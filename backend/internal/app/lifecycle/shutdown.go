package lifecycle

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// ShutdownConfig holds shutdown configuration
type ShutdownConfig struct {
	Timeout time.Duration
}

// GracefulShutdown handles graceful shutdown of the application
func GracefulShutdown(
	ctx context.Context,
	logger *zap.Logger,
	shutdownFunc func(context.Context) error,
	cfg ShutdownConfig,
) {
	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := shutdownFunc(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	} else {
		logger.Info("Server exited gracefully")
	}
}

