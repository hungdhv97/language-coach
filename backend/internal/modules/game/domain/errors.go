package domain

import "errors"

// Game domain errors - sentinel errors using errors.New()
var (
	ErrInsufficientWords      = errors.New("INSUFFICIENT_WORDS")
	ErrSessionNotFound        = errors.New("SESSION_NOT_FOUND")
	ErrSessionEnded           = errors.New("SESSION_ENDED")
	ErrQuestionNotFound       = errors.New("QUESTION_NOT_FOUND")
	ErrQuestionNotInSession   = errors.New("QUESTION_NOT_IN_SESSION")
	ErrOptionNotFound         = errors.New("OPTION_NOT_FOUND")
	ErrAnswerAlreadySubmitted = errors.New("ANSWER_ALREADY_SUBMITTED")
	ErrInvalidMode            = errors.New("INVALID_MODE")
	ErrSessionNotOwned        = errors.New("SESSION_NOT_OWNED")
)
