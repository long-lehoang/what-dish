package engagement

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Helpers (mocks are defined in service_test.go and shared across test files)
// ---------------------------------------------------------------------------

// setUserID is a Gin middleware that injects a user ID into the context,
// simulating the auth middleware for tests.
func setUserID(userID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}
}

func setupEngagementRouter(h *Handler, userID *uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	if userID != nil {
		authed := r.Group("", setUserID(*userID))
		authed.POST("/favorites", h.HandleAddFavorite)
		authed.DELETE("/favorites/:recipe_id", h.HandleRemoveFavorite)
		authed.GET("/favorites", h.HandleListFavorites)
		authed.GET("/favorites/check", h.HandleCheckFavorites)
		authed.POST("/views", h.HandleRecordView)
	} else {
		r.POST("/favorites", h.HandleAddFavorite)
		r.DELETE("/favorites/:recipe_id", h.HandleRemoveFavorite)
		r.GET("/favorites", h.HandleListFavorites)
		r.GET("/favorites/check", h.HandleCheckFavorites)
		r.POST("/views", h.HandleRecordView)
	}

	return r
}

func parseEngagementBody(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var result map[string]any
	err := json.NewDecoder(w.Body).Decode(&result)
	assert.NoError(t, err, "response body should be valid JSON")
	return result
}

// ---------------------------------------------------------------------------
// Tests — POST /views (no auth required)
// ---------------------------------------------------------------------------

func TestHandleRecordView_OK(t *testing.T) {
	recipeID := uuid.New()
	var recorded bool
	viewRepo := &mockViewRepo{
		recordFn: func(_ context.Context, view *ViewHistory) error {
			recorded = true
			assert.Equal(t, recipeID, view.RecipeID)
			assert.Equal(t, "search", view.Source)
			assert.NotEmpty(t, view.SessionID)
			return nil
		},
	}
	svc := NewEngagementService(&mockFavoriteRepo{}, viewRepo, &mockRatingRepo{})
	h := NewHandler(svc)

	// No auth for views — anonymous user.
	router := setupEngagementRouter(h, nil)

	reqBody := RecordViewRequest{
		RecipeID: recipeID.String(),
		Source:   "search",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/views", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.True(t, recorded, "view should have been recorded")

	body := parseEngagementBody(t, w)
	data := body["data"].(map[string]any)
	assert.Equal(t, true, data["recorded"])
}

func TestHandleRecordView_WithSessionHeader(t *testing.T) {
	recipeID := uuid.New()
	customSession := "my-custom-session-123"

	var capturedSessionID string
	viewRepo := &mockViewRepo{
		recordFn: func(_ context.Context, view *ViewHistory) error {
			capturedSessionID = view.SessionID
			return nil
		},
	}
	svc := NewEngagementService(&mockFavoriteRepo{}, viewRepo, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, nil)

	reqBody := RecordViewRequest{
		RecipeID: recipeID.String(),
		Source:   "direct",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/views", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Session-ID", customSession)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, customSession, capturedSessionID)
}

