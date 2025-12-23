package logger

import (
	"time"
)

// ILogger defines the logger interface for structured logging
// Fields are represented as map[string]interface{} to be independent of any logging library
type ILogger interface {
	// Debug logs a message at debug level
	Debug(msg string, fields ...map[string]interface{})
	// Info logs a message at info level
	Info(msg string, fields ...map[string]interface{})
	// Warn logs a message at warn level
	Warn(msg string, fields ...map[string]interface{})
	// Error logs a message at error level
	Error(msg string, fields ...map[string]interface{})
	// Fatal logs a message at fatal level and exits
	Fatal(msg string, fields ...map[string]interface{})
	// With creates a child logger with the given fields
	With(fields ...map[string]interface{}) ILogger
	// Sync flushes any buffered log entries
	Sync() error
}

// Field helpers for common logging scenarios
// These return map[string]interface{} to be independent of any logging library

// String creates a string field
func String(key, value string) map[string]interface{} {
	return map[string]interface{}{key: value}
}

// Int creates an int field
func Int(key string, value int) map[string]interface{} {
	return map[string]interface{}{key: value}
}

// Int64 creates an int64 field
func Int64(key string, value int64) map[string]interface{} {
	return map[string]interface{}{key: value}
}

// Float64 creates a float64 field
func Float64(key string, value float64) map[string]interface{} {
	return map[string]interface{}{key: value}
}

// Bool creates a bool field
func Bool(key string, value bool) map[string]interface{} {
	return map[string]interface{}{key: value}
}

// Error creates an error field
func Error(err error) map[string]interface{} {
	if err == nil {
		return map[string]interface{}{"error": nil}
	}
	return map[string]interface{}{"error": err.Error()}
}

// Any creates a field with any value
func Any(key string, value interface{}) map[string]interface{} {
	return map[string]interface{}{key: value}
}

// Duration creates a duration field
func Duration(key string, value time.Duration) map[string]interface{} {
	return map[string]interface{}{key: value}
}

// Strings creates a string slice field
func Strings(key string, values []string) map[string]interface{} {
	return map[string]interface{}{key: values}
}

// Ints creates an int slice field
func Ints(key string, values []int) map[string]interface{} {
	return map[string]interface{}{key: values}
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
func WithFields(l ILogger, fields ...map[string]interface{}) ILogger {
	return l.With(fields...)
}
