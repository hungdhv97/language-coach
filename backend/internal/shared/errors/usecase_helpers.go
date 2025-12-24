package errors

import (
	dictionarydomain "github.com/english-coach/backend/internal/modules/dictionary/domain"
	gamedomain "github.com/english-coach/backend/internal/modules/game/domain"
	userdomain "github.com/english-coach/backend/internal/modules/user/domain"
)

// MapDomainErrorToAppError maps a domain error to an AppError
// This is used by usecase layer to convert domain errors to standardized AppErrors
func MapDomainErrorToAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	// Check if it's already an AppError
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	// Map user domain errors
	if err := mapUserDomainErrorToAppError(err); err != nil {
		return err
	}

	// Map game domain errors
	if err := mapGameDomainErrorToAppError(err); err != nil {
		return err
	}

	// Map dictionary domain errors
	if err := mapDictionaryDomainErrorToAppError(err); err != nil {
		return err
	}

	// For unexpected errors, return internal error
	// This should rarely happen if error flow is correct
	return ErrInternalError.WithCause(err)
}

// mapUserDomainErrorToAppError maps user domain errors to AppError
func mapUserDomainErrorToAppError(err error) *AppError {
	switch err {
	case userdomain.ErrEmailRequired:
		return ErrEmailRequired
	case userdomain.ErrEmailExists:
		return ErrEmailExists
	case userdomain.ErrUsernameExists:
		return ErrUsernameExists
	case userdomain.ErrInvalidPassword:
		return ErrInvalidPassword
	case userdomain.ErrInvalidCredentials:
		return ErrInvalidCredentials
	case userdomain.ErrUserInactive:
		return ErrUserInactive
	case userdomain.ErrProfileNotFound:
		return ErrProfileNotFound
	case userdomain.ErrUserNotFound:
		return ErrUserNotFound
	default:
		return nil
	}
}

// mapGameDomainErrorToAppError maps game domain errors to AppError
func mapGameDomainErrorToAppError(err error) *AppError {
	switch err {
	case gamedomain.ErrInsufficientWords:
		return ErrInsufficientWords
	case gamedomain.ErrSessionNotFound:
		return ErrSessionNotFound
	case gamedomain.ErrSessionEnded:
		return ErrSessionEnded
	case gamedomain.ErrQuestionNotFound:
		return ErrQuestionNotFound
	case gamedomain.ErrQuestionNotInSession:
		return ErrQuestionNotInSession
	case gamedomain.ErrOptionNotFound:
		return ErrOptionNotFound
	case gamedomain.ErrAnswerAlreadySubmitted:
		return ErrAnswerAlreadySubmitted
	case gamedomain.ErrInvalidMode:
		return ErrInvalidMode
	case gamedomain.ErrSessionNotOwned:
		return ErrSessionNotOwned
	default:
		return nil
	}
}

// mapDictionaryDomainErrorToAppError maps dictionary domain errors to AppError
func mapDictionaryDomainErrorToAppError(err error) *AppError {
	switch err {
	case dictionarydomain.ErrWordNotFound:
		return ErrWordNotFound
	default:
		return nil
	}
}
