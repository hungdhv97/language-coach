package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS creates a Gin middleware for CORS handling
func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowedOriginsMap := make(map[string]bool)
	for _, origin := range allowedOrigins {
		allowedOriginsMap[origin] = true
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		if origin != "" {
			// Support wildcard "*" (allow any origin) by echoing the request origin.
			// This works with credentials, unlike setting "*" directly.
			if allowedOriginsMap["*"] {
				c.Header("Access-Control-Allow-Origin", origin)
			} else if allowedOriginsMap[origin] {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
