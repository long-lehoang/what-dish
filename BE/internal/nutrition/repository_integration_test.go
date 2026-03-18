package nutrition_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lehoanglong/whatdish/internal/nutrition"
)

func seedRecipeForNutrition(t *testing.T, name string) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	var id uuid.UUID
	err := testDB.Pool.QueryRow(ctx,
		`INSERT INTO recipes (name, slug, status) VALUES ($1, $2, 'PUBLISHED') RETURNING id`,
		name, "slug-"+uuid.New().String()[:8],
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func TestIntegration_NutritionRepo_UpsertAndGet(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := nutrition.NewNutritionRepo(testDB.Pool)

	recipeID := seedRecipeForNutrition(t, "Phở Bò")
	cal := 450.0
	protein := 30.0
	carbs := 55.0
	fat := 10.0

	n := &nutrition.RecipeNutrition{
		ID:        uuid.New(),
		RecipeID:  recipeID,
		Calories:  &cal,
		Protein:   &protein,
		Carbs:     &carbs,
		Fat:       &fat,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := repo.Upsert(ctx, n)
	require.NoError(t, err)

	got, err := repo.GetByRecipeID(ctx, recipeID)
	require.NoError(t, err)
	assert.InDelta(t, 450.0, *got.Calories, 0.01)
	assert.InDelta(t, 30.0, *got.Protein, 0.01)

	newCal := 500.0
	n.Calories = &newCal
	n.UpdatedAt = time.Now().UTC()
	err = repo.Upsert(ctx, n)
	require.NoError(t, err)

	got, err = repo.GetByRecipeID(ctx, recipeID)
	require.NoError(t, err)
	assert.InDelta(t, 500.0, *got.Calories, 0.01)
}

func TestIntegration_NutritionRepo_GetByRecipeID_NotFound(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := nutrition.NewNutritionRepo(testDB.Pool)

	_, err := repo.GetByRecipeID(ctx, uuid.New())
	assert.Error(t, err)
}

func TestIntegration_NutritionRepo_GetByRecipeIDs(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := nutrition.NewNutritionRepo(testDB.Pool)

	r1 := seedRecipeForNutrition(t, "Dish 1")
	r2 := seedRecipeForNutrition(t, "Dish 2")
	r3 := seedRecipeForNutrition(t, "Dish 3")

	cal1 := 400.0
	cal2 := 600.0
	now := time.Now().UTC()

	require.NoError(t, repo.Upsert(ctx, &nutrition.RecipeNutrition{
		ID: uuid.New(), RecipeID: r1, Calories: &cal1, CreatedAt: now, UpdatedAt: now,
	}))
	require.NoError(t, repo.Upsert(ctx, &nutrition.RecipeNutrition{
		ID: uuid.New(), RecipeID: r2, Calories: &cal2, CreatedAt: now, UpdatedAt: now,
	}))

	results, err := repo.GetByRecipeIDs(ctx, []uuid.UUID{r1, r2, r3})
	require.NoError(t, err)
	assert.Len(t, results, 2)

	results, err = repo.GetByRecipeIDs(ctx, nil)
	require.NoError(t, err)
	assert.Nil(t, results)
}

func TestIntegration_NutritionRepo_GetIDsByCalorieRange(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := nutrition.NewNutritionRepo(testDB.Pool)

	r1 := seedRecipeForNutrition(t, "Low Cal")
	r2 := seedRecipeForNutrition(t, "Mid Cal")
	r3 := seedRecipeForNutrition(t, "High Cal")

	cal1 := 200.0
	cal2 := 500.0
	cal3 := 900.0
	now := time.Now().UTC()

	require.NoError(t, repo.Upsert(ctx, &nutrition.RecipeNutrition{
		ID: uuid.New(), RecipeID: r1, Calories: &cal1, CreatedAt: now, UpdatedAt: now,
	}))
	require.NoError(t, repo.Upsert(ctx, &nutrition.RecipeNutrition{
		ID: uuid.New(), RecipeID: r2, Calories: &cal2, CreatedAt: now, UpdatedAt: now,
	}))
	require.NoError(t, repo.Upsert(ctx, &nutrition.RecipeNutrition{
		ID: uuid.New(), RecipeID: r3, Calories: &cal3, CreatedAt: now, UpdatedAt: now,
	}))

	ids, err := repo.GetIDsByCalorieRange(ctx, 400, 600)
	require.NoError(t, err)
	assert.Len(t, ids, 1)
	assert.Equal(t, r2, ids[0])

	ids, err = repo.GetIDsByCalorieRange(ctx, 100, 1000)
	require.NoError(t, err)
	assert.Len(t, ids, 3)
}

func TestIntegration_GoalRepo_List(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := nutrition.NewGoalRepo(testDB.Pool)

	_, err := testDB.Pool.Exec(ctx,
		`INSERT INTO nutrition_goals (name, description, is_active, sort_order)
		 VALUES ('Lose Weight', 'Cut calories', true, 1),
		        ('Maintain', 'Stay the same', true, 2),
		        ('Inactive Goal', 'Hidden', false, 3)`)
	require.NoError(t, err)

	goals, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, goals, 2)
	assert.Equal(t, "Lose Weight", goals[0].Name)
	assert.Equal(t, "Maintain", goals[1].Name)
}
