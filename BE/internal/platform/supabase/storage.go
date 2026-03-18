package supabase

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// StorageClient provides access to Supabase Storage.
type StorageClient struct {
	baseURL    string
	serviceKey string
	httpClient *http.Client
}

// NewStorageClient creates a new Supabase Storage client.
func NewStorageClient(baseURL, serviceKey string) *StorageClient {
	return &StorageClient{
		baseURL:    baseURL,
		serviceKey: serviceKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Upload uploads a file to the specified bucket and path.
func (c *StorageClient) Upload(ctx context.Context, bucket, path string, body io.Reader, contentType string) (string, error) {
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", c.baseURL, bucket, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return "", fmt.Errorf("supabase.Upload: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.serviceKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "true")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("supabase.Upload: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("supabase.Upload: status %d: %s", resp.StatusCode, string(respBody))
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", c.baseURL, bucket, path)
	return publicURL, nil
}

// GetPublicURL returns the public URL for a file in storage.
func (c *StorageClient) GetPublicURL(bucket, path string) string {
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", c.baseURL, bucket, path)
}
