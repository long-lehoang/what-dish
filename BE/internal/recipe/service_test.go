package recipe

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Mock implementations
// ---------------------------------------------------------------------------

type mockDishRepo struct {
	listFn               func(ctx context.Context, filter DishFilter) ([]Dish, int64, error)
	getByIDFn            func(ctx context.Context, id uuid.UUID) (*DishDetail, error)
	getBySlugFn          func(ctx context.Context, slug string) (*DishDetail, error)
	getRandom            func(ctx context.Context, filter DishFilter, excludeIDs []uuid.UUID) (*DishDetail, error)
	searchFn             func(ctx context.Context, query string, page, pageSize int) ([]Dish, int64, error)
	getAllPublishedIDsFn func(ctx context.Context) ([]uuid.UUID, error)
	upsertFromSyncFn     func(ctx context.Context, dish *Dish, ingredients []Ingredient, steps []Step, tagIDs []uuid.UUID) (bool, error)
	softDeleteFn         func(ctx context.Context, id uuid.UUID) error
	incrementViewFn      func(ctx context.Context, id uuid.UUID) error
}

func (m *mockDishRepo) List(ctx context.Context, filter DishFilter) ([]Dish, int64, error) {
	if m.listFn != nil {
		return m.listFn(ctx, filter)
	}
	return nil, 0, nil
}

func (m *mockDishRepo) GetByID(ctx context.Context, id uuid.UUID) (*DishDetail, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockDishRepo) GetBySlug(ctx context.Context, slug string) (*DishDetail, error) {
	if m.getBySlugFn != nil {
		return m.getBySlugFn(ctx, slug)
	}
	return nil, nil
}

func (m *mockDishRepo) GetRandom(ctx context.Context, filter DishFilter, excludeIDs []uuid.UUID) (*DishDetail, error) {
	if m.getRandom != nil {
		return m.getRandom(ctx, filter, excludeIDs)
	}
	return nil, nil
}

func (m *mockDishRepo) Search(ctx context.Context, query string, page, pageSize int) ([]Dish, int64, error) {
	if m.searchFn != nil {
		return m.searchFn(ctx, query, page, pageSize)
	}
	return nil, 0, nil
}

func (m *mockDishRepo) GetAllPublishedIDs(ctx context.Context) ([]uuid.UUID, error) {
	if m.getAllPublishedIDsFn != nil {
		return m.getAllPublishedIDsFn(ctx)
	}
	return nil, nil
}

func (m *mockDishRepo) UpsertFromSync(ctx context.Context, dish *Dish, ingredients []Ingredient, steps []Step, tagIDs []uuid.UUID) (bool, error) {
	if m.upsertFromSyncFn != nil {
		return m.upsertFromSyncFn(ctx, dish, ingredients, steps, tagIDs)
	}
	return false, nil
}

func (m *mockDishRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if m.softDeleteFn != nil {
		return m.softDeleteFn(ctx, id)
	}
	return nil
}

func (m *mockDishRepo) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	if m.incrementViewFn != nil {
		return m.incrementViewFn(ctx, id)
	}
	return nil
}

type mockCategoryRepo struct {
	listFn      func(ctx context.Context, categoryType string) ([]Category, error)
	getBySlugFn func(ctx context.Context, slug string) (*Category, error)
}

func (m *mockCategoryRepo) List(ctx context.Context, categoryType string) ([]Category, error) {
	if m.listFn != nil {
		return m.listFn(ctx, categoryType)
	}
	return nil, nil
}

func (m *mockCategoryRepo) GetBySlug(ctx context.Context, slug string) (*Category, error) {
	if m.getBySlugFn != nil {
		return m.getBySlugFn(ctx, slug)
	}
	return nil, nil
}

type mockTagRepo struct {
	listFn      func(ctx context.Context) ([]Tag, error)
	getBySlugFn func(ctx context.Context, slug string) (*Tag, error)
}

func (m *mockTagRepo) List(ctx context.Context) ([]Tag, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}
	return nil, nil
}

