package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ILogger defines the logger interface for structured logging
type ILogger interface {
	// Debug logs a message at debug level
	Debug(msg string, fields ...zap.Field)
	// Info logs a message at info level
	Info(msg string, fields ...zap.Field)
	// Warn logs a message at warn level
	Warn(msg string, fields ...zap.Field)
	// Error logs a message at error level
	Error(msg string, fields ...zap.Field)
	// Fatal logs a message at fatal level and exits
	Fatal(msg string, fields ...zap.Field)
	// With creates a child logger with the given fields
	With(fields ...zap.Field) ILogger
	// Sync flushes any buffered log entries
	Sync() error
}

// Field helpers for common logging scenarios

// String creates a zap.String field
func String(key, value string) zap.Field {
	return zap.String(key, value)
}

// Int creates a zap.Int field
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Int64 creates a zap.Int64 field
func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

// Float64 creates a zap.Float64 field
func Float64(key string, value float64) zap.Field {
	return zap.Float64(key, value)
}

// Bool creates a zap.Bool field
func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

// Error creates a zap.Error field
func Error(err error) zap.Field {
	return zap.Error(err)
}

// Any creates a zap.Any field
func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// Duration creates a zap.Duration field
func Duration(key string, value time.Duration) zap.Field {
	return zap.Duration(key, value)
}

// Strings creates a zap.Strings field
func Strings(key string, values []string) zap.Field {
	return zap.Strings(key, values)
}

// Ints creates a zap.Ints field
func Ints(key string, values []int) zap.Field {
	return zap.Ints(key, values)
}

// Object creates a zap.Object field
func Object(key string, obj zapcore.ObjectMarshaler) zap.Field {
	return zap.Object(key, obj)
}

// Namespace creates a zap.Namespace field
func Namespace(key string) zap.Field {
	return zap.Namespace(key)
}

// Contextual field helpers for common use cases

// WithRequestID creates a logger with request ID field
func WithRequestID(l ILogger, requestID string) ILogger {
	return l.With(String("request_id", requestID))
}

// WithUserID creates a logger with user ID field
func WithUserID(l ILogger, userID int64) ILogger {
	return l.With(Int64("user_id", userID))
}

// WithTraceID creates a logger with trace ID field
func WithTraceID(l ILogger, traceID string) ILogger {
	return l.With(String("trace_id", traceID))
}

// WithMethod creates a logger with HTTP method field
func WithMethod(l ILogger, method string) ILogger {
	return l.With(String("method", method))
}

// WithPath creates a logger with HTTP path field
func WithPath(l ILogger, path string) ILogger {
	return l.With(String("path", path))
}

// WithRemoteAddr creates a logger with remote address field
func WithRemoteAddr(l ILogger, remoteAddr string) ILogger {
	return l.With(String("remote_addr", remoteAddr))
}

// WithStatus creates a logger with HTTP status code field
func WithStatus(l ILogger, statusCode int) ILogger {
	return l.With(Int("status_code", statusCode))
}

// WithDuration creates a logger with duration field
func WithDuration(l ILogger, duration time.Duration) ILogger {
	return l.With(Duration("duration", duration))
}

// WithError creates a logger with error field
func WithError(l ILogger, err error) ILogger {
	return l.With(Error(err))
}

// WithFields creates a logger with multiple fields
func WithFields(l ILogger, fields ...zap.Field) ILogger {
	return l.With(fields...)
}
