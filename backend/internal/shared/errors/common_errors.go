package errors

// Common application errors - pre-defined AppError instances with Vietnamese messages
// All error messages that are shown to end users are in Vietnamese
// These can be used directly or with WithDetails() to add additional context

var (
	// Common errors
	ErrInvalidRequest   = NewAppError(CodeInvalidRequest, "Yêu cầu không hợp lệ")
	ErrInvalidParameter = NewAppError(CodeInvalidParameter, "Tham số không hợp lệ")
	ErrValidationError  = NewAppError(CodeValidationError, "Dữ liệu không hợp lệ")
	ErrUnauthorized     = NewAppError(CodeUnauthorized, "Chưa xác thực")
	ErrForbidden        = NewAppError(CodeForbidden, "Không có quyền truy cập")
	ErrNotFound         = NewAppError(CodeNotFound, "Không tìm thấy")
	ErrConflict         = NewAppError(CodeConflict, "Xung đột dữ liệu")
	ErrInternalError    = NewAppError(CodeInternalError, "Đã xảy ra lỗi hệ thống")

	// User domain errors
	ErrEmailRequired      = NewAppError(CodeEmailRequired, "Email hoặc tên đăng nhập là bắt buộc")
	ErrEmailExists        = NewAppError(CodeEmailExists, "Email đã tồn tại")
	ErrUsernameExists     = NewAppError(CodeUsernameExists, "Tên đăng nhập đã tồn tại")
	ErrInvalidPassword    = NewAppError(CodeInvalidPassword, "Mật khẩu phải có ít nhất 6 ký tự")
	ErrInvalidCredentials = NewAppError(CodeInvalidCredentials, "Email hoặc tên đăng nhập không hợp lệ")
	ErrUserInactive       = NewAppError(CodeUserInactive, "Tài khoản người dùng đã bị vô hiệu hóa")
	ErrProfileNotFound    = NewAppError(CodeProfileNotFound, "Không tìm thấy hồ sơ người dùng")
	ErrUserNotFound       = NewAppError(CodeUserNotFound, "Không tìm thấy người dùng")

	// Game domain errors
	ErrInsufficientWords      = NewAppError(CodeInsufficientWords, "Không đủ từ vựng để tạo phiên chơi. Vui lòng chọn chủ đề hoặc cấp độ khác")
	ErrSessionNotFound        = NewAppError(CodeSessionNotFound, "Không tìm thấy phiên chơi")
	ErrSessionEnded           = NewAppError(CodeSessionEnded, "Phiên chơi đã kết thúc")
	ErrQuestionNotFound       = NewAppError(CodeQuestionNotFound, "Không tìm thấy câu hỏi")
	ErrQuestionNotInSession   = NewAppError(CodeQuestionNotInSession, "Câu hỏi không thuộc về phiên chơi này")
	ErrOptionNotFound         = NewAppError(CodeOptionNotFound, "Không tìm thấy lựa chọn đã chọn")
	ErrAnswerAlreadySubmitted = NewAppError(CodeAnswerAlreadySubmitted, "Đã gửi câu trả lời cho câu hỏi này")
	ErrInvalidMode            = NewAppError(CodeInvalidMode, "Chế độ không hợp lệ")
	ErrSessionNotOwned        = NewAppError(CodeSessionNotOwned, "Phiên chơi không thuộc về người dùng này")

	// Dictionary domain errors
	ErrWordNotFound = NewAppError(CodeWordNotFound, "Không tìm thấy từ")
)
