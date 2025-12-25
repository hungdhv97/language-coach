package game

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/english-coach/backend/internal/modules/game/domain"
	db "github.com/english-coach/backend/internal/platform/db/sqlc/gen/game"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
)

// gameQuestionRepo implements GameQuestionRepository using sqlc
type gameQuestionRepo struct {
	*GameRepository
}

// CreateBatch creates multiple questions and their options in a transaction
func (r *gameQuestionRepo) CreateBatch(ctx context.Context, questions []*domain.GameQuestion, options []*domain.GameQuestionOption) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return sharederrors.MapGameRepositoryError(err, "CreateBatch")
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Insert questions
	for _, question := range questions {
		var sourceSenseID pgtype.Int8
		if question.SourceSenseID != nil {
			sourceSenseID = pgtype.Int8{Int64: *question.SourceSenseID, Valid: true}
		}
		createdAt := pgtype.Timestamp{Time: time.Now(), Valid: true}

		result, err := qtx.CreateGameQuestion(ctx, db.CreateGameQuestionParams{
			SessionID:           question.SessionID,
			QuestionOrder:       question.QuestionOrder,
			QuestionType:        question.QuestionType,
			SourceWordID:        question.SourceWordID,
			SourceSenseID:       sourceSenseID,
			CorrectTargetWordID: question.CorrectTargetWordID,
			SourceLanguageID:    question.SourceLanguageID,
			TargetLanguageID:    question.TargetLanguageID,
			CreatedAt:           createdAt,
		})
		if err != nil {
			return sharederrors.MapGameRepositoryError(err, "CreateBatch")
		}
		question.ID = result.ID
		question.CreatedAt = result.CreatedAt.Time
	}

	// Insert options
	// Each question has exactly 4 options (A, B, C, D), so we map options to questions by index
	const optionsPerQuestion = 4
	expectedOptionsCount := len(questions) * optionsPerQuestion
	if len(options) != expectedOptionsCount {
		return sharederrors.NewAppError(
			sharederrors.CodeInternalError,
			"mismatch between questions and options count",
		).WithMetadata("expected", expectedOptionsCount).
			WithMetadata("questions", len(questions)).
			WithMetadata("options", len(options))
	}

	for i, question := range questions {
		// Calculate the start index for this question's options
		optionStartIndex := i * optionsPerQuestion
		optionEndIndex := optionStartIndex + optionsPerQuestion

		// Insert options for this question
		for j := optionStartIndex; j < optionEndIndex; j++ {
			option := options[j]
			// Update the option's QuestionID to the actual question ID
			option.QuestionID = question.ID

			optionID, err := qtx.CreateGameQuestionOption(ctx, db.CreateGameQuestionOptionParams{
				QuestionID:   option.QuestionID,
				OptionLabel:  option.OptionLabel,
				TargetWordID: option.TargetWordID,
				IsCorrect:    option.IsCorrect,
			})
			if err != nil {
				return sharederrors.MapGameRepositoryError(err, "CreateBatch")
			}
			option.ID = optionID
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return sharederrors.MapGameRepositoryError(err, "CreateBatch")
	}
	return nil
}

// FindGameQuestionsBySessionID returns all questions for a session with their options
func (r *gameQuestionRepo) FindGameQuestionsBySessionID(ctx context.Context, sessionID int64) (*domain.GameQuestionsResult, error) {
	questionRows, err := r.queries.FindGameQuestionsBySessionID(ctx, sessionID)
	if err != nil {
		return nil, sharederrors.MapGameRepositoryError(err, "FindGameQuestionsBySessionID")
	}

	questions := make([]*domain.GameQuestion, 0, len(questionRows))
	questionIDs := make([]int64, 0, len(questionRows))
	for _, row := range questionRows {
		var sourceSenseID *int64
		if row.SourceSenseID.Valid {
			val := row.SourceSenseID.Int64
			sourceSenseID = &val
		}

		question := &domain.GameQuestion{
			ID:                  row.ID,
			SessionID:           row.SessionID,
			QuestionOrder:       row.QuestionOrder,
			QuestionType:        row.QuestionType,
			SourceWordID:        row.SourceWordID,
			SourceSenseID:       sourceSenseID,
			CorrectTargetWordID: row.CorrectTargetWordID,
			SourceLanguageID:    row.SourceLanguageID,
			TargetLanguageID:    row.TargetLanguageID,
			CreatedAt:           row.CreatedAt.Time,
		}
		questions = append(questions, question)
		questionIDs = append(questionIDs, question.ID)
	}

	if len(questionIDs) == 0 {
		return &domain.GameQuestionsResult{
			Questions: []*domain.GameQuestion{},
			Options:   []*domain.GameQuestionOption{},
		}, nil
	}

	optionRows, err := r.queries.FindGameQuestionOptionsByQuestionIDs(ctx, questionIDs)
	if err != nil {
		return nil, sharederrors.MapGameRepositoryError(err, "FindGameQuestionsBySessionID")
	}

	options := make([]*domain.GameQuestionOption, 0, len(optionRows))
	for _, row := range optionRows {
		options = append(options, &domain.GameQuestionOption{
			ID:           row.ID,
			QuestionID:   row.QuestionID,
			OptionLabel:  row.OptionLabel,
			TargetWordID: row.TargetWordID,
			IsCorrect:    row.IsCorrect,
		})
	}

	return &domain.GameQuestionsResult{
		Questions: questions,
		Options:   options,
	}, nil
}

// FindGameQuestionByID returns a question by ID with its options
func (r *gameQuestionRepo) FindGameQuestionByID(ctx context.Context, questionID int64) (*domain.GameQuestionWithOptions, error) {
	questionRow, err := r.queries.FindGameQuestionByID(ctx, questionID)
	if err != nil {
		return nil, sharederrors.MapGameRepositoryError(err, "FindGameQuestionByID")
	}

	var sourceSenseID *int64
	if questionRow.SourceSenseID.Valid {
		val := questionRow.SourceSenseID.Int64
		sourceSenseID = &val
	}

	question := &domain.GameQuestion{
		ID:                  questionRow.ID,
		SessionID:           questionRow.SessionID,
		QuestionOrder:       questionRow.QuestionOrder,
		QuestionType:        questionRow.QuestionType,
		SourceWordID:        questionRow.SourceWordID,
		SourceSenseID:       sourceSenseID,
		CorrectTargetWordID: questionRow.CorrectTargetWordID,
		SourceLanguageID:    questionRow.SourceLanguageID,
		TargetLanguageID:    questionRow.TargetLanguageID,
		CreatedAt:           questionRow.CreatedAt.Time,
	}

	optionRows, err := r.queries.FindGameQuestionOptionsByQuestionID(ctx, questionID)
	if err != nil {
		return nil, sharederrors.MapGameRepositoryError(err, "FindGameQuestionByID")
	}

	options := make([]*domain.GameQuestionOption, 0, len(optionRows))
	for _, row := range optionRows {
		options = append(options, &domain.GameQuestionOption{
			ID:           row.ID,
			QuestionID:   row.QuestionID,
			OptionLabel:  row.OptionLabel,
			TargetWordID: row.TargetWordID,
			IsCorrect:    row.IsCorrect,
		})
	}

	return &domain.GameQuestionWithOptions{
		Question: question,
		Options:  options,
	}, nil
}
