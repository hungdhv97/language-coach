package app

import (
	"github.com/english-coach/backend/internal/interface/http/handler"
	"github.com/gin-gonic/gin"
)

// RegisterDictionaryRoutes registers dictionary-related routes
func RegisterDictionaryRoutes(router *gin.RouterGroup, dictHandler *handler.DictionaryHandler) {
	// Reference data endpoints
	ref := router.Group("/reference")
	{
		ref.GET("/languages", dictHandler.GetLanguages)
		ref.GET("/topics", dictHandler.GetTopics)
		ref.GET("/levels", dictHandler.GetLevels)
	}

	// Dictionary search and word detail endpoints
	router.GET("/search", dictHandler.SearchWords)
	router.GET("/words/:wordId", dictHandler.GetWordDetail)
}

// RegisterGameRoutes registers game-related routes
func RegisterGameRoutes(router *gin.RouterGroup, gameHandler *handler.GameHandler) {
	// Game session endpoints
	sessions := router.Group("/sessions")
	{
		sessions.POST("", gameHandler.CreateSession)
		sessions.GET("/:sessionId", gameHandler.GetSession)
		sessions.POST("/:sessionId/answers", gameHandler.SubmitAnswer)
	}
}

// RegisterStatisticsRoutes registers statistics-related routes
func RegisterStatisticsRoutes(router *gin.RouterGroup, statisticsHandler *handler.StatisticsHandler) {
	// Statistics endpoints
	statistics := router.Group("/statistics")
	{
		statistics.GET("/sessions/:sessionId", statisticsHandler.GetSessionStatistics)
	}
}
