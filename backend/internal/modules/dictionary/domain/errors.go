package domain

import "errors"

// Dictionary domain errors - sentinel errors using errors.New()
var (
	ErrWordNotFound = errors.New("Word not found")
)
