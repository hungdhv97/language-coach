package errors

import (
	"net/http"
)

// HTTPErrorResponse represents a standardized HTTP error response structure
type HTTPErrorResponse struct {
	Code     string                 `json:"code"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// MapToHTTPResponse maps an AppError to HTTP status code and standardized error response
// This function is called by the error handler middleware
func MapToHTTPResponse(err *AppError) (int, *HTTPErrorResponse) {
	if err == nil {
		return http.StatusInternalServerError, &HTTPErrorResponse{
			Code:    ErrInternalError.Code,
			Message: ErrInternalError.Message,
		}
	}

	// Map error code to HTTP status code
	statusCode := mapCodeToHTTPStatus(err.Code)

	return statusCode, &HTTPErrorResponse{
		Code:     err.Code,
		Message:  err.Message,
		Metadata: err.Metadata,
	}
}

// mapCodeToHTTPStatus maps error codes to HTTP status codes
// This is the only place where HTTP status codes are determined
func mapCodeToHTTPStatus(code string) int {
	switch code {
	// 400 Bad Request
	case CodeInvalidRequest, CodeInvalidParameter, CodeValidationError,
		CodeEmailRequired, CodeInvalidPassword, CodeInvalidMode,
		CodeInsufficientWords, CodeSessionEnded, CodeQuestionNotInSession,
		CodeAnswerAlreadySubmitted:
		return http.StatusBadRequest

	// 401 Unauthorized
	case CodeUnauthorized, CodeInvalidCredentials:
		return http.StatusUnauthorized

	// 403 Forbidden
	case CodeForbidden, CodeUserInactive, CodeSessionNotOwned:
		return http.StatusForbidden

	// 404 Not Found
	case CodeNotFound, CodeUserNotFound, CodeProfileNotFound,
		CodeSessionNotFound, CodeQuestionNotFound, CodeOptionNotFound,
		CodeWordNotFound:
		return http.StatusNotFound

	// 409 Conflict
	case CodeConflict, CodeEmailExists, CodeUsernameExists:
		return http.StatusConflict

	// 500 Internal Server Error (default)
	default:
		return http.StatusInternalServerError
	}
}
