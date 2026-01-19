package middleware

import (
	"net/http"
	"strings"

	"github.com/fayzzzm/go-bro/pkg/auth"
	"github.com/gin-gonic/gin"
)

const (
	AuthUserIDKey    = "auth_user_id"
	AuthUserEmailKey = "auth_user_email"
	AuthCookieName   = "auth_token"
)

// AuthMiddleware validates JWT tokens from cookies (primary) or Authorization header (fallback)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Try to get token from cookie first (preferred)
		cookie, err := c.Cookie(AuthCookieName)
		if err == nil && cookie != "" {
			tokenString = cookie
		}

		// 2. Fallback to Authorization header (for API clients)
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		// No token found
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in context for handlers to use
		c.Set(AuthUserIDKey, claims.UserID)
		c.Set(AuthUserEmailKey, claims.Email)

		c.Next()
	}
}

// GetUserID extracts user ID from context (set by AuthMiddleware)
func GetUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get(AuthUserIDKey)
	if !exists {
		return 0, false
	}
	id, ok := userID.(int)
	return id, ok
}

// GetUserEmail extracts user email from context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(AuthUserEmailKey)
	if !exists {
		return "", false
	}
	e, ok := email.(string)
	return e, ok
}
