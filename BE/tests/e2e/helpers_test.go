package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// httpClient is a shared HTTP client for E2E tests.
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// --- HTTP helpers ---

func doGet(t *testing.T, path string) *http.Response {
	t.Helper()
	url := infra.appBaseURL + path
	resp, err := httpClient.Get(url)
	require.NoError(t, err, "GET %s", path)
	return resp
}

func doPost(t *testing.T, path string, body any) *http.Response {
	t.Helper()
	url := infra.appBaseURL + path
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(data)
	}
	resp, err := httpClient.Post(url, "application/json", reader)
	require.NoError(t, err, "POST %s", path)
	return resp
}

func doAuthedGet(t *testing.T, path string, userID uuid.UUID) *http.Response {
	t.Helper()
	url := infra.appBaseURL + path
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer test-token-"+userID.String())
	resp, err := httpClient.Do(req)
	require.NoError(t, err, "authed GET %s", path)
	return resp
}

func doAuthedPost(t *testing.T, path string, body any, userID uuid.UUID) *http.Response {
	t.Helper()
	url := infra.appBaseURL + path
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequest("POST", url, reader)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token-"+userID.String())
	resp, err := httpClient.Do(req)
	require.NoError(t, err, "authed POST %s", path)
	return resp
}

func doAuthedPut(t *testing.T, path string, body any, userID uuid.UUID) *http.Response {
	t.Helper()
	url := infra.appBaseURL + path
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequest("PUT", url, reader)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token-"+userID.String())
	resp, err := httpClient.Do(req)
	require.NoError(t, err, "authed PUT %s", path)
	return resp
}

func doAuthedDelete(t *testing.T, path string, userID uuid.UUID) *http.Response {
	t.Helper()
	url := infra.appBaseURL + path
	req, err := http.NewRequest("DELETE", url, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer test-token-"+userID.String())
	resp, err := httpClient.Do(req)
	require.NoError(t, err, "authed DELETE %s", path)
	return resp
}

// --- Response parsing ---

type apiDataResponse struct {
	Data json.RawMessage `json:"data"`
}

type apiListResponse struct {
	Data       json.RawMessage `json:"data"`
	Pagination struct {
		Page       int   `json:"page"`
		PageSize   int   `json:"pageSize"`
		Total      int64 `json:"total"`
		TotalPages int   `json:"totalPages"`
	} `json:"pagination"`
}

type apiErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func parseData(t *testing.T, resp *http.Response) apiDataResponse {
	t.Helper()
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var r apiDataResponse
	require.NoError(t, json.Unmarshal(body, &r), "body: %s", string(body))
	return r
}

func parseList(t *testing.T, resp *http.Response) apiListResponse {
	t.Helper()
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var r apiListResponse
	require.NoError(t, json.Unmarshal(body, &r), "body: %s", string(body))
	return r
}

func parseError(t *testing.T, resp *http.Response) apiErrorResponse {
	t.Helper()
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var r apiErrorResponse
	require.NoError(t, json.Unmarshal(body, &r), "body: %s", string(body))
	return r
}

func dataAsMap(t *testing.T, raw json.RawMessage) map[string]any {
	t.Helper()
	var m map[string]any
	require.NoError(t, json.Unmarshal(raw, &m))
	return m
}

func dataAsSlice(t *testing.T, raw json.RawMessage) []map[string]any {
	t.Helper()
	var s []map[string]any
	require.NoError(t, json.Unmarshal(raw, &s))
	return s
}

// --- Seed helpers ---

// seedTestRecipes inserts test recipes with ingredients, nutrition, and tags.
// Returns a slice of recipe IDs for use in tests.
func seedTestRecipes(ctx context.Context, connStr string) []uuid.UUID {
	pool, err := pgxPoolFromConnStr(ctx, connStr)
	if err != nil {
		fmt.Printf("e2e: seed recipes pool: %v\n", err)
		return nil
	}
	defer pool.Close()

	recipes := []struct {
		name       string
		slug       string
		difficulty string
		prepTime   int
		cookTime   int
		servings   int
	}{
		{"Phở Bò", "pho-bo", "MEDIUM", 30, 120, 4},
		{"Bún Chả", "bun-cha", "EASY", 20, 15, 2},
		{"Cơm Tấm", "com-tam", "EASY", 15, 30, 1},
		{"Bánh Mì", "banh-mi", "EASY", 10, 5, 1},
		{"Gỏi Cuốn", "goi-cuon", "EASY", 20, 0, 4},
	}

	ids := make([]uuid.UUID, 0, len(recipes))
	for _, r := range recipes {
		var id uuid.UUID
		err := pool.QueryRow(ctx,
			`INSERT INTO recipes (name, slug, status, difficulty, prep_time, cook_time, servings)
			 VALUES ($1, $2, 'PUBLISHED', $3, $4, $5, $6)
			 ON CONFLICT (slug) DO UPDATE SET
				status = 'PUBLISHED',
				difficulty = EXCLUDED.difficulty,
				prep_time = EXCLUDED.prep_time,
				cook_time = EXCLUDED.cook_time,
				servings = EXCLUDED.servings
			 RETURNING id`,
			r.name, r.slug, r.difficulty, r.prepTime, r.cookTime, r.servings,
		).Scan(&id)
		if err != nil {
			fmt.Printf("e2e: seed recipe %s: %v\n", r.name, err)
			continue
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		fmt.Println("e2e: no recipes seeded, skipping nutrition/ingredient seeding")
		return nil
	}

	// Add nutrition for the first 3 recipes.
	nutritionData := []struct {
		idx      int
		calories float64
		protein  float64
		carbs    float64
		fat      float64
	}{
		{0, 450, 30, 55, 10}, // Phở Bò
		{1, 550, 25, 40, 20}, // Bún Chả
		{2, 600, 20, 70, 15}, // Cơm Tấm
	}

	for _, n := range nutritionData {
		if n.idx >= len(ids) {
			continue
		}
		_, err := pool.Exec(ctx,
			`INSERT INTO nutrition_recipe (id, recipe_id, calories, protein, carbs, fat, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
			 ON CONFLICT (recipe_id) DO UPDATE SET
				calories = EXCLUDED.calories,
				protein = EXCLUDED.protein,
				carbs = EXCLUDED.carbs,
				fat = EXCLUDED.fat,
				updated_at = EXCLUDED.updated_at`,
			uuid.New(), ids[n.idx], n.calories, n.protein, n.carbs, n.fat, time.Now().UTC(),
		)
		if err != nil {
			fmt.Printf("e2e: seed nutrition for recipe %d: %v\n", n.idx, err)
		}
	}

	// Add some ingredients for the first recipe (replace existing).
	_, _ = pool.Exec(ctx, `DELETE FROM recipe_ingredients WHERE recipe_id = $1`, ids[0])
	ingredients := []struct {
		name   string
		amount string
		unit   string
	}{
		{"Bánh phở", "500", "g"},
		{"Thịt bò", "300", "g"},
		{"Hành tây", "1", "củ"},
	}
	for i, ing := range ingredients {
		_, err := pool.Exec(ctx,
			`INSERT INTO recipe_ingredients (id, recipe_id, name, amount, unit, sort_order)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			uuid.New(), ids[0], ing.name, ing.amount, ing.unit, i+1,
		)
		if err != nil {
			fmt.Printf("e2e: seed ingredient: %v\n", err)
		}
	}

	return ids
}

// --- Utility ---

func pgxPoolFromConnStr(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, connStr)
}

func readFileContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
