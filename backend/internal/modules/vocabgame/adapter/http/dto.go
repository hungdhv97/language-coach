package http

import (
	"time"
)

// CreateSessionRequest represents the request body for creating a vocabgame session
type CreateSessionRequest struct {
	Mode             string  `json:"mode" binding:"required"`
	SourceLanguageID int16   `json:"source_language_id" binding:"required"`
	TargetLanguageID int16   `json:"target_language_id" binding:"required"`
	LevelID          int64   `json:"level_id" binding:"required"`
	TopicIDs         []int64 `json:"topic_ids,omitempty"`
}

// CreateSessionResponse represents the response body for creating a vocabgame session
type CreateSessionResponse struct {
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
	Mode             string    `json:"mode"`
	SourceLanguageID int16     `json:"source_language_id"`
	TargetLanguageID int16     `json:"target_language_id"`
	TopicID          *int64    `json:"topic_id,omitempty"`
	LevelID          *int64    `json:"level_id,omitempty"`
	TotalQuestions   int16     `json:"total_questions"`
	CorrectQuestions int16     `json:"correct_questions"`
	StartedAt        time.Time `json:"started_at"`
	EndedAt          *time.Time `json:"ended_at,omitempty"`
}

// SubmitAnswerRequest represents the request body for submitting an answer
type SubmitAnswerRequest struct {
	QuestionID       int64 `json:"question_id" binding:"required"`
	SelectedOptionID int64 `json:"selected_option_id" binding:"required"`
	ResponseTimeMs   *int  `json:"response_time_ms,omitempty"`
}

// SubmitAnswerResponse represents the response body for submitting an answer
type SubmitAnswerResponse struct {
	ID               int64     `json:"id"`
	QuestionID       int64     `json:"question_id"`
	SessionID        int64     `json:"session_id"`
	UserID           int64     `json:"user_id"`
	SelectedOptionID *int64    `json:"selected_option_id,omitempty"`
	IsCorrect        bool      `json:"is_correct"`
	ResponseTimeMs   *int      `json:"response_time_ms,omitempty"`
	AnsweredAt       time.Time `json:"answered_at"`
}

// GetSessionRequest represents the path parameter for getting a session
type GetSessionRequest struct {
	SessionID int64 `uri:"sessionId" binding:"required"`
}

// GameSessionResponse represents a vocabgame session for HTTP response
type GameSessionResponse struct {
	ID               int64      `json:"id"`
	UserID           int64      `json:"user_id"`
	Mode             string     `json:"mode"`
	SourceLanguageID int16      `json:"source_language_id"`
	TargetLanguageID int16      `json:"target_language_id"`
	TopicID          *int64     `json:"topic_id,omitempty"`
	LevelID          *int64     `json:"level_id,omitempty"`
	TotalQuestions   int16      `json:"total_questions"`
	CorrectQuestions int16      `json:"correct_questions"`
	StartedAt        time.Time  `json:"started_at"`
	EndedAt          *time.Time `json:"ended_at,omitempty"`
}

// GameQuestionResponse represents a vocabgame question for HTTP response
type GameQuestionResponse struct {
	ID                  int64     `json:"id"`
	SessionID           int64     `json:"session_id"`
	QuestionOrder       int16     `json:"question_order"`
	QuestionType        string    `json:"question_type"`
	SourceWordID        int64     `json:"source_word_id"`
	SourceSenseID       *int64    `json:"source_sense_id,omitempty"`
	CorrectTargetWordID int64     `json:"correct_target_word_id"`
	SourceLanguageID    int16     `json:"source_language_id"`
	TargetLanguageID    int16     `json:"target_language_id"`
	CreatedAt           time.Time `json:"created_at"`
}

// OptionResponse represents an option for a question (without is_correct for security)
type OptionResponse struct {
	ID           int64  `json:"id"`
	QuestionID   int64  `json:"question_id"`
	OptionLabel  string `json:"option_label"`
	TargetWordID int64  `json:"target_word_id"`
	WordText     string `json:"word_text"`
}

// QuestionWithOptions represents a question with its options for the response
type QuestionWithOptions struct {
	GameQuestionResponse
	SourceWordText string          `json:"source_word_text"`
	Options        []OptionResponse `json:"options"`
}

// GetSessionResponse represents the response for getting a session
type GetSessionResponse struct {
	Session   GameSessionResponse   `json:"session"`
	Questions []QuestionWithOptions `json:"questions"`
}

// ListSessionsResponse represents the response for listing sessions
type ListSessionsResponse struct {
	Sessions []GameSessionResponse `json:"sessions"`
}