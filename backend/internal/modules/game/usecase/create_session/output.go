package create_session

import "time"

// Output represents the output for creating a game session
type Output struct {
	ID                int64     `json:"id"`
	UserID            int64     `json:"user_id"`
	Mode              string    `json:"mode"`
	SourceLanguageID  int16     `json:"source_language_id"`
	TargetLanguageID  int16     `json:"target_language_id"`
	TopicID           *int64    `json:"topic_id,omitempty"`
	LevelID           *int64    `json:"level_id,omitempty"`
	TotalQuestions    int16     `json:"total_questions"`
	CorrectQuestions  int16     `json:"correct_questions"`
	StartedAt         time.Time `json:"started_at"`
	EndedAt           *time.Time `json:"ended_at,omitempty"`
}

