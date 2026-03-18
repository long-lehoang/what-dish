package e2e_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_Health(t *testing.T) {
	requireE2E(t)

	resp := doGet(t, "/health")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var m map[string]any
	require.NoError(t, json.Unmarshal(body, &m))
	assert.Equal(t, "ok", m["status"])
	require.NotEmpty(t, m["timestamp"])
}
