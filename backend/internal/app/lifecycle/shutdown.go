package lifecycle

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/english-coach/backend/internal/shared/logger"
)

// ShutdownConfig holds shutdown configuration
type ShutdownConfig struct {
	Timeout time.Duration
}

// GracefulShutdown handles graceful shutdown of the application
func GracefulShutdown(
	ctx context.Context,
	appLogger logger.ILogger,
	shutdownFunc func(context.Context) error,
	cfg ShutdownConfig,
) {
	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := shutdownFunc(shutdownCtx); err != nil {
		appLogger.Error("Server forced to shutdown", logger.Error(err))
	} else {
		appLogger.Info("Server exited gracefully")
	}
}
