package dictionary

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/english-coach/backend/internal/domain/dictionary/port"
	db "github.com/english-coach/backend/internal/infrastructure/db/sqlc/gen/dictionary"
)

// DictionaryRepository implements dictionary repository interfaces using sqlc
type DictionaryRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewDictionaryRepository creates a new dictionary repository
func NewDictionaryRepository(pool *pgxpool.Pool) *DictionaryRepository {
	return &DictionaryRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// LanguageRepository returns a LanguageRepository implementation
func (r *DictionaryRepository) LanguageRepository() port.LanguageRepository {
	return &languageRepository{
		DictionaryRepository: r,
	}
}

// TopicRepository returns a TopicRepository implementation
func (r *DictionaryRepository) TopicRepository() port.TopicRepository {
	return &topicRepository{
		DictionaryRepository: r,
	}
}

// LevelRepository returns a LevelRepository implementation
func (r *DictionaryRepository) LevelRepository() port.LevelRepository {
	return &levelRepository{
		DictionaryRepository: r,
	}
}

// WordRepository returns a WordRepository implementation
func (r *DictionaryRepository) WordRepository() port.WordRepository {
	return &wordRepository{
		DictionaryRepository: r,
	}
}

// SenseRepository returns a SenseRepository implementation
func (r *DictionaryRepository) SenseRepository() port.SenseRepository {
	return &senseRepository{
		DictionaryRepository: r,
	}
}
