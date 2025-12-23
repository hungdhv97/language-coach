package get_word_detail

import (
	"context"
	"fmt"

	"github.com/english-coach/backend/internal/modules/dictionary/domain"
	"github.com/english-coach/backend/internal/shared/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler provides dictionary lookup functionality
type Handler struct {
	wordRepo         domain.WordRepository
	senseRepo        domain.SenseRepository
	languageRepo     domain.LanguageRepository
	levelRepo        domain.LevelRepository
	partOfSpeechRepo domain.PartOfSpeechRepository
	pool             *pgxpool.Pool
	logger           logger.ILogger
}

// NewHandler creates a new dictionary handler
func NewHandler(
	wordRepo domain.WordRepository,
	senseRepo domain.SenseRepository,
	languageRepo domain.LanguageRepository,
	levelRepo domain.LevelRepository,
	partOfSpeechRepo domain.PartOfSpeechRepository,
	pool *pgxpool.Pool,
	logger logger.ILogger,
) *Handler {
	return &Handler{
		wordRepo:         wordRepo,
		senseRepo:        senseRepo,
		languageRepo:     languageRepo,
		levelRepo:        levelRepo,
		partOfSpeechRepo: partOfSpeechRepo,
		pool:             pool,
		logger:           logger,
	}
}

// Execute retrieves detailed information about a word including senses, translations, examples, and pronunciations
func (h *Handler) Execute(ctx context.Context, input Input) (*Output, error) {
	// Get word
	word, err := h.wordRepo.FindByID(ctx, input.WordID)
	if err != nil {
		return nil, err
	}

	// Check if word is nil
	if word == nil {
		return nil, fmt.Errorf("word not found")
	}

	// Get senses
	senses, err := h.senseRepo.FindByWordID(ctx, input.WordID)
	if err != nil {
		return nil, err
	}

	// Get translations for each sense
	senseIDs := make([]int64, len(senses))
	for i, sense := range senses {
		senseIDs[i] = sense.ID
	}

	// Get sense translations
	translations, err := h.getSenseTranslations(ctx, senseIDs)
	if err != nil {
		h.logger.Warn("failed to fetch sense translations", logger.Error(err))
		translations = make(map[int64][]*domain.Word)
	}
	if translations == nil {
		translations = make(map[int64][]*domain.Word)
	}

	// Get examples for senses
	examples, err := h.getExamples(ctx, senseIDs)
	if err != nil {
		h.logger.Warn("failed to fetch examples", logger.Error(err))
		examples = make(map[int64][]*domain.Example)
	}
	if examples == nil {
		examples = make(map[int64][]*domain.Example)
	}

	// Get pronunciations
	pronunciations, err := h.getPronunciations(ctx, input.WordID)
	if err != nil {
		h.logger.Warn("failed to fetch pronunciations", logger.Error(err))
		pronunciations = []*domain.Pronunciation{}
	}
	if pronunciations == nil {
		pronunciations = []*domain.Pronunciation{}
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
		posData, err := h.partOfSpeechRepo.FindByIDs(ctx, posIDs)
		if err != nil {
			h.logger.Warn("failed to fetch part of speech names", logger.Error(err))
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
		level, err := h.levelRepo.FindByID(ctx, levelID)
		if err != nil {
			h.logger.Warn("failed to fetch level name", logger.Int64("level_id", levelID), logger.Error(err))
		} else {
			levelMap[levelID] = &level.Name
		}
	}

	// Fetch language names
	langMap := make(map[int16]*string)
	for _, langID := range langIDs {
		lang, err := h.languageRepo.FindByID(ctx, langID)
		if err != nil {
			h.logger.Warn("failed to fetch language name", logger.Int("language_id", int(langID)), logger.Error(err))
		} else {
			name := lang.Name
			langMap[langID] = &name
		}
	}

	// Build sense details
	senseDetails := make([]SenseDetail, len(senses))
	for i, sense := range senses {
		var levelName *string
		if sense.LevelID != nil {
			if name, ok := levelMap[*sense.LevelID]; ok {
				levelName = name
			}
		}

		// Ensure translations and examples are not nil
		senseTranslations := translations[sense.ID]
		if senseTranslations == nil {
			senseTranslations = []*domain.Word{}
		}
		senseExamples := examples[sense.ID]
		if senseExamples == nil {
			senseExamples = []*domain.Example{}
		}

		senseDetails[i] = SenseDetail{
			ID:                   sense.ID,
			SenseOrder:           sense.SenseOrder,
			PartOfSpeechID:       sense.PartOfSpeechID,
			PartOfSpeechName:     posMap[sense.PartOfSpeechID],
			Definition:           sense.Definition,
			DefinitionLanguageID: sense.DefinitionLanguageID,
			LevelID:              sense.LevelID,
			LevelName:            levelName,
			Note:                 sense.Note,
			Translations:         senseTranslations,
			Examples:             senseExamples,
		}
	}

	// Get topics for word
	topics, err := h.getWordTopics(ctx, input.WordID)
	if err != nil {
		h.logger.Warn("failed to fetch word topics", logger.Error(err))
	} else {
		word.Topics = topics
	}

	// Get relations for word
	relations, err := h.getWordRelations(ctx, input.WordID)
	if err != nil {
		h.logger.Warn("failed to fetch word relations", logger.Error(err))
		relations = []*domain.WordRelation{}
	}
	if relations == nil {
		relations = []*domain.WordRelation{}
	}

	return &Output{
		Word:           word,
		Senses:         senseDetails,
		Pronunciations: pronunciations,
		Relations:      relations,
	}, nil
}

// getSenseTranslations retrieves translations for given sense IDs
func (h *Handler) getSenseTranslations(ctx context.Context, senseIDs []int64) (map[int64][]*domain.Word, error) {
	if len(senseIDs) == 0 {
		return make(map[int64][]*domain.Word), nil
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
	rows, err := h.pool.Query(ctx, query, senseIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]*domain.Word)
	for rows.Next() {
		var senseID int64
		var word domain.Word
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
func (h *Handler) getExamples(ctx context.Context, senseIDs []int64) (map[int64][]*domain.Example, error) {
	if len(senseIDs) == 0 {
		return make(map[int64][]*domain.Example), nil
	}

	query := `
		SELECT e.id, e.source_sense_id, e.language_id, e.content, e.audio_url, e.source
		FROM examples e
		WHERE e.source_sense_id = ANY($1)
		ORDER BY e.source_sense_id, e.id
	`
	rows, err := h.pool.Query(ctx, query, senseIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]*domain.Example)
	exampleMap := make(map[int64]*domain.Example)
	for rows.Next() {
		var example domain.Example
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
		example.Translations = []domain.ExampleTranslationSimple{}
		exampleMap[example.ID] = &example
		result[example.SourceSenseID] = append(result[example.SourceSenseID], &example)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Fetch example translations
	if len(exampleMap) > 0 {
		exampleIDs := make([]int64, 0, len(exampleMap))
		for id := range exampleMap {
			exampleIDs = append(exampleIDs, id)
		}

		transQuery := `
			SELECT et.example_id, l.code, et.content
			FROM example_translations et
			INNER JOIN languages l ON et.language_id = l.id
			WHERE et.example_id = ANY($1)
			ORDER BY et.example_id
		`
		transRows, err := h.pool.Query(ctx, transQuery, exampleIDs)
		if err != nil {
			h.logger.Warn("failed to fetch example translations", logger.Error(err))
		} else {
			defer transRows.Close()
			for transRows.Next() {
				var exampleID int64
				var langCode, content string
				if err := transRows.Scan(&exampleID, &langCode, &content); err != nil {
					h.logger.Warn("failed to scan example translation", logger.Error(err))
					continue
				}
				if example, ok := exampleMap[exampleID]; ok {
					example.Translations = append(example.Translations, domain.ExampleTranslationSimple{
						Language: langCode,
						Content:  content,
					})
				}
			}
		}
	}

	return result, nil
}

// getPronunciations retrieves pronunciations for a word
func (h *Handler) getPronunciations(ctx context.Context, wordID int64) ([]*domain.Pronunciation, error) {
	query := `
		SELECT id, word_id, dialect, ipa, phonetic, audio_url
		FROM pronunciations
		WHERE word_id = $1
		ORDER BY id
	`
	rows, err := h.pool.Query(ctx, query, wordID)
	if err != nil {
		return []*domain.Pronunciation{}, err
	}
	defer rows.Close()

	var pronunciations []*domain.Pronunciation
	for rows.Next() {
		var pron domain.Pronunciation
		var dialect, ipa, phonetic, audioURL *string
		if err := rows.Scan(
			&pron.ID,
			&pron.WordID,
			&dialect,
			&ipa,
			&phonetic,
			&audioURL,
		); err != nil {
			return []*domain.Pronunciation{}, err
		}
		pron.Dialect = dialect
		pron.IPA = ipa
		pron.Phonetic = phonetic
		pron.AudioURL = audioURL
		pronunciations = append(pronunciations, &pron)
	}

	if err := rows.Err(); err != nil {
		return []*domain.Pronunciation{}, err
	}

	if pronunciations == nil {
		return []*domain.Pronunciation{}, nil
	}

	return pronunciations, nil
}

// getWordTopics retrieves topic objects (with code and name) for a word
func (h *Handler) getWordTopics(ctx context.Context, wordID int64) ([]*domain.Topic, error) {
	query := `
		SELECT t.id, t.code, t.name
		FROM word_topics wt
		INNER JOIN topics t ON wt.topic_id = t.id
		WHERE wt.word_id = $1
		ORDER BY t.code
	`
	rows, err := h.pool.Query(ctx, query, wordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []*domain.Topic
	for rows.Next() {
		var topic domain.Topic
		if err := rows.Scan(&topic.ID, &topic.Code, &topic.Name); err != nil {
			return nil, err
		}
		topics = append(topics, &topic)
	}

	if err := rows.Err(); err != nil {
		return []*domain.Topic{}, err
	}

	if topics == nil {
		return []*domain.Topic{}, nil
	}

	return topics, nil
}

// getWordRelations retrieves relations for a word
func (h *Handler) getWordRelations(ctx context.Context, wordID int64) ([]*domain.WordRelation, error) {
	query := `
		SELECT wr.relation_type, wr.note, 
		       tw.id, tw.language_id, tw.lemma, tw.lemma_normalized, tw.search_key,
		       tw.romanization, tw.script_code, tw.frequency_rank, tw.note,
		       tw.created_at, tw.updated_at
		FROM word_relations wr
		INNER JOIN words tw ON wr.to_word_id = tw.id
		WHERE wr.from_word_id = $1
		ORDER BY wr.relation_type, tw.lemma
	`
	rows, err := h.pool.Query(ctx, query, wordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relations []*domain.WordRelation
	for rows.Next() {
		var relation domain.WordRelation
		var targetWord domain.Word
		var lemmaNormalized, searchKey, romanization, scriptCode, note, targetNote *string
		var frequencyRank *int
		if err := rows.Scan(
			&relation.RelationType,
			&note,
			&targetWord.ID,
			&targetWord.LanguageID,
			&targetWord.Lemma,
			&lemmaNormalized,
			&searchKey,
			&romanization,
			&scriptCode,
			&frequencyRank,
			&targetNote,
			&targetWord.CreatedAt,
			&targetWord.UpdatedAt,
		); err != nil {
			return nil, err
		}
		relation.Note = note
		targetWord.LemmaNormalized = lemmaNormalized
		targetWord.SearchKey = searchKey
		targetWord.Romanization = romanization
		targetWord.ScriptCode = scriptCode
		targetWord.FrequencyRank = frequencyRank
		targetWord.Note = targetNote

		// Get topics for target word
		targetTopics, err := h.getWordTopics(ctx, targetWord.ID)
		if err != nil {
			h.logger.Warn("failed to fetch topics for related word", logger.Int64("word_id", targetWord.ID), logger.Error(err))
		} else {
			targetWord.Topics = targetTopics
		}

		relation.TargetWord = &targetWord
		relations = append(relations, &relation)
	}

	if err := rows.Err(); err != nil {
		return []*domain.WordRelation{}, err
	}

	if relations == nil {
		return []*domain.WordRelation{}, nil
	}

	return relations, nil
}
