package domain

import "errors"

// Game domain errors - sentinel errors using errors.New()
var (
	ErrInsufficientWords      = errors.New("Insufficient words available")
	ErrSessionNotFound        = errors.New("Session not found")
	ErrSessionEnded           = errors.New("Session has ended")
	ErrQuestionNotFound       = errors.New("Question not found")
	ErrQuestionNotInSession   = errors.New("Question does not belong to this session")
	ErrOptionNotFound         = errors.New("Option not found")
	ErrAnswerAlreadySubmitted = errors.New("Answer has already been submitted")
	ErrInvalidMode            = errors.New("Invalid mode")
	ErrSessionNotOwned        = errors.New("Session is not owned by this user")
)
