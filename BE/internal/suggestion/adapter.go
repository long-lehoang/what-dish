package suggestion

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DishReaderAdapter implements DishReader by querying the recipes table directly.
// This is a pragmatic choice for a monolith — in a microservice world, this would
// be an HTTP client calling the recipe service.
type DishReaderAdapter struct {
	pool *pgxpool.Pool
}

func NewDishReaderAdapter(pool *pgxpool.Pool) *DishReaderAdapter {
	return &DishReaderAdapter{pool: pool}
}

func (a *DishReaderAdapter) GetAllPublishedIDs(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := a.pool.Query(ctx,
		"SELECT id FROM recipes WHERE status = 'PUBLISHED' AND deleted_at IS NULL")
	if err != nil {
		return nil, fmt.Errorf("DishReaderAdapter.GetAllPublishedIDs: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("DishReaderAdapter.GetAllPublishedIDs scan: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (a *DishReaderAdapter) GetByID(ctx context.Context, id uuid.UUID) (*RecipeSummary, error) {
	var s RecipeSummary
	err := a.pool.QueryRow(ctx, `
		SELECT r.id, r.name, r.slug, COALESCE(r.image_url, ''), COALESCE(r.difficulty, ''), COALESCE(r.cook_time, 0), n.calories
		FROM recipes r
		LEFT JOIN nutrition_recipe n ON n.recipe_id = r.id
		WHERE r.id = $1 AND r.deleted_at IS NULL
	`, id).Scan(&s.ID, &s.Name, &s.Slug, &s.ImageURL, &s.Difficulty, &s.CookTime, &s.Calories)
	if err != nil {
		return nil, fmt.Errorf("DishReaderAdapter.GetByID: %w", err)
	}
	return &s, nil
}

func (a *DishReaderAdapter) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]RecipeSummary, error) {
	if len(ids) == 0 {
		return []RecipeSummary{}, nil
	}

	rows, err := a.pool.Query(ctx, `
		SELECT r.id, r.name, r.slug, COALESCE(r.image_url, ''), COALESCE(r.difficulty, ''), COALESCE(r.cook_time, 0), n.calories
		FROM recipes r
		LEFT JOIN nutrition_recipe n ON n.recipe_id = r.id
		WHERE r.id = ANY($1) AND r.deleted_at IS NULL
	`, ids)
	if err != nil {
		return nil, fmt.Errorf("DishReaderAdapter.GetByIDs: %w", err)
	}
	defer rows.Close()

	var results []RecipeSummary
	for rows.Next() {
		var s RecipeSummary
		if err := rows.Scan(&s.ID, &s.Name, &s.Slug, &s.ImageURL, &s.Difficulty, &s.CookTime, &s.Calories); err != nil {
			return nil, fmt.Errorf("DishReaderAdapter.GetByIDs scan: %w", err)
		}
		results = append(results, s)
	}
	return results, rows.Err()
}

func (a *DishReaderAdapter) GetIDsByFilter(ctx context.Context, filter RecipeFilter) ([]uuid.UUID, error) {
	query := "SELECT id FROM recipes WHERE status = 'PUBLISHED' AND deleted_at IS NULL"
	var args []any
	argIdx := 1

	if filter.DishTypeID != nil {
		query += fmt.Sprintf(" AND dish_type_id = $%d", argIdx)
		args = append(args, *filter.DishTypeID)
		argIdx++
	}
	if filter.RegionID != nil {
		query += fmt.Sprintf(" AND region_id = $%d", argIdx)
		args = append(args, *filter.RegionID)
		argIdx++
	}
	if filter.MainIngredientID != nil {
		query += fmt.Sprintf(" AND main_ingredient_id = $%d", argIdx)
		args = append(args, *filter.MainIngredientID)
		argIdx++
	}
	if filter.MealTypeID != nil {
		query += fmt.Sprintf(" AND meal_type_id = $%d", argIdx)
		args = append(args, *filter.MealTypeID)
		argIdx++
	}
	if filter.Difficulty != nil {
		query += fmt.Sprintf(" AND difficulty = $%d", argIdx)
		args = append(args, *filter.Difficulty)
		argIdx++
	}
	if filter.MaxCookTime != nil {
		query += fmt.Sprintf(" AND cook_time <= $%d", argIdx)
		args = append(args, *filter.MaxCookTime)
	}

	rows, err := a.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("DishReaderAdapter.GetIDsByFilter: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("DishReaderAdapter.GetIDsByFilter scan: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// CalorieProviderAdapter implements CalorieProvider by querying the nutrition_recipe table.
type CalorieProviderAdapter struct {
	pool *pgxpool.Pool
}

func NewCalorieProviderAdapter(pool *pgxpool.Pool) *CalorieProviderAdapter {
	return &CalorieProviderAdapter{pool: pool}
}

func (a *CalorieProviderAdapter) GetRecipeIDsByCalorieRange(ctx context.Context, min, max float64) ([]uuid.UUID, error) {
	rows, err := a.pool.Query(ctx,
		"SELECT recipe_id FROM nutrition_recipe WHERE calories >= $1 AND calories <= $2",
		min, max)
	if err != nil {
		return nil, fmt.Errorf("CalorieProviderAdapter.GetRecipeIDsByCalorieRange: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("CalorieProviderAdapter.GetRecipeIDsByCalorieRange scan: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (a *CalorieProviderAdapter) GetCalories(ctx context.Context, recipeID uuid.UUID) (*float64, error) {
	var cal *float64
	err := a.pool.QueryRow(ctx,
		"SELECT calories FROM nutrition_recipe WHERE recipe_id = $1",
		recipeID).Scan(&cal)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // no nutrition data is not an error
		}
		return nil, fmt.Errorf("suggestion.CalorieProviderAdapter.GetCalories: %w", err)
	}
	return cal, nil
}
