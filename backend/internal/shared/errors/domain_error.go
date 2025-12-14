package errors

import (
	"fmt"
	"net/http"
)

// DomainError represents a domain-specific error with code, message, and details
type DomainError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	StatusCode int         `json:"-"` // HTTP status code, not exposed in JSON
}

// Error implements the error interface
func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new domain error
func NewDomainError(code, message string, statusCode int) *DomainError {
	return &DomainError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// WithDetails adds details to the error
func (e *DomainError) WithDetails(details interface{}) *DomainError {
	e.Details = details
	return e
}

// IsDomainError checks if an error is a DomainError
func IsDomainError(err error) (*DomainError, bool) {
	if err == nil {
		return nil, false
	}

	// Check if error is directly a DomainError
	if de, ok := err.(*DomainError); ok {
		return de, true
	}

	return nil, false
}

// WrapError wraps a standard error with domain error information
// This is used for fmt.Errorf cases where we want to preserve the original error
// but still log it properly
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Common HTTP status codes
const (
	StatusBadRequest          = http.StatusBadRequest          // 400
	StatusUnauthorized        = http.StatusUnauthorized        // 401
	StatusForbidden           = http.StatusForbidden           // 403
	StatusNotFound            = http.StatusNotFound            // 404
	StatusConflict            = http.StatusConflict            // 409
	StatusInternalServerError = http.StatusInternalServerError // 500
)
