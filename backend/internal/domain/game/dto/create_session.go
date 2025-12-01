package dto

import (
	"errors"
	"fmt"
)

// CreateGameSessionRequest represents the request to create a game session
type CreateGameSessionRequest struct {
	SourceLanguageID int16  `json:"source_language_id" validate:"required,gt=0"`
	TargetLanguageID int16  `json:"target_language_id" validate:"required,gt=0"`
	Mode             string `json:"mode" validate:"required,oneof=topic level"`
	TopicID          *int64 `json:"topic_id,omitempty" validate:"omitempty,gt=0"`
	LevelID          *int64 `json:"level_id,omitempty" validate:"omitempty,gt=0"`
}

// Validate validates the CreateGameSessionRequest
func (r *CreateGameSessionRequest) Validate() error {
	// Source and target languages must be different
	if r.SourceLanguageID == r.TargetLanguageID {
		return errors.New("source and target languages must be different")
	}

	// Mode must be either 'topic' or 'level'
	if r.Mode != "topic" && r.Mode != "level" {
		return errors.New("mode must be either 'topic' or 'level'")
	}

	// Topic XOR Level required (exactly one must be set)
	if r.Mode == "topic" {
		if r.TopicID == nil || *r.TopicID <= 0 {
			return errors.New("topic_id is required when mode is 'topic'")
		}
		if r.LevelID != nil {
			return errors.New("cannot specify both topic_id and level_id at the same time")
		}
	} else if r.Mode == "level" {
		if r.LevelID == nil || *r.LevelID <= 0 {
			return errors.New("level_id is required when mode is 'level'")
		}
		if r.TopicID != nil {
			return errors.New("cannot specify both topic_id and level_id at the same time")
		}
	}

	return nil
}

// Error implements error interface for validation errors
func (r *CreateGameSessionRequest) Error() string {
	if err := r.Validate(); err != nil {
		return err.Error()
	}
	return ""
}

// GetValidationError returns a formatted validation error message
func (r *CreateGameSessionRequest) GetValidationError() error {
	if err := r.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}

