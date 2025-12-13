package dto

// SessionStatistics represents statistics for a game session
type SessionStatistics struct {
	SessionID             int64   `json:"session_id"`
	TotalQuestions        int     `json:"total_questions"`
	CorrectAnswers        int     `json:"correct_answers"`
	WrongAnswers          int     `json:"wrong_answers"`
	AccuracyPercentage    float64 `json:"accuracy_percentage"`      // 0-100
	DurationSeconds       int     `json:"duration_seconds"`         // Total session duration
	AverageResponseTimeMs int     `json:"average_response_time_ms"` // Average response time in milliseconds
}
