package middleware

import (
	"time"

	"github.com/english-coach/backend/internal/shared/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggerMiddleware creates a Gin middleware for request logging
func LoggerMiddleware(appLogger logger.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate or extract request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Request.Header.Set("X-Request-ID", requestID)
		}
		c.Header("X-Request-ID", requestID)

		// Add request ID to logger context
		requestLogger := appLogger.With(
			logger.String("request_id", requestID),
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.URL.Path),
			logger.String("remote_addr", c.ClientIP()),
		)

		// Store logger in context
		c.Set("logger", requestLogger)

		// Log request start
		requestLogger.Info("request started")

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request completion
		requestLogger.Info("request completed",
			logger.Int("status_code", c.Writer.Status()),
			logger.Duration("duration", duration),
			logger.Int64("duration_ms", duration.Milliseconds()),
		)
	}
}
