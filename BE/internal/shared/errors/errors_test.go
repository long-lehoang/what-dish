package errors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPStatus_DirectErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"ErrValidation", ErrValidation, http.StatusBadRequest},
		{"ErrUnauthorized", ErrUnauthorized, http.StatusUnauthorized},
		{"ErrForbidden", ErrForbidden, http.StatusForbidden},
		{"ErrNotFound", ErrNotFound, http.StatusNotFound},
		{"ErrConflict", ErrConflict, http.StatusConflict},
		{"ErrInternal", ErrInternal, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantStatus, HTTPStatus(tt.err))
		})
	}
}

func TestHTTPStatus_UnknownErrorReturns500(t *testing.T) {
	unknown := fmt.Errorf("something unexpected")
	assert.Equal(t, http.StatusInternalServerError, HTTPStatus(unknown))
}

func TestHTTPStatus_WrappedErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			"WrappedNotFound",
			fmt.Errorf("dish lookup: %w", ErrNotFound),
			http.StatusNotFound,
		},
		{
			"WrappedValidation",
			fmt.Errorf("input check: %w", ErrValidation),
			http.StatusBadRequest,
		},
		{
			"WrappedUnauthorized",
			fmt.Errorf("auth: %w", ErrUnauthorized),
			http.StatusUnauthorized,
		},
		{
			"WrappedForbidden",
			fmt.Errorf("access: %w", ErrForbidden),
			http.StatusForbidden,
		},
		{
			"WrappedConflict",
			fmt.Errorf("duplicate: %w", ErrConflict),
			http.StatusConflict,
		},
		{
			"DoubleWrappedNotFound",
			fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", ErrNotFound)),
			http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantStatus, HTTPStatus(tt.err))
		})
	}
}
