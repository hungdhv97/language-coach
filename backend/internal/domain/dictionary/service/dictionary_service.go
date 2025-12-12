package service

import (
	"context"

	"github.com/english-coach/backend/internal/domain/dictionary/dto"
	"github.com/english-coach/backend/internal/domain/dictionary/model"
	"github.com/english-coach/backend/internal/domain/dictionary/port"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// DictionaryService provides dictionary lookup functionality
type DictionaryService struct {
	wordRepo         port.WordRepository
	senseRepo        port.SenseRepository
	languageRepo     port.LanguageRepository
	levelRepo        port.LevelRepository
	partOfSpeechRepo port.PartOfSpeechRepository
	pool             *pgxpool.Pool
	logger           *zap.Logger
}

// NewDictionaryService creates a new dictionary service
func NewDictionaryService(
	wordRepo port.WordRepository,
	senseRepo port.SenseRepository,
	languageRepo port.LanguageRepository,
	levelRepo port.LevelRepository,
	partOfSpeechRepo port.PartOfSpeechRepository,
	pool *pgxpool.Pool,
	logger *zap.Logger,
) *DictionaryService {
	return &DictionaryService{
		wordRepo:         wordRepo,
		senseRepo:        senseRepo,
		languageRepo:     languageRepo,
		levelRepo:        levelRepo,
		partOfSpeechRepo: partOfSpeechRepo,
		pool:             pool,
		logger:           logger,
	}
}

// GetWordDetail retrieves detailed information about a word including senses, translations, examples, and pronunciations
func (s *DictionaryService) GetWordDetail(ctx context.Context, wordID int64) (*dto.WordDetail, error) {
	// Get word
	word, err := s.wordRepo.FindByID(ctx, wordID)
	if err != nil {
		return nil, err
	}

	// Get senses
	senses, err := s.senseRepo.FindByWordID(ctx, wordID)
	if err != nil {
		return nil, err
	}

	// Get translations for each sense
	senseIDs := make([]int64, len(senses))
	for i, sense := range senses {
		senseIDs[i] = sense.ID
	}

	// Get sense translations
	translations, err := s.getSenseTranslations(ctx, senseIDs)
	if err != nil {
		return nil, err
	}

	// Get examples for senses
	examples, err := s.getExamples(ctx, senseIDs)
	if err != nil {
		return nil, err
	}

	// Get pronunciations
	pronunciations, err := s.getPronunciations(ctx, wordID)
	if err != nil {
		return nil, err
	}

	// Get part of speech IDs and language IDs for lookup
	posIDs := make([]int16, 0)
	levelIDs := make([]int64, 0)
	langIDs := make([]int16, 0)
	langIDSet := make(map[int16]bool)

	for _, sense := range senses {
		posIDs = append(posIDs, sense.PartOfSpeechID)
		if sense.LevelID != nil {
			levelIDs = append(levelIDs, *sense.LevelID)
		}
		if !langIDSet[sense.DefinitionLanguageID] {
			langIDs = append(langIDs, sense.DefinitionLanguageID)
			langIDSet[sense.DefinitionLanguageID] = true
		}
	}

	// Fetch part of speech names
	posMap := make(map[int16]*string)
	if len(posIDs) > 0 {
		posData, err := s.partOfSpeechRepo.FindByIDs(ctx, posIDs)
		if err != nil {
			s.logger.Warn("failed to fetch part of speech names", zap.Error(err))
		} else {
			for id, pos := range posData {
				name := pos.Name
				posMap[id] = &name
			}
		}
	}

	// Fetch level names
	levelMap := make(map[int64]*string)
	for _, levelID := range levelIDs {
		level, err := s.levelRepo.FindByID(ctx, levelID)
		if err != nil {
			s.logger.Warn("failed to fetch level name", zap.Int64("level_id", levelID), zap.Error(err))
		} else {
			levelMap[levelID] = &level.Name
		}
	}

	// Fetch language names
	langMap := make(map[int16]*string)
	for _, langID := range langIDs {
		lang, err := s.languageRepo.FindByID(ctx, langID)
		if err != nil {
			s.logger.Warn("failed to fetch language name", zap.Int16("language_id", langID), zap.Error(err))
		} else {
			name := lang.Name
			if lang.NativeName != nil {
				name = *lang.NativeName
			}
			langMap[langID] = &name
		}
	}

	// Build sense details
	senseDetails := make([]dto.SenseDetail, len(senses))
	for i, sense := range senses {
		var levelName *string
		if sense.LevelID != nil {
			if name, ok := levelMap[*sense.LevelID]; ok {
				levelName = name
			}
		}

		senseDetails[i] = dto.SenseDetail{
			ID:                     sense.ID,
			SenseOrder:             sense.SenseOrder,
			PartOfSpeechID:         sense.PartOfSpeechID,
			PartOfSpeechName:       posMap[sense.PartOfSpeechID],
			Definition:             sense.Definition,
			DefinitionLanguageID:   sense.DefinitionLanguageID,
			DefinitionLanguageName: langMap[sense.DefinitionLanguageID],
			UsageLabel:             sense.UsageLabel,
			LevelID:                sense.LevelID,
			LevelName:              levelName,
			Note:                   sense.Note,
			Translations:           translations[sense.ID],
			Examples:               examples[sense.ID],
		}
	}

	return &dto.WordDetail{
		Word:           word,
		Senses:         senseDetails,
		Pronunciations: pronunciations,
	}, nil
}

// getSenseTranslations retrieves translations for given sense IDs
func (s *DictionaryService) getSenseTranslations(ctx context.Context, senseIDs []int64) (map[int64][]*model.Word, error) {
	if len(senseIDs) == 0 {
		return make(map[int64][]*model.Word), nil
	}

	query := `
		SELECT st.source_sense_id, tw.id, tw.language_id, tw.lemma, tw.lemma_normalized, tw.search_key,
		       tw.romanization, tw.script_code, tw.frequency_rank,
		       tw.note, tw.created_at, tw.updated_at
		FROM sense_translations st
		INNER JOIN words tw ON st.target_word_id = tw.id
		WHERE st.source_sense_id = ANY($1)
		ORDER BY st.source_sense_id, st.priority, tw.frequency_rank NULLS LAST
	`
	rows, err := s.pool.Query(ctx, query, senseIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]*model.Word)
	for rows.Next() {
		var senseID int64
		var word model.Word
		var lemmaNormalized, searchKey, romanization, scriptCode, note *string
		var frequencyRank *int
		if err := rows.Scan(
			&senseID,
			&word.ID,
			&word.LanguageID,
			&word.Lemma,
			&lemmaNormalized,
			&searchKey,
			&romanization,
			&scriptCode,
			&frequencyRank,
			&note,
			&word.CreatedAt,
			&word.UpdatedAt,
		); err != nil {
			return nil, err
		}
		word.LemmaNormalized = lemmaNormalized
		word.SearchKey = searchKey
		word.Romanization = romanization
		word.ScriptCode = scriptCode
		word.FrequencyRank = frequencyRank
		word.Note = note
		result[senseID] = append(result[senseID], &word)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// getExamples retrieves examples for given sense IDs
func (s *DictionaryService) getExamples(ctx context.Context, senseIDs []int64) (map[int64][]*model.Example, error) {
	if len(senseIDs) == 0 {
		return make(map[int64][]*model.Example), nil
	}

	query := `
		SELECT e.id, e.source_sense_id, e.language_id, e.content, e.audio_url, e.source
		FROM examples e
		WHERE e.source_sense_id = ANY($1)
		ORDER BY e.source_sense_id, e.id
	`
	rows, err := s.pool.Query(ctx, query, senseIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]*model.Example)
	for rows.Next() {
		var example model.Example
		var audioURL, source *string
		if err := rows.Scan(
			&example.ID,
			&example.SourceSenseID,
			&example.LanguageID,
			&example.Content,
			&audioURL,
			&source,
		); err != nil {
			return nil, err
		}
		example.AudioURL = audioURL
		example.Source = source
		result[example.SourceSenseID] = append(result[example.SourceSenseID], &example)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// getPronunciations retrieves pronunciations for a word
func (s *DictionaryService) getPronunciations(ctx context.Context, wordID int64) ([]*model.Pronunciation, error) {
	query := `
		SELECT id, word_id, dialect, ipa, phonetic, audio_url
		FROM pronunciations
		WHERE word_id = $1
		ORDER BY id
	`
	rows, err := s.pool.Query(ctx, query, wordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pronunciations []*model.Pronunciation
	for rows.Next() {
		var pron model.Pronunciation
		var dialect, ipa, phonetic, audioURL *string
		if err := rows.Scan(
			&pron.ID,
			&pron.WordID,
			&dialect,
			&ipa,
			&phonetic,
			&audioURL,
		); err != nil {
			return nil, err
		}
		pron.Dialect = dialect
		pron.IPA = ipa
		pron.Phonetic = phonetic
		pron.AudioURL = audioURL
		pronunciations = append(pronunciations, &pron)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pronunciations, nil
}
