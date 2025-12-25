package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers vocabgame-related HTTP routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, authMiddleware gin.HandlerFunc) {
	// VocabGame routes: /api/v1/vocabgames/... (protected - requires login)
	vocabGameGroup := router.Group("/vocabgames")
	vocabGameGroup.Use(authMiddleware)
	{
		sessionsGroup := vocabGameGroup.Group("/sessions")
		{
			sessionsGroup.POST("", handler.CreateSession)
			sessionsGroup.GET("", handler.ListSessions) // Must be before /:sessionId to avoid route conflict
			sessionsGroup.GET("/:sessionId", handler.GetSession)
			sessionsGroup.POST("/:sessionId/answers", handler.SubmitAnswer)
		}
	}
}
