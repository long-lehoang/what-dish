package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "user_id"

// AuthProvider verifies tokens and returns the user ID.
type AuthProvider interface {
	VerifyToken(ctx context.Context, token string) (uuid.UUID, error)
}

// RequireAuth middleware requires a valid JWT.
func RequireAuth(auth AuthProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "missing authorization token",
			})
			return
		}

		userID, err := auth.VerifyToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "invalid or expired token",
			})
			return
		}

		c.Set(string(userIDKey), userID)
		c.Next()
	}
}

// OptionalAuth middleware extracts user ID if present but doesn't require it.
func OptionalAuth(auth AuthProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token != "" {
			userID, err := auth.VerifyToken(c.Request.Context(), token)
			if err == nil {
				c.Set(string(userIDKey), userID)
			}
		}
		c.Next()
	}
}

// RequireAdmin middleware requires the authenticated user to be in the admin list.
// Must be used after RequireAuth.
func RequireAdmin(adminIDs []uuid.UUID) gin.HandlerFunc {
	adminSet := make(map[uuid.UUID]struct{}, len(adminIDs))
	for _, id := range adminIDs {
		adminSet[id] = struct{}{}
	}

	return func(c *gin.Context) {
		userID, ok := GetUserID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "admin access required",
			})
			return
		}

		if _, isAdmin := adminSet[userID]; !isAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "admin access required",
			})
			return
		}

		c.Next()
	}
}

// GetUserID extracts the user ID from the gin context.
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(string(userIDKey))
	if !exists {
		return uuid.UUID{}, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

func extractToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}
