package middleware

import (
	"net/http"
	"strings"

	sharederrors "github.com/english-coach/backend/internal/shared/errors"
	"github.com/english-coach/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler is a centralized error handling middleware for Gin
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Try to extract domain error
			domainErr, isDomainErr := sharederrors.IsDomainError(err)
			if isDomainErr {
				// Log domain error (message is already in Vietnamese for user)
				logger.Info("domain error",
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("error_code", domainErr.Code),
					zap.String("error_message", domainErr.Message),
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
			logger.Error("request error",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("error", errorMsg),
				zap.Error(err),
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
