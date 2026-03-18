package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const (
	baseURL    = "https://api.notion.com/v1"
	apiVersion = "2022-06-28"
)

// Client is an HTTP client for the Notion API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Notion API client.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// QueryDatabase queries a Notion database with optional filters and pagination.
func (c *Client) QueryDatabase(ctx context.Context, databaseID string, startCursor *string) (*DatabaseQueryResponse, error) {
	body := map[string]any{
		"filter": map[string]any{
			"property": "Status",
			"select": map[string]string{
				"equals": "Published",
			},
		},
		"page_size": 100,
	}
	if startCursor != nil {
		body["start_cursor"] = *startCursor
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("notion.QueryDatabase: marshal body: %w", err)
	}

	url := fmt.Sprintf("%s/databases/%s/query", baseURL, databaseID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("notion.QueryDatabase: create request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("notion.QueryDatabase: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		slog.Error("notion API error", "status", resp.StatusCode, "body", string(respBody))
		return nil, fmt.Errorf("notion.QueryDatabase: API returned status %d", resp.StatusCode)
	}

	var result DatabaseQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("notion.QueryDatabase: decode response: %w", err)
	}

	return &result, nil
}

// GetBlockChildren retrieves all child blocks of a page or block.
func (c *Client) GetBlockChildren(ctx context.Context, blockID string, startCursor *string) (*BlocksResponse, error) {
	url := fmt.Sprintf("%s/blocks/%s/children?page_size=100", baseURL, blockID)
	if startCursor != nil {
		url += "&start_cursor=" + *startCursor
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("notion.GetBlockChildren: create request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("notion.GetBlockChildren: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		slog.Error("notion API error", "status", resp.StatusCode, "body", string(respBody))
		return nil, fmt.Errorf("notion.GetBlockChildren: API returned status %d", resp.StatusCode)
	}

	var result BlocksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("notion.GetBlockChildren: decode response: %w", err)
	}

	return &result, nil
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Notion-Version", apiVersion)
	req.Header.Set("Content-Type", "application/json")
}
