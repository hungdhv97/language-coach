package main

import (
	"log"

	config "github.com/english-coach/backend/configs"
	"github.com/english-coach/backend/internal/app/bootstrap"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Wire all dependencies and create application
	app, err := bootstrap.Wire(cfg)
	if err != nil {
		log.Fatalf("Failed to wire application: %v", err)
	}
	defer app.Container.Close()

	// Log application startup
	app.Logger.Info("Starting English Coach Backend API",
		zap.String("env", cfg.App.Env),
		zap.String("name", cfg.App.Name),
	)

	// Run application (handles graceful shutdown)
	app.Run()
}
