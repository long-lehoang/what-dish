package e2e_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_Register(t *testing.T) {
	requireE2E(t)

	body := map[string]any{
		"email":    "newuser@example.com",
		"password": "securepassword123",
		"name":     "New User",
	}

	resp := doPost(t, "/api/v1/auth/register", body)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	data := parseData(t, resp)
	m := dataAsMap(t, data.Data)
	require.NotNil(t, m["user"])
	require.NotNil(t, m["tokens"])

	user := m["user"].(map[string]any)
	assert.Equal(t, "newuser@example.com", user["email"])
	assert.Equal(t, "New User", user["name"])

	tokens := m["tokens"].(map[string]any)
	require.NotEmpty(t, tokens["accessToken"])
	require.NotEmpty(t, tokens["refreshToken"])
}

func TestE2E_Register_InvalidBody(t *testing.T) {
	requireE2E(t)

	// Missing password.
	body := map[string]any{
		"email": "bad@example.com",
		"name":  "Bad User",
	}

	resp := doPost(t, "/api/v1/auth/register", body)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()
}

func TestE2E_Login(t *testing.T) {
	requireE2E(t)

	// Register first.
	regBody := map[string]any{
		"email":    "logintest@example.com",
		"password": "testpassword123",
		"name":     "Login User",
	}
	regResp := doPost(t, "/api/v1/auth/register", regBody)
	require.Equal(t, http.StatusCreated, regResp.StatusCode)
	regResp.Body.Close()

	// Login.
	loginBody := map[string]any{
		"email":    "logintest@example.com",
		"password": "testpassword123",
	}
	resp := doPost(t, "/api/v1/auth/login", loginBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	m := dataAsMap(t, data.Data)
	require.NotNil(t, m["tokens"])

	tokens := m["tokens"].(map[string]any)
	require.NotEmpty(t, tokens["accessToken"])
}

func TestE2E_Refresh(t *testing.T) {
	requireE2E(t)

	body := map[string]any{
		"refreshToken": "refresh-" + "00000000-0000-0000-0000-000000000001",
	}

	resp := doPost(t, "/api/v1/auth/refresh", body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data := parseData(t, resp)
	m := dataAsMap(t, data.Data)
	require.NotEmpty(t, m["accessToken"])
	require.NotEmpty(t, m["refreshToken"])
}
