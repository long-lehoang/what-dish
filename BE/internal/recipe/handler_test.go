package recipe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Helpers (mocks are defined in service_test.go and shared across test files)
// ---------------------------------------------------------------------------

func setupTestRouter(h *DishHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/recipes", h.HandleListDishes)
	r.GET("/recipes/random", h.HandleGetRandomDish)
	r.GET("/recipes/search", h.HandleSearchDishes)
	r.GET("/recipes/:id", h.HandleGetDish)
	r.GET("/categories", h.HandleListCategories)
	r.GET("/tags", h.HandleListTags)
	return r
}

func sampleDish() Dish {
	now := time.Now().UTC()
	diff := "EASY"
	cookTime := 30
	return Dish{
		ID:         uuid.New(),
		Name:       "Pho Bo",
		Slug:       "pho-bo",
		Difficulty: &diff,
		CookTime:   &cookTime,
		Status:     "PUBLISHED",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func sampleDishDetail() *DishDetail {
	d := sampleDish()
	return &DishDetail{
		Dish: d,
		Ingredients: []Ingredient{
			{
				ID:       uuid.New(),
				RecipeID: d.ID,
				Name:     "Beef",
			},
		},
		Steps: []Step{
			{
				ID:          uuid.New(),
				RecipeID:    d.ID,
				StepNumber:  1,
				Description: "Boil water",
			},
		},
		Tags: []Tag{
			{ID: uuid.New(), Name: "quick", Slug: "quick"},
		},
	}
}

// parseBody decodes a JSON response body into a generic map.
func parseBody(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var result map[string]any
	err := json.NewDecoder(w.Body).Decode(&result)
	assert.NoError(t, err, "response body should be valid JSON")
	return result
}

// ---------------------------------------------------------------------------
// Tests — GET /recipes
// ---------------------------------------------------------------------------

func TestHandleListDishes_OK(t *testing.T) {
	dishes := []Dish{sampleDish(), sampleDish()}
	dishRepo := &mockDishRepo{
		listFn: func(_ context.Context, f DishFilter) ([]Dish, int64, error) {
			assert.Equal(t, 1, f.Page)
			assert.Equal(t, 20, f.PageSize)
			return dishes, 2, nil
		},
	}
	svc := NewDishService(dishRepo, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/recipes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseBody(t, w)
	data, ok := body["data"].([]any)
	assert.True(t, ok, "response should contain data array")
	assert.Len(t, data, 2)

	pagination, ok := body["pagination"].(map[string]any)
	assert.True(t, ok, "response should contain pagination object")
	assert.Equal(t, float64(2), pagination["total"])
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(20), pagination["pageSize"])
}

func TestHandleListDishes_WithDifficultyFilter(t *testing.T) {
	var capturedFilter DishFilter
	dishRepo := &mockDishRepo{
		listFn: func(_ context.Context, f DishFilter) ([]Dish, int64, error) {
			capturedFilter = f
			return []Dish{}, 0, nil
		},
	}
	svc := NewDishService(dishRepo, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/recipes?difficulty=EASY", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, capturedFilter.Difficulty)
	assert.Equal(t, "EASY", *capturedFilter.Difficulty)
}

func TestHandleListDishes_WithMultipleFilters(t *testing.T) {
	dishTypeID := uuid.New()
	regionID := uuid.New()

	var capturedFilter DishFilter
	dishRepo := &mockDishRepo{
		listFn: func(_ context.Context, f DishFilter) ([]Dish, int64, error) {
			capturedFilter = f
			return []Dish{}, 0, nil
		},
	}
	svc := NewDishService(dishRepo, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	url := fmt.Sprintf("/recipes?dish_type=%s&region=%s&difficulty=HARD&max_cook_time=45&tags=quick,healthy", dishTypeID, regionID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, capturedFilter.DishTypeID)
	assert.Equal(t, dishTypeID, *capturedFilter.DishTypeID)
	assert.NotNil(t, capturedFilter.RegionID)
	assert.Equal(t, regionID, *capturedFilter.RegionID)
	assert.NotNil(t, capturedFilter.Difficulty)
	assert.Equal(t, "HARD", *capturedFilter.Difficulty)
	assert.NotNil(t, capturedFilter.MaxCookTime)
	assert.Equal(t, 45, *capturedFilter.MaxCookTime)
	assert.Equal(t, []string{"quick", "healthy"}, capturedFilter.Tags)
}

func TestHandleListDishes_Pagination(t *testing.T) {
	var capturedFilter DishFilter
	dishRepo := &mockDishRepo{
		listFn: func(_ context.Context, f DishFilter) ([]Dish, int64, error) {
			capturedFilter = f
			return []Dish{}, 0, nil
		},
	}
	svc := NewDishService(dishRepo, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/recipes?page=3&pageSize=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 3, capturedFilter.Page)
	assert.Equal(t, 10, capturedFilter.PageSize)
}

func TestHandleListDishes_ServiceError(t *testing.T) {
	dishRepo := &mockDishRepo{
		listFn: func(_ context.Context, _ DishFilter) ([]Dish, int64, error) {
			return nil, 0, fmt.Errorf("database connection failed")
		},
	}
	svc := NewDishService(dishRepo, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/recipes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	body := parseBody(t, w)
	_, hasError := body["error"]
	assert.True(t, hasError)
}

// ---------------------------------------------------------------------------
// Tests — GET /recipes/:id
// ---------------------------------------------------------------------------

func TestHandleGetDish_ByUUID(t *testing.T) {
	detail := sampleDishDetail()
	dishRepo := &mockDishRepo{
		getByIDFn: func(_ context.Context, id uuid.UUID) (*DishDetail, error) {
			assert.Equal(t, detail.ID, id)
			return detail, nil
		},
	}
	svc := NewDishService(dishRepo, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/recipes/"+detail.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseBody(t, w)
	data, ok := body["data"].(map[string]any)
	assert.True(t, ok, "response should contain data object")
	assert.Equal(t, detail.Name, data["name"])
	assert.Equal(t, detail.Slug, data["slug"])

	ingredients, ok := data["ingredients"].([]any)
	assert.True(t, ok)
	assert.Len(t, ingredients, 1)

	steps, ok := data["steps"].([]any)
	assert.True(t, ok)
	assert.Len(t, steps, 1)

	tags, ok := data["tags"].([]any)
	assert.True(t, ok)
	assert.Len(t, tags, 1)
}

func TestHandleGetDish_BySlug(t *testing.T) {
	detail := sampleDishDetail()
	var calledSlug string
	dishRepo := &mockDishRepo{
		getBySlugFn: func(_ context.Context, slug string) (*DishDetail, error) {
			calledSlug = slug
			return detail, nil
		},
	}
	svc := NewDishService(dishRepo, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/recipes/pho-bo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pho-bo", calledSlug)

	body := parseBody(t, w)
	data, ok := body["data"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, detail.Name, data["name"])
}

func TestHandleGetDish_NotFound(t *testing.T) {
	dishRepo := &mockDishRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*DishDetail, error) {
			return nil, fmt.Errorf("recipe.GetDish: %w", fmt.Errorf("resource not found"))
		},
	}
	svc := NewDishService(dishRepo, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/recipes/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Without wrapping with shared/errors sentinel, this maps to 500.
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	body := parseBody(t, w)
	_, hasError := body["error"]
	assert.True(t, hasError, "error response should contain error field")
}

// ---------------------------------------------------------------------------
// Tests — GET /recipes/random
// ---------------------------------------------------------------------------

func TestHandleGetRandomDish_OK(t *testing.T) {
	detail := sampleDishDetail()
	dishRepo := &mockDishRepo{
		getRandom: func(_ context.Context, _ DishFilter, excludeIDs []uuid.UUID) (*DishDetail, error) {
			assert.Empty(t, excludeIDs)
			return detail, nil
		},
	}
	svc := NewDishService(dishRepo, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/recipes/random", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseBody(t, w)
	data, ok := body["data"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, detail.Name, data["name"])
}

func TestHandleGetRandomDish_WithExcludeIDs(t *testing.T) {
	excludeID1 := uuid.New()
	excludeID2 := uuid.New()
	detail := sampleDishDetail()

	var capturedExcludeIDs []uuid.UUID
	dishRepo := &mockDishRepo{
		getRandom: func(_ context.Context, _ DishFilter, excludeIDs []uuid.UUID) (*DishDetail, error) {
			capturedExcludeIDs = excludeIDs
			return detail, nil
		},
	}
	svc := NewDishService(dishRepo, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	url := fmt.Sprintf("/recipes/random?exclude_ids=%s,%s", excludeID1, excludeID2)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, capturedExcludeIDs, 2)
	assert.Contains(t, capturedExcludeIDs, excludeID1)
	assert.Contains(t, capturedExcludeIDs, excludeID2)
}

func TestHandleGetRandomDish_WithFilters(t *testing.T) {
	detail := sampleDishDetail()
	var capturedFilter DishFilter
	dishRepo := &mockDishRepo{
		getRandom: func(_ context.Context, f DishFilter, _ []uuid.UUID) (*DishDetail, error) {
			capturedFilter = f
			return detail, nil
		},
	}
	svc := NewDishService(dishRepo, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/recipes/random?difficulty=EASY&max_cook_time=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, capturedFilter.Difficulty)
	assert.Equal(t, "EASY", *capturedFilter.Difficulty)
	assert.NotNil(t, capturedFilter.MaxCookTime)
	assert.Equal(t, 20, *capturedFilter.MaxCookTime)
}

// ---------------------------------------------------------------------------
// Tests — GET /recipes/search
// ---------------------------------------------------------------------------

func TestHandleSearchDishes_OK(t *testing.T) {
	dishes := []Dish{sampleDish()}
	var capturedQuery string
	dishRepo := &mockDishRepo{
		searchFn: func(_ context.Context, query string, page, pageSize int) ([]Dish, int64, error) {
			capturedQuery = query
			assert.Equal(t, 1, page)
			assert.Equal(t, 20, pageSize)
			return dishes, 1, nil
		},
	}
	svc := NewDishService(dishRepo, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/recipes/search?q=pho", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pho", capturedQuery)

	body := parseBody(t, w)
	data, ok := body["data"].([]any)
	assert.True(t, ok)
	assert.Len(t, data, 1)

	pagination, ok := body["pagination"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, float64(1), pagination["total"])
}

func TestHandleSearchDishes_MissingQuery(t *testing.T) {
	svc := NewDishService(&mockDishRepo{}, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/recipes/search", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	body := parseBody(t, w)
	assert.Equal(t, "query parameter 'q' is required", body["message"])
}

func TestHandleSearchDishes_EmptyQuery(t *testing.T) {
	svc := NewDishService(&mockDishRepo{}, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/recipes/search?q=", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleSearchDishes_WithPagination(t *testing.T) {
	var capturedPage, capturedPageSize int
	dishRepo := &mockDishRepo{
		searchFn: func(_ context.Context, _ string, page, pageSize int) ([]Dish, int64, error) {
			capturedPage = page
			capturedPageSize = pageSize
			return []Dish{}, 0, nil
		},
	}
	svc := NewDishService(dishRepo, &mockCategoryRepo{}, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/recipes/search?q=bun&page=2&limit=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, capturedPage)
	assert.Equal(t, 5, capturedPageSize)
}

// ---------------------------------------------------------------------------
// Tests — GET /categories
// ---------------------------------------------------------------------------

func TestHandleListCategories_OK(t *testing.T) {
	cats := []Category{
		{ID: uuid.New(), Name: "Main Course", Slug: "main-course", Type: "DISH_TYPE", IsActive: true},
		{ID: uuid.New(), Name: "Soup", Slug: "soup", Type: "DISH_TYPE", IsActive: true},
	}
	catRepo := &mockCategoryRepo{
		listFn: func(_ context.Context, categoryType string) ([]Category, error) {
			assert.Equal(t, "", categoryType)
			return cats, nil
		},
	}
	svc := NewDishService(&mockDishRepo{}, catRepo, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseBody(t, w)
	data, ok := body["data"].([]any)
	assert.True(t, ok, "response should contain data array")
	assert.Len(t, data, 2)

	first := data[0].(map[string]any)
	assert.Equal(t, "Main Course", first["name"])
	assert.Equal(t, "DISH_TYPE", first["type"])
}

func TestHandleListCategories_WithTypeFilter(t *testing.T) {
	var capturedType string
	catRepo := &mockCategoryRepo{
		listFn: func(_ context.Context, categoryType string) ([]Category, error) {
			capturedType = categoryType
			return []Category{}, nil
		},
	}
	svc := NewDishService(&mockDishRepo{}, catRepo, &mockTagRepo{})
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/categories?type=DISH_TYPE", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "DISH_TYPE", capturedType)
}

// ---------------------------------------------------------------------------
// Tests — GET /tags
// ---------------------------------------------------------------------------

func TestHandleListTags_OK(t *testing.T) {
	tags := []Tag{
		{ID: uuid.New(), Name: "quick", Slug: "quick"},
		{ID: uuid.New(), Name: "healthy", Slug: "healthy"},
	}
	tagRepo := &mockTagRepo{
		listFn: func(_ context.Context) ([]Tag, error) {
			return tags, nil
		},
	}
	svc := NewDishService(&mockDishRepo{}, &mockCategoryRepo{}, tagRepo)
	h := NewDishHandler(svc)
	router := setupTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/tags", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseBody(t, w)
	data, ok := body["data"].([]any)
	assert.True(t, ok, "response should contain data array")
	assert.Len(t, data, 2)

	first := data[0].(map[string]any)
	assert.Equal(t, "quick", first["name"])
	assert.Equal(t, "quick", first["slug"])
}
