package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/lehoanglong/whatdish/internal/shared/middleware"
	"github.com/lehoanglong/whatdish/internal/shared/response"
)

// Handler exposes HTTP endpoints for the user bounded context.
type Handler struct {
	authSvc    *AuthService
	profileSvc *ProfileService
	validate   *validator.Validate
}

// NewHandler creates a new user Handler.
func NewHandler(authSvc *AuthService, profileSvc *ProfileService) *Handler {
	return &Handler{
		authSvc:    authSvc,
		profileSvc: profileSvc,
		validate:   validator.New(),
	}
}

// HandleRegister handles POST /auth/register.
func (h *Handler) HandleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.authSvc.Register(c.Request.Context(), req)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.Created(c, result)
}

// HandleLogin handles POST /auth/login.
func (h *Handler) HandleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.authSvc.Login(c.Request.Context(), req)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}

// HandleRefresh handles POST /auth/refresh.
func (h *Handler) HandleRefresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.authSvc.RefreshToken(c.Request.Context(), req)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, tokens)
}

// HandleGetProfile handles GET /users/me.
func (h *Handler) HandleGetProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.ErrMsg(c, http.StatusUnauthorized, "missing user identity")
		return
	}

	result, err := h.profileSvc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}

// HandleUpdateProfile handles PUT /users/me/profile.
func (h *Handler) HandleUpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.ErrMsg(c, http.StatusUnauthorized, "missing user identity")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ErrMsg(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.profileSvc.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}
