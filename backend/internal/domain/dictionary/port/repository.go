package port

import (
	"context"

	"github.com/english-coach/backend/internal/modules/dictionary/domain"
)

// LanguageRepository defines operations for language data access
type LanguageRepository interface {
	// FindAll returns all languages
	FindAll(ctx context.Context) ([]*domain.Language, error)
	// FindByID returns a language by ID
	FindByID(ctx context.Context, id int16) (*domain.Language, error)
	// FindByCode returns a language by code
	FindByCode(ctx context.Context, code string) (*domain.Language, error)
}

// TopicRepository defines operations for topic data access
type TopicRepository interface {
	// FindAll returns all topics
	FindAll(ctx context.Context) ([]*domain.Topic, error)
	// FindByID returns a topic by ID
	FindByID(ctx context.Context, id int64) (*domain.Topic, error)
	// FindByCode returns a topic by code
	FindByCode(ctx context.Context, code string) (*domain.Topic, error)
}

// LevelRepository defines operations for level data access
type LevelRepository interface {
	// FindAll returns all levels
	FindAll(ctx context.Context) ([]*domain.Level, error)
	// FindByID returns a level by ID
	FindByID(ctx context.Context, id int64) (*domain.Level, error)
	// FindByCode returns a level by code
	FindByCode(ctx context.Context, code string) (*domain.Level, error)
	// FindByLanguageID returns all levels for a specific language
	FindByLanguageID(ctx context.Context, languageID int16) ([]*domain.Level, error)
}

// WordRepository defines operations for word data access
type WordRepository interface {
	// FindByID returns a word by ID
	FindByID(ctx context.Context, id int64) (*domain.Word, error)
	// FindByIDs returns multiple words by their IDs
	FindByIDs(ctx context.Context, ids []int64) ([]*domain.Word, error)
	// FindWordsByTopicAndLanguages finds words filtered by topic and language pair
	FindWordsByTopicAndLanguages(ctx context.Context, topicID int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*domain.Word, error)
	// FindWordsByLevelAndLanguages finds words filtered by level and language pair
	FindWordsByLevelAndLanguages(ctx context.Context, levelID int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*domain.Word, error)
	// FindWordsByLevelAndTopicsAndLanguages finds words filtered by level, optional topics, and language pair
	// If topicIDs is nil or empty, returns all words for the level (no topic filter)
	FindWordsByLevelAndTopicsAndLanguages(ctx context.Context, levelID int64, topicIDs []int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*domain.Word, error)
	// FindTranslationsForWord finds translation words for a given source word and target language
	FindTranslationsForWord(ctx context.Context, sourceWordID int64, targetLanguageID int16, limit int) ([]*domain.Word, error)
	// SearchWords searches for words using multiple strategies (lemma, normalized, search_key)
	SearchWords(ctx context.Context, query string, languageID int16, limit, offset int) ([]*domain.Word, error)
	// CountSearchWords returns the total count of words matching the search query
	CountSearchWords(ctx context.Context, query string, languageID int16) (int, error)
}

// SenseRepository defines operations for sense data access
type SenseRepository interface {
	// FindByWordID returns all senses for a word, ordered by sense_order
	FindByWordID(ctx context.Context, wordID int64) ([]*domain.Sense, error)
	// FindByWordIDs returns senses for multiple words
	FindByWordIDs(ctx context.Context, wordIDs []int64) (map[int64][]*domain.Sense, error)
}

// PartOfSpeechRepository defines operations for part of speech data access
type PartOfSpeechRepository interface {
	// FindAll returns all parts of speech
	FindAll(ctx context.Context) ([]*domain.PartOfSpeech, error)
	// FindByID returns a part of speech by ID
	FindByID(ctx context.Context, id int16) (*domain.PartOfSpeech, error)
	// FindByCode returns a part of speech by code
	FindByCode(ctx context.Context, code string) (*domain.PartOfSpeech, error)
	// FindByIDs returns parts of speech by their IDs
	FindByIDs(ctx context.Context, ids []int16) (map[int16]*domain.PartOfSpeech, error)
}
