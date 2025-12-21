package port

import (
	"context"

	"github.com/english-coach/backend/internal/modules/game/domain"
)

// GameSessionRepository defines operations for game session data access
type GameSessionRepository interface {
	// Create creates a new game session
	Create(ctx context.Context, session *domain.GameSession) error
	// FindByID returns a game session by ID
	FindByID(ctx context.Context, id int64) (*domain.GameSession, error)
	// Update updates a game session
	Update(ctx context.Context, session *domain.GameSession) error
	// EndSession marks a session as ended
	EndSession(ctx context.Context, sessionID int64, endedAt interface{}) error
}

// GameQuestionRepository defines operations for game question data access
type GameQuestionRepository interface {
	// CreateBatch creates multiple questions and their options in a transaction
	CreateBatch(ctx context.Context, questions []*domain.GameQuestion, options []*domain.GameQuestionOption) error
	// FindBySessionID returns all questions for a session
	FindBySessionID(ctx context.Context, sessionID int64) ([]*domain.GameQuestion, []*domain.GameQuestionOption, error)
	// FindByID returns a question by ID with its options
	FindByID(ctx context.Context, questionID int64) (*domain.GameQuestion, []*domain.GameQuestionOption, error)
}

// GameAnswerRepository defines operations for game answer data access
type GameAnswerRepository interface {
	// Create creates a new answer
	Create(ctx context.Context, answer *domain.GameAnswer) error
	// FindByQuestionID returns the answer for a specific question
	FindByQuestionID(ctx context.Context, questionID, sessionID, userID int64) (*domain.GameAnswer, error)
	// FindBySessionID returns all answers for a session
	FindBySessionID(ctx context.Context, sessionID, userID int64) ([]*domain.GameAnswer, error)
}

