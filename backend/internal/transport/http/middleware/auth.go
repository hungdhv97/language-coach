package middleware

import (
	"net/http"
	"strings"

	"github.com/english-coach/backend/internal/shared/auth"
	"github.com/english-coach/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a Gin middleware for JWT authentication
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.NewError(
				"UNAUTHORIZED",
				"Yêu cầu header Authorization",
				nil,
			))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, response.NewError(
				"UNAUTHORIZED",
				"Định dạng header Authorization không hợp lệ",
				nil,
			))
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateToken(parts[1])
		if err != nil {
			statusCode := http.StatusUnauthorized
			code := "UNAUTHORIZED"
			if err == auth.ErrExpiredToken {
				code = "TOKEN_EXPIRED"
			}

			c.JSON(statusCode, response.NewError(
				code,
				err.Error(),
				nil,
			))
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("claims", claims)

		c.Next()
	}
}
