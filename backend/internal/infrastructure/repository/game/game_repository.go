package game

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/english-coach/backend/internal/domain/game/port"
	db "github.com/english-coach/backend/internal/infrastructure/db/sqlc/gen/game"
)

// GameRepository implements game repository interfaces using sqlc
type GameRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewGameRepository creates a new game repository
func NewGameRepository(pool *pgxpool.Pool) *GameRepository {
	return &GameRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// GameSessionRepo returns a GameSessionRepository implementation
func (r *GameRepository) GameSessionRepo() port.GameSessionRepository {
	return &gameSessionRepo{
		GameRepository: r,
	}
}

// GameQuestionRepo returns a GameQuestionRepository implementation
func (r *GameRepository) GameQuestionRepo() port.GameQuestionRepository {
	return &gameQuestionRepo{
		GameRepository: r,
	}
}

// GameAnswerRepo returns a GameAnswerRepository implementation
func (r *GameRepository) GameAnswerRepo() port.GameAnswerRepository {
	return &gameAnswerRepo{
		GameRepository: r,
	}
}
