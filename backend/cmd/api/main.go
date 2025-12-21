package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/english-coach/backend/configs"
	"github.com/english-coach/backend/internal/domain/dictionary/service"
	gamesvc "github.com/english-coach/backend/internal/domain/game/service"
	gamecmd "github.com/english-coach/backend/internal/domain/game/usecase/command"
	userregister "github.com/english-coach/backend/internal/modules/user/usecase/register"
	userlogin "github.com/english-coach/backend/internal/modules/user/usecase/login"
	usergetprofile "github.com/english-coach/backend/internal/modules/user/usecase/get_profile"
	userupdateprofile "github.com/english-coach/backend/internal/modules/user/usecase/update_profile"
	"github.com/english-coach/backend/internal/shared/auth"
	"github.com/english-coach/backend/internal/platform/db"
	"github.com/english-coach/backend/internal/shared/logger"
	"github.com/english-coach/backend/internal/modules/dictionary/infra/persistence/postgres/dictionary"
	gamerepo "github.com/english-coach/backend/internal/modules/game/infra/persistence/postgres"
	userrepo "github.com/english-coach/backend/internal/modules/user/infra/persistence/postgres/user"
	httpserver "github.com/english-coach/backend/internal/transport/http"
	"github.com/english-coach/backend/internal/transport/http/handler"
	"github.com/english-coach/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	appLogger, err := logger.NewLogger(cfg.App.Env)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	appLogger.Info("Starting English Coach Backend API",
		zap.String("env", cfg.App.Env),
		zap.String("name", cfg.App.Name),
	)

	// Initialize database connection
	ctx := context.Background()
	pool, err := db.NewPostgres(ctx, db.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		Database:        cfg.Database.Database,
		SSLMode:         cfg.Database.SSLMode,
		MaxConns:        cfg.Database.MaxConns,
		MinConns:        cfg.Database.MinConns,
		MaxConnLifetime: cfg.Database.MaxConnLifetime,
		MaxConnIdleTime: cfg.Database.MaxConnIdleTime,
	})
	if err != nil {
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	appLogger.Info("Database connection established")

	// Log CORS configuration
	appLogger.Info("CORS configuration loaded",
		zap.Strings("allowed_origins", cfg.CORS.AllowedOrigins),
	)

	// Setup HTTP server
	server := httpserver.NewServer(
		httpserver.Config{
			Port:            cfg.Server.Port,
			ReadTimeout:     cfg.Server.ReadTimeout,
			WriteTimeout:    cfg.Server.WriteTimeout,
			IdleTimeout:     cfg.Server.IdleTimeout,
			ShutdownTimeout: cfg.Server.ShutdownTimeout,
		},
		appLogger.Logger,
		middleware.CORS(cfg.CORS.AllowedOrigins),
		middleware.ErrorHandler(appLogger.Logger),
		middleware.LoggerMiddleware(appLogger.Logger),
	)

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)

	// Initialize repositories
	dictRepo := dictionary.NewDictionaryRepository(pool)
	gameRepository := gamerepo.NewGameRepository(pool)
	userRepository := userrepo.NewUserRepository(pool)

	// Initialize services
	dictService := service.NewDictionaryService(
		dictRepo.WordRepository(),
		dictRepo.SenseRepository(),
		dictRepo.LanguageRepository(),
		dictRepo.LevelRepository(),
		dictRepo.PartOfSpeechRepository(),
		pool,
		appLogger.Logger,
	)

	// Initialize game question generator service
	questionGenerator := gamesvc.NewQuestionGeneratorService(
		dictRepo.WordRepository(),
		appLogger.Logger,
	)

	// Initialize game use cases
	createSessionUC := gamecmd.NewCreateGameSessionUseCase(
		gameRepository.GameSessionRepo(),
		gameRepository.GameQuestionRepo(),
		questionGenerator,
		appLogger.Logger,
	)

	submitAnswerUC := gamecmd.NewSubmitAnswerUseCase(
		gameRepository.GameAnswerRepo(),
		gameRepository.GameQuestionRepo(),
		gameRepository.GameSessionRepo(),
		appLogger.Logger,
	)

	// Initialize user use cases
	registerUC := userregister.NewHandler(
		userRepository.UserRepository(),
		appLogger.Logger,
	)

	loginUC := userlogin.NewHandler(
		userRepository.UserRepository(),
		jwtManager,
		appLogger.Logger,
	)

	getProfileUC := usergetprofile.NewHandler(
		userRepository.UserProfileRepository(),
		appLogger.Logger,
	)

	updateProfileUC := userupdateprofile.NewHandler(
		userRepository.UserProfileRepository(),
		appLogger.Logger,
	)

	// Initialize handlers
	dictHandler := handler.NewDictionaryHandler(
		dictRepo.LanguageRepository(),
		dictRepo.TopicRepository(),
		dictRepo.LevelRepository(),
		dictRepo.WordRepository(),
		dictService,
		appLogger.Logger,
	)

	gameHandler := handler.NewGameHandler(
		createSessionUC,
		submitAnswerUC,
		gameRepository.GameQuestionRepo(),
		gameRepository.GameSessionRepo(),
		appLogger.Logger,
	)

	userHandler := handler.NewUserHandler(
		registerUC,
		loginUC,
		getProfileUC,
		updateProfileUC,
		userRepository.UserRepository(),
		userRepository.UserProfileRepository(),
		appLogger.Logger,
	)

	// Register routes
	router := server.Router()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// OpenAPI documentation endpoint (Swagger UI)
	openapiHandler := handler.NewOpenAPIHandler(appLogger.Logger, "docs/openapi/openapi.yaml")
	router.GET("/docs", openapiHandler.GetSwaggerUI)

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		// Auth routes: /api/v1/auth/... (public)
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/register", userHandler.Register)
			authGroup.POST("/login", userHandler.Login)
			authGroup.GET("/check-email", userHandler.CheckEmailAvailability)
			authGroup.GET("/check-username", userHandler.CheckUsernameAvailability)
		}

		// Reference routes: /api/v1/reference/... (public)
		referenceGroup := apiV1.Group("/reference")
		{
			referenceGroup.GET("/languages", dictHandler.GetLanguages)
			referenceGroup.GET("/topics", dictHandler.GetTopics)
			referenceGroup.GET("/levels", dictHandler.GetLevels)
		}

		// Dictionary routes: /api/v1/dictionary/... (public)
		dictionaryGroup := apiV1.Group("/dictionary")
		{
			dictionaryGroup.GET("/search", dictHandler.SearchWords)
			dictionaryGroup.GET("/words/:wordId", dictHandler.GetWordDetail)
		}

		// User routes: /api/v1/users/... (protected)
		userGroup := apiV1.Group("/users")
		userGroup.Use(middleware.AuthMiddleware(jwtManager))
		{
			userGroup.GET("/profile", userHandler.GetProfile)
			userGroup.PUT("/profile", userHandler.UpdateProfile)
		}

		// Game routes: /api/v1/games/... (protected - requires login)
		gameGroup := apiV1.Group("/games")
		gameGroup.Use(middleware.AuthMiddleware(jwtManager))
		{
			sessionsGroup := gameGroup.Group("/sessions")
			{
				sessionsGroup.POST("", gameHandler.CreateSession)
				sessionsGroup.GET("/:sessionId", gameHandler.GetSession)
				sessionsGroup.POST("/:sessionId/answers", gameHandler.SubmitAnswer)
			}
		}
	}

	// Start server in goroutine
	go func() {
		appLogger.Info("Starting HTTP server", zap.Int("port", cfg.Server.Port))
		if err := server.Start(); err != nil {
			appLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("Server forced to shutdown", zap.Error(err))
	}

	appLogger.Info("Server exited")
}