func (m *mockTagRepo) GetBySlug(ctx context.Context, slug string) (*Tag, error) {
	if m.getBySlugFn != nil {
		return m.getBySlugFn(ctx, slug)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// Tests: DishService.ListDishes
// ---------------------------------------------------------------------------

func TestDishService_ListDishes_Success(t *testing.T) {
	now := time.Now().UTC()
	dishes := []Dish{
		{ID: uuid.New(), Name: "Pho Bo", Slug: "pho-bo", Status: "published", CreatedAt: now, UpdatedAt: now},
		{ID: uuid.New(), Name: "Bun Cha", Slug: "bun-cha", Status: "published", CreatedAt: now, UpdatedAt: now},
	}

	svc := NewDishService(
		&mockDishRepo{
			listFn: func(_ context.Context, _ DishFilter) ([]Dish, int64, error) {
				return dishes, 2, nil
			},
		},
		&mockCategoryRepo{},
		&mockTagRepo{},
	)

	got, total, err := svc.ListDishes(context.Background(), DishFilter{Page: 1, PageSize: 20})

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, got, 2)
	assert.Equal(t, "Pho Bo", got[0].Name)
	assert.Equal(t, "Bun Cha", got[1].Name)
}

func TestDishService_ListDishes_RepoError(t *testing.T) {
	repoErr := errors.New("database connection lost")

	svc := NewDishService(
		&mockDishRepo{
			listFn: func(_ context.Context, _ DishFilter) ([]Dish, int64, error) {
				return nil, 0, repoErr
			},
		},
		&mockCategoryRepo{},
		&mockTagRepo{},
	)

	got, total, err := svc.ListDishes(context.Background(), DishFilter{Page: 1, PageSize: 20})

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Contains(t, err.Error(), "recipe.ListDishes")
	assert.Nil(t, got)
	assert.Equal(t, int64(0), total)
}

func TestDishService_ListDishes_EmptyResult(t *testing.T) {
	svc := NewDishService(
		&mockDishRepo{
			listFn: func(_ context.Context, _ DishFilter) ([]Dish, int64, error) {
				return []Dish{}, 0, nil
			},
		},
		&mockCategoryRepo{},
		&mockTagRepo{},
	)

	got, total, err := svc.ListDishes(context.Background(), DishFilter{Page: 1, PageSize: 20})

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, got)
}

// ---------------------------------------------------------------------------
// Tests: DishService.GetDish
// ---------------------------------------------------------------------------

func TestDishService_GetDish_Success(t *testing.T) {
	dishID := uuid.New()
	now := time.Now().UTC()
	detail := &DishDetail{
		Dish: Dish{
			ID:        dishID,
			Name:      "Pho Bo",
			Slug:      "pho-bo",
			Status:    "published",
			CreatedAt: now,
			UpdatedAt: now,
		},
		Ingredients: []Ingredient{
			{ID: uuid.New(), RecipeID: dishID, Name: "Banh pho", SortOrder: 1},
		},
		Steps: []Step{
			{ID: uuid.New(), RecipeID: dishID, StepNumber: 1, Description: "Boil broth", SortOrder: 1},
		},
		Tags: []Tag{
			{ID: uuid.New(), Name: "Soup", Slug: "soup"},
		},
	}

	svc := NewDishService(
		&mockDishRepo{
			getByIDFn: func(_ context.Context, id uuid.UUID) (*DishDetail, error) {
				assert.Equal(t, dishID, id)
				return detail, nil
			},
		},
		&mockCategoryRepo{},
		&mockTagRepo{},
	)

	got, err := svc.GetDish(context.Background(), dishID)

	assert.NoError(t, err)
	assert.Equal(t, dishID, got.ID)
	assert.Equal(t, "Pho Bo", got.Name)
	assert.Len(t, got.Ingredients, 1)
	assert.Len(t, got.Steps, 1)
	assert.Len(t, got.Tags, 1)
}

func TestDishService_GetDish_NotFound(t *testing.T) {
	notFoundErr := errors.New("not found")
	dishID := uuid.New()

	svc := NewDishService(
		&mockDishRepo{
			getByIDFn: func(_ context.Context, _ uuid.UUID) (*DishDetail, error) {
				return nil, notFoundErr
			},
		},
		&mockCategoryRepo{},
		&mockTagRepo{},
	)

	got, err := svc.GetDish(context.Background(), dishID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, notFoundErr)
	assert.Contains(t, err.Error(), "recipe.GetDish")
	assert.Nil(t, got)
}

// ---------------------------------------------------------------------------
// Tests: DishService.GetDishBySlug
// ---------------------------------------------------------------------------