func TestHandleRecordView_AuthenticatedUser(t *testing.T) {
	recipeID := uuid.New()
	userID := uuid.New()

	var capturedUserID *uuid.UUID
	viewRepo := &mockViewRepo{
		recordFn: func(_ context.Context, view *ViewHistory) error {
			capturedUserID = view.UserID
			return nil
		},
	}
	svc := NewEngagementService(&mockFavoriteRepo{}, viewRepo, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, &userID)

	reqBody := RecordViewRequest{
		RecipeID: recipeID.String(),
		Source:   "suggestion",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/views", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotNil(t, capturedUserID)
	assert.Equal(t, userID, *capturedUserID)
}

func TestHandleRecordView_InvalidBody(t *testing.T) {
	svc := NewEngagementService(&mockFavoriteRepo{}, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, nil)

	req := httptest.NewRequest(http.MethodPost, "/views", bytes.NewReader([]byte(`{bad json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	body := parseEngagementBody(t, w)
	assert.Equal(t, "invalid request body", body["message"])
}

func TestHandleRecordView_MissingRecipeID(t *testing.T) {
	svc := NewEngagementService(&mockFavoriteRepo{}, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, nil)

	// Missing recipeId field.
	req := httptest.NewRequest(http.MethodPost, "/views", bytes.NewReader([]byte(`{"source":"search"}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleRecordView_InvalidRecipeID(t *testing.T) {
	svc := NewEngagementService(&mockFavoriteRepo{}, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, nil)

	reqBody := RecordViewRequest{
		RecipeID: "not-a-uuid",
		Source:   "search",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/views", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleRecordView_ServiceError(t *testing.T) {
	recipeID := uuid.New()
	viewRepo := &mockViewRepo{
		recordFn: func(_ context.Context, _ *ViewHistory) error {
			return fmt.Errorf("database write failed")
		},
	}
	svc := NewEngagementService(&mockFavoriteRepo{}, viewRepo, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, nil)

	reqBody := RecordViewRequest{
		RecipeID: recipeID.String(),
		Source:   "search",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/views", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ---------------------------------------------------------------------------
// Tests — POST /favorites (auth required)
// ---------------------------------------------------------------------------

func TestHandleAddFavorite_OK(t *testing.T) {
	userID := uuid.New()
	recipeID := uuid.New()

	favRepo := &mockFavoriteRepo{
		addFn: func(_ context.Context, uid, rid uuid.UUID) error {
			assert.Equal(t, userID, uid)
			assert.Equal(t, recipeID, rid)
			return nil
		},
	}
	svc := NewEngagementService(favRepo, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, &userID)

	reqBody := AddFavoriteRequest{RecipeID: recipeID.String()}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/favorites", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	body := parseEngagementBody(t, w)
	data := body["data"].(map[string]any)
	assert.Equal(t, recipeID.String(), data["recipeId"])
	assert.Contains(t, data, "createdAt")
}

func TestHandleAddFavorite_NoAuth(t *testing.T) {
	svc := NewEngagementService(&mockFavoriteRepo{}, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	// No user ID — simulates unauthenticated request.
	router := setupEngagementRouter(h, nil)

	recipeID := uuid.New()
	reqBody := AddFavoriteRequest{RecipeID: recipeID.String()}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/favorites", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	body := parseEngagementBody(t, w)
	assert.Equal(t, "missing user identity", body["message"])
}

func TestHandleAddFavorite_InvalidRecipeID(t *testing.T) {
	userID := uuid.New()
	svc := NewEngagementService(&mockFavoriteRepo{}, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, &userID)

	reqBody := AddFavoriteRequest{RecipeID: "not-a-uuid"}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/favorites", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleAddFavorite_InvalidBody(t *testing.T) {
	userID := uuid.New()
	svc := NewEngagementService(&mockFavoriteRepo{}, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, &userID)

	req := httptest.NewRequest(http.MethodPost, "/favorites", bytes.NewReader([]byte(`{bad`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---------------------------------------------------------------------------
// Tests — DELETE /favorites/:recipe_id (auth required)
// ---------------------------------------------------------------------------

func TestHandleRemoveFavorite_OK(t *testing.T) {
	userID := uuid.New()
	recipeID := uuid.New()

	var removeCalled bool
	favRepo := &mockFavoriteRepo{
		removeFn: func(_ context.Context, uid, rid uuid.UUID) error {
			removeCalled = true
			assert.Equal(t, userID, uid)
			assert.Equal(t, recipeID, rid)
			return nil
		},
	}
	svc := NewEngagementService(favRepo, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, &userID)

	req := httptest.NewRequest(http.MethodDelete, "/favorites/"+recipeID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.True(t, removeCalled)
}

func TestHandleRemoveFavorite_NoAuth(t *testing.T) {
	svc := NewEngagementService(&mockFavoriteRepo{}, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, nil)

	req := httptest.NewRequest(http.MethodDelete, "/favorites/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandleRemoveFavorite_InvalidID(t *testing.T) {
	userID := uuid.New()
	svc := NewEngagementService(&mockFavoriteRepo{}, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, &userID)

	req := httptest.NewRequest(http.MethodDelete, "/favorites/bad-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---------------------------------------------------------------------------
// Tests — GET /favorites (auth required)
// ---------------------------------------------------------------------------

func TestHandleListFavorites_OK(t *testing.T) {
	userID := uuid.New()
	recipeID := uuid.New()

	favRepo := &mockFavoriteRepo{
		listByUserFn: func(_ context.Context, uid uuid.UUID, page, pageSize int) ([]Favorite, int64, error) {
			assert.Equal(t, userID, uid)
			assert.Equal(t, 1, page)
			assert.Equal(t, 20, pageSize)
			return []Favorite{
				{ID: uuid.New(), UserID: userID, RecipeID: recipeID, CreatedAt: time.Now().UTC()},
			}, 1, nil
		},
	}
	svc := NewEngagementService(favRepo, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, &userID)

	req := httptest.NewRequest(http.MethodGet, "/favorites", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseEngagementBody(t, w)
	data, ok := body["data"].([]any)
	assert.True(t, ok)
	assert.Len(t, data, 1)

	first := data[0].(map[string]any)
	assert.Equal(t, recipeID.String(), first["recipeId"])

	pagination := body["pagination"].(map[string]any)
	assert.Equal(t, float64(1), pagination["total"])
}

func TestHandleListFavorites_NoAuth(t *testing.T) {
	svc := NewEngagementService(&mockFavoriteRepo{}, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, nil)

	req := httptest.NewRequest(http.MethodGet, "/favorites", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ---------------------------------------------------------------------------
// Tests — GET /favorites/check (auth required)
// ---------------------------------------------------------------------------

func TestHandleCheckFavorites_OK(t *testing.T) {
	userID := uuid.New()
	recipeID1 := uuid.New()
	recipeID2 := uuid.New()

	favRepo := &mockFavoriteRepo{
		checkFn: func(_ context.Context, uid uuid.UUID, ids []uuid.UUID) (map[uuid.UUID]bool, error) {
			assert.Equal(t, userID, uid)
			assert.Len(t, ids, 2)
			return map[uuid.UUID]bool{
				recipeID1: true,
				recipeID2: false,
			}, nil
		},
	}
	svc := NewEngagementService(favRepo, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, &userID)

	url := fmt.Sprintf("/favorites/check?recipe_ids=%s,%s", recipeID1, recipeID2)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	body := parseEngagementBody(t, w)
	data := body["data"].(map[string]any)
	assert.Equal(t, true, data[recipeID1.String()])
	assert.Equal(t, false, data[recipeID2.String()])
}

func TestHandleCheckFavorites_MissingRecipeIDs(t *testing.T) {
	userID := uuid.New()
	svc := NewEngagementService(&mockFavoriteRepo{}, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, &userID)

	req := httptest.NewRequest(http.MethodGet, "/favorites/check", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	body := parseEngagementBody(t, w)
	assert.Equal(t, "recipe_ids query parameter is required", body["message"])
}

func TestHandleCheckFavorites_InvalidUUID(t *testing.T) {
	userID := uuid.New()
	svc := NewEngagementService(&mockFavoriteRepo{}, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, &userID)

	req := httptest.NewRequest(http.MethodGet, "/favorites/check?recipe_ids=not-a-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCheckFavorites_NoAuth(t *testing.T) {
	svc := NewEngagementService(&mockFavoriteRepo{}, &mockViewRepo{}, &mockRatingRepo{})
	h := NewHandler(svc)
	router := setupEngagementRouter(h, nil)

	req := httptest.NewRequest(http.MethodGet, "/favorites/check?recipe_ids="+uuid.New().String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
