package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoggerMiddleware creates a Gin middleware for request logging with zap
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
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
		requestLogger := logger.With(
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("remote_addr", c.ClientIP()),
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
			zap.Int("status_code", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.Int64("duration_ms", duration.Milliseconds()),
		)
	}
}
