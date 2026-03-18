package e2e_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_GetProfile_Unauthorized(t *testing.T) {
	requireE2E(t)

	resp := doGet(t, "/api/v1/users/me")
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()
}

func TestE2E_UpdateAndGetProfile(t *testing.T) {
	requireE2E(t)

	userID := uuid.New()

	// Update profile.
	body := map[string]any{
		"gender":        "MALE",
		"age":           30,
		"heightCm":      175.0,
		"weightKg":      70.0,
		"activityLevel": "MODERATE",
		"goal":          "MAINTAIN",
	}

	resp := doAuthedPut(t, "/api/v1/users/me/profile", body, userID)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	m := dataAsMap(t, data.Data)
	assert.Equal(t, "MALE", m["gender"])
	assert.Equal(t, float64(30), m["age"])

	// Verify TDEE was calculated.
	require.NotNil(t, m["bmr"])
	require.NotNil(t, m["tdee"])
	require.NotNil(t, m["dailyTarget"])

	bmr := m["bmr"].(float64)
	assert.Greater(t, bmr, float64(0))

	// GET profile.
	resp = doAuthedGet(t, "/api/v1/users/me", userID)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data = parseData(t, resp)
	m = dataAsMap(t, data.Data)
	assert.Equal(t, "MALE", m["gender"])
	assert.Equal(t, float64(30), m["age"])
}

func TestE2E_UpdateProfile_WithAllergies(t *testing.T) {
	requireE2E(t)

	userID := uuid.New()

	body := map[string]any{
		"gender":        "FEMALE",
		"age":           25,
		"heightCm":      165.0,
		"weightKg":      55.0,
		"activityLevel": "LIGHT",
		"goal":          "LOSE_WEIGHT",
		"allergies": []map[string]string{
			{"ingredientName": "Peanut", "allergyType": "ALLERGY"},
			{"ingredientName": "Shrimp", "allergyType": "DISLIKE"},
		},
	}

	resp := doAuthedPut(t, "/api/v1/users/me/profile", body, userID)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	m := dataAsMap(t, data.Data)

	allergies, ok := m["allergies"].([]any)
	require.True(t, ok)
	assert.Len(t, allergies, 2)
}

func TestE2E_UpdateProfile_InvalidBody(t *testing.T) {
	requireE2E(t)

	userID := uuid.New()

	// Invalid gender value.
	body := map[string]any{
		"gender": "INVALID",
	}

	resp := doAuthedPut(t, "/api/v1/users/me/profile", body, userID)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()
}
