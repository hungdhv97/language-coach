package domain

import (
	"context"
)

// GameSessionRepository defines operations for vocabgame session data access
type GameSessionRepository interface {
	// Create creates a new vocabgame session
	Create(ctx context.Context, session *GameSession) error
	// FindGameSessionByID returns a vocabgame session by ID
	FindGameSessionByID(ctx context.Context, id int64) (*GameSession, error)
	// FindGameSessionsByUserID returns a list of game sessions for a user with pagination
	FindGameSessionsByUserID(ctx context.Context, userID int64, limit, offset int) ([]*GameSession, error)
	// CountGameSessionsByUserID returns the total count of game sessions for a user
	CountGameSessionsByUserID(ctx context.Context, userID int64) (int64, error)
	// Update updates a vocabgame session
	Update(ctx context.Context, session *GameSession) error
	// EndSession marks a session as ended
	EndSession(ctx context.Context, sessionID int64, endedAt interface{}) error
}

// GameQuestionRepository defines operations for vocabgame question data access
type GameQuestionRepository interface {
	// CreateBatch creates multiple questions and their options in a transaction
	CreateBatch(ctx context.Context, questions []*GameQuestion, options []*GameQuestionOption) error
	// FindGameQuestionsBySessionID returns all questions for a session with their options
	FindGameQuestionsBySessionID(ctx context.Context, sessionID int64) (*GameQuestionsResult, error)
	// FindGameQuestionByID returns a question by ID with its options
	FindGameQuestionByID(ctx context.Context, questionID int64) (*GameQuestionWithOptions, error)
}

// GameAnswerRepository defines operations for vocabgame answer data access
type GameAnswerRepository interface {
	// Create creates a new answer
	Create(ctx context.Context, answer *GameAnswer) error
	// FindGameAnswerByQuestionID returns the answer for a specific question in a session
	FindGameAnswerByQuestionID(ctx context.Context, questionID, sessionID, userID int64) (*GameAnswer, error)
	// FindGameAnswersBySessionID returns all answers for a session
	FindGameAnswersBySessionID(ctx context.Context, sessionID, userID int64) ([]*GameAnswer, error)
}
