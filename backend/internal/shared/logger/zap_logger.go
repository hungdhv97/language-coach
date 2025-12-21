package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger wraps zap logger
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new structured logger using zap
func NewLogger(env string) (*Logger, error) {
	// Map environment to log level
	var level zapcore.Level
	switch env {
	case "development":
		level = zapcore.DebugLevel
	case "staging", "production":
		level = zapcore.InfoLevel
	default:
		level = zapcore.InfoLevel
	}

	// Base encoder config
	encoderCfg := zap.NewProductionEncoderConfig()
	if env == "development" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}

	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	// Console encoder (with color)
	consoleEncoderCfg := encoderCfg
	consoleEncoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderCfg)

	// File encoder (no color)
	fileEncoderCfg := encoderCfg
	fileEncoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderCfg)

	// Ensure log directory exists
	if err := os.MkdirAll("logs", 0o755); err != nil {
		return nil, err
	}

	// Lumberjack rotation config
	appWriter := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    50,   // MB per file
		MaxBackups: 7,    // number of old files to keep
		MaxAge:     30,   // days to keep a log file
		Compress:   true, // gzip old log files
	}
	errorWriter := &lumberjack.Logger{
		Filename:   "logs/error.log",
		MaxSize:    50,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	}

	appSyncer := zapcore.AddSync(appWriter)
	errorSyncer := zapcore.AddSync(errorWriter)
	consoleSyncer := zapcore.AddSync(os.Stdout)

	core := zapcore.NewTee(
		// All logs -> app.log (rotated daily)
		zapcore.NewCore(fileEncoder, appSyncer, zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= level
		})),
		// Error and above -> error.log (rotated daily)
		zapcore.NewCore(fileEncoder, errorSyncer, zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.ErrorLevel
		})),
		// Console output
		zapcore.NewCore(consoleEncoder, consoleSyncer, zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= level
		})),
	)

	logger := zap.New(core, zap.AddCaller())
	return &Logger{Logger: logger}, nil
}
