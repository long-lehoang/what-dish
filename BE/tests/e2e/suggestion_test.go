package e2e_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_RandomSuggestion(t *testing.T) {
	requireE2E(t)

	resp := doPost(t, "/api/v1/suggestions/random", map[string]any{})
	if resp.StatusCode != http.StatusOK {
		errResp := parseError(t, resp)
		t.Fatalf("expected 200 but got %d: %s", resp.StatusCode, errResp.Message)
	}

	data := parseData(t, resp)
	m := dataAsMap(t, data.Data)
	require.NotEmpty(t, m["sessionId"])
	assert.Equal(t, "RANDOM", m["type"])
	require.NotNil(t, m["recipes"])
}

func TestE2E_RandomSuggestion_EmptyBody(t *testing.T) {
	requireE2E(t)

	// Send an empty JSON object (not nil).
	resp := doPost(t, "/api/v1/suggestions/random", map[string]any{})
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestE2E_CalorieSuggestion_InvalidBody(t *testing.T) {
	requireE2E(t)

	// Missing required targetCalories field.
	resp := doPost(t, "/api/v1/suggestions/by-calories", map[string]any{})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()
}

func TestE2E_SuggestionHistory_Unauthorized(t *testing.T) {
	requireE2E(t)

	// No auth header → 401.
	resp := doGet(t, "/api/v1/suggestions/history")
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()
}

func TestE2E_SuggestionHistory(t *testing.T) {
	requireE2E(t)

	userID := uuid.New()

	// Seed sessions directly into the DB since the suggestion endpoint
	// is public and doesn't attach a user_id.
	ctx := context.Background()
	pool, err := pgxPoolFromConnStr(ctx, infra.pgConnStr)
	require.NoError(t, err)
	defer pool.Close()

	now := time.Now().UTC()
	for i := 0; i < 2; i++ {
		_, err := pool.Exec(ctx,
			`INSERT INTO suggestion_sessions (id, user_id, session_type, input_params, result_recipe_ids, created_at)
			 VALUES ($1, $2, 'RANDOM', '{}', '{}', $3)`,
			uuid.New(), userID, now.Add(time.Duration(i)*time.Second),
		)
		require.NoError(t, err)
	}

	// Get history.
	resp := doAuthedGet(t, "/api/v1/suggestions/history?page=1&pageSize=10", userID)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := parseList(t, resp)
	assert.Equal(t, int64(2), result.Pagination.Total)
}
