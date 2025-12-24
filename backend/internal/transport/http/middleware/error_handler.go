package middleware

import (
	"net/http"

	sharederrors "github.com/english-coach/backend/internal/shared/errors"
	"github.com/english-coach/backend/internal/shared/logger"
	"github.com/english-coach/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// ErrorHandler is a centralized error handling middleware for Gin
// This is the ONLY place where errors are logged and converted to HTTP responses
func ErrorHandler(appLogger logger.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Try to extract AppError (from usecase layer)
			appErr, isAppErr := sharederrors.IsAppError(err)
			if isAppErr {
				// Log AppError with full context (for internal debugging)
				appLogger.Error("application error",
					logger.String("method", c.Request.Method),
					logger.String("path", c.Request.URL.Path),
					logger.String("error_code", appErr.Code),
					logger.String("error_message", appErr.Message),
					logger.Any("metadata", appErr.Metadata),
					logger.Error(appErr.Cause), // Include cause for debugging
				)

				// Map AppError to HTTP response
				statusCode, httpErr := sharederrors.MapToHTTPResponse(appErr)
				response.ErrorResponse(
					c,
					statusCode,
					httpErr.Code,
					httpErr.Message,
					httpErr.Metadata,
				)
				return
			}

			// For unexpected errors (not AppError), log and return generic internal error
			// This should rarely happen if error flow is correct
			appLogger.Error("unexpected error type",
				logger.String("method", c.Request.Method),
				logger.String("path", c.Request.URL.Path),
				logger.Error(err),
			)

			// Return generic internal error (don't leak error details)
			response.ErrorResponse(
				c,
				http.StatusInternalServerError,
				sharederrors.CodeInternalError,
				"Đã xảy ra lỗi hệ thống",
				nil,
			)
		}
	}
}

// SetError sets an error in the Gin context
// This is used by adapters to pass errors to the error handler middleware
func SetError(c *gin.Context, err error) {
	c.Error(err)
}
