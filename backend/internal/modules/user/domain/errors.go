package domain

import (
	"github.com/english-coach/backend/internal/shared/errors"
)

// User domain error codes
const (
	CodeEmailRequired      = "EMAIL_REQUIRED"
	CodeEmailExists        = "EMAIL_EXISTS"
	CodeUsernameExists     = "USERNAME_EXISTS"
	CodeInvalidPassword    = "INVALID_PASSWORD"
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeUserInactive       = "USER_INACTIVE"
	CodeProfileNotFound    = "PROFILE_NOT_FOUND"
	CodeUserNotFound       = "USER_NOT_FOUND"
)

// User domain errors
var (
	ErrEmailRequired = errors.NewDomainError(
		CodeEmailRequired,
		"Email hoặc tên đăng nhập là bắt buộc",
		errors.StatusBadRequest,
	)

	ErrEmailExists = errors.NewDomainError(
		CodeEmailExists,
		"Email đã tồn tại",
		errors.StatusConflict,
	)

	ErrUsernameExists = errors.NewDomainError(
		CodeUsernameExists,
		"Tên đăng nhập đã tồn tại",
		errors.StatusConflict,
	)

	ErrInvalidPassword = errors.NewDomainError(
		CodeInvalidPassword,
		"Mật khẩu phải có ít nhất 6 ký tự",
		errors.StatusBadRequest,
	)

	ErrInvalidCredentials = errors.NewDomainError(
		CodeInvalidCredentials,
		"Email hoặc tên đăng nhập không hợp lệ",
		errors.StatusUnauthorized,
	)

	ErrUserInactive = errors.NewDomainError(
		CodeUserInactive,
		"Tài khoản người dùng đã bị vô hiệu hóa",
		errors.StatusForbidden,
	)

	ErrProfileNotFound = errors.NewDomainError(
		CodeProfileNotFound,
		"Không tìm thấy hồ sơ người dùng",
		errors.StatusNotFound,
	)

	ErrUserNotFound = errors.NewDomainError(
		CodeUserNotFound,
		"Không tìm thấy người dùng",
		errors.StatusNotFound,
	)
)
