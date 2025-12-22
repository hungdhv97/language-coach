package response

import (
	"github.com/english-coach/backend/internal/shared/pagination"
	"github.com/gin-gonic/gin"
)

// Success sends a success response
func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
}

// Paginated sends a paginated success response
func Paginated(c *gin.Context, statusCode int, data interface{}, params *pagination.Params, total int64) {
	metadata := pagination.CalculateMetadata(params, total)
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
		"pagination": gin.H{
			"page":       metadata.Page,
			"pageSize":   metadata.PageSize,
			"total":      metadata.Total,
			"totalPages": metadata.TotalPages,
			"limit":      metadata.Limit,
			"offset":     metadata.Offset,
			"hasNext":    metadata.HasNext,
			"hasPrev":    metadata.HasPrev,
		},
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, code, message string, details interface{}) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error": gin.H{
			"code":    code,
			"message": message,
			"details": details,
		},
	})
}
