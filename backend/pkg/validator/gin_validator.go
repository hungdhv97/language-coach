package validator

import (
	"net/http"

	"github.com/english-coach/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validators here if needed
	// validate.RegisterValidation("custom_tag", customValidatorFunc)
}

// RegisterGinValidator registers go-playground/validator with Gin
func RegisterGinValidator() error {
	if _, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Use the same validator instance
		_ = validate
		return nil
	}
	return nil
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}

// ValidateStruct validates a struct and returns errors in Gin format
func ValidateStruct(c *gin.Context, s interface{}) bool {
	if err := validate.Struct(s); err != nil {
		// Handle validation errors
		var errs []string

		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				errs = append(errs, getValidationErrorMessage(e))
			}
		} else {
			errs = []string{err.Error()}
		}

		response.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", errs)
		return false
	}
	return true
}

// ShouldBindJSON validates and binds JSON body
func ShouldBindJSON(c *gin.Context, s interface{}) bool {
	if err := c.ShouldBindJSON(s); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON format", err.Error())
		return false
	}
	return ValidateStruct(c, s)
}

// getValidationErrorMessage returns a human-readable error message
func getValidationErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "email":
		return e.Field() + " must be a valid email"
	case "min":
		return e.Field() + " must be at least " + e.Param() + " characters"
	case "max":
		return e.Field() + " must be at most " + e.Param() + " characters"
	default:
		return e.Field() + " is invalid"
	}
}
