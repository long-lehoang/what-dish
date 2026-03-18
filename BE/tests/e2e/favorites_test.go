package e2e_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_AddFavorite(t *testing.T) {
	requireE2E(t)
	require.NotEmpty(t, recipeIDs)

	userID := uuid.New()
	body := map[string]any{"recipeId": recipeIDs[0].String()}

	resp := doAuthedPost(t, "/api/v1/favorites", body, userID)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	data := parseData(t, resp)
	m := dataAsMap(t, data.Data)
	assert.Equal(t, recipeIDs[0].String(), m["recipeId"])
	require.NotEmpty(t, m["createdAt"])
}

func TestE2E_AddFavorite_Unauthorized(t *testing.T) {
	requireE2E(t)

	body := map[string]any{"recipeId": uuid.New().String()}
	resp := doPost(t, "/api/v1/favorites", body)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()
}

func TestE2E_ListFavorites(t *testing.T) {
	requireE2E(t)
	require.True(t, len(recipeIDs) >= 2)

	userID := uuid.New()

	// Add 2 favorites.
	doAuthedPost(t, "/api/v1/favorites", map[string]any{"recipeId": recipeIDs[0].String()}, userID).Body.Close()
	doAuthedPost(t, "/api/v1/favorites", map[string]any{"recipeId": recipeIDs[1].String()}, userID).Body.Close()

	// List.
	resp := doAuthedGet(t, "/api/v1/favorites?page=1&pageSize=10", userID)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := parseList(t, resp)
	items := dataAsSlice(t, result.Data)
	assert.Len(t, items, 2)
	assert.Equal(t, int64(2), result.Pagination.Total)
}

func TestE2E_RemoveFavorite(t *testing.T) {
	requireE2E(t)
	require.NotEmpty(t, recipeIDs)

	userID := uuid.New()
	recipeID := recipeIDs[0]

	// Add.
	doAuthedPost(t, "/api/v1/favorites", map[string]any{"recipeId": recipeID.String()}, userID).Body.Close()

	// Remove.
	resp := doAuthedDelete(t, "/api/v1/favorites/"+recipeID.String(), userID)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()

	// Verify removed.
	resp = doAuthedGet(t, "/api/v1/favorites?page=1&pageSize=10", userID)
	result := parseList(t, resp)
	items := dataAsSlice(t, result.Data)
	assert.Len(t, items, 0)
}

func TestE2E_CheckFavorites(t *testing.T) {
	requireE2E(t)
	require.True(t, len(recipeIDs) >= 2)

	userID := uuid.New()
	r1 := recipeIDs[0]
	r2 := recipeIDs[1]

	// Only favorite r1.
	doAuthedPost(t, "/api/v1/favorites", map[string]any{"recipeId": r1.String()}, userID).Body.Close()

	resp := doAuthedGet(t, "/api/v1/favorites/check?recipe_ids="+r1.String()+","+r2.String(), userID)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	m := dataAsMap(t, data.Data)
	assert.Equal(t, true, m[r1.String()])
	assert.Equal(t, false, m[r2.String()])
}
