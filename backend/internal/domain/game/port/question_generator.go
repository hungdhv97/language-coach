package port

import (
	"context"

	"github.com/english-coach/backend/internal/modules/game/domain"
)

// QuestionGenerator defines operations for generating game questions
type QuestionGenerator interface {
	// GenerateQuestions generates questions for a game session
	// Returns questions with their options (4 options per question, 1 correct)
	// mode is always "level" now
	// topicIDs is optional array - if nil or empty, includes all topics
	GenerateQuestions(
		ctx context.Context,
		sessionID int64,
		sourceLanguageID, targetLanguageID int16,
		mode string, // Always "level" now
		topicIDs []int64, // Optional array of topic IDs (nil/empty means all topics)
		levelID int64, // Required level ID
		questionCount int,
	) ([]*domain.GameQuestion, []*domain.GameQuestionOption, error)
}
