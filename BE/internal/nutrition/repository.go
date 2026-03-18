package nutrition

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrors "github.com/lehoanglong/whatdish/internal/shared/errors"
)

// NutritionRepo is the PostgreSQL implementation of NutritionRepository.
type NutritionRepo struct {
	pool *pgxpool.Pool
}

// NewNutritionRepo creates a new NutritionRepo.
func NewNutritionRepo(pool *pgxpool.Pool) *NutritionRepo {
	return &NutritionRepo{pool: pool}
}

// GetByRecipeID returns nutrition data for a single recipe.
func (r *NutritionRepo) GetByRecipeID(ctx context.Context, recipeID uuid.UUID) (*RecipeNutrition, error) {
	query := `
		SELECT id, recipe_id, calories, protein, carbs, fat, fiber, sugar, sodium,
		       serving_size, data_source, is_verified, created_at, updated_at
		FROM nutrition_recipe
		WHERE recipe_id = $1`

	var n RecipeNutrition
	err := r.pool.QueryRow(ctx, query, recipeID).Scan(
		&n.ID, &n.RecipeID, &n.Calories, &n.Protein, &n.Carbs, &n.Fat,
		&n.Fiber, &n.Sugar, &n.Sodium, &n.ServingSize, &n.DataSource,
		&n.IsVerified, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("NutritionRepo.GetByRecipeID: %w", apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("NutritionRepo.GetByRecipeID: %w", err)
	}

	return &n, nil
}

// GetByRecipeIDs returns nutrition data for multiple recipes.
func (r *NutritionRepo) GetByRecipeIDs(ctx context.Context, recipeIDs []uuid.UUID) ([]RecipeNutrition, error) {
	if len(recipeIDs) == 0 {
		return nil, nil
	}

	query := `
		SELECT id, recipe_id, calories, protein, carbs, fat, fiber, sugar, sodium,
		       serving_size, data_source, is_verified, created_at, updated_at
		FROM nutrition_recipe
		WHERE recipe_id = ANY($1)`

	rows, err := r.pool.Query(ctx, query, recipeIDs)
	if err != nil {
		return nil, fmt.Errorf("NutritionRepo.GetByRecipeIDs: %w", err)
	}
	defer rows.Close()

	var items []RecipeNutrition
	for rows.Next() {
		var n RecipeNutrition
		if err := rows.Scan(
			&n.ID, &n.RecipeID, &n.Calories, &n.Protein, &n.Carbs, &n.Fat,
			&n.Fiber, &n.Sugar, &n.Sodium, &n.ServingSize, &n.DataSource,
			&n.IsVerified, &n.CreatedAt, &n.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("NutritionRepo.GetByRecipeIDs: %w", err)
		}
		items = append(items, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("NutritionRepo.GetByRecipeIDs: %w", err)
	}

	return items, nil
}

// Upsert inserts or updates nutrition data for a recipe.
func (r *NutritionRepo) Upsert(ctx context.Context, nutrition *RecipeNutrition) error {
	query := `
		INSERT INTO nutrition_recipe (id, recipe_id, calories, protein, carbs, fat, fiber,
		                              sugar, sodium, serving_size, data_source, is_verified,
		                              created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (recipe_id) DO UPDATE SET
			calories = EXCLUDED.calories,
			protein = EXCLUDED.protein,
			carbs = EXCLUDED.carbs,
			fat = EXCLUDED.fat,
			fiber = EXCLUDED.fiber,
			sugar = EXCLUDED.sugar,
			sodium = EXCLUDED.sodium,
			serving_size = EXCLUDED.serving_size,
			data_source = EXCLUDED.data_source,
			is_verified = EXCLUDED.is_verified,
			updated_at = EXCLUDED.updated_at`

	_, err := r.pool.Exec(ctx, query,
		nutrition.ID, nutrition.RecipeID, nutrition.Calories, nutrition.Protein,
		nutrition.Carbs, nutrition.Fat, nutrition.Fiber, nutrition.Sugar, nutrition.Sodium,
		nutrition.ServingSize, nutrition.DataSource, nutrition.IsVerified,
		nutrition.CreatedAt, nutrition.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("NutritionRepo.Upsert: %w", err)
	}

	slog.Debug("nutrition upserted", "recipe_id", nutrition.RecipeID)

	return nil
}

// GetIDsByCalorieRange returns recipe IDs within the specified calorie range.
func (r *NutritionRepo) GetIDsByCalorieRange(ctx context.Context, min, max float64) ([]uuid.UUID, error) {
	query := `
		SELECT recipe_id
		FROM nutrition_recipe
		WHERE calories >= $1 AND calories <= $2`

	rows, err := r.pool.Query(ctx, query, min, max)
	if err != nil {
		return nil, fmt.Errorf("NutritionRepo.GetIDsByCalorieRange: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("NutritionRepo.GetIDsByCalorieRange: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("NutritionRepo.GetIDsByCalorieRange: %w", err)
	}

	return ids, nil
}

// GoalRepo is the PostgreSQL implementation of GoalRepository.
type GoalRepo struct {
	pool *pgxpool.Pool
}

// NewGoalRepo creates a new GoalRepo.
func NewGoalRepo(pool *pgxpool.Pool) *GoalRepo {
	return &GoalRepo{pool: pool}
}

// List returns all active nutrition goals ordered by sort_order.
func (r *GoalRepo) List(ctx context.Context) ([]NutritionGoal, error) {
	query := `
		SELECT id, name, description, meal_calories_min, meal_calories_max,
		       daily_calories_min, daily_calories_max, protein_pct, carbs_pct, fat_pct,
		       is_active, sort_order
		FROM nutrition_goals
		WHERE is_active = true
		ORDER BY sort_order`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GoalRepo.List: %w", err)
	}
	defer rows.Close()

	var goals []NutritionGoal
	for rows.Next() {
		var g NutritionGoal
		if err := rows.Scan(
			&g.ID, &g.Name, &g.Description, &g.MealCaloriesMin, &g.MealCaloriesMax,
			&g.DailyCaloriesMin, &g.DailyCaloriesMax, &g.ProteinPct, &g.CarbsPct, &g.FatPct,
			&g.IsActive, &g.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("GoalRepo.List: %w", err)
		}
		goals = append(goals, g)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GoalRepo.List: %w", err)
	}

	return goals, nil
}
