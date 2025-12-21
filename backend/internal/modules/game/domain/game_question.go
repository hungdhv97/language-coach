package domain

import "time"

// GameQuestion represents a single question within a game session
type GameQuestion struct {
	ID                  int64     `json:"id"`
	SessionID           int64     `json:"session_id"`
	QuestionOrder       int16     `json:"question_order"`
	QuestionType        string    `json:"question_type"`
	SourceWordID        int64     `json:"source_word_id"`
	SourceSenseID       *int64    `json:"source_sense_id,omitempty"`
	CorrectTargetWordID int64    `json:"correct_target_word_id"`
	SourceLanguageID    int16     `json:"source_language_id"`
	TargetLanguageID    int16     `json:"target_language_id"`
	CreatedAt           time.Time `json:"created_at"`
}

// GameQuestionOption represents one of the four multiple-choice answers (A, B, C, D)
type GameQuestionOption struct {
	ID            int64  `json:"id"`
	QuestionID    int64  `json:"question_id"`
	OptionLabel   string `json:"option_label"` // 'A', 'B', 'C', 'D'
	TargetWordID  int64  `json:"target_word_id"`
	IsCorrect     bool   `json:"is_correct"`
}

