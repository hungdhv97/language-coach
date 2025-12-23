package middleware

import (
	"net/http"
	"strings"

	sharederrors "github.com/english-coach/backend/internal/shared/errors"
	"github.com/english-coach/backend/internal/shared/logger"
	"github.com/english-coach/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// ErrorHandler is a centralized error handling middleware for Gin
func ErrorHandler(appLogger logger.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Try to extract domain error
			domainErr, isDomainErr := sharederrors.IsDomainError(err)
			if isDomainErr {
				// Log domain error (message is already in Vietnamese for user)
				appLogger.Info("domain error",
					logger.String("method", c.Request.Method),
					logger.String("path", c.Request.URL.Path),
					logger.String("error_code", domainErr.Code),
					logger.String("error_message", domainErr.Message),
				)

				// Use domain error details if available
				details := domainErr.Details
				if details == nil {
					details = nil
				}

				c.JSON(domainErr.StatusCode, response.NewError(
					domainErr.Code,
					domainErr.Message,
					details,
				))
				return
			}

			// For non-domain errors (fmt.Errorf, etc.), log in English lowercase
			// and return generic internal error
			errorMsg := strings.ToLower(err.Error())
			appLogger.Error("request error",
				logger.String("method", c.Request.Method),
				logger.String("path", c.Request.URL.Path),
				logger.String("error", errorMsg),
				logger.Error(err),
			)

			// Determine status code
			statusCode := http.StatusInternalServerError
			code := sharederrors.CodeInternalError
			message := sharederrors.ErrInternalError.Message

			// Check if status code was already set
			if c.Writer.Status() != http.StatusOK {
				statusCode = c.Writer.Status()
			}

			c.JSON(statusCode, response.NewError(
				code,
				message,
				nil,
			))
		}
	}
}

// Helper function to set error in Gin context
func SetError(c *gin.Context, err error) {
	c.Error(err)
}
