package submit_answer

import (
	"context"
	"time"

	"github.com/english-coach/backend/internal/modules/game/domain"
	"github.com/english-coach/backend/internal/shared/errors"
	"github.com/english-coach/backend/internal/shared/logger"
)

// Handler handles answer submission
type Handler struct {
	answerRepo   domain.GameAnswerRepository
	questionRepo domain.GameQuestionRepository
	sessionRepo  domain.GameSessionRepository
	logger       logger.ILogger
}

// NewHandler creates a new use case
func NewHandler(
	answerRepo domain.GameAnswerRepository,
	questionRepo domain.GameQuestionRepository,
	sessionRepo domain.GameSessionRepository,
	logger logger.ILogger,
) *Handler {
	return &Handler{
		answerRepo:   answerRepo,
		questionRepo: questionRepo,
		sessionRepo:  sessionRepo,
		logger:       logger,
	}
}

// Execute submits an answer to a question
func (h *Handler) Execute(ctx context.Context, input Input, sessionID, userID int64) (*Output, error) {
	// Get question and options to verify the answer
	question, options, err := h.questionRepo.FindByID(ctx, input.QuestionID)
	if err != nil {
		h.logger.Error("failed to find question",
			logger.Error(err),
			logger.Int64("question_id", input.QuestionID),
		)
		return nil, errors.WrapError(err, "failed to find question")
	}

	// Check if question is nil
	if question == nil {
		return nil, domain.ErrQuestionNotFound
	}

	// Verify question belongs to session
	if question.SessionID != sessionID {
		return nil, domain.ErrQuestionNotInSession
	}

	// Find the selected option
	var selectedOption *domain.GameQuestionOption
	var isCorrect bool
	for _, opt := range options {
		if opt.ID == input.SelectedOptionID {
			selectedOption = opt
			isCorrect = opt.IsCorrect
			break
		}
	}

	if selectedOption == nil {
		return nil, domain.ErrOptionNotFound
	}

	// Check if answer already exists
	existingAnswer, _ := h.answerRepo.FindByQuestionID(ctx, input.QuestionID, sessionID, userID)
	if existingAnswer != nil {
		return nil, domain.ErrAnswerAlreadySubmitted
	}

	// Create answer
	answer := &domain.GameAnswer{
		QuestionID:       input.QuestionID,
		SessionID:        sessionID,
		UserID:           userID,
		SelectedOptionID: &input.SelectedOptionID,
		IsCorrect:        isCorrect,
		ResponseTimeMs:   input.ResponseTimeMs,
		AnsweredAt:       time.Now(),
	}

	if err := h.answerRepo.Create(ctx, answer); err != nil {
		h.logger.Error("failed to create answer",
			logger.Error(err),
			logger.Int64("question_id", input.QuestionID),
		)
		return nil, errors.WrapError(err, "failed to create answer")
	}

	// Update session correct count if answer is correct
	if isCorrect {
		session, err := h.sessionRepo.FindByID(ctx, sessionID)
		if err == nil {
			session.CorrectQuestions++
			if err := h.sessionRepo.Update(ctx, session); err != nil {
				h.logger.Warn("failed to update session correct count",
					logger.Error(err),
					logger.Int64("session_id", sessionID),
				)
			}
		}
	}

	// Log answer submission
	fields := []map[string]interface{}{
		logger.Int64("answer_id", answer.ID),
		logger.Int64("question_id", input.QuestionID),
		logger.Int64("session_id", sessionID),
		logger.Int64("user_id", userID),
		logger.Bool("is_correct", isCorrect),
	}
	if input.ResponseTimeMs != nil {
		fields = append(fields, logger.Int("response_time_ms", *input.ResponseTimeMs))
	}
	h.logger.Info("answer submitted", fields...)

	return &Output{
		ID:               answer.ID,
		QuestionID:       answer.QuestionID,
		SessionID:        answer.SessionID,
		UserID:           answer.UserID,
		SelectedOptionID: answer.SelectedOptionID,
		IsCorrect:        answer.IsCorrect,
		ResponseTimeMs:   answer.ResponseTimeMs,
		AnsweredAt:       answer.AnsweredAt,
	}, nil
}
