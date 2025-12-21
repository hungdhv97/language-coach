package game

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/english-coach/backend/internal/modules/game/domain"
	db "github.com/english-coach/backend/internal/platform/db/sqlc/gen/game"
	"github.com/english-coach/backend/internal/shared/errors"
)

// gameSessionRepo implements GameSessionRepository using sqlc
type gameSessionRepo struct {
	*GameRepository
}

// Create creates a new game session
func (r *gameSessionRepo) Create(ctx context.Context, session *domain.GameSession) error {
	var topicID, levelID pgtype.Int8
	if session.TopicID != nil {
		topicID = pgtype.Int8{Int64: *session.TopicID, Valid: true}
	}
	if session.LevelID != nil {
		levelID = pgtype.Int8{Int64: *session.LevelID, Valid: true}
	}

	totalQuestions := pgtype.Int2{Int16: session.TotalQuestions, Valid: true}
	correctQuestions := pgtype.Int2{Int16: session.CorrectQuestions, Valid: true}
	startedAt := pgtype.Timestamp{Time: time.Now(), Valid: true}

	result, err := r.queries.CreateGameSession(ctx, db.CreateGameSessionParams{
		UserID:           session.UserID,
		Mode:             session.Mode,
		SourceLanguageID: session.SourceLanguageID,
		TargetLanguageID: session.TargetLanguageID,
		TopicID:          topicID,
		LevelID:          levelID,
		TotalQuestions:   totalQuestions,
		CorrectQuestions: correctQuestions,
		StartedAt:        startedAt,
	})
	if err != nil {
		return errors.MapPgError(err)
	}

	session.ID = result.ID
	session.StartedAt = result.StartedAt.Time
	return nil
}

// FindByID returns a game session by ID
func (r *gameSessionRepo) FindByID(ctx context.Context, id int64) (*domain.GameSession, error) {
	row, err := r.queries.FindGameSessionByID(ctx, id)
	if err != nil {
		return nil, errors.MapPgError(err)
	}

	var topicID, levelID *int64
	var endedAt *time.Time

	if row.TopicID.Valid {
		val := row.TopicID.Int64
		topicID = &val
	}
	if row.LevelID.Valid {
		val := row.LevelID.Int64
		levelID = &val
	}
	if row.EndedAt.Valid {
		endedAt = &row.EndedAt.Time
	}

	return &domain.GameSession{
		ID:               row.ID,
		UserID:           row.UserID,
		Mode:             row.Mode,
		SourceLanguageID: row.SourceLanguageID,
		TargetLanguageID: row.TargetLanguageID,
		TopicID:          topicID,
		LevelID:          levelID,
		TotalQuestions:   int16(row.TotalQuestions.Int16),
		CorrectQuestions: int16(row.CorrectQuestions.Int16),
		StartedAt:        row.StartedAt.Time,
		EndedAt:          endedAt,
	}, nil
}

// Update updates a game session
func (r *gameSessionRepo) Update(ctx context.Context, session *domain.GameSession) error {
	totalQuestions := pgtype.Int2{Int16: session.TotalQuestions, Valid: true}
	correctQuestions := pgtype.Int2{Int16: session.CorrectQuestions, Valid: true}
	var endedAt pgtype.Timestamp
	if session.EndedAt != nil {
		endedAt = pgtype.Timestamp{Time: *session.EndedAt, Valid: true}
	}

	err := r.queries.UpdateGameSession(ctx, db.UpdateGameSessionParams{
		ID:               session.ID,
		TotalQuestions:   totalQuestions,
		CorrectQuestions: correctQuestions,
		EndedAt:          endedAt,
	})
	return errors.MapPgError(err)
}

// EndSession marks a session as ended
func (r *gameSessionRepo) EndSession(ctx context.Context, sessionID int64, endedAt interface{}) error {
	var endTime time.Time
	if endedAt != nil {
		if t, ok := endedAt.(time.Time); ok {
			endTime = t
		} else {
			endTime = time.Now()
		}
	} else {
		endTime = time.Now()
	}

	endedAtPg := pgtype.Timestamp{Time: endTime, Valid: true}
	err := r.queries.EndGameSession(ctx, db.EndGameSessionParams{
		ID:      sessionID,
		EndedAt: endedAtPg,
	})
	return errors.MapPgError(err)
}
