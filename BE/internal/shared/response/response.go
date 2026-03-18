package response

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	apperrors "github.com/lehoanglong/whatdish/internal/shared/errors"
)

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

type ListResponse struct {
	Data       any        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type DataResponse struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, DataResponse{Data: data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, DataResponse{Data: data})
}

func List(c *gin.Context, data any, page, pageSize int, total int64) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	c.JSON(http.StatusOK, ListResponse{
		Data: data,
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func Err(c *gin.Context, err error) {
	status := apperrors.HTTPStatus(err)
	message := err.Error()
	if status >= 500 {
		slog.Error("internal error", "error", err, "path", c.Request.URL.Path)
		message = "an internal error occurred"
	}
	c.JSON(status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}

func ErrMsg(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}

func ParsePagination(c *gin.Context) (page, pageSize int) {
	page = 1
	pageSize = 20

	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}
	if ps := c.Query("limit"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}
	return page, pageSize
}

func Offset(page, pageSize int) int {
	return (page - 1) * pageSize
}
