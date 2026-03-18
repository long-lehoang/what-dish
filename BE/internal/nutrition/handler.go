package nutrition

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/lehoanglong/whatdish/internal/shared/response"
)

// Handler exposes HTTP endpoints for the nutrition bounded context.
type Handler struct {
	svc      *NutritionService
	validate *validator.Validate
}

// NewHandler creates a new nutrition Handler.
func NewHandler(svc *NutritionService) *Handler {
	return &Handler{
		svc:      svc,
		validate: validator.New(),
	}
}

// HandleGetRecipeNutrition handles GET /nutrition/recipes/:id.
func (h *Handler) HandleGetRecipeNutrition(c *gin.Context) {
	recipeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrMsg(c, http.StatusBadRequest, "invalid recipe id")
		return
	}

	result, err := h.svc.GetByRecipeID(c.Request.Context(), recipeID)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}

// HandleGetBatchNutrition handles GET /nutrition/recipes?ids=uuid1,uuid2.
func (h *Handler) HandleGetBatchNutrition(c *gin.Context) {
	idsParam := c.Query("ids")
	if idsParam == "" {
		response.ErrMsg(c, http.StatusBadRequest, "ids query parameter is required")
		return
	}

	parts := strings.Split(idsParam, ",")
	ids := make([]uuid.UUID, 0, len(parts))
	for _, p := range parts {
		id, err := uuid.Parse(strings.TrimSpace(p))
		if err != nil {
			response.ErrMsg(c, http.StatusBadRequest, "invalid uuid in ids: "+p)
			return
		}
		ids = append(ids, id)
	}

	results, err := h.svc.GetByRecipeIDs(c.Request.Context(), ids)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, results)
}

// HandleCalculateTDEE handles POST /nutrition/calculate-tdee.
func (h *Handler) HandleCalculateTDEE(c *gin.Context) {
	var req CalculateTDEERequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.svc.CalculateTDEE(c.Request.Context(), req)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}

// HandleListGoals handles GET /nutrition/goals.
func (h *Handler) HandleListGoals(c *gin.Context) {
	goals, err := h.svc.ListGoals(c.Request.Context())
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, goals)
}

// HandleUpsertNutrition handles POST /nutrition/recipes/:id (admin endpoint).
func (h *Handler) HandleUpsertNutrition(c *gin.Context) {
	recipeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ErrMsg(c, http.StatusBadRequest, "invalid recipe id")
		return
	}

	var req UpsertNutritionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.svc.Upsert(c.Request.Context(), recipeID, req)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.Created(c, result)
}
