package error

import (
	"github.com/english-coach/backend/internal/shared/errors"
)

// Dictionary domain error codes
const (
	CodeWordNotFound = "WORD_NOT_FOUND"
)

// Dictionary domain errors
var (
	ErrWordNotFound = errors.NewDomainError(
		CodeWordNotFound,
		"Không tìm thấy từ",
		errors.StatusNotFound,
	)
)
