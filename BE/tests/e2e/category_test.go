package e2e_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_ListCategories(t *testing.T) {
	requireE2E(t)

	resp := doGet(t, "/api/v1/categories")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	items := dataAsSlice(t, data.Data)
	assert.Greater(t, len(items), 0, "should have seeded categories")
}

func TestE2E_ListCategories_ByType(t *testing.T) {
	requireE2E(t)

	resp := doGet(t, "/api/v1/categories?type=DISH_TYPE")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	items := dataAsSlice(t, data.Data)
	for _, cat := range items {
		assert.Equal(t, "DISH_TYPE", cat["type"])
	}
}

func TestE2E_ListTags(t *testing.T) {
	requireE2E(t)

	resp := doGet(t, "/api/v1/tags")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	items := dataAsSlice(t, data.Data)
	assert.Greater(t, len(items), 0, "should have seeded tags")
}
