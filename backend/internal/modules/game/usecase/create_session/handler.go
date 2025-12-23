package create_session

import (
	"context"
	"strings"
	"time"

	"github.com/english-coach/backend/internal/modules/game/domain"
	"github.com/english-coach/backend/internal/shared/constants"
	"github.com/english-coach/backend/internal/shared/errors"
	"github.com/english-coach/backend/internal/shared/logger"
)

// Handler handles game session creation
type Handler struct {
	sessionRepo       domain.GameSessionRepository
	questionRepo      domain.GameQuestionRepository
	questionGenerator domain.QuestionGenerator
	logger            logger.ILogger
}

// NewHandler creates a new use case
func NewHandler(
	sessionRepo domain.GameSessionRepository,
	questionRepo domain.GameQuestionRepository,
	questionGenerator domain.QuestionGenerator,
	logger logger.ILogger,
) *Handler {
	return &Handler{
		sessionRepo:       sessionRepo,
		questionRepo:      questionRepo,
		questionGenerator: questionGenerator,
		logger:            logger,
	}
}

// Execute creates a new game session
func (h *Handler) Execute(ctx context.Context, input Input, userID int64) (*Output, error) {
	// Validate request
	if err := input.Validate(); err != nil {
		return nil, errors.ErrValidationError.WithDetails(err.Error())
	}

	// Create game session model
	// Note: TopicID is kept for backward compatibility with DB schema, but we use TopicIDs array for filtering
	var topicID *int64
	if len(input.TopicIDs) > 0 {
		// Store first topic ID for DB compatibility (schema still has single topic_id)
		topicID = &input.TopicIDs[0]
	}
	levelID := &input.LevelID

	session := &domain.GameSession{
		UserID:           userID,
		Mode:             input.Mode,
		SourceLanguageID: input.SourceLanguageID,
		TargetLanguageID: input.TargetLanguageID,
		TopicID:          topicID,
		LevelID:          levelID,
		TotalQuestions:   0, // Will be set when questions are generated
		CorrectQuestions: 0,
		StartedAt:        time.Now(),
	}

	// Save session to database first (needed for question generation)
	if err := h.sessionRepo.Create(ctx, session); err != nil {
		h.logger.Error("failed to create game session",
			logger.Error(err),
			logger.Int64("user_id", userID),
			logger.String("mode", input.Mode),
		)
		return nil, errors.WrapError(err, "failed to create game session")
	}

	// Generate questions upfront - request up to MaxGameQuestionCount (20)
	questions, options, err := h.questionGenerator.GenerateQuestions(
		ctx,
		session.ID,
		input.SourceLanguageID,
		input.TargetLanguageID,
		input.Mode,
		input.TopicIDs,
		input.LevelID,
		constants.MaxGameQuestionCount,
	)
	if err != nil {
		h.logger.Error("failed to generate questions",
			logger.Error(err),
			logger.Int64("session_id", session.ID),
			logger.String("mode", input.Mode),
			logger.Int("source_language_id", int(input.SourceLanguageID)),
			logger.Int("target_language_id", int(input.TargetLanguageID)),
			logger.Any("topic_ids", input.TopicIDs),
			logger.Any("level_id", input.LevelID),
		)
		// Check for insufficient words error (FR-026)
		// Error message format: "Không đủ từ: cần X, có Y"
		if strings.Contains(err.Error(), "Không đủ từ") {
			return nil, domain.ErrInsufficientWords
		}
		return nil, errors.WrapError(err, "failed to generate questions")
	}

	// Check if we have at least the minimum required questions (1)
	if len(questions) < constants.MinGameQuestionCount {
		return nil, domain.ErrInsufficientWords
	}

	// Save questions and options
	if err := h.questionRepo.CreateBatch(ctx, questions, options); err != nil {
		h.logger.Error("failed to save questions",
			logger.Error(err),
			logger.Int64("session_id", session.ID),
		)
		return nil, errors.WrapError(err, "failed to save questions")
	}

	// Update session with question count
	session.TotalQuestions = int16(len(questions))
	if err := h.sessionRepo.Update(ctx, session); err != nil {
		h.logger.Error("failed to update session with question count",
			logger.Error(err),
			logger.Int64("session_id", session.ID),
		)
		// Non-fatal error, continue
	}

	// Log session creation
	h.logger.Info("game session created with questions",
		logger.Int64("session_id", session.ID),
		logger.Int64("user_id", userID),
		logger.String("mode", input.Mode),
		logger.Int("source_language_id", int(input.SourceLanguageID)),
		logger.Int("target_language_id", int(input.TargetLanguageID)),
		logger.Int("question_count", len(questions)),
	)

	return &Output{
		ID:               session.ID,
		UserID:           session.UserID,
		Mode:             session.Mode,
		SourceLanguageID: session.SourceLanguageID,
		TargetLanguageID: session.TargetLanguageID,
		TopicID:          session.TopicID,
		LevelID:          session.LevelID,
		TotalQuestions:   session.TotalQuestions,
		CorrectQuestions: session.CorrectQuestions,
		StartedAt:        session.StartedAt,
		EndedAt:          session.EndedAt,
	}, nil
}
