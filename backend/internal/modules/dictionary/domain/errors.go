package domain

import "errors"

// Dictionary domain errors - sentinel errors using errors.New()
var (
	ErrWordNotFound        = errors.New("Word not found")
	ErrTopicNotFound       = errors.New("Topic not found")
	ErrLevelNotFound       = errors.New("Level not found")
	ErrLanguageNotFound    = errors.New("Language not found")
	ErrPartOfSpeechNotFound = errors.New("Part of speech not found")
	ErrSenseNotFound       = errors.New("Sense not found")
)
