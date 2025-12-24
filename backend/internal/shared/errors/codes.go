package errors

// Error codes catalog - used across the application
// These codes are stable and should not change frequently
// Client applications can hardcode these codes

// Common error codes
const (
	CodeInvalidRequest   = "INVALID_REQUEST"
	CodeInvalidParameter = "INVALID_PARAMETER"
	CodeValidationError  = "VALIDATION_ERROR"
	CodeUnauthorized     = "UNAUTHORIZED"
	CodeForbidden        = "FORBIDDEN"
	CodeInternalError    = "INTERNAL_ERROR"
	CodeNotFound         = "NOT_FOUND"
	CodeConflict         = "CONFLICT"
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

// Game domain error codes
const (
	CodeInsufficientWords      = "INSUFFICIENT_WORDS"
	CodeSessionNotFound        = "SESSION_NOT_FOUND"
	CodeSessionEnded           = "SESSION_ENDED"
	CodeQuestionNotFound       = "QUESTION_NOT_FOUND"
	CodeQuestionNotInSession   = "QUESTION_NOT_IN_SESSION"
	CodeOptionNotFound         = "OPTION_NOT_FOUND"
	CodeAnswerAlreadySubmitted = "ANSWER_ALREADY_SUBMITTED"
	CodeInvalidMode            = "INVALID_MODE"
	CodeSessionNotOwned        = "SESSION_NOT_OWNED"
)

// Dictionary domain error codes
const (
	CodeWordNotFound = "WORD_NOT_FOUND"
)
