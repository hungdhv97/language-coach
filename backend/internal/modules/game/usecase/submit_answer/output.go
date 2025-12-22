package submit_answer

import "time"

// Output represents the output for submitting an answer
type Output struct {
	ID               int64     `json:"id"`
	QuestionID       int64     `json:"question_id"`
	SessionID        int64     `json:"session_id"`
	UserID           int64     `json:"user_id"`
	SelectedOptionID *int64    `json:"selected_option_id,omitempty"`
	IsCorrect        bool      `json:"is_correct"`
	ResponseTimeMs   *int      `json:"response_time_ms,omitempty"`
	AnsweredAt       time.Time `json:"answered_at"`
}

