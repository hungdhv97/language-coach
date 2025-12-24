package domain

import "errors"

// User domain errors - sentinel errors using errors.New()
var (
	ErrEmailRequired      = errors.New("Email or username is required")
	ErrEmailExists        = errors.New("Email already exists")
	ErrUsernameExists     = errors.New("Username already exists")
	ErrInvalidPassword    = errors.New("Password must be at least 6 characters")
	ErrInvalidCredentials = errors.New("Invalid email or username")
	ErrUserInactive       = errors.New("User account is inactive")
	ErrProfileNotFound    = errors.New("User profile not found")
	ErrUserNotFound       = errors.New("User not found")
)
