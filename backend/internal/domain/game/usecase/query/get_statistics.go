package query

import (
	"context"
	"fmt"
	"time"

	"github.com/english-coach/backend/internal/domain/game/dto"
	"github.com/english-coach/backend/internal/domain/game/port"
	"go.uber.org/zap"
)

// GetSessionStatisticsUseCase handles retrieving session statistics
type GetSessionStatisticsUseCase struct {
	sessionRepo port.GameSessionRepository
	answerRepo  port.GameAnswerRepository
	logger      *zap.Logger
}

// NewGetSessionStatisticsUseCase creates a new use case
func NewGetSessionStatisticsUseCase(
	sessionRepo port.GameSessionRepository,
	answerRepo port.GameAnswerRepository,
	logger *zap.Logger,
) *GetSessionStatisticsUseCase {
	return &GetSessionStatisticsUseCase{
		sessionRepo: sessionRepo,
		answerRepo:  answerRepo,
		logger:      logger,
	}
}

// Execute retrieves statistics for a game session
func (uc *GetSessionStatisticsUseCase) Execute(ctx context.Context, sessionID, userID int64) (*dto.SessionStatistics, error) {
	// Get session
	session, err := uc.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		uc.logger.Error("failed to find session",
			zap.Error(err),
			zap.Int64("session_id", sessionID),
		)
		return nil, fmt.Errorf("failed to find session: %w", err)
	}

	// Check if session is nil
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	// Verify session belongs to user
	if session.UserID != userID {
		return nil, fmt.Errorf("session does not belong to user")
	}

	// Get all answers for the session
	answers, err := uc.answerRepo.FindBySessionID(ctx, sessionID, userID)
	if err != nil {
		uc.logger.Error("failed to find answers",
			zap.Error(err),
			zap.Int64("session_id", sessionID),
		)
		return nil, fmt.Errorf("failed to find answers: %w", err)
	}

	// Calculate statistics
	stats := &dto.SessionStatistics{
		SessionID:      sessionID,
		TotalQuestions: int(session.TotalQuestions),
		CorrectAnswers: int(session.CorrectQuestions),
	}

	// Calculate wrong answers
	stats.WrongAnswers = stats.TotalQuestions - stats.CorrectAnswers

	// Calculate accuracy percentage
	if stats.TotalQuestions > 0 {
		stats.AccuracyPercentage = float64(stats.CorrectAnswers) / float64(stats.TotalQuestions) * 100.0
	} else {
		stats.AccuracyPercentage = 0.0
	}

	// Calculate duration
	if session.EndedAt != nil {
		duration := session.EndedAt.Sub(session.StartedAt)
		stats.DurationSeconds = int(duration.Seconds())
	} else {
		// If session hasn't ended, calculate from current time
		duration := time.Since(session.StartedAt)
		stats.DurationSeconds = int(duration.Seconds())
	}

	// Calculate average response time
	totalResponseTime := 0
	responseTimeCount := 0
	for _, answer := range answers {
		if answer.ResponseTimeMs != nil {
			totalResponseTime += *answer.ResponseTimeMs
			responseTimeCount++
		}
	}

	if responseTimeCount > 0 {
		stats.AverageResponseTimeMs = totalResponseTime / responseTimeCount
	} else {
		stats.AverageResponseTimeMs = 0
	}

	// Log statistics retrieval
	uc.logger.Info("session statistics retrieved",
		zap.Int64("session_id", sessionID),
		zap.Int64("user_id", userID),
		zap.Int("total_questions", stats.TotalQuestions),
		zap.Int("correct_answers", stats.CorrectAnswers),
		zap.Float64("accuracy", stats.AccuracyPercentage),
	)

	return stats, nil
}
