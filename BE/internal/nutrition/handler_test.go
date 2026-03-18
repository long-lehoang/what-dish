package nutrition

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Helpers (mocks are defined in service_test.go and shared across test files)
// ---------------------------------------------------------------------------

func setupNutritionRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/nutrition/calculate-tdee", h.HandleCalculateTDEE)
	r.GET("/nutrition/goals", h.HandleListGoals)
	r.GET("/nutrition/recipes/:id", h.HandleGetRecipeNutrition)
	r.GET("/nutrition/recipes", h.HandleGetBatchNutrition)
	return r
}

func parseNutritionBody(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var result map[string]any
	err := json.NewDecoder(w.Body).Decode(&result)
	assert.NoError(t, err, "response body should be valid JSON")
	return result
}

// ---------------------------------------------------------------------------
// Tests — POST /nutrition/calculate-tdee
// ---------------------------------------------------------------------------

func TestHandleCalculateTDEE_ValidRequest(t *testing.T) {
	svc := NewNutritionService(&mockNutritionRepo{}, &mockGoalRepo{})
	h := NewHandler(svc)
	router := setupNutritionRouter(h)

	reqBody := CalculateTDEERequest{
		Gender:        "MALE",
		Age:           25,
		HeightCm:      175,
		WeightKg:      70,
		ActivityLevel: "MODERATE",
		Goal:          "MAINTAIN",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/nutrition/calculate-tdee", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseNutritionBody(t, w)
	data, ok := body["data"].(map[string]any)
	assert.True(t, ok, "response should contain data object")

	assert.Contains(t, data, "bmr")
	assert.Contains(t, data, "tdee")
	assert.Contains(t, data, "dailyTarget")
	assert.Contains(t, data, "mealBreakdown")

	bmr, ok := data["bmr"].(float64)
	assert.True(t, ok)
	assert.Greater(t, bmr, 0.0)

	tdee, ok := data["tdee"].(float64)
	assert.True(t, ok)
	assert.Greater(t, tdee, bmr, "TDEE should be greater than BMR for MODERATE activity")

	mealBreakdown, ok := data["mealBreakdown"].(map[string]any)
	assert.True(t, ok)
	assert.Contains(t, mealBreakdown, "breakfast")
	assert.Contains(t, mealBreakdown, "lunch")
	assert.Contains(t, mealBreakdown, "dinner")
}

func TestHandleCalculateTDEE_InvalidBody(t *testing.T) {
	svc := NewNutritionService(&mockNutritionRepo{}, &mockGoalRepo{})
	h := NewHandler(svc)
	router := setupNutritionRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/nutrition/calculate-tdee", bytes.NewReader([]byte(`{invalid`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	body := parseNutritionBody(t, w)
	assert.Equal(t, "invalid request body", body["message"])
}

func TestHandleCalculateTDEE_MissingRequiredFields(t *testing.T) {
	svc := NewNutritionService(&mockNutritionRepo{}, &mockGoalRepo{})
	h := NewHandler(svc)
	router := setupNutritionRouter(h)

	// Missing all required fields — sends empty object.
	req := httptest.NewRequest(http.MethodPost, "/nutrition/calculate-tdee", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCalculateTDEE_InvalidGender(t *testing.T) {
	svc := NewNutritionService(&mockNutritionRepo{}, &mockGoalRepo{})
	h := NewHandler(svc)
	router := setupNutritionRouter(h)

	reqBody := CalculateTDEERequest{
		Gender:        "INVALID",
		Age:           25,
		HeightCm:      175,
		WeightKg:      70,
		ActivityLevel: "MODERATE",
		Goal:          "MAINTAIN",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/nutrition/calculate-tdee", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCalculateTDEE_InvalidActivityLevel(t *testing.T) {
	svc := NewNutritionService(&mockNutritionRepo{}, &mockGoalRepo{})
	h := NewHandler(svc)
	router := setupNutritionRouter(h)

	reqBody := CalculateTDEERequest{
		Gender:        "MALE",
		Age:           25,
		HeightCm:      175,
		WeightKg:      70,
		ActivityLevel: "SUPER_ACTIVE",
		Goal:          "MAINTAIN",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/nutrition/calculate-tdee", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCalculateTDEE_FemaleRequest(t *testing.T) {
	svc := NewNutritionService(&mockNutritionRepo{}, &mockGoalRepo{})
	h := NewHandler(svc)
	router := setupNutritionRouter(h)

	reqBody := CalculateTDEERequest{
		Gender:        "FEMALE",
		Age:           30,
		HeightCm:      165,
		WeightKg:      60,
		ActivityLevel: "LIGHT",
		Goal:          "LOSE_WEIGHT",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/nutrition/calculate-tdee", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseNutritionBody(t, w)
	data := body["data"].(map[string]any)

	dailyTarget := data["dailyTarget"].(float64)
	tdee := data["tdee"].(float64)
	assert.Less(t, dailyTarget, tdee, "LOSE_WEIGHT goal should reduce daily target below TDEE")
}

// ---------------------------------------------------------------------------
// Tests — GET /nutrition/goals
// ---------------------------------------------------------------------------

func TestHandleListGoals_OK(t *testing.T) {
	goals := []NutritionGoal{
		{ID: uuid.New(), Name: "Weight Loss", Description: "Reduce calorie intake", IsActive: true, SortOrder: 1},
		{ID: uuid.New(), Name: "Maintain", Description: "Keep current weight", IsActive: true, SortOrder: 2},
	}
	gRepo := &mockGoalRepo{
		listFn: func(_ context.Context) ([]NutritionGoal, error) {
			return goals, nil
		},
	}
	svc := NewNutritionService(&mockNutritionRepo{}, gRepo)
	h := NewHandler(svc)
	router := setupNutritionRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/nutrition/goals", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseNutritionBody(t, w)
	data, ok := body["data"].([]any)
	assert.True(t, ok, "response should contain data array")
	assert.Len(t, data, 2)

	first := data[0].(map[string]any)
	assert.Equal(t, "Weight Loss", first["name"])
}

func TestHandleListGoals_ServiceError(t *testing.T) {
	gRepo := &mockGoalRepo{
		listFn: func(_ context.Context) ([]NutritionGoal, error) {
			return nil, fmt.Errorf("database error")
		},
	}
	svc := NewNutritionService(&mockNutritionRepo{}, gRepo)
	h := NewHandler(svc)
	router := setupNutritionRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/nutrition/goals", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ---------------------------------------------------------------------------
// Tests — GET /nutrition/recipes/:id
// ---------------------------------------------------------------------------

func TestHandleGetRecipeNutrition_OK(t *testing.T) {
	recipeID := uuid.New()
	cal := 350.0
	prot := 25.0
	nutrition := &RecipeNutrition{
		ID:       uuid.New(),
		RecipeID: recipeID,
		Calories: &cal,
		Protein:  &prot,
	}
	nRepo := &mockNutritionRepo{
		getByRecipeIDFn: func(_ context.Context, id uuid.UUID) (*RecipeNutrition, error) {
			assert.Equal(t, recipeID, id)
			return nutrition, nil
		},
	}
	svc := NewNutritionService(nRepo, &mockGoalRepo{})
	h := NewHandler(svc)
	router := setupNutritionRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/nutrition/recipes/"+recipeID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseNutritionBody(t, w)
	data := body["data"].(map[string]any)
	assert.Equal(t, 350.0, data["calories"])
	assert.Equal(t, 25.0, data["protein"])
}

func TestHandleGetRecipeNutrition_InvalidID(t *testing.T) {
	svc := NewNutritionService(&mockNutritionRepo{}, &mockGoalRepo{})
	h := NewHandler(svc)
	router := setupNutritionRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/nutrition/recipes/not-a-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	body := parseNutritionBody(t, w)
	assert.Equal(t, "invalid recipe id", body["message"])
}

// ---------------------------------------------------------------------------
// Tests — GET /nutrition/recipes?ids=...
// ---------------------------------------------------------------------------

func TestHandleGetBatchNutrition_OK(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	cal1 := 400.0
	cal2 := 500.0
	nRepo := &mockNutritionRepo{
		getByRecipeIDsFn: func(_ context.Context, ids []uuid.UUID) ([]RecipeNutrition, error) {
			assert.Len(t, ids, 2)
			return []RecipeNutrition{
				{ID: uuid.New(), RecipeID: id1, Calories: &cal1},
				{ID: uuid.New(), RecipeID: id2, Calories: &cal2},
			}, nil
		},
	}
	svc := NewNutritionService(nRepo, &mockGoalRepo{})
	h := NewHandler(svc)
	router := setupNutritionRouter(h)

	url := fmt.Sprintf("/nutrition/recipes?ids=%s,%s", id1, id2)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseNutritionBody(t, w)
	data, ok := body["data"].([]any)
	assert.True(t, ok)
	assert.Len(t, data, 2)
}

func TestHandleGetBatchNutrition_MissingIDs(t *testing.T) {
	svc := NewNutritionService(&mockNutritionRepo{}, &mockGoalRepo{})
	h := NewHandler(svc)
	router := setupNutritionRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/nutrition/recipes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	body := parseNutritionBody(t, w)
	assert.Equal(t, "ids query parameter is required", body["message"])
}

func TestHandleGetBatchNutrition_InvalidUUID(t *testing.T) {
	svc := NewNutritionService(&mockNutritionRepo{}, &mockGoalRepo{})
	h := NewHandler(svc)
	router := setupNutritionRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/nutrition/recipes?ids=bad-uuid,also-bad", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
