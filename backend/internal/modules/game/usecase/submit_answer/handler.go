package submit_answer

import (
	"context"
	"time"

	"github.com/english-coach/backend/internal/modules/game/domain"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
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
func (h *Handler) Execute(ctx context.Context, input SubmitAnswerInput, sessionID, userID int64) (*SubmitAnswerOutput, error) {
	// Get question and options to verify the answer
	questionWithOptions, err := h.questionRepo.FindGameQuestionByID(ctx, input.QuestionID)
	if err != nil {
		h.logger.Error("failed to find question",
			logger.Error(err),
			logger.Int64("question_id", input.QuestionID),
		)
		return nil, sharederrors.MapDomainErrorToAppError(err)
	}

	// Check if question is nil
	if questionWithOptions == nil || questionWithOptions.Question == nil {
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrQuestionNotFound)
	}

	question := questionWithOptions.Question
	options := questionWithOptions.Options

	// Verify question belongs to session
	if question.SessionID != sessionID {
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrQuestionNotInSession)
	}

	// Check if session exists and has not ended
	session, err := h.sessionRepo.FindGameSessionByID(ctx, sessionID)
	if err != nil {
		h.logger.Error("failed to find session",
			logger.Error(err),
			logger.Int64("session_id", sessionID),
		)
		return nil, sharederrors.MapDomainErrorToAppError(err)
	}
	if session == nil {
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrSessionNotFound)
	}
	if session.EndedAt != nil {
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrSessionEnded)
	}

	// Verify user owns session
	if session.UserID != userID {
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrSessionNotOwned)
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
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrOptionNotFound)
	}

	// Check if answer already exists
	existingAnswer, err := h.answerRepo.FindGameAnswerByQuestionID(ctx, input.QuestionID, sessionID, userID)
	if err != nil {
		// If not found, it's a normal case (answer doesn't exist yet, allow submission)
		if sharederrors.IsNotFound(err) {
			existingAnswer = nil
		} else {
			// Real error occurred
			h.logger.Error("failed to check existing answer",
				logger.Error(err),
				logger.Int64("question_id", input.QuestionID),
				logger.Int64("session_id", sessionID),
				logger.Int64("user_id", userID),
			)
			return nil, sharederrors.MapDomainErrorToAppError(err)
		}
	}
	if existingAnswer != nil {
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrAnswerAlreadySubmitted)
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
		return nil, sharederrors.MapDomainErrorToAppError(err)
	}

	// Update session correct count if answer is correct
	if isCorrect {
		session.CorrectQuestions++
		if err := h.sessionRepo.Update(ctx, session); err != nil {
			h.logger.Error("failed to update session correct count",
				logger.Error(err),
				logger.Int64("session_id", sessionID),
			)
			return nil, sharederrors.MapDomainErrorToAppError(err)
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

	return &SubmitAnswerOutput{
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
