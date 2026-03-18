package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrors "github.com/lehoanglong/whatdish/internal/shared/errors"
)

// ProfileRepo is the PostgreSQL implementation of ProfileRepository.
type ProfileRepo struct {
	pool *pgxpool.Pool
}

// NewProfileRepo creates a new ProfileRepo.
func NewProfileRepo(pool *pgxpool.Pool) *ProfileRepo {
	return &ProfileRepo{pool: pool}
}

// GetByUserID returns the user profile for the given user ID.
func (r *ProfileRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	query := `
		SELECT id, user_id, gender, age, height_cm, weight_kg,
		       activity_level, goal, bmr, tdee, daily_target,
		       created_at, updated_at
		FROM user_profiles
		WHERE user_id = $1`

	var p UserProfile
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&p.ID, &p.UserID, &p.Gender, &p.Age, &p.HeightCm, &p.WeightKg,
		&p.ActivityLevel, &p.Goal, &p.BMR, &p.TDEE, &p.DailyTarget,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("ProfileRepo.GetByUserID: %w", apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("ProfileRepo.GetByUserID: %w", err)
	}

	return &p, nil
}

// Upsert inserts or updates a user profile.
func (r *ProfileRepo) Upsert(ctx context.Context, profile *UserProfile) error {
	query := `
		INSERT INTO user_profiles (id, user_id, gender, age, height_cm, weight_kg,
		                           activity_level, goal, bmr, tdee, daily_target,
		                           created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (user_id) DO UPDATE SET
			gender = EXCLUDED.gender,
			age = EXCLUDED.age,
			height_cm = EXCLUDED.height_cm,
			weight_kg = EXCLUDED.weight_kg,
			activity_level = EXCLUDED.activity_level,
			goal = EXCLUDED.goal,
			bmr = EXCLUDED.bmr,
			tdee = EXCLUDED.tdee,
			daily_target = EXCLUDED.daily_target,
			updated_at = EXCLUDED.updated_at`

	_, err := r.pool.Exec(ctx, query,
		profile.ID, profile.UserID, profile.Gender, profile.Age,
		profile.HeightCm, profile.WeightKg, profile.ActivityLevel, profile.Goal,
		profile.BMR, profile.TDEE, profile.DailyTarget,
		profile.CreatedAt, profile.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("ProfileRepo.Upsert: %w", err)
	}

	slog.Debug("profile upserted", "user_id", profile.UserID)

	return nil
}

// AllergyRepo is the PostgreSQL implementation of AllergyRepository.
type AllergyRepo struct {
	pool *pgxpool.Pool
}

// NewAllergyRepo creates a new AllergyRepo.
func NewAllergyRepo(pool *pgxpool.Pool) *AllergyRepo {
	return &AllergyRepo{pool: pool}
}

// ListByUser returns all allergies for the given user.
func (r *AllergyRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]UserAllergy, error) {
	query := `
		SELECT id, user_id, ingredient_id, ingredient_name, allergy_type
		FROM user_allergies
		WHERE user_id = $1
		ORDER BY ingredient_name`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("AllergyRepo.ListByUser: %w", err)
	}
	defer rows.Close()

	var allergies []UserAllergy
	for rows.Next() {
		var a UserAllergy
		if err := rows.Scan(&a.ID, &a.UserID, &a.IngredientID, &a.IngredientName, &a.AllergyType); err != nil {
			return nil, fmt.Errorf("AllergyRepo.ListByUser: %w", err)
		}
		allergies = append(allergies, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("AllergyRepo.ListByUser: %w", err)
	}

	return allergies, nil
}

// Set replaces all allergies for a user with the provided list.
// This deletes existing allergies and inserts the new set in a transaction.
func (r *AllergyRepo) Set(ctx context.Context, userID uuid.UUID, allergies []UserAllergy) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("AllergyRepo.Set: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err.Error() != "tx is closed" {
			slog.Error("AllergyRepo.Set: rollback failed", "error", err)
		}
	}()

	// Delete existing allergies.
	_, err = tx.Exec(ctx, `DELETE FROM user_allergies WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("AllergyRepo.Set: %w", err)
	}

	// Insert new allergies.
	for _, a := range allergies {
		_, err = tx.Exec(ctx,
			`INSERT INTO user_allergies (id, user_id, ingredient_id, ingredient_name, allergy_type)
			 VALUES ($1, $2, $3, $4, $5)`,
			a.ID, a.UserID, a.IngredientID, a.IngredientName, a.AllergyType,
		)
		if err != nil {
			return fmt.Errorf("AllergyRepo.Set: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("AllergyRepo.Set: %w", err)
	}

	slog.Debug("allergies set", "user_id", userID, "count", len(allergies))

	return nil
}
