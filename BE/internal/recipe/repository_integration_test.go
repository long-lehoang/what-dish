package recipe_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lehoanglong/whatdish/internal/recipe"
)

// seedCategory inserts a category and returns it.
func seedCategory(t *testing.T, name, slug, catType string) recipe.Category {
	t.Helper()
	ctx := context.Background()
	var cat recipe.Category
	err := testDB.Pool.QueryRow(ctx,
		`INSERT INTO categories (name, slug, type, sort_order, is_active)
		 VALUES ($1, $2, $3, 0, true)
		 RETURNING id, name, slug, type, icon_url, sort_order, is_active, created_at, updated_at`,
		name, slug, catType,
	).Scan(&cat.ID, &cat.Name, &cat.Slug, &cat.Type, &cat.IconURL,
		&cat.SortOrder, &cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt)
	require.NoError(t, err)
	return cat
}

// seedTag inserts a tag and returns it.
func seedTag(t *testing.T, name, slug string) recipe.Tag {
	t.Helper()
	ctx := context.Background()
	var tag recipe.Tag
	err := testDB.Pool.QueryRow(ctx,
		`INSERT INTO tags (name, slug) VALUES ($1, $2) RETURNING id, name, slug`,
		name, slug,
	).Scan(&tag.ID, &tag.Name, &tag.Slug)
	require.NoError(t, err)
	return tag
}

