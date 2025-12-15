package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap logger
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new structured logger using zap
func NewLogger(env string) (*Logger, error) {
	var zapLevel zapcore.Level
	var config zap.Config

	// Map environment to log level
	switch env {
	case "development":
		zapLevel = zapcore.DebugLevel
		config = zap.NewDevelopmentConfig()
	case "staging", "production":
		zapLevel = zapcore.InfoLevel
		config = zap.NewProductionConfig()
	default:
		zapLevel = zapcore.InfoLevel
		config = zap.NewDevelopmentConfig()
	}

	config.Level = zap.NewAtomicLevelAt(zapLevel)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: logger}, nil
}
