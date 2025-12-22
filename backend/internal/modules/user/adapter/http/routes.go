package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers user-related HTTP routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, authMiddleware gin.HandlerFunc) {
	// Auth routes: /api/v1/auth/... (public)
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", handler.Register)
		authGroup.POST("/login", handler.Login)
		authGroup.GET("/check-email", handler.CheckEmailAvailability)
		authGroup.GET("/check-username", handler.CheckUsernameAvailability)
	}

	// User routes: /api/v1/users/... (protected)
	userGroup := router.Group("/users")
	userGroup.Use(authMiddleware)
	{
		userGroup.GET("/profile", handler.GetProfile)
		userGroup.PUT("/profile", handler.UpdateProfile)
	}
}