// seedRecipe inserts a published recipe and returns its ID.
func seedRecipe(t *testing.T, name, slug string, dishTypeID *uuid.UUID) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	var id uuid.UUID
	err := testDB.Pool.QueryRow(ctx,
		`INSERT INTO recipes (name, slug, status, difficulty, servings, dish_type_id)
		 VALUES ($1, $2, 'PUBLISHED', 'EASY', 2, $3) RETURNING id`,
		name, slug, dishTypeID,
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func TestIntegration_DishRepo_List(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	// Empty state.
	dishes, total, err := repo.List(ctx, recipe.DishFilter{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, dishes)

	// Seed recipes.
	seedRecipe(t, "Phở Bò", "pho-bo", nil)
	seedRecipe(t, "Bún Chả", "bun-cha", nil)
	seedRecipe(t, "Cơm Tấm", "com-tam", nil)

	dishes, total, err = repo.List(ctx, recipe.DishFilter{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, dishes, 3)
}

func TestIntegration_DishRepo_ListWithFilter(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	cat := seedCategory(t, "Noodle Soup", "noodle-soup", "DISH_TYPE")
	seedRecipe(t, "Phở Bò", "pho-bo", &cat.ID)
	seedRecipe(t, "Bún Chả", "bun-cha", nil)

	dishes, total, err := repo.List(ctx, recipe.DishFilter{
		DishTypeID: &cat.ID,
		Page:       1,
		PageSize:   10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "Phở Bò", dishes[0].Name)
}

func TestIntegration_DishRepo_ListPagination(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	for i := 0; i < 5; i++ {
		seedRecipe(t, "Dish "+uuid.New().String()[:4], "dish-"+uuid.New().String()[:8], nil)
	}

	dishes, total, err := repo.List(ctx, recipe.DishFilter{Page: 1, PageSize: 2})
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, dishes, 2)

	dishes2, _, err := repo.List(ctx, recipe.DishFilter{Page: 2, PageSize: 2})
	require.NoError(t, err)
	assert.Len(t, dishes2, 2)
	assert.NotEqual(t, dishes[0].ID, dishes2[0].ID)
}

func TestIntegration_DishRepo_GetByID(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	id := seedRecipe(t, "Phở Bò", "pho-bo", nil)

	detail, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Phở Bò", detail.Name)
	assert.Equal(t, "pho-bo", detail.Slug)
	assert.NotNil(t, detail.Ingredients)
	assert.NotNil(t, detail.Steps)
	assert.NotNil(t, detail.Tags)
}

func TestIntegration_DishRepo_GetByID_NotFound(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	_, err := repo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
}

func TestIntegration_DishRepo_GetBySlug(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	seedRecipe(t, "Phở Bò", "pho-bo", nil)

	detail, err := repo.GetBySlug(ctx, "pho-bo")
	require.NoError(t, err)
	assert.Equal(t, "Phở Bò", detail.Name)
}

func TestIntegration_DishRepo_GetRandom(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	seedRecipe(t, "Phở Bò", "pho-bo", nil)
	seedRecipe(t, "Bún Chả", "bun-cha", nil)

	detail, err := repo.GetRandom(ctx, recipe.DishFilter{Page: 1, PageSize: 10}, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, detail.Name)
}

func TestIntegration_DishRepo_GetRandom_WithExclusion(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	id1 := seedRecipe(t, "Phở Bò", "pho-bo", nil)
	id2 := seedRecipe(t, "Bún Chả", "bun-cha", nil)

	detail, err := repo.GetRandom(ctx, recipe.DishFilter{Page: 1, PageSize: 10}, []uuid.UUID{id1})
	require.NoError(t, err)
	assert.Equal(t, id2, detail.ID)
}

func TestIntegration_DishRepo_Search(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	seedRecipe(t, "Phở Bò", "pho-bo", nil)
	seedRecipe(t, "Bún Chả", "bun-cha", nil)

	dishes, total, err := repo.Search(ctx, "pho", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "Phở Bò", dishes[0].Name)
}

func TestIntegration_DishRepo_UpsertFromSync(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	extID := "notion-page-123"
	servings := 4
	difficulty := "MEDIUM"
	dish := recipe.Dish{
		ExternalID: &extID,
		Name:       "Canh Chua",
		Slug:       "canh-chua",
		Status:     "PUBLISHED",
		Servings:   &servings,
		Difficulty: &difficulty,
	}

	ingredients := []recipe.Ingredient{
		{Name: "Cá lóc", SortOrder: 0},
		{Name: "Cà chua", SortOrder: 1},
	}
	steps := []recipe.Step{
		{StepNumber: 1, Description: "Sơ chế cá", SortOrder: 1},
		{StepNumber: 2, Description: "Nấu canh", SortOrder: 2},
	}

	// First upsert — insert.
	isNew, err := repo.UpsertFromSync(ctx, &dish, ingredients, steps, nil)
	require.NoError(t, err)
	assert.True(t, isNew)

	detail, err := repo.GetBySlug(ctx, "canh-chua")
	require.NoError(t, err)
	assert.Equal(t, "Canh Chua", detail.Name)
	assert.Len(t, detail.Ingredients, 2)
	assert.Len(t, detail.Steps, 2)

	// Second upsert — update name.
	dish.Name = "Canh Chua Cá Lóc"
	isNew, err = repo.UpsertFromSync(ctx, &dish, ingredients, steps, nil)
	require.NoError(t, err)
	assert.False(t, isNew)

	detail, err = repo.GetBySlug(ctx, "canh-chua")
	require.NoError(t, err)
	assert.Equal(t, "Canh Chua Cá Lóc", detail.Name)
}

func TestIntegration_DishRepo_UpsertFromSync_WithTags(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	tag := seedTag(t, "Quick", "quick")

	extID := "notion-page-456"
	servings := 2
	difficulty := "EASY"
	dish := recipe.Dish{
		ExternalID: &extID,
		Name:       "Mì Xào",
		Slug:       "mi-xao",
		Status:     "PUBLISHED",
		Servings:   &servings,
		Difficulty: &difficulty,
	}

	_, err := repo.UpsertFromSync(ctx, &dish, nil, nil, []uuid.UUID{tag.ID})
	require.NoError(t, err)

	detail, err := repo.GetBySlug(ctx, "mi-xao")
	require.NoError(t, err)
	assert.Len(t, detail.Tags, 1)
	assert.Equal(t, "Quick", detail.Tags[0].Name)
}

func TestIntegration_DishRepo_SoftDelete(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	id := seedRecipe(t, "To Delete", "to-delete", nil)

	err := repo.SoftDelete(ctx, id)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, id)
	assert.Error(t, err)
}

func TestIntegration_DishRepo_SoftDelete_NotFound(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	err := repo.SoftDelete(ctx, uuid.New())
	assert.Error(t, err)
}

func TestIntegration_DishRepo_IncrementViewCount(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	id := seedRecipe(t, "Popular Dish", "popular-dish", nil)

	err := repo.IncrementViewCount(ctx, id)
	require.NoError(t, err)

	detail, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, 1, detail.ViewCount)
}

func TestIntegration_DishRepo_GetAllPublishedIDs(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewDishRepository(testDB.Pool)

	id1 := seedRecipe(t, "Dish 1", "dish-1", nil)
	id2 := seedRecipe(t, "Dish 2", "dish-2", nil)

	// Insert a DRAFT dish (should not appear).
	_, err := testDB.Pool.Exec(ctx,
		`INSERT INTO recipes (name, slug, status) VALUES ('Draft', 'draft', 'DRAFT')`)
	require.NoError(t, err)

	ids, err := repo.GetAllPublishedIDs(ctx)
	require.NoError(t, err)
	assert.Len(t, ids, 2)
	assert.Contains(t, ids, id1)
	assert.Contains(t, ids, id2)
}

func TestIntegration_CategoryRepo_List(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewCategoryRepository(testDB.Pool)

	seedCategory(t, "Soup", "soup", "DISH_TYPE")
	seedCategory(t, "Stir Fry", "stir-fry", "DISH_TYPE")
	seedCategory(t, "Northern", "northern", "REGION")

	cats, err := repo.List(ctx, "")
	require.NoError(t, err)
	assert.Len(t, cats, 3)

	cats, err = repo.List(ctx, "DISH_TYPE")
	require.NoError(t, err)
	assert.Len(t, cats, 2)
}

func TestIntegration_CategoryRepo_GetBySlug(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewCategoryRepository(testDB.Pool)

	seedCategory(t, "Soup", "soup", "DISH_TYPE")

	cat, err := repo.GetBySlug(ctx, "soup")
	require.NoError(t, err)
	assert.Equal(t, "Soup", cat.Name)
	assert.Equal(t, "DISH_TYPE", cat.Type)
}

func TestIntegration_TagRepo_List(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewTagRepository(testDB.Pool)

	seedTag(t, "Quick", "quick")
	seedTag(t, "Healthy", "healthy")

	tags, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, tags, 2)
}

func TestIntegration_TagRepo_GetBySlug(t *testing.T) {
	requireIntegrationDB(t)
	ctx := context.Background()
	repo := recipe.NewTagRepository(testDB.Pool)

	seedTag(t, "Quick", "quick")

	tag, err := repo.GetBySlug(ctx, "quick")
	require.NoError(t, err)
	assert.Equal(t, "Quick", tag.Name)
}
