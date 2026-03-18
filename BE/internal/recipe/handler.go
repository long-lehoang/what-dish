package recipe

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lehoanglong/whatdish/internal/shared/response"
)

type DishHandler struct {
	service DishServicePort
}

func NewDishHandler(service DishServicePort) *DishHandler {
	return &DishHandler{service: service}
}

// HandleListDishes handles GET /recipes
func (h *DishHandler) HandleListDishes(c *gin.Context) {
	filter := parseDishFilter(c)

	dishes, total, err := h.service.ListDishes(c.Request.Context(), filter)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.List(c, ToDishListResponse(dishes), filter.Page, filter.PageSize, total)
}

// HandleGetDish handles GET /recipes/:id
// Tries UUID parse first; falls back to slug lookup.
func (h *DishHandler) HandleGetDish(c *gin.Context) {
	param := c.Param("id")

	id, err := uuid.Parse(param)
	if err == nil {
		detail, err := h.service.GetDish(c.Request.Context(), id)
		if err != nil {
			response.Err(c, err)
			return
		}
		resp := ToDishDetailResponse(detail)
		response.OK(c, resp)
		return
	}

	// Param is not a valid UUID; treat as slug.
	detail, err := h.service.GetDishBySlug(c.Request.Context(), param)
	if err != nil {
		response.Err(c, err)
		return
	}
	resp := ToDishDetailResponse(detail)
	response.OK(c, resp)
}

// HandleGetRandomDish handles GET /recipes/random
func (h *DishHandler) HandleGetRandomDish(c *gin.Context) {
	filter := parseDishFilter(c)

	var excludeIDs []uuid.UUID
	if v := c.Query("exclude_ids"); v != "" {
		for _, raw := range strings.Split(v, ",") {
			if id, err := uuid.Parse(strings.TrimSpace(raw)); err == nil {
				excludeIDs = append(excludeIDs, id)
			}
		}
	}

	detail, err := h.service.GetRandomDish(c.Request.Context(), filter, excludeIDs)
	if err != nil {
		response.Err(c, err)
		return
	}
	resp := ToDishDetailResponse(detail)
	response.OK(c, resp)
}

// HandleSearchDishes handles GET /recipes/search
func (h *DishHandler) HandleSearchDishes(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.ErrMsg(c, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	page, pageSize := response.ParsePagination(c)

	dishes, total, err := h.service.SearchDishes(c.Request.Context(), q, page, pageSize)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.List(c, ToDishListResponse(dishes), page, pageSize, total)
}

// HandleListCategories handles GET /categories
func (h *DishHandler) HandleListCategories(c *gin.Context) {
	categoryType := c.Query("type")

	cats, err := h.service.ListCategories(c.Request.Context(), categoryType)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, ToCategoryListResponse(cats))
}

// HandleListTags handles GET /tags
func (h *DishHandler) HandleListTags(c *gin.Context) {
	tags, err := h.service.ListTags(c.Request.Context())
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, ToTagListResponse(tags))
}

func parseDishFilter(c *gin.Context) DishFilter {
	page, pageSize := response.ParsePagination(c)
	filter := DishFilter{
		Page:     page,
		PageSize: pageSize,
	}

	if v := c.Query("dish_type"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.DishTypeID = &id
		}
	}
	if v := c.Query("region"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.RegionID = &id
		}
	}
	if v := c.Query("main_ingredient"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.MainIngredientID = &id
		}
	}
	if v := c.Query("meal_type"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.MealTypeID = &id
		}
	}
	if v := c.Query("difficulty"); v != "" {
		filter.Difficulty = &v
	}
	if v := c.Query("max_cook_time"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			filter.MaxCookTime = &n
		}
	}
	if v := c.Query("tags"); v != "" {
		filter.Tags = strings.Split(v, ",")
	}

	return filter
}
