package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger wraps zap logger and implements ILogger interface
type Logger struct {
	*zap.Logger
}

// Ensure Logger implements ILogger interface
var _ ILogger = (*Logger)(nil)

// Debug logs a message at debug level
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

// Info logs a message at info level
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

// Warn logs a message at warn level
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Error logs a message at error level
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

// Fatal logs a message at fatal level and exits
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}

// With creates a child logger with the given fields
func (l *Logger) With(fields ...zap.Field) ILogger {
	return &Logger{Logger: l.Logger.With(fields...)}
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
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
