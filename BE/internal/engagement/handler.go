package engagement

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/lehoanglong/whatdish/internal/shared/middleware"
	"github.com/lehoanglong/whatdish/internal/shared/response"
)

// Handler exposes HTTP endpoints for the engagement bounded context.
type Handler struct {
	svc      *EngagementService
	validate *validator.Validate
}

// NewHandler creates a new engagement Handler.
func NewHandler(svc *EngagementService) *Handler {
	return &Handler{
		svc:      svc,
		validate: validator.New(),
	}
}

// HandleAddFavorite handles POST /favorites.
func (h *Handler) HandleAddFavorite(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.ErrMsg(c, http.StatusUnauthorized, "missing user identity")
		return
	}

	var req AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, err.Error())
		return
	}

	recipeID, err := uuid.Parse(req.RecipeID)
	if err != nil {
		response.ErrMsg(c, http.StatusBadRequest, "invalid recipe id")
		return
	}

	result, err := h.svc.AddFavorite(c.Request.Context(), userID, recipeID)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.Created(c, result)
}

// HandleRemoveFavorite handles DELETE /favorites/:recipe_id.
func (h *Handler) HandleRemoveFavorite(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.ErrMsg(c, http.StatusUnauthorized, "missing user identity")
		return
	}

	recipeID, err := uuid.Parse(c.Param("recipe_id"))
	if err != nil {
		response.ErrMsg(c, http.StatusBadRequest, "invalid recipe id")
		return
	}

	if err := h.svc.RemoveFavorite(c.Request.Context(), userID, recipeID); err != nil {
		response.Err(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// HandleListFavorites handles GET /favorites.
func (h *Handler) HandleListFavorites(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.ErrMsg(c, http.StatusUnauthorized, "missing user identity")
		return
	}

	page, pageSize := response.ParsePagination(c)

	results, total, err := h.svc.ListFavorites(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.List(c, results, page, pageSize, total)
}

// HandleCheckFavorites handles GET /favorites/check?recipe_ids=uuid1,uuid2.
func (h *Handler) HandleCheckFavorites(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.ErrMsg(c, http.StatusUnauthorized, "missing user identity")
		return
	}

	idsParam := c.Query("recipe_ids")
	if idsParam == "" {
		response.ErrMsg(c, http.StatusBadRequest, "recipe_ids query parameter is required")
		return
	}

	parts := strings.Split(idsParam, ",")
	ids := make([]uuid.UUID, 0, len(parts))
	for _, p := range parts {
		id, err := uuid.Parse(strings.TrimSpace(p))
		if err != nil {
			response.ErrMsg(c, http.StatusBadRequest, "invalid uuid in recipe_ids: "+p)
			return
		}
		ids = append(ids, id)
	}

	result, err := h.svc.CheckFavorites(c.Request.Context(), userID, ids)
	if err != nil {
		response.Err(c, err)
		return
	}

	// Convert map[uuid.UUID]bool to map[string]bool for JSON.
	out := make(map[string]bool, len(result))
	for id, v := range result {
		out[id.String()] = v
	}

	response.OK(c, out)
}

// HandleRecordView handles POST /views.
// Supports both authenticated and anonymous users (session-based).
func (h *Handler) HandleRecordView(c *gin.Context) {
	var req RecordViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, err.Error())
		return
	}

	recipeID, err := uuid.Parse(req.RecipeID)
	if err != nil {
		response.ErrMsg(c, http.StatusBadRequest, "invalid recipe id")
		return
	}

	// Extract optional user ID (may not be present for anonymous views).
	var userID *uuid.UUID
	if uid, ok := middleware.GetUserID(c); ok {
		userID = &uid
	}

	// Use a session ID from header or generate one.
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	if err := h.svc.RecordView(c.Request.Context(), userID, sessionID, recipeID, req.Source); err != nil {
		response.Err(c, err)
		return
	}

	response.Created(c, gin.H{"recorded": true})
}
