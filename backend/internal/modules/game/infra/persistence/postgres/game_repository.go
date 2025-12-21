package game

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/english-coach/backend/internal/modules/game/domain"
	db "github.com/english-coach/backend/internal/platform/db/sqlc/gen/game"
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
func (r *GameRepository) GameSessionRepo() domain.GameSessionRepository {
	return &gameSessionRepo{
		GameRepository: r,
	}
}

// GameQuestionRepo returns a GameQuestionRepository implementation
func (r *GameRepository) GameQuestionRepo() domain.GameQuestionRepository {
	return &gameQuestionRepo{
		GameRepository: r,
	}
}

// GameAnswerRepo returns a GameAnswerRepository implementation
func (r *GameRepository) GameAnswerRepo() domain.GameAnswerRepository {
	return &gameAnswerRepo{
		GameRepository: r,
	}
}
