package e2e_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_ListRecipes(t *testing.T) {
	requireE2E(t)

	resp := doGet(t, "/api/v1/recipes")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := parseList(t, resp)
	items := dataAsSlice(t, result.Data)
	assert.GreaterOrEqual(t, len(items), 5, "should have at least 5 seeded recipes")
	assert.Greater(t, result.Pagination.Total, int64(0))
}

func TestE2E_ListRecipes_Pagination(t *testing.T) {
	requireE2E(t)

	resp := doGet(t, "/api/v1/recipes?page=1&pageSize=2")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := parseList(t, resp)
	items := dataAsSlice(t, result.Data)
	assert.Len(t, items, 2)
	assert.Equal(t, 1, result.Pagination.Page)
	assert.Equal(t, 2, result.Pagination.PageSize)
	assert.GreaterOrEqual(t, result.Pagination.Total, int64(5))
}

func TestE2E_GetRecipeByID(t *testing.T) {
	requireE2E(t)
	require.NotEmpty(t, recipeIDs)

	resp := doGet(t, "/api/v1/recipes/"+recipeIDs[0].String())
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	m := dataAsMap(t, data.Data)
	assert.Equal(t, "Phở Bò", m["name"])
	assert.Equal(t, "pho-bo", m["slug"])
}

func TestE2E_GetRecipeBySlug(t *testing.T) {
	requireE2E(t)

	resp := doGet(t, "/api/v1/recipes/bun-cha")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	m := dataAsMap(t, data.Data)
	assert.Equal(t, "Bún Chả", m["name"])
}

func TestE2E_GetRecipe_NotFound(t *testing.T) {
	requireE2E(t)

	resp := doGet(t, "/api/v1/recipes/"+uuid.New().String())
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()
}

func TestE2E_GetRandomRecipe(t *testing.T) {
	requireE2E(t)

	resp := doGet(t, "/api/v1/recipes/random")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	m := dataAsMap(t, data.Data)
	require.NotEmpty(t, m["name"])
	require.NotEmpty(t, m["slug"])
}

func TestE2E_SearchRecipes(t *testing.T) {
	requireE2E(t)

	resp := doGet(t, "/api/v1/recipes/search?q=pho")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	result := parseList(t, resp)
	items := dataAsSlice(t, result.Data)
	assert.GreaterOrEqual(t, len(items), 1)
	assert.Equal(t, "Phở Bò", items[0]["name"])
}
