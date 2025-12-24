package errors

import (
	dictionarydomain "github.com/english-coach/backend/internal/modules/dictionary/domain"
	gamedomain "github.com/english-coach/backend/internal/modules/game/domain"
	userdomain "github.com/english-coach/backend/internal/modules/user/domain"
)

// MapToDomainError translates technical errors (pgx, etc.) to domain errors
// This is used by infrastructure layer to translate technical errors to business errors
// Each domain should provide its own mapping function

// MapUserRepositoryError translates technical errors to user domain errors
func MapUserRepositoryError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Check for "not found" errors
	if IsNotFound(err) {
		switch operation {
		case "FindByID", "FindByEmail", "FindByUsername":
			return userdomain.ErrUserNotFound
		case "GetByUserID", "GetProfile":
			return userdomain.ErrProfileNotFound
		default:
			return userdomain.ErrUserNotFound
		}
	}

	// Check for unique violation errors
	if IsUniqueViolation(err) {
		field := GetUniqueConstraintField(err)
		switch field {
		case "users_email_key", "users_email_unique":
			return userdomain.ErrEmailExists
		case "users_username_key", "users_username_unique":
			return userdomain.ErrUsernameExists
		default:
			// Generic conflict - let usecase decide based on context
			return err // Return as-is, usecase will handle
		}
	}

	// For other errors, return as-is (let usecase handle unexpected errors)
	return err
}

// MapGameRepositoryError translates technical errors to game domain errors
func MapGameRepositoryError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Check for "not found" errors
	if IsNotFound(err) {
		switch operation {
		case "FindSessionByID", "FindByID":
			return gamedomain.ErrSessionNotFound
		case "FindQuestionByID":
			return gamedomain.ErrQuestionNotFound
		case "FindOptionByID":
			return gamedomain.ErrOptionNotFound
		default:
			return err // Return as-is, let usecase handle
		}
	}

	// For other errors, return as-is
	return err
}

// MapDictionaryRepositoryError translates technical errors to dictionary domain errors
func MapDictionaryRepositoryError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Check for "not found" errors
	if IsNotFound(err) {
		switch operation {
		case "FindWordByID", "FindByID":
			return dictionarydomain.ErrWordNotFound
		default:
			return err // Return as-is, let usecase handle
		}
	}

	// For other errors, return as-is
	return err
}
