package domain

import (
	"context"
)

// LanguageRepository defines operations for language data access
type LanguageRepository interface {
	// FindAllLanguages returns all languages
	FindAllLanguages(ctx context.Context) ([]*Language, error)
	// FindLanguageByID returns a language by ID
	FindLanguageByID(ctx context.Context, id int16) (*Language, error)
	// FindLanguageByCode returns a language by code
	FindLanguageByCode(ctx context.Context, code string) (*Language, error)
}

// TopicRepository defines operations for topic data access
type TopicRepository interface {
	// FindAllTopics returns all topics
	FindAllTopics(ctx context.Context) ([]*Topic, error)
	// FindTopicByID returns a topic by ID
	FindTopicByID(ctx context.Context, id int64) (*Topic, error)
	// FindTopicByCode returns a topic by code
	FindTopicByCode(ctx context.Context, code string) (*Topic, error)
}

// LevelRepository defines operations for level data access
type LevelRepository interface {
	// FindAllLevels returns all levels
	FindAllLevels(ctx context.Context) ([]*Level, error)
	// FindLevelByID returns a level by ID
	FindLevelByID(ctx context.Context, id int64) (*Level, error)
	// FindLevelByCode returns a level by code
	FindLevelByCode(ctx context.Context, code string) (*Level, error)
	// FindLevelsByLanguageID returns all levels for a specific language
	FindLevelsByLanguageID(ctx context.Context, languageID int16) ([]*Level, error)
}

// WordRepository defines operations for word data access
type WordRepository interface {
	// FindWordByID returns a word by ID
	FindWordByID(ctx context.Context, id int64) (*Word, error)
	// FindWordsByIDs returns multiple words by their IDs
	FindWordsByIDs(ctx context.Context, ids []int64) ([]*Word, error)
	// FindWordsByTopicAndLanguages finds words filtered by topic and language pair
	FindWordsByTopicAndLanguages(ctx context.Context, topicID int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*Word, error)
	// FindWordsByLevelAndLanguages finds words filtered by level and language pair
	FindWordsByLevelAndLanguages(ctx context.Context, levelID int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*Word, error)
	// FindWordsByLevelAndTopicsAndLanguages finds words filtered by level, optional topics, and language pair
	// If topicIDs is nil or empty, returns all words for the level (no topic filter)
	FindWordsByLevelAndTopicsAndLanguages(ctx context.Context, levelID int64, topicIDs []int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*Word, error)
	// FindTranslationsForWord finds translation words for a given source word and target language
	FindTranslationsForWord(ctx context.Context, sourceWordID int64, targetLanguageID int16, limit int) ([]*Word, error)
	// SearchWords searches for words using multiple strategies (lemma, normalized, search_key)
	SearchWords(ctx context.Context, query string, languageID int16, limit, offset int) ([]*Word, error)
	// CountSearchWords returns the total count of words matching the search query
	CountSearchWords(ctx context.Context, query string, languageID int16) (int, error)
}

// SenseRepository defines operations for sense data access
type SenseRepository interface {
	// FindSensesByWordID returns all senses for a word, ordered by sense_order
	FindSensesByWordID(ctx context.Context, wordID int64) ([]*Sense, error)
	// FindSensesByWordIDs returns senses for multiple words
	FindSensesByWordIDs(ctx context.Context, wordIDs []int64) (map[int64][]*Sense, error)
}

// PartOfSpeechRepository defines operations for part of speech data access
type PartOfSpeechRepository interface {
	// FindAllPartsOfSpeech returns all parts of speech
	FindAllPartsOfSpeech(ctx context.Context) ([]*PartOfSpeech, error)
	// FindPartOfSpeechByID returns a part of speech by ID
	FindPartOfSpeechByID(ctx context.Context, id int16) (*PartOfSpeech, error)
	// FindPartOfSpeechByCode returns a part of speech by code
	FindPartOfSpeechByCode(ctx context.Context, code string) (*PartOfSpeech, error)
	// FindPartsOfSpeechByIDs returns parts of speech by their IDs
	FindPartsOfSpeechByIDs(ctx context.Context, ids []int16) (map[int16]*PartOfSpeech, error)
}

