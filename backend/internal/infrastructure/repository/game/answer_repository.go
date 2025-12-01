package game

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/english-coach/backend/internal/domain/game/model"
	db "github.com/english-coach/backend/internal/infrastructure/db/sqlc/gen/game"
	"github.com/english-coach/backend/internal/infrastructure/repository/common"
)

// gameAnswerRepo implements GameAnswerRepository using sqlc
type gameAnswerRepo struct {
	*GameRepository
}

// Create creates a new answer
func (r *gameAnswerRepo) Create(ctx context.Context, answer *model.GameAnswer) error {
	var selectedOptionID pgtype.Int8
	if answer.SelectedOptionID != nil {
		selectedOptionID = pgtype.Int8{Int64: *answer.SelectedOptionID, Valid: true}
	}
	var responseTimeMs pgtype.Int4
	if answer.ResponseTimeMs != nil {
		responseTimeMs = pgtype.Int4{Int32: int32(*answer.ResponseTimeMs), Valid: true}
	}
	answeredAt := pgtype.Timestamp{Time: time.Now(), Valid: true}

	result, err := r.queries.CreateGameAnswer(ctx, db.CreateGameAnswerParams{
		QuestionID:       answer.QuestionID,
		SessionID:        answer.SessionID,
		UserID:           answer.UserID,
		SelectedOptionID: selectedOptionID,
		IsCorrect:        answer.IsCorrect,
		ResponseTimeMs:   responseTimeMs,
		AnsweredAt:       answeredAt,
	})
	if err != nil {
		return common.MapPgError(err)
	}

	answer.ID = result.ID
	answer.AnsweredAt = result.AnsweredAt.Time
	return nil
}

// FindByQuestionID returns the answer for a specific question
func (r *gameAnswerRepo) FindByQuestionID(ctx context.Context, questionID, sessionID, userID int64) (*model.GameAnswer, error) {
	row, err := r.queries.FindGameAnswerByQuestionID(ctx, db.FindGameAnswerByQuestionIDParams{
		QuestionID: questionID,
		SessionID:  sessionID,
		UserID:     userID,
	})
	if err != nil {
		return nil, common.MapPgError(err)
	}

	var selectedOptionID *int64
	var responseTimeMs *int

	if row.SelectedOptionID.Valid {
		val := row.SelectedOptionID.Int64
		selectedOptionID = &val
	}
	if row.ResponseTimeMs.Valid {
		val := int(row.ResponseTimeMs.Int32)
		responseTimeMs = &val
	}

	return &model.GameAnswer{
		ID:               row.ID,
		QuestionID:       row.QuestionID,
		SessionID:        row.SessionID,
		UserID:           row.UserID,
		SelectedOptionID: selectedOptionID,
		IsCorrect:        row.IsCorrect,
		ResponseTimeMs:   responseTimeMs,
		AnsweredAt:       row.AnsweredAt.Time,
	}, nil
}

// FindBySessionID returns all answers for a session
func (r *gameAnswerRepo) FindBySessionID(ctx context.Context, sessionID, userID int64) ([]*model.GameAnswer, error) {
	rows, err := r.queries.FindGameAnswersBySessionID(ctx, db.FindGameAnswersBySessionIDParams{
		SessionID: sessionID,
		UserID:    userID,
	})
	if err != nil {
		return nil, common.MapPgError(err)
	}

	answers := make([]*model.GameAnswer, 0, len(rows))
	for _, row := range rows {
		var selectedOptionID *int64
		var responseTimeMs *int

		if row.SelectedOptionID.Valid {
			val := row.SelectedOptionID.Int64
			selectedOptionID = &val
		}
		if row.ResponseTimeMs.Valid {
			val := int(row.ResponseTimeMs.Int32)
			responseTimeMs = &val
		}

		answers = append(answers, &model.GameAnswer{
			ID:               row.ID,
			QuestionID:       row.QuestionID,
			SessionID:        row.SessionID,
			UserID:           row.UserID,
			SelectedOptionID: selectedOptionID,
			IsCorrect:        row.IsCorrect,
			ResponseTimeMs:   responseTimeMs,
			AnsweredAt:       row.AnsweredAt.Time,
		})
	}

	return answers, nil
}
