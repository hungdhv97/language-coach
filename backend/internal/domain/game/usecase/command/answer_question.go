package command

import (
	"context"
	"time"

	"github.com/english-coach/backend/internal/domain/game/dto"
	"github.com/english-coach/backend/internal/modules/game/domain"
	"github.com/english-coach/backend/internal/shared/errors"
	"go.uber.org/zap"
)

// SubmitAnswerUseCase handles answer submission
type SubmitAnswerUseCase struct {
	answerRepo   domain.GameAnswerRepository
	questionRepo domain.GameQuestionRepository
	sessionRepo  domain.GameSessionRepository
	logger       *zap.Logger
}

// NewSubmitAnswerUseCase creates a new use case
func NewSubmitAnswerUseCase(
	answerRepo domain.GameAnswerRepository,
	questionRepo domain.GameQuestionRepository,
	sessionRepo domain.GameSessionRepository,
	logger *zap.Logger,
) *SubmitAnswerUseCase {
	return &SubmitAnswerUseCase{
		answerRepo:   answerRepo,
		questionRepo: questionRepo,
		sessionRepo:  sessionRepo,
		logger:       logger,
	}
}

// Execute submits an answer to a question
func (uc *SubmitAnswerUseCase) Execute(ctx context.Context, req *dto.SubmitAnswerRequest, sessionID, userID int64) (*domain.GameAnswer, error) {
	// Get question and options to verify the answer
	question, options, err := uc.questionRepo.FindByID(ctx, req.QuestionID)
	if err != nil {
		uc.logger.Error("failed to find question",
			zap.Error(err),
			zap.Int64("question_id", req.QuestionID),
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
		if opt.ID == req.SelectedOptionID {
			selectedOption = opt
			isCorrect = opt.IsCorrect
			break
		}
	}

	if selectedOption == nil {
		return nil, domain.ErrOptionNotFound
	}

	// Check if answer already exists
	existingAnswer, _ := uc.answerRepo.FindByQuestionID(ctx, req.QuestionID, sessionID, userID)
	if existingAnswer != nil {
		return nil, domain.ErrAnswerAlreadySubmitted
	}

	// Create answer
	answer := &domain.GameAnswer{
		QuestionID:       req.QuestionID,
		SessionID:        sessionID,
		UserID:           userID,
		SelectedOptionID: &req.SelectedOptionID,
		IsCorrect:        isCorrect,
		ResponseTimeMs:   req.ResponseTimeMs,
		AnsweredAt:       time.Now(),
	}

	if err := uc.answerRepo.Create(ctx, answer); err != nil {
		uc.logger.Error("failed to create answer",
			zap.Error(err),
			zap.Int64("question_id", req.QuestionID),
		)
		return nil, errors.WrapError(err, "failed to create answer")
	}

	// Update session correct count if answer is correct
	if isCorrect {
		session, err := uc.sessionRepo.FindByID(ctx, sessionID)
		if err == nil {
			session.CorrectQuestions++
			if err := uc.sessionRepo.Update(ctx, session); err != nil {
				uc.logger.Warn("failed to update session correct count",
					zap.Error(err),
					zap.Int64("session_id", sessionID),
				)
			}
		}
	}

	// Log answer submission
	uc.logger.Info("answer submitted",
		zap.Int64("answer_id", answer.ID),
		zap.Int64("question_id", req.QuestionID),
		zap.Int64("session_id", sessionID),
		zap.Int64("user_id", userID),
		zap.Bool("is_correct", isCorrect),
		zap.Intp("response_time_ms", req.ResponseTimeMs),
	)

	return answer, nil
}
