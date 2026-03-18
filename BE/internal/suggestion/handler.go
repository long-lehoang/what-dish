package suggestion

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lehoanglong/whatdish/internal/shared/middleware"
	"github.com/lehoanglong/whatdish/internal/shared/response"
)

// SuggestionHandler handles HTTP requests for the suggestion bounded context.
type SuggestionHandler struct {
	service *SuggestionService
}

// NewSuggestionHandler creates a new SuggestionHandler.
func NewSuggestionHandler(service *SuggestionService) *SuggestionHandler {
	return &SuggestionHandler{service: service}
}

// HandleRandomSuggestion handles POST /suggestions/random.
func (h *SuggestionHandler) HandleRandomSuggestion(c *gin.Context) {
	var req RandomSuggestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Allow empty body for random suggestions.
		req = RandomSuggestionRequest{}
	}

	params := make(map[string]any)
	if req.Filters != nil {
		params["filters"] = req.Filters
	}

	// Attach user ID if authenticated.
	if userID, ok := middleware.GetUserID(c); ok {
		params["userId"] = userID.String()
	}

	result, err := h.service.Suggest(c.Request.Context(), "RANDOM", params)
	if err != nil {
		slog.Error("suggestion.HandleRandomSuggestion", "error", err)
		response.Err(c, err)
		return
	}

	response.OK(c, toSuggestionResponse(result))
}

// HandleCalorieSuggestion handles POST /suggestions/by-calories.
func (h *SuggestionHandler) HandleCalorieSuggestion(c *gin.Context) {
	var req CalorieSuggestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, err.Error())
		return
	}

	params := map[string]any{
		"target_calories": req.TargetCalories,
		"meal_type":       req.MealType,
	}
	if req.TolerancePct != nil {
		params["tolerance_pct"] = *req.TolerancePct
	}
	if req.Filters != nil {
		params["filters"] = req.Filters
	}
	if userID, ok := middleware.GetUserID(c); ok {
		params["userId"] = userID.String()
	}

	result, err := h.service.Suggest(c.Request.Context(), "BY_CALORIES", params)
	if err != nil {
		slog.Error("suggestion.HandleCalorieSuggestion", "error", err)
		response.Err(c, err)
		return
	}

	response.OK(c, toSuggestionResponse(result))
}

// HandleGroupSuggestion handles POST /suggestions/by-group.
func (h *SuggestionHandler) HandleGroupSuggestion(c *gin.Context) {
	var req GroupSuggestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, err.Error())
		return
	}

	params := map[string]any{
		"group_size": req.GroupSize,
		"group_type": req.GroupType,
		"meal_type":  req.MealType,
	}
	if req.Filters != nil {
		params["filters"] = req.Filters
	}
	if userID, ok := middleware.GetUserID(c); ok {
		params["userId"] = userID.String()
	}

	result, err := h.service.Suggest(c.Request.Context(), "BY_GROUP", params)
	if err != nil {
		slog.Error("suggestion.HandleGroupSuggestion", "error", err)
		response.Err(c, err)
		return
	}

	response.OK(c, toSuggestionResponse(result))
}

// HandleGetHistory handles GET /suggestions/history.
// Requires authentication.
func (h *SuggestionHandler) HandleGetHistory(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.ErrMsg(c, http.StatusUnauthorized, "authentication required")
		return
	}

	sessionType := c.Query("sessionType")
	page, pageSize := response.ParsePagination(c)

	sessions, total, err := h.service.GetHistory(c.Request.Context(), userID, sessionType, page, pageSize)
	if err != nil {
		slog.Error("suggestion.HandleGetHistory", "error", err)
		response.Err(c, err)
		return
	}

	response.List(c, sessions, page, pageSize, total)
}

// toSuggestionResponse converts a SuggestionResult to the API response DTO.
func toSuggestionResponse(result *SuggestionResult) SuggestionResponse {
	return SuggestionResponse{
		SessionID:     result.Session.ID,
		Type:          result.Session.SessionType,
		Recipes:       result.Recipes,
		TotalCalories: result.Session.TotalCalories,
	}
}
