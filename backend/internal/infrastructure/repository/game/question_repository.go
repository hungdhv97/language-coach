package game

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/english-coach/backend/internal/domain/game/model"
	db "github.com/english-coach/backend/internal/infrastructure/db/sqlc/gen/game"
	"github.com/english-coach/backend/internal/infrastructure/repository/common"
)

// gameQuestionRepo implements GameQuestionRepository using sqlc
type gameQuestionRepo struct {
	*GameRepository
}

// CreateBatch creates multiple questions and their options in a transaction
func (r *gameQuestionRepo) CreateBatch(ctx context.Context, questions []*model.GameQuestion, options []*model.GameQuestionOption) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return common.MapPgError(err)
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
			return common.MapPgError(err)
		}
		question.ID = result.ID
		question.CreatedAt = result.CreatedAt.Time
	}

	// Insert options
	for _, option := range options {
		optionID, err := qtx.CreateGameQuestionOption(ctx, db.CreateGameQuestionOptionParams{
			QuestionID:   option.QuestionID,
			OptionLabel:  option.OptionLabel,
			TargetWordID: option.TargetWordID,
			IsCorrect:    option.IsCorrect,
		})
		if err != nil {
			return common.MapPgError(err)
		}
		option.ID = optionID
	}

	return tx.Commit(ctx)
}

// FindBySessionID returns all questions for a session
func (r *gameQuestionRepo) FindBySessionID(ctx context.Context, sessionID int64) ([]*model.GameQuestion, []*model.GameQuestionOption, error) {
	questionRows, err := r.queries.FindGameQuestionsBySessionID(ctx, sessionID)
	if err != nil {
		return nil, nil, common.MapPgError(err)
	}

	questions := make([]*model.GameQuestion, 0, len(questionRows))
	questionIDs := make([]int64, 0, len(questionRows))
	for _, row := range questionRows {
		var sourceSenseID *int64
		if row.SourceSenseID.Valid {
			val := row.SourceSenseID.Int64
			sourceSenseID = &val
		}

		question := &model.GameQuestion{
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
		return questions, []*model.GameQuestionOption{}, nil
	}

	optionRows, err := r.queries.FindGameQuestionOptionsByQuestionIDs(ctx, questionIDs)
	if err != nil {
		return nil, nil, common.MapPgError(err)
	}

	options := make([]*model.GameQuestionOption, 0, len(optionRows))
	for _, row := range optionRows {
		options = append(options, &model.GameQuestionOption{
			ID:           row.ID,
			QuestionID:   row.QuestionID,
			OptionLabel:  row.OptionLabel,
			TargetWordID: row.TargetWordID,
			IsCorrect:    row.IsCorrect,
		})
	}

	return questions, options, nil
}

// FindByID returns a question by ID with its options
func (r *gameQuestionRepo) FindByID(ctx context.Context, questionID int64) (*model.GameQuestion, []*model.GameQuestionOption, error) {
	questionRow, err := r.queries.FindGameQuestionByID(ctx, questionID)
	if err != nil {
		return nil, nil, common.MapPgError(err)
	}

	var sourceSenseID *int64
	if questionRow.SourceSenseID.Valid {
		val := questionRow.SourceSenseID.Int64
		sourceSenseID = &val
	}

	question := &model.GameQuestion{
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
		return nil, nil, common.MapPgError(err)
	}

	options := make([]*model.GameQuestionOption, 0, len(optionRows))
	for _, row := range optionRows {
		options = append(options, &model.GameQuestionOption{
			ID:           row.ID,
			QuestionID:   row.QuestionID,
			OptionLabel:  row.OptionLabel,
			TargetWordID: row.TargetWordID,
			IsCorrect:    row.IsCorrect,
		})
	}

	return question, options, nil
}
