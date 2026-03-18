package suggestion

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// --- SessionRepo (PostgreSQL) ---

// SessionRepo implements SessionRepository using PostgreSQL.
type SessionRepo struct {
	pool *pgxpool.Pool
}

// NewSessionRepo creates a new SessionRepo.
func NewSessionRepo(pool *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{pool: pool}
}

// Create inserts a new suggestion session.
func (r *SessionRepo) Create(ctx context.Context, session *SuggestionSession) error {
	inputJSON, err := json.Marshal(session.InputParams)
	if err != nil {
		return fmt.Errorf("suggestion.SessionRepo.Create: marshal input_params: %w", err)
	}

	query := `
		INSERT INTO suggestion_sessions (id, user_id, session_type, input_params, result_recipe_ids, total_calories, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = r.pool.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.SessionType,
		inputJSON,
		session.ResultRecipeIDs,
		session.TotalCalories,
		session.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("suggestion.SessionRepo.Create: %w", err)
	}
	return nil
}

// ListByUser returns paginated suggestion sessions for a user, optionally filtered by type.
func (r *SessionRepo) ListByUser(ctx context.Context, userID uuid.UUID, sessionType string, page, pageSize int) ([]SuggestionSession, int64, error) {
	offset := (page - 1) * pageSize

	// Count query.
	countQuery := `SELECT COUNT(*) FROM suggestion_sessions WHERE user_id = $1`
	countArgs := []any{userID}

	if sessionType != "" {
		countQuery += ` AND session_type = $2`
		countArgs = append(countArgs, sessionType)
	}

	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("suggestion.SessionRepo.ListByUser: count: %w", err)
	}

	// Data query.
	dataQuery := `
		SELECT id, user_id, session_type, input_params, result_recipe_ids, total_calories, created_at
		FROM suggestion_sessions
		WHERE user_id = $1`
	dataArgs := []any{userID}

	if sessionType != "" {
		dataQuery += ` AND session_type = $2`
		dataArgs = append(dataArgs, sessionType)
	}

	dataQuery += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(dataArgs)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(dataArgs)+2)
	dataArgs = append(dataArgs, pageSize, offset)

	rows, err := r.pool.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("suggestion.SessionRepo.ListByUser: query: %w", err)
	}
	defer rows.Close()

	sessions := make([]SuggestionSession, 0)
	for rows.Next() {
		var s SuggestionSession
		var inputJSON []byte

		if err := rows.Scan(&s.ID, &s.UserID, &s.SessionType, &inputJSON, &s.ResultRecipeIDs, &s.TotalCalories, &s.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("suggestion.SessionRepo.ListByUser: scan: %w", err)
		}
		if inputJSON != nil {
			if err := json.Unmarshal(inputJSON, &s.InputParams); err != nil {
				return nil, 0, fmt.Errorf("suggestion.SessionRepo.ListByUser: unmarshal input_params: %w", err)
			}
		}
		sessions = append(sessions, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("suggestion.SessionRepo.ListByUser: rows: %w", err)
	}

	return sessions, total, nil
}

// --- ConfigRepo (PostgreSQL) ---

// ConfigRepo implements ConfigRepository using PostgreSQL.
type ConfigRepo struct {
	pool *pgxpool.Pool
}

// NewConfigRepo creates a new ConfigRepo.
func NewConfigRepo(pool *pgxpool.Pool) *ConfigRepo {
	return &ConfigRepo{pool: pool}
}

// ListByType returns active suggestion configs filtered by type, ordered by sort_order.
func (r *ConfigRepo) ListByType(ctx context.Context, configType string) ([]SuggestionConfig, error) {
	query := `
		SELECT id, config_type, name, description, params, is_active, sort_order, created_at, updated_at
		FROM suggestion_configs
		WHERE is_active = true`
	args := make([]any, 0)

	if configType != "" {
		query += ` AND config_type = $1`
		args = append(args, configType)
	}
	query += ` ORDER BY sort_order ASC`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("suggestion.ConfigRepo.ListByType: %w", err)
	}
	defer rows.Close()

	configs := make([]SuggestionConfig, 0)
	for rows.Next() {
		var c SuggestionConfig
		var paramsJSON []byte

		if err := rows.Scan(&c.ID, &c.ConfigType, &c.Name, &c.Description, &paramsJSON, &c.IsActive, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("suggestion.ConfigRepo.ListByType: scan: %w", err)
		}
		if paramsJSON != nil {
			if err := json.Unmarshal(paramsJSON, &c.Params); err != nil {
				return nil, fmt.Errorf("suggestion.ConfigRepo.ListByType: unmarshal params: %w", err)
			}
		}
		configs = append(configs, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("suggestion.ConfigRepo.ListByType: rows: %w", err)
	}

	return configs, nil
}

// --- ExclusionRepo (PostgreSQL) ---

// ExclusionRepo implements ExclusionRepository using PostgreSQL.
type ExclusionRepo struct {
	pool *pgxpool.Pool
}

// NewExclusionRepo creates a new ExclusionRepo.
func NewExclusionRepo(pool *pgxpool.Pool) *ExclusionRepo {
	return &ExclusionRepo{pool: pool}
}

// GetExcludedRecipeIDs returns the recipe IDs currently excluded for a user.
func (r *ExclusionRepo) GetExcludedRecipeIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT recipe_id
		FROM suggestion_exclusions
		WHERE user_id = $1 AND excluded_until > $2`

	rows, err := r.pool.Query(ctx, query, userID, time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("suggestion.ExclusionRepo.GetExcludedRecipeIDs: %w", err)
	}
	defer rows.Close()

	ids, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return nil, fmt.Errorf("suggestion.ExclusionRepo.GetExcludedRecipeIDs: collect: %w", err)
	}
	return ids, nil
}

// Add inserts a new exclusion rule.
func (r *ExclusionRepo) Add(ctx context.Context, rule *ExclusionRule) error {
	query := `
		INSERT INTO suggestion_exclusions (id, user_id, recipe_id, excluded_until, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, recipe_id) DO UPDATE SET excluded_until = EXCLUDED.excluded_until`

	_, err := r.pool.Exec(ctx, query,
		rule.ID,
		rule.UserID,
		rule.RecipeID,
		rule.ExcludedUntil,
		rule.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("suggestion.ExclusionRepo.Add: %w", err)
	}
	return nil
}

// CleanExpired removes exclusion rules that have passed their excluded_until time.
func (r *ExclusionRepo) CleanExpired(ctx context.Context) error {
	query := `DELETE FROM suggestion_exclusions WHERE excluded_until <= $1`

	_, err := r.pool.Exec(ctx, query, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("suggestion.ExclusionRepo.CleanExpired: %w", err)
	}
	return nil
}
