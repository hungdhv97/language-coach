package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers dictionary-related HTTP routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	// Reference routes: /api/v1/reference/... (public)
	referenceGroup := router.Group("/reference")
	{
		referenceGroup.GET("/languages", handler.GetLanguages)
		referenceGroup.GET("/topics", handler.GetTopics)
		referenceGroup.GET("/levels", handler.GetLevels)
	}

	// Dictionary routes: /api/v1/dictionary/... (public)
	dictionaryGroup := router.Group("/dictionary")
	{
		dictionaryGroup.GET("/search", handler.SearchWords)
		dictionaryGroup.GET("/words/:wordId", handler.GetWordDetail)
	}
}