func TestDishService_GetDishBySlug_Success(t *testing.T) {
	slug := "pho-bo"
	now := time.Now().UTC()
	detail := &DishDetail{
		Dish: Dish{
			ID:        uuid.New(),
			Name:      "Pho Bo",
			Slug:      slug,
			Status:    "published",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	svc := NewDishService(
		&mockDishRepo{
			getBySlugFn: func(_ context.Context, s string) (*DishDetail, error) {
				assert.Equal(t, slug, s)
				return detail, nil
			},
		},
		&mockCategoryRepo{},
		&mockTagRepo{},
	)

	got, err := svc.GetDishBySlug(context.Background(), slug)

	assert.NoError(t, err)
	assert.Equal(t, slug, got.Slug)
	assert.Equal(t, "Pho Bo", got.Name)
}

func TestDishService_GetDishBySlug_NotFound(t *testing.T) {
	notFoundErr := errors.New("not found")

	svc := NewDishService(
		&mockDishRepo{
			getBySlugFn: func(_ context.Context, _ string) (*DishDetail, error) {
				return nil, notFoundErr
			},
		},
		&mockCategoryRepo{},
		&mockTagRepo{},
	)

	got, err := svc.GetDishBySlug(context.Background(), "nonexistent")

	assert.Error(t, err)
	assert.ErrorIs(t, err, notFoundErr)
	assert.Nil(t, got)
}

// ---------------------------------------------------------------------------
// Tests: DishService.GetRandomDish
// ---------------------------------------------------------------------------

func TestDishService_GetRandomDish_Success(t *testing.T) {
	now := time.Now().UTC()
	detail := &DishDetail{
		Dish: Dish{
			ID:        uuid.New(),
			Name:      "Bun Bo Hue",
			Slug:      "bun-bo-hue",
			Status:    "published",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	excludeIDs := []uuid.UUID{uuid.New()}

	svc := NewDishService(
		&mockDishRepo{
			getRandom: func(_ context.Context, f DishFilter, excl []uuid.UUID) (*DishDetail, error) {
				assert.Len(t, excl, 1)
				return detail, nil
			},
		},
		&mockCategoryRepo{},
		&mockTagRepo{},
	)

	got, err := svc.GetRandomDish(context.Background(), DishFilter{}, excludeIDs)

	assert.NoError(t, err)
	assert.Equal(t, "Bun Bo Hue", got.Name)
}

func TestDishService_GetRandomDish_RepoError(t *testing.T) {
	repoErr := errors.New("no dishes available")

	svc := NewDishService(
		&mockDishRepo{
			getRandom: func(_ context.Context, _ DishFilter, _ []uuid.UUID) (*DishDetail, error) {
				return nil, repoErr
			},
		},
		&mockCategoryRepo{},
		&mockTagRepo{},
	)

	got, err := svc.GetRandomDish(context.Background(), DishFilter{}, nil)

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, got)
}

// ---------------------------------------------------------------------------
// Tests: DishService.SearchDishes
// ---------------------------------------------------------------------------

func TestDishService_SearchDishes_Success(t *testing.T) {
	now := time.Now().UTC()
	dishes := []Dish{
		{ID: uuid.New(), Name: "Pho Bo", Slug: "pho-bo", Status: "published", CreatedAt: now, UpdatedAt: now},
	}

	svc := NewDishService(
		&mockDishRepo{
			searchFn: func(_ context.Context, q string, page, pageSize int) ([]Dish, int64, error) {
				assert.Equal(t, "pho", q)
				assert.Equal(t, 1, page)
				assert.Equal(t, 20, pageSize)
				return dishes, 1, nil
			},
		},
		&mockCategoryRepo{},
		&mockTagRepo{},
	)

	got, total, err := svc.SearchDishes(context.Background(), "pho", 1, 20)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, got, 1)
	assert.Equal(t, "Pho Bo", got[0].Name)
}

func TestDishService_SearchDishes_RepoError(t *testing.T) {
	repoErr := errors.New("search index error")

	svc := NewDishService(
		&mockDishRepo{
			searchFn: func(_ context.Context, _ string, _, _ int) ([]Dish, int64, error) {
				return nil, 0, repoErr
			},
		},
		&mockCategoryRepo{},
		&mockTagRepo{},
	)

	got, total, err := svc.SearchDishes(context.Background(), "pho", 1, 20)

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, got)
	assert.Equal(t, int64(0), total)
}

func TestDishService_SearchDishes_EmptyResult(t *testing.T) {
	svc := NewDishService(
		&mockDishRepo{
			searchFn: func(_ context.Context, _ string, _, _ int) ([]Dish, int64, error) {
				return []Dish{}, 0, nil
			},
		},
		&mockCategoryRepo{},
		&mockTagRepo{},
	)

	got, total, err := svc.SearchDishes(context.Background(), "nonexistent", 1, 20)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, got)
}

