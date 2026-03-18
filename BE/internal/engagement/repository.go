package engagement

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrors "github.com/lehoanglong/whatdish/internal/shared/errors"
)

// FavoriteRepo is the PostgreSQL implementation of FavoriteRepository.
type FavoriteRepo struct {
	pool *pgxpool.Pool
}

// NewFavoriteRepo creates a new FavoriteRepo.
func NewFavoriteRepo(pool *pgxpool.Pool) *FavoriteRepo {
	return &FavoriteRepo{pool: pool}
}

// Add inserts a favorite. Returns a conflict error if the favorite already exists.
func (r *FavoriteRepo) Add(ctx context.Context, userID, recipeID uuid.UUID) error {
	query := `
		INSERT INTO engagement_favorites (id, user_id, recipe_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, recipe_id) DO NOTHING`

	_, err := r.pool.Exec(ctx, query, uuid.New(), userID, recipeID, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("FavoriteRepo.Add: %w", err)
	}

	slog.Debug("favorite added", "user_id", userID, "recipe_id", recipeID)

	return nil
}

// Remove deletes a favorite by user and recipe ID.
func (r *FavoriteRepo) Remove(ctx context.Context, userID, recipeID uuid.UUID) error {
	query := `DELETE FROM engagement_favorites WHERE user_id = $1 AND recipe_id = $2`

	tag, err := r.pool.Exec(ctx, query, userID, recipeID)
	if err != nil {
		return fmt.Errorf("FavoriteRepo.Remove: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("FavoriteRepo.Remove: %w", apperrors.ErrNotFound)
	}

	slog.Debug("favorite removed", "user_id", userID, "recipe_id", recipeID)

	return nil
}

// ListByUser returns a paginated list of a user's favorites.
func (r *FavoriteRepo) ListByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]Favorite, int64, error) {
	// Count total.
	var total int64
	countQuery := `SELECT COUNT(*) FROM engagement_favorites WHERE user_id = $1`
	if err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("FavoriteRepo.ListByUser: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT id, user_id, recipe_id, created_at
		FROM engagement_favorites
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("FavoriteRepo.ListByUser: %w", err)
	}
	defer rows.Close()

	var favorites []Favorite
	for rows.Next() {
		var f Favorite
		if err := rows.Scan(&f.ID, &f.UserID, &f.RecipeID, &f.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("FavoriteRepo.ListByUser: %w", err)
		}
		favorites = append(favorites, f)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("FavoriteRepo.ListByUser: %w", err)
	}

	return favorites, total, nil
}

// Check returns a map indicating which of the given recipe IDs are favorited.
func (r *FavoriteRepo) Check(ctx context.Context, userID uuid.UUID, recipeIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	if len(recipeIDs) == 0 {
		return make(map[uuid.UUID]bool), nil
	}

	query := `
		SELECT recipe_id
		FROM engagement_favorites
		WHERE user_id = $1 AND recipe_id = ANY($2)`

	rows, err := r.pool.Query(ctx, query, userID, recipeIDs)
	if err != nil {
		return nil, fmt.Errorf("FavoriteRepo.Check: %w", err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID]bool, len(recipeIDs))
	for _, id := range recipeIDs {
		result[id] = false
	}

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("FavoriteRepo.Check: %w", err)
		}
		result[id] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("FavoriteRepo.Check: %w", err)
	}

	return result, nil
}

// ViewRepo is the PostgreSQL implementation of ViewRepository.
type ViewRepo struct {
	pool *pgxpool.Pool
}

// NewViewRepo creates a new ViewRepo.
func NewViewRepo(pool *pgxpool.Pool) *ViewRepo {
	return &ViewRepo{pool: pool}
}

// Record inserts a view history entry.
func (r *ViewRepo) Record(ctx context.Context, view *ViewHistory) error {
	query := `
		INSERT INTO engagement_views (id, user_id, session_id, recipe_id, source, viewed_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.pool.Exec(ctx, query,
		view.ID, view.UserID, view.SessionID, view.RecipeID, view.Source, view.ViewedAt,
	)
	if err != nil {
		return fmt.Errorf("ViewRepo.Record: %w", err)
	}

	return nil
}

// RatingRepo is the PostgreSQL implementation of RatingRepository.
type RatingRepo struct {
	pool *pgxpool.Pool
}

// NewRatingRepo creates a new RatingRepo.
func NewRatingRepo(pool *pgxpool.Pool) *RatingRepo {
	return &RatingRepo{pool: pool}
}

// Upsert creates or updates a rating.
func (r *RatingRepo) Upsert(ctx context.Context, rating *Rating) error {
	query := `
		INSERT INTO engagement_ratings (id, user_id, recipe_id, score, comment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, recipe_id) DO UPDATE SET
			score = EXCLUDED.score,
			comment = EXCLUDED.comment,
			updated_at = EXCLUDED.updated_at`

	_, err := r.pool.Exec(ctx, query,
		rating.ID, rating.UserID, rating.RecipeID, rating.Score, rating.Comment,
		rating.CreatedAt, rating.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("RatingRepo.Upsert: %w", err)
	}

	slog.Debug("rating upserted", "user_id", rating.UserID, "recipe_id", rating.RecipeID, "score", rating.Score)

	return nil
}

// GetByRecipe returns all ratings for a recipe.
func (r *RatingRepo) GetByRecipe(ctx context.Context, recipeID uuid.UUID) ([]Rating, error) {
	query := `
		SELECT id, user_id, recipe_id, score, comment, created_at, updated_at
		FROM engagement_ratings
		WHERE recipe_id = $1
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, recipeID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("RatingRepo.GetByRecipe: %w", err)
	}
	defer rows.Close()

	var ratings []Rating
	for rows.Next() {
		var rt Rating
		if err := rows.Scan(&rt.ID, &rt.UserID, &rt.RecipeID, &rt.Score, &rt.Comment, &rt.CreatedAt, &rt.UpdatedAt); err != nil {
			return nil, fmt.Errorf("RatingRepo.GetByRecipe: %w", err)
		}
		ratings = append(ratings, rt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("RatingRepo.GetByRecipe: %w", err)
	}

	return ratings, nil
}
