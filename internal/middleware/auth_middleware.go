package middleware

import (
	"net/http"
	"renew-guard/pkg/jwt"
	"renew-guard/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	UserIDKey           = "userID"
	UserEmailKey        = "userEmail"
)

// AuthMiddleware validates JWT tokens and adds user information to context
func AuthMiddleware(jwtUtil *jwt.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		// Extract bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := jwtUtil.ValidateToken(tokenString)
		if err != nil {
			if err == jwt.ErrExpiredToken {
				utils.ErrorResponse(c, http.StatusUnauthorized, "Token has expired")
			} else {
				utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token")
			}
			c.Abort()
			return
		}

		// Add user information to context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)

		c.Next()
	}
}

// GetUserID retrieves the user ID from the context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetUserEmail retrieves the user email from the context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(UserEmailKey)
	if !exists {
		return "", false
	}
	return email.(string), true
}