// ---------------------------------------------------------------------------
// Tests: DishService.ListCategories
// ---------------------------------------------------------------------------

func TestDishService_ListCategories_Success(t *testing.T) {
	now := time.Now().UTC()
	categories := []Category{
		{ID: uuid.New(), Name: "Soup", Slug: "soup", Type: "DISH_TYPE", IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: uuid.New(), Name: "Stir Fry", Slug: "stir-fry", Type: "DISH_TYPE", IsActive: true, CreatedAt: now, UpdatedAt: now},
	}

	svc := NewDishService(
		&mockDishRepo{},
		&mockCategoryRepo{
			listFn: func(_ context.Context, catType string) ([]Category, error) {
				assert.Equal(t, "DISH_TYPE", catType)
				return categories, nil
			},
		},
		&mockTagRepo{},
	)

	got, err := svc.ListCategories(context.Background(), "DISH_TYPE")

	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "Soup", got[0].Name)
	assert.Equal(t, "Stir Fry", got[1].Name)
}

func TestDishService_ListCategories_RepoError(t *testing.T) {
	repoErr := errors.New("db error")

	svc := NewDishService(
		&mockDishRepo{},
		&mockCategoryRepo{
			listFn: func(_ context.Context, _ string) ([]Category, error) {
				return nil, repoErr
			},
		},
		&mockTagRepo{},
	)

	got, err := svc.ListCategories(context.Background(), "DISH_TYPE")

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Contains(t, err.Error(), "recipe.ListCategories")
	assert.Nil(t, got)
}

// ---------------------------------------------------------------------------
// Tests: DishService.ListTags
// ---------------------------------------------------------------------------

func TestDishService_ListTags_Success(t *testing.T) {
	tags := []Tag{
		{ID: uuid.New(), Name: "Quick", Slug: "quick"},
		{ID: uuid.New(), Name: "Healthy", Slug: "healthy"},
		{ID: uuid.New(), Name: "Kid Friendly", Slug: "kid-friendly"},
	}

	svc := NewDishService(
		&mockDishRepo{},
		&mockCategoryRepo{},
		&mockTagRepo{
			listFn: func(_ context.Context) ([]Tag, error) {
				return tags, nil
			},
		},
	)

	got, err := svc.ListTags(context.Background())

	assert.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Equal(t, "Quick", got[0].Name)
	assert.Equal(t, "Healthy", got[1].Name)
	assert.Equal(t, "Kid Friendly", got[2].Name)
}

func TestDishService_ListTags_RepoError(t *testing.T) {
	repoErr := errors.New("db error")

	svc := NewDishService(
		&mockDishRepo{},
		&mockCategoryRepo{},
		&mockTagRepo{
			listFn: func(_ context.Context) ([]Tag, error) {
				return nil, repoErr
			},
		},
	)

	got, err := svc.ListTags(context.Background())

	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Contains(t, err.Error(), "recipe.ListTags")
	assert.Nil(t, got)
}

// ---------------------------------------------------------------------------
// Tests: DishService.ListDishes with filters (table-driven)
// ---------------------------------------------------------------------------

func TestDishService_ListDishes_PassesFilterToRepo(t *testing.T) {
	dishTypeID := uuid.New()
	difficulty := "EASY"
	maxCookTime := 30

	var capturedFilter DishFilter
	svc := NewDishService(
		&mockDishRepo{
			listFn: func(_ context.Context, f DishFilter) ([]Dish, int64, error) {
				capturedFilter = f
				return []Dish{}, 0, nil
			},
		},
		&mockCategoryRepo{},
		&mockTagRepo{},
	)

	filter := DishFilter{
		DishTypeID:  &dishTypeID,
		Difficulty:  &difficulty,
		MaxCookTime: &maxCookTime,
		Tags:        []string{"quick", "healthy"},
		Page:        2,
		PageSize:    10,
	}

	_, _, err := svc.ListDishes(context.Background(), filter)

	assert.NoError(t, err)
	assert.Equal(t, &dishTypeID, capturedFilter.DishTypeID)
	assert.Equal(t, &difficulty, capturedFilter.Difficulty)
	assert.Equal(t, &maxCookTime, capturedFilter.MaxCookTime)
	assert.Equal(t, []string{"quick", "healthy"}, capturedFilter.Tags)
	assert.Equal(t, 2, capturedFilter.Page)
	assert.Equal(t, 10, capturedFilter.PageSize)
}
