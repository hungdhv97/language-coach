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
		// Game Session operations
		case "FindSessionByID", "FindByID":
			return gamedomain.ErrSessionNotFound
		// Game Question operations
		case "FindQuestionByID":
			return gamedomain.ErrQuestionNotFound
		case "FindQuestionsBySessionID":
			// FindQuestionsBySessionID returns empty slice if not found, not an error
			// But if there's a DB error, return as-is
			return err
		// Game Question Option operations
		case "FindOptionByID":
			return gamedomain.ErrOptionNotFound
		// Game Answer operations
		case "FindAnswerByQuestionID":
			// Answer not found is not necessarily an error - might be first time answering
			// Return as-is, let usecase decide
			return err
		case "FindAnswersBySessionID":
			// FindAnswersBySessionID returns empty slice if not found, not an error
			// But if there's a DB error, return as-is
			return err
		// Create/Update operations
		case "Create", "CreateBatch", "Update", "EndSession":
			// These operations should not return "not found" errors
			// If they do, it's likely a constraint violation or other issue
			return err
		default:
			return err // Return as-is, let usecase handle
		}
	}

	// Check for unique violation errors (if any unique constraints exist)
	if IsUniqueViolation(err) {
		// Game domain doesn't have unique constraints that need special handling
		// Return as-is, let usecase handle
		return err
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
		// Word operations - specific operation name
		switch operation {
		case "FindWordByID":
			return dictionarydomain.ErrWordNotFound
		}

		// Operations that return collections (empty slice/map if not found, not an error)
		// These should not return "not found" errors, but if they do, it's a DB error
		switch operation {
		case "FindByIDs", "FindWordsByTopicAndLanguages", "FindWordsByLevelAndLanguages",
			"FindWordsByLevelAndTopicsAndLanguages", "FindTranslationsForWord",
			"SearchWords", "CountSearchWords", "FindAll", "FindByLanguageID",
			"FindByWordID", "FindByWordIDs":
			// These operations return empty results if not found, not an error
			// But if there's a DB error, return as-is
			return err
		}

		// Generic operations (FindByID, FindByCode) used by multiple entities
		// Since we can't distinguish which entity, return error as-is
		// Usecase layer should check IsNotFound and map to appropriate domain error
		// based on context (e.g., if calling topicRepo.FindByID, map to ErrTopicNotFound)
		switch operation {
		case "FindByID", "FindByCode":
			// Return as-is, let usecase layer handle based on context
			return err
		default:
			return err // Return as-is, let usecase handle
		}
	}

	// Check for unique violation errors (if any unique constraints exist)
	if IsUniqueViolation(err) {
		// Dictionary domain doesn't have unique constraints that need special handling
		// Return as-is, let usecase handle
		return err
	}

	// For other errors, return as-is
	return err
}
