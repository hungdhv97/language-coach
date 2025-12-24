package domain

import "errors"

// User domain errors - sentinel errors using errors.New()
var (
	ErrEmailRequired      = errors.New("EMAIL_REQUIRED")
	ErrEmailExists        = errors.New("EMAIL_EXISTS")
	ErrUsernameExists     = errors.New("USERNAME_EXISTS")
	ErrInvalidPassword    = errors.New("INVALID_PASSWORD")
	ErrInvalidCredentials = errors.New("INVALID_CREDENTIALS")
	ErrUserInactive       = errors.New("USER_INACTIVE")
	ErrProfileNotFound    = errors.New("PROFILE_NOT_FOUND")
	ErrUserNotFound       = errors.New("USER_NOT_FOUND")
)
