package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lehoanglong/whatdish/internal/user"
)

// AuthClient implements user.AuthProvider by delegating to Supabase Auth REST API.
type AuthClient struct {
	baseURL    string
	anonKey    string
	serviceKey string
	httpClient *http.Client
}

// NewAuthClient creates a new Supabase Auth client.
func NewAuthClient(baseURL, anonKey, serviceKey string) *AuthClient {
	return &AuthClient{
		baseURL:    baseURL,
		anonKey:    anonKey,
		serviceKey: serviceKey,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *AuthClient) Register(ctx context.Context, email, password, name string) (*user.AuthUser, *user.AuthTokens, error) {
	body := map[string]any{
		"email":    email,
		"password": password,
		"data": map[string]string{
			"name": name,
		},
	}

	resp, err := c.post(ctx, "/auth/v1/signup", body, c.anonKey)
	if err != nil {
		return nil, nil, fmt.Errorf("supabase.Register: %w", err)
	}

	var result authResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, nil, fmt.Errorf("supabase.Register: decode: %w", err)
	}

	if result.Error != "" {
		return nil, nil, fmt.Errorf("supabase.Register: %s", result.ErrorDescription)
	}

	authUser := &user.AuthUser{
		ID:    result.User.ID,
		Email: result.User.Email,
		Name:  name,
	}
	tokens := &user.AuthTokens{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}

	return authUser, tokens, nil
}

func (c *AuthClient) Login(ctx context.Context, email, password string) (*user.AuthUser, *user.AuthTokens, error) {
	body := map[string]string{
		"email":    email,
		"password": password,
	}

	resp, err := c.post(ctx, "/auth/v1/token?grant_type=password", body, c.anonKey)
	if err != nil {
		return nil, nil, fmt.Errorf("supabase.Login: %w", err)
	}

	var result authResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, nil, fmt.Errorf("supabase.Login: decode: %w", err)
	}

	if result.Error != "" {
		return nil, nil, fmt.Errorf("supabase.Login: %s", result.ErrorDescription)
	}

	authUser := &user.AuthUser{
		ID:    result.User.ID,
		Email: result.User.Email,
		Name:  result.User.UserMetadata.Name,
	}
	tokens := &user.AuthTokens{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}

	return authUser, tokens, nil
}

func (c *AuthClient) RefreshToken(ctx context.Context, refreshToken string) (*user.AuthTokens, error) {
	body := map[string]string{
		"refresh_token": refreshToken,
	}

	resp, err := c.post(ctx, "/auth/v1/token?grant_type=refresh_token", body, c.anonKey)
	if err != nil {
		return nil, fmt.Errorf("supabase.RefreshToken: %w", err)
	}

	var result authResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("supabase.RefreshToken: decode: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("supabase.RefreshToken: %s", result.ErrorDescription)
	}

	return &user.AuthTokens{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (c *AuthClient) VerifyToken(ctx context.Context, token string) (uuid.UUID, error) {
	u, err := c.GetUser(ctx, token)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("supabase.VerifyToken: %w", err)
	}
	return u.ID, nil
}

func (c *AuthClient) GetUser(ctx context.Context, token string) (*user.AuthUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/auth/v1/user", nil)
	if err != nil {
		return nil, fmt.Errorf("supabase.GetUser: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", c.anonKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("supabase.GetUser: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("supabase.GetUser: status %d", resp.StatusCode)
	}

	var result supabaseUser
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("supabase.GetUser: decode: %w", err)
	}

	return &user.AuthUser{
		ID:    result.ID,
		Email: result.Email,
		Name:  result.UserMetadata.Name,
	}, nil
}

func (c *AuthClient) post(ctx context.Context, path string, body any, apiKey string) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("supabase auth %s returned %d: %s", path, resp.StatusCode, respBody)
	}
	return respBody, nil
}

// Internal types for Supabase API responses

type authResponse struct {
	AccessToken      string       `json:"access_token"`
	RefreshToken     string       `json:"refresh_token"`
	ExpiresIn        int          `json:"expires_in"`
	User             supabaseUser `json:"user"`
	Error            string       `json:"error"`
	ErrorDescription string       `json:"error_description"`
}

type supabaseUser struct {
	ID           uuid.UUID    `json:"id"`
	Email        string       `json:"email"`
	UserMetadata userMetadata `json:"user_metadata"`
}

type userMetadata struct {
	Name string `json:"name"`
}
