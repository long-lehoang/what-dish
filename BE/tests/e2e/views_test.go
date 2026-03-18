package e2e_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_RecordView(t *testing.T) {
	requireE2E(t)
	require.NotEmpty(t, recipeIDs)

	body := map[string]any{
		"recipeId": recipeIDs[0].String(),
		"source":   "search",
	}

	resp := doPost(t, "/api/v1/views", body)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	data := parseData(t, resp)
	m := dataAsMap(t, data.Data)
	assert.Equal(t, true, m["recorded"])
}

func TestE2E_RecordView_InvalidBody(t *testing.T) {
	requireE2E(t)

	resp := doPost(t, "/api/v1/views", map[string]any{})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()
}
