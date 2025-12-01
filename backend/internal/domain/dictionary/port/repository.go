package port

import (
	"context"

	"github.com/english-coach/backend/internal/domain/dictionary/model"
)

// LanguageRepository defines operations for language data access
type LanguageRepository interface {
	// FindAll returns all languages
	FindAll(ctx context.Context) ([]*model.Language, error)
	// FindByID returns a language by ID
	FindByID(ctx context.Context, id int16) (*model.Language, error)
	// FindByCode returns a language by code
	FindByCode(ctx context.Context, code string) (*model.Language, error)
}

// TopicRepository defines operations for topic data access
type TopicRepository interface {
	// FindAll returns all topics
	FindAll(ctx context.Context) ([]*model.Topic, error)
	// FindByID returns a topic by ID
	FindByID(ctx context.Context, id int64) (*model.Topic, error)
	// FindByCode returns a topic by code
	FindByCode(ctx context.Context, code string) (*model.Topic, error)
}

// LevelRepository defines operations for level data access
type LevelRepository interface {
	// FindAll returns all levels
	FindAll(ctx context.Context) ([]*model.Level, error)
	// FindByID returns a level by ID
	FindByID(ctx context.Context, id int64) (*model.Level, error)
	// FindByCode returns a level by code
	FindByCode(ctx context.Context, code string) (*model.Level, error)
	// FindByLanguageID returns all levels for a specific language
	FindByLanguageID(ctx context.Context, languageID int16) ([]*model.Level, error)
}

// WordRepository defines operations for word data access
type WordRepository interface {
	// FindByID returns a word by ID
	FindByID(ctx context.Context, id int64) (*model.Word, error)
	// FindByIDs returns multiple words by their IDs
	FindByIDs(ctx context.Context, ids []int64) ([]*model.Word, error)
	// FindWordsByTopicAndLanguages finds words filtered by topic and language pair
	FindWordsByTopicAndLanguages(ctx context.Context, topicID int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*model.Word, error)
	// FindWordsByLevelAndLanguages finds words filtered by level and language pair
	FindWordsByLevelAndLanguages(ctx context.Context, levelID int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*model.Word, error)
	// FindTranslationsForWord finds translation words for a given source word and target language
	FindTranslationsForWord(ctx context.Context, sourceWordID int64, targetLanguageID int16, limit int) ([]*model.Word, error)
	// SearchWords searches for words using multiple strategies (lemma, normalized, search_key)
	SearchWords(ctx context.Context, query string, languageID *int16, limit, offset int) ([]*model.Word, error)
	// CountSearchWords returns the total count of words matching the search query
	CountSearchWords(ctx context.Context, query string, languageID *int16) (int, error)
}

// SenseRepository defines operations for sense data access
type SenseRepository interface {
	// FindByWordID returns all senses for a word, ordered by sense_order
	FindByWordID(ctx context.Context, wordID int64) ([]*model.Sense, error)
	// FindByWordIDs returns senses for multiple words
	FindByWordIDs(ctx context.Context, wordIDs []int64) (map[int64][]*model.Sense, error)
}
