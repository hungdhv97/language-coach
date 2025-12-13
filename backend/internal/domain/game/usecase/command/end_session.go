package command

import (
	"context"
	"fmt"
	"time"

	"github.com/english-coach/backend/internal/domain/game/port"
	"go.uber.org/zap"
)

// EndGameSessionUseCase handles ending a game session
type EndGameSessionUseCase struct {
	sessionRepo port.GameSessionRepository
	logger      *zap.Logger
}

// NewEndGameSessionUseCase creates a new use case
func NewEndGameSessionUseCase(
	sessionRepo port.GameSessionRepository,
	logger *zap.Logger,
) *EndGameSessionUseCase {
	return &EndGameSessionUseCase{
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

// Execute ends a game session
func (uc *EndGameSessionUseCase) Execute(ctx context.Context, sessionID int64) error {
	// Get session to verify it exists
	session, err := uc.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		uc.logger.Error("failed to find session",
			zap.Error(err),
			zap.Int64("session_id", sessionID),
		)
		return fmt.Errorf("failed to find session: %w", err)
	}

	// Check if session is nil
	if session == nil {
		return fmt.Errorf("session not found")
	}

	// Check if already ended
	if session.EndedAt != nil {
		return fmt.Errorf("session already ended")
	}

	// End session
	if err := uc.sessionRepo.EndSession(ctx, sessionID, time.Now()); err != nil {
		uc.logger.Error("failed to end session",
			zap.Error(err),
			zap.Int64("session_id", sessionID),
		)
		return fmt.Errorf("failed to end session: %w", err)
	}

	// Log session end
	uc.logger.Info("game session ended",
		zap.Int64("session_id", sessionID),
		zap.Int16("total_questions", session.TotalQuestions),
		zap.Int16("correct_questions", session.CorrectQuestions),
	)

	return nil
}
