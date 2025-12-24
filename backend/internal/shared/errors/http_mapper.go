package errors

import (
	"net/http"
)

// HTTPErrorResponse represents a standardized HTTP error response structure
// This is used by MapToHTTPResponse to return error details without importing the response package
// (to avoid import cycles)
type HTTPErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// MapToHTTPResponse maps a DomainError to HTTP status code and standardized error response
// This function ensures consistent error response format across the application
//
// Returns:
//   - HTTP status code (int)
//   - HTTPErrorResponse struct with code, message, and details
//
// The caller should use these values to construct the HTTP response (e.g., via gin.JSON)
// The HTTP status code is taken directly from DomainError.StatusCode which is set when
// the error is created using NewDomainError()
func MapToHTTPResponse(err *DomainError) (int, *HTTPErrorResponse) {
	if err == nil {
		return http.StatusInternalServerError, &HTTPErrorResponse{
			Code:    CodeInternalError,
			Message: "Đã xảy ra lỗi hệ thống",
			Details: nil,
		}
	}

	// Use the status code from DomainError (set when error is created)
	statusCode := err.StatusCode
	if statusCode == 0 {
		// Fallback to 500 if status code is not set
		statusCode = http.StatusInternalServerError
	}

	return statusCode, &HTTPErrorResponse{
		Code:    err.Code,
		Message: err.Message,
		Details: err.Details,
	}
}
