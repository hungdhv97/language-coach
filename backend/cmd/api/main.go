package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/english-coach/backend/internal/app"
	"github.com/english-coach/backend/internal/config"
	"github.com/english-coach/backend/internal/domain/dictionary/service"
	"github.com/english-coach/backend/internal/infrastructure/db"
	"github.com/english-coach/backend/internal/infrastructure/logger"
	"github.com/english-coach/backend/internal/infrastructure/repository/dictionary"
	httpServer "github.com/english-coach/backend/internal/interface/http"
	"github.com/english-coach/backend/internal/interface/http/handler"
	"github.com/english-coach/backend/internal/interface/http/middleware"
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
	appLogger, err := logger.NewLogger(cfg.Logging.Level)
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
	server := httpServer.NewServer(
		httpServer.Config{
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

	// Initialize repositories
	dictRepo := dictionary.NewDictionaryRepository(pool)

	// Initialize services
	dictService := service.NewDictionaryService(
		dictRepo.WordRepository(),
		dictRepo.SenseRepository(),
		pool,
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

	// Register routes
	router := server.Router()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		// Dictionary routes: /api/v1/dictionary/...
		dictionaryGroup := apiV1.Group("/dictionary")
		app.RegisterDictionaryRoutes(dictionaryGroup, dictHandler)
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
