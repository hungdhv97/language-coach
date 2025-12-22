package create_session

import (
	"errors"
)

// Input represents the input to create a game session
type Input struct {
	SourceLanguageID int16   `json:"source_language_id" validate:"required,gt=0"`
	TargetLanguageID int16   `json:"target_language_id" validate:"required,gt=0"`
	Mode             string  `json:"mode" validate:"required,oneof=level"`               // Always 'level' now
	LevelID          int64   `json:"level_id" validate:"required,gt=0"`                  // Required
	TopicIDs         []int64 `json:"topic_ids,omitempty" validate:"omitempty,dive,gt=0"` // Optional array (empty/null means all topics)
}

// Validate validates the Input
func (r *Input) Validate() error {
	// Source and target languages must be different
	if r.SourceLanguageID == r.TargetLanguageID {
		return errors.New("Ngôn ngữ nguồn và ngôn ngữ đích phải khác nhau")
	}

	// Mode must be 'level'
	if r.Mode != "level" {
		return errors.New("Chế độ phải là 'level'")
	}

	// Level ID is required
	if r.LevelID <= 0 {
		return errors.New("Level_id là bắt buộc và phải lớn hơn 0")
	}

	// TopicIDs is optional (empty array or nil means all topics)
	// If provided, all topic IDs must be valid
	for _, topicID := range r.TopicIDs {
		if topicID <= 0 {
			return errors.New("Tất cả topic_ids phải lớn hơn 0")
		}
	}

	return nil
}

