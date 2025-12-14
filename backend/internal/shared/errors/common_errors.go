package errors

// Common/shared error codes (used across multiple domains or in handlers)
const (
	CodeInvalidRequest   = "INVALID_REQUEST"
	CodeInvalidParameter = "INVALID_PARAMETER"
	CodeValidationError  = "VALIDATION_ERROR"
	CodeUnauthorized     = "UNAUTHORIZED"
	CodeForbidden        = "FORBIDDEN"
	CodeInternalError    = "INTERNAL_ERROR"
	CodeNotFound         = "NOT_FOUND"
)

// Common errors
var (
	ErrInvalidRequest = NewDomainError(
		CodeInvalidRequest,
		"Dữ liệu yêu cầu không hợp lệ",
		StatusBadRequest,
	)

	ErrInvalidParameter = NewDomainError(
		CodeInvalidParameter,
		"Tham số không hợp lệ",
		StatusBadRequest,
	)

	ErrValidationError = NewDomainError(
		CodeValidationError,
		"Lỗi xác thực dữ liệu",
		StatusBadRequest,
	)

	ErrUnauthorized = NewDomainError(
		CodeUnauthorized,
		"Người dùng chưa được xác thực",
		StatusUnauthorized,
	)

	ErrForbidden = NewDomainError(
		CodeForbidden,
		"Bạn không có quyền thực hiện hành động này",
		StatusForbidden,
	)

	ErrInternalError = NewDomainError(
		CodeInternalError,
		"Đã xảy ra lỗi hệ thống",
		StatusInternalServerError,
	)

	ErrNotFound = NewDomainError(
		CodeNotFound,
		"Không tìm thấy tài nguyên",
		StatusNotFound,
	)
)
