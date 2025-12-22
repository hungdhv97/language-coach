package bootstrap

import (
	"context"

	config "github.com/english-coach/backend/configs"
	"github.com/english-coach/backend/internal/app/di"
	"github.com/english-coach/backend/internal/app/lifecycle"
	"github.com/english-coach/backend/internal/transport/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Application represents the fully wired application
type Application struct {
	Container *di.Container
	Server    *http.Server
	Logger    *zap.Logger
}

// Wire wires all dependencies and returns a fully configured application
func Wire(cfg *config.Config) (*Application, error) {
	// Create dependency injection container
	container, err := di.NewContainer(cfg)
	if err != nil {
		return nil, err
	}

	// Register routes
	registerRoutes(container)

	app := &Application{
		Container: container,
		Server:    container.HTTPServer,
		Logger:    container.Logger.Logger,
	}

	return app, nil
}

// Start starts the application
func (app *Application) Start() error {
	app.Logger.Info("Starting HTTP server",
		zap.Int("port", app.Container.Config.Server.Port),
	)
	return app.Server.Start()
}

// Shutdown gracefully shuts down the application
func (app *Application) Shutdown(ctx context.Context) error {
	app.Logger.Info("Shutting down server...")
	return app.Server.Shutdown(ctx)
}

// Run starts the application and handles graceful shutdown
func (app *Application) Run() {
	// Start server in goroutine
	go func() {
		if err := app.Start(); err != nil {
			app.Logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Handle graceful shutdown
	lifecycle.GracefulShutdown(
		context.Background(),
		app.Logger,
		app.Shutdown,
		lifecycle.ShutdownConfig{
			Timeout: app.Container.Config.Server.ShutdownTimeout,
		},
	)

	// Cleanup resources
	if err := app.Container.Close(); err != nil {
		app.Logger.Error("Error closing container", zap.Error(err))
	}
}

// registerRoutes registers all HTTP routes
func registerRoutes(container *di.Container) {
	router := container.HTTPServer.Router()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// OpenAPI documentation endpoint (Swagger UI)
	router.GET("/docs", container.OpenAPIHandler.GetSwaggerUI)

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		// Auth routes: /api/v1/auth/... (public)
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/register", container.UserHandler.Register)
			authGroup.POST("/login", container.UserHandler.Login)
			authGroup.GET("/check-email", container.UserHandler.CheckEmailAvailability)
			authGroup.GET("/check-username", container.UserHandler.CheckUsernameAvailability)
		}

		// Reference routes: /api/v1/reference/... (public)
		referenceGroup := apiV1.Group("/reference")
		{
			referenceGroup.GET("/languages", container.DictionaryHandler.GetLanguages)
			referenceGroup.GET("/topics", container.DictionaryHandler.GetTopics)
			referenceGroup.GET("/levels", container.DictionaryHandler.GetLevels)
		}

		// Dictionary routes: /api/v1/dictionary/... (public)
		dictionaryGroup := apiV1.Group("/dictionary")
		{
			dictionaryGroup.GET("/search", container.DictionaryHandler.SearchWords)
			dictionaryGroup.GET("/words/:wordId", container.DictionaryHandler.GetWordDetail)
		}

		// User routes: /api/v1/users/... (protected)
		userGroup := apiV1.Group("/users")
		userGroup.Use(container.AuthMiddleware)
		{
			userGroup.GET("/profile", container.UserHandler.GetProfile)
			userGroup.PUT("/profile", container.UserHandler.UpdateProfile)
		}

		// Game routes: /api/v1/games/... (protected - requires login)
		gameGroup := apiV1.Group("/games")
		gameGroup.Use(container.AuthMiddleware)
		{
			sessionsGroup := gameGroup.Group("/sessions")
			{
				sessionsGroup.POST("", container.GameHandler.CreateSession)
				sessionsGroup.GET("/:sessionId", container.GameHandler.GetSession)
				sessionsGroup.POST("/:sessionId/answers", container.GameHandler.SubmitAnswer)
			}
		}
	}
}

