package e2e_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_GetRecipeNutrition(t *testing.T) {
	requireE2E(t)
	require.NotEmpty(t, recipeIDs)

	// First recipe (Phở Bò) has nutrition seeded.
	resp := doGet(t, "/api/v1/nutrition/recipes/"+recipeIDs[0].String())
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	m := dataAsMap(t, data.Data)
	assert.Equal(t, float64(450), m["calories"])
	assert.Equal(t, float64(30), m["protein"])
}

func TestE2E_GetRecipeNutrition_NotFound(t *testing.T) {
	requireE2E(t)

	// Last recipe (Gỏi Cuốn) has no nutrition data.
	resp := doGet(t, "/api/v1/nutrition/recipes/"+recipeIDs[len(recipeIDs)-1].String())
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()
}

func TestE2E_GetBatchNutrition(t *testing.T) {
	requireE2E(t)
	require.True(t, len(recipeIDs) >= 3)

	ids := recipeIDs[0].String() + "," + recipeIDs[1].String() + "," + recipeIDs[2].String()
	resp := doGet(t, "/api/v1/nutrition/recipes?ids="+ids)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	items := dataAsSlice(t, data.Data)
	assert.Len(t, items, 3)
}

func TestE2E_CalculateTDEE(t *testing.T) {
	requireE2E(t)

	body := map[string]any{
		"gender":        "MALE",
		"age":           30,
		"heightCm":      175,
		"weightKg":      70,
		"activityLevel": "MODERATE",
		"goal":          "MAINTAIN",
	}

	resp := doPost(t, "/api/v1/nutrition/calculate-tdee", body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	m := dataAsMap(t, data.Data)

	bmr, ok := m["bmr"].(float64)
	require.True(t, ok)
	assert.Greater(t, bmr, float64(0))

	tdee, ok := m["tdee"].(float64)
	require.True(t, ok)
	assert.Greater(t, tdee, bmr)

	dailyTarget, ok := m["dailyTarget"].(float64)
	require.True(t, ok)
	// MAINTAIN goal → dailyTarget == tdee
	assert.InDelta(t, tdee, dailyTarget, 0.01)

	require.NotNil(t, m["mealBreakdown"])
}

func TestE2E_CalculateTDEE_InvalidBody(t *testing.T) {
	requireE2E(t)

	resp := doPost(t, "/api/v1/nutrition/calculate-tdee", map[string]any{"gender": "MALE"})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()
}

func TestE2E_ListNutritionGoals(t *testing.T) {
	requireE2E(t)

	resp := doGet(t, "/api/v1/nutrition/goals")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	items := dataAsSlice(t, data.Data)
	assert.Greater(t, len(items), 0, "should have seeded nutrition goals")
}
