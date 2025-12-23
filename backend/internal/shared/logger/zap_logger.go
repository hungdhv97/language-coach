package logger

import (
	"os"
	"time"

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

// convertFields converts map[string]interface{} fields to zap.Field slice
func convertFields(fields ...map[string]interface{}) []zap.Field {
	if len(fields) == 0 {
		return nil
	}

	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		for key, value := range field {
			zapFields = append(zapFields, convertValue(key, value))
		}
	}
	return zapFields
}

// convertValue converts a single key-value pair to zap.Field
func convertValue(key string, value interface{}) zap.Field {
	switch v := value.(type) {
	case string:
		return zap.String(key, v)
	case int:
		return zap.Int(key, v)
	case int64:
		return zap.Int64(key, v)
	case float64:
		return zap.Float64(key, v)
	case bool:
		return zap.Bool(key, v)
	case time.Duration:
		return zap.Duration(key, v)
	case []string:
		return zap.Strings(key, v)
	case []int:
		return zap.Ints(key, v)
	case error:
		return zap.Error(v)
	case nil:
		return zap.Any(key, nil)
	default:
		return zap.Any(key, v)
	}
}

// Debug logs a message at debug level
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	l.Logger.Debug(msg, convertFields(fields...)...)
}

// Info logs a message at info level
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	l.Logger.Info(msg, convertFields(fields...)...)
}

// Warn logs a message at warn level
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	l.Logger.Warn(msg, convertFields(fields...)...)
}

// Error logs a message at error level
func (l *Logger) Error(msg string, fields ...map[string]interface{}) {
	l.Logger.Error(msg, convertFields(fields...)...)
}

// Fatal logs a message at fatal level and exits
func (l *Logger) Fatal(msg string, fields ...map[string]interface{}) {
	l.Logger.Fatal(msg, convertFields(fields...)...)
}

// With creates a child logger with the given fields
func (l *Logger) With(fields ...map[string]interface{}) ILogger {
	return &Logger{Logger: l.Logger.With(convertFields(fields...)...)}
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
