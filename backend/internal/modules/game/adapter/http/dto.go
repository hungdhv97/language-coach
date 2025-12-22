package http

import (
	"github.com/english-coach/backend/internal/modules/game/domain"
)

// CreateSessionRequest represents the request body for creating a game session
type CreateSessionRequest struct {
	Mode             string  `json:"mode" binding:"required"`
	SourceLanguageID int16   `json:"source_language_id" binding:"required"`
	TargetLanguageID int16   `json:"target_language_id" binding:"required"`
	LevelID          int64   `json:"level_id" binding:"required"`
	TopicIDs         []int64 `json:"topic_ids,omitempty"`
}

// SubmitAnswerRequest represents the request body for submitting an answer
type SubmitAnswerRequest struct {
	QuestionID       int64 `json:"question_id" binding:"required"`
	SelectedOptionID int64 `json:"selected_option_id" binding:"required"`
	ResponseTimeMs   *int  `json:"response_time_ms,omitempty"`
}

// GetSessionRequest represents the path parameter for getting a session
type GetSessionRequest struct {
	SessionID int64 `uri:"sessionId" binding:"required"`
}

// QuestionWithOptions represents a question with its options for the response
type QuestionWithOptions struct {
	*domain.GameQuestion
	Options []*domain.GameQuestionOption `json:"options"`
}

// GetSessionResponse represents the response for getting a session
type GetSessionResponse struct {
	Session   *domain.GameSession `json:"session"`
	Questions []QuestionWithOptions `json:"questions"`
}

