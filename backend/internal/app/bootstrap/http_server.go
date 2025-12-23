package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/english-coach/backend/internal/shared/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// HTTPServerConfig holds HTTP server configuration
type HTTPServerConfig struct {
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// HTTPServer represents the HTTP server
type HTTPServer struct {
	router *gin.Engine
	server *http.Server
}

// NewHTTPServer creates a new HTTP server using Gin
func NewHTTPServer(
	cfg HTTPServerConfig,
	appLogger logger.ILogger,
	corsMiddleware, errorMiddleware, loggerMiddleware gin.HandlerFunc,
) *HTTPServer {
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
		if appLogger != nil {
			appLogger.Error("panic recovered",
				logger.Any("error", recovered),
				logger.String("path", c.Request.URL.Path),
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

	return &HTTPServer{
		router: router,
		server: server,
	}
}

// Router returns the Gin router
func (s *HTTPServer) Router() *gin.Engine {
	return s.router
}

// Start starts the HTTP server
func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
