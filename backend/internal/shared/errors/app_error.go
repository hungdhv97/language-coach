package errors

// AppError represents a standardized application error that flows from usecase to adapter
// This is the only error type that should be returned from usecases
type AppError struct {
	// Code is a machine-readable error code (e.g., "USER_NOT_FOUND", "INVALID_REQUEST")
	Code string `json:"code"`

	// Message is a human-readable error message (safe to show to client)
	Message string `json:"message"`

	// Metadata contains additional context (field errors, params, etc.)
	// This is optional and only included when relevant
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Cause is the underlying error (for internal logging/debugging only)
	// This should NEVER be exposed to clients
	Cause error `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Unwrap returns the underlying cause error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithMetadata adds metadata to the error
// Returns a new AppError instance to avoid mutating shared error instances
func (e *AppError) WithMetadata(key string, value interface{}) *AppError {
	newErr := &AppError{
		Code:    e.Code,
		Message: e.Message,
		Cause:   e.Cause,
	}
	if e.Metadata != nil {
		// Copy existing metadata
		newErr.Metadata = make(map[string]interface{})
		for k, v := range e.Metadata {
			newErr.Metadata[k] = v
		}
	} else {
		newErr.Metadata = make(map[string]interface{})
	}
	newErr.Metadata[key] = value
	return newErr
}

// WithCause adds a cause error (for internal debugging)
// Returns a new AppError instance to avoid mutating shared error instances
func (e *AppError) WithCause(cause error) *AppError {
	newErr := &AppError{
		Code:    e.Code,
		Message: e.Message,
		Cause:   cause,
	}
	if e.Metadata != nil {
		// Copy existing metadata
		newErr.Metadata = make(map[string]interface{})
		for k, v := range e.Metadata {
			newErr.Metadata[k] = v
		}
	}
	return newErr
}

// WithDetails adds details to the error metadata
// This is a convenience method that adds a "details" key to the metadata
func (e *AppError) WithDetails(details string) *AppError {
	return e.WithMetadata("details", details)
}

// NewAppError creates a new AppError
func NewAppError(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}

	return nil, false
}
