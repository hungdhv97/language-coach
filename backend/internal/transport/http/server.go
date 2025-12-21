package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/requestid"
	"go.uber.org/zap"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	server *http.Server
}

// Config holds server configuration
type Config struct {
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// NewServer creates a new HTTP server using Gin with Sonic JSON binding
func NewServer(cfg Config, logger *zap.Logger, corsMiddleware, errorMiddleware, loggerMiddleware gin.HandlerFunc) *Server {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	// Use Sonic for JSON binding (faster JSON parsing)
	gin.EnableJsonDecoderUseNumber()
	
	router := gin.New()

	// Add request ID middleware
	router.Use(requestid.New())

	// Add logger middleware
	if loggerMiddleware != nil {
		router.Use(loggerMiddleware)
	}

	// Add CORS middleware
	if corsMiddleware != nil {
		router.Use(corsMiddleware)
	}

	// Add recovery middleware
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if logger != nil {
			logger.Error("panic recovered",
				zap.Any("error", recovered),
				zap.String("path", c.Request.URL.Path),
			)
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "An internal error occurred",
		})
	}))

	// Add error handler middleware
	if errorMiddleware != nil {
		router.Use(errorMiddleware)
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		router: router,
		server: server,
	}
}

// Router returns the Gin router
func (s *Server) Router() *gin.Engine {
	return s.router
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
