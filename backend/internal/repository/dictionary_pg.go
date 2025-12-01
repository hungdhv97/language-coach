package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/english-coach/backend/internal/domain/dictionary/model"
	"github.com/english-coach/backend/internal/domain/dictionary/port"
)

// DictionaryRepository implements dictionary repository interfaces
type DictionaryRepository struct {
	pool *pgxpool.Pool
}

// NewDictionaryRepository creates a new dictionary repository
func NewDictionaryRepository(pool *pgxpool.Pool) *DictionaryRepository {
	return &DictionaryRepository{
		pool: pool,
	}
}

// LanguageRepository returns a LanguageRepository implementation
func (r *DictionaryRepository) LanguageRepository() port.LanguageRepository {
	return &languageRepository{r}
}

// TopicRepository returns a TopicRepository implementation
func (r *DictionaryRepository) TopicRepository() port.TopicRepository {
	return &topicRepository{r}
}

// LevelRepository returns a LevelRepository implementation
func (r *DictionaryRepository) LevelRepository() port.LevelRepository {
	return &levelRepository{r}
}

// languageRepository implements LanguageRepository
type languageRepository struct {
	*DictionaryRepository
}

// FindAll returns all languages
func (r *languageRepository) FindAll(ctx context.Context) ([]*model.Language, error) {
	query := `SELECT id, code, name, native_name FROM languages ORDER BY code`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var languages []*model.Language
	for rows.Next() {
		var lang model.Language
		var nativeName *string
		if err := rows.Scan(&lang.ID, &lang.Code, &lang.Name, &nativeName); err != nil {
			return nil, err
		}
		lang.NativeName = nativeName
		languages = append(languages, &lang)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return languages, nil
}

// FindByID returns a language by ID
func (r *languageRepository) FindByID(ctx context.Context, id int16) (*model.Language, error) {
	query := `SELECT id, code, name, native_name FROM languages WHERE id = $1`
	var lang model.Language
	var nativeName *string
	err := r.pool.QueryRow(ctx, query, id).Scan(&lang.ID, &lang.Code, &lang.Name, &nativeName)
	if err != nil {
		return nil, err
	}
	lang.NativeName = nativeName
	return &lang, nil
}

// FindByCode returns a language by code
func (r *languageRepository) FindByCode(ctx context.Context, code string) (*model.Language, error) {
	query := `SELECT id, code, name, native_name FROM languages WHERE code = $1`
	var lang model.Language
	var nativeName *string
	err := r.pool.QueryRow(ctx, query, code).Scan(&lang.ID, &lang.Code, &lang.Name, &nativeName)
	if err != nil {
		return nil, err
	}
	lang.NativeName = nativeName
	return &lang, nil
}

// topicRepository implements TopicRepository
type topicRepository struct {
	*DictionaryRepository
}

// FindAll returns all topics
func (r *topicRepository) FindAll(ctx context.Context) ([]*model.Topic, error) {
	query := `SELECT id, code, name FROM topics ORDER BY code`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []*model.Topic
	for rows.Next() {
		var topic model.Topic
		if err := rows.Scan(&topic.ID, &topic.Code, &topic.Name); err != nil {
			return nil, err
		}
		topics = append(topics, &topic)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return topics, nil
}

// FindByID returns a topic by ID
func (r *topicRepository) FindByID(ctx context.Context, id int64) (*model.Topic, error) {
	query := `SELECT id, code, name FROM topics WHERE id = $1`
	var topic model.Topic
	err := r.pool.QueryRow(ctx, query, id).Scan(&topic.ID, &topic.Code, &topic.Name)
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

// FindByCode returns a topic by code
func (r *topicRepository) FindByCode(ctx context.Context, code string) (*model.Topic, error) {
	query := `SELECT id, code, name FROM topics WHERE code = $1`
	var topic model.Topic
	err := r.pool.QueryRow(ctx, query, code).Scan(&topic.ID, &topic.Code, &topic.Name)
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

// levelRepository implements LevelRepository
type levelRepository struct {
	*DictionaryRepository
}

// FindAll returns all levels
func (r *levelRepository) FindAll(ctx context.Context) ([]*model.Level, error) {
	query := `SELECT id, code, name, description, language_id, difficulty_order FROM levels ORDER BY language_id, difficulty_order NULLS LAST, code`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []*model.Level
	for rows.Next() {
		var level model.Level
		var description *string
		var languageID *int16
		var difficultyOrder *int16
		if err := rows.Scan(&level.ID, &level.Code, &level.Name, &description, &languageID, &difficultyOrder); err != nil {
			return nil, err
		}
		level.Description = description
		level.LanguageID = languageID
		level.DifficultyOrder = difficultyOrder
		levels = append(levels, &level)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return levels, nil
}

// FindByID returns a level by ID
func (r *levelRepository) FindByID(ctx context.Context, id int64) (*model.Level, error) {
	query := `SELECT id, code, name, description, language_id, difficulty_order FROM levels WHERE id = $1`
	var level model.Level
	var description *string
	var languageID *int16
	var difficultyOrder *int16
	err := r.pool.QueryRow(ctx, query, id).Scan(&level.ID, &level.Code, &level.Name, &description, &languageID, &difficultyOrder)
	if err != nil {
		return nil, err
	}
	level.Description = description
	level.LanguageID = languageID
	level.DifficultyOrder = difficultyOrder
	return &level, nil
}

// FindByCode returns a level by code
func (r *levelRepository) FindByCode(ctx context.Context, code string) (*model.Level, error) {
	query := `SELECT id, code, name, description, language_id, difficulty_order FROM levels WHERE code = $1`
	var level model.Level
	var description *string
	var languageID *int16
	var difficultyOrder *int16
	err := r.pool.QueryRow(ctx, query, code).Scan(&level.ID, &level.Code, &level.Name, &description, &languageID, &difficultyOrder)
	if err != nil {
		return nil, err
	}
	level.Description = description
	level.LanguageID = languageID
	level.DifficultyOrder = difficultyOrder
	return &level, nil
}

// FindByLanguageID returns all levels for a specific language
func (r *levelRepository) FindByLanguageID(ctx context.Context, languageID int16) ([]*model.Level, error) {
	query := `SELECT id, code, name, description, language_id, difficulty_order FROM levels WHERE language_id = $1 ORDER BY difficulty_order, code`
	rows, err := r.pool.Query(ctx, query, languageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []*model.Level
	for rows.Next() {
		var level model.Level
		var description *string
		var langID *int16
		var difficultyOrder *int16
		if err := rows.Scan(&level.ID, &level.Code, &level.Name, &description, &langID, &difficultyOrder); err != nil {
			return nil, err
		}
		level.Description = description
		level.LanguageID = langID
		level.DifficultyOrder = difficultyOrder
		levels = append(levels, &level)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return levels, nil
}

// WordRepository returns a WordRepository implementation
func (r *DictionaryRepository) WordRepository() port.WordRepository {
	return &wordRepository{r}
}

// wordRepository implements WordRepository
type wordRepository struct {
	*DictionaryRepository
}

// FindByID returns a word by ID
func (r *wordRepository) FindByID(ctx context.Context, id int64) (*model.Word, error) {
	query := `
		SELECT id, language_id, lemma, lemma_normalized, search_key,
		       part_of_speech_id, romanization, script_code, frequency_rank,
		       notes, created_at, updated_at
		FROM words
		WHERE id = $1
	`
	var word model.Word
	var lemmaNormalized, searchKey, romanization, scriptCode, notes *string
	var partOfSpeechID *int16
	var frequencyRank *int
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&word.ID,
		&word.LanguageID,
		&word.Lemma,
		&lemmaNormalized,
		&searchKey,
		&partOfSpeechID,
		&romanization,
		&scriptCode,
		&frequencyRank,
		&notes,
		&word.CreatedAt,
		&word.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	word.LemmaNormalized = lemmaNormalized
	word.SearchKey = searchKey
	word.PartOfSpeechID = partOfSpeechID
	word.Romanization = romanization
	word.ScriptCode = scriptCode
	word.FrequencyRank = frequencyRank
	word.Notes = notes
	return &word, nil
}

// FindByIDs returns multiple words by their IDs
func (r *wordRepository) FindByIDs(ctx context.Context, ids []int64) ([]*model.Word, error) {
	if len(ids) == 0 {
		return []*model.Word{}, nil
	}

	query := `
		SELECT id, language_id, lemma, lemma_normalized, search_key,
		       part_of_speech_id, romanization, script_code, frequency_rank,
		       notes, created_at, updated_at
		FROM words
		WHERE id = ANY($1)
		ORDER BY id
	`
	rows, err := r.pool.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*model.Word
	for rows.Next() {
		var word model.Word
		var lemmaNormalized, searchKey, romanization, scriptCode, notes *string
		var partOfSpeechID *int16
		var frequencyRank *int
		if err := rows.Scan(
			&word.ID,
			&word.LanguageID,
			&word.Lemma,
			&lemmaNormalized,
			&searchKey,
			&partOfSpeechID,
			&romanization,
			&scriptCode,
			&frequencyRank,
			&notes,
			&word.CreatedAt,
			&word.UpdatedAt,
		); err != nil {
			return nil, err
		}
		word.LemmaNormalized = lemmaNormalized
		word.SearchKey = searchKey
		word.PartOfSpeechID = partOfSpeechID
		word.Romanization = romanization
		word.ScriptCode = scriptCode
		word.FrequencyRank = frequencyRank
		word.Notes = notes
		words = append(words, &word)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// FindWordsByTopicAndLanguages finds words filtered by topic and language pair
func (r *wordRepository) FindWordsByTopicAndLanguages(ctx context.Context, topicID int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*model.Word, error) {
	query := `
		SELECT DISTINCT w.id, w.language_id, w.lemma, w.lemma_normalized, w.search_key,
		       w.part_of_speech_id, w.romanization, w.script_code, w.frequency_rank,
		       w.notes, w.created_at, w.updated_at
		FROM words w
		INNER JOIN word_topics wt ON w.id = wt.word_id
		WHERE wt.topic_id = $1
		  AND w.language_id = $2
		  AND EXISTS (
		      SELECT 1
		      FROM senses s
		      INNER JOIN sense_translations st ON s.id = st.source_sense_id
		      INNER JOIN words tw ON st.target_word_id = tw.id
		      WHERE s.word_id = w.id
		        AND tw.language_id = $3
		  )
		ORDER BY w.frequency_rank NULLS LAST, w.id
		LIMIT $4
	`
	rows, err := r.pool.Query(ctx, query, topicID, sourceLanguageID, targetLanguageID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*model.Word
	for rows.Next() {
		var word model.Word
		var lemmaNormalized, searchKey, romanization, scriptCode, notes *string
		var partOfSpeechID *int16
		var frequencyRank *int
		if err := rows.Scan(
			&word.ID,
			&word.LanguageID,
			&word.Lemma,
			&lemmaNormalized,
			&searchKey,
			&partOfSpeechID,
			&romanization,
			&scriptCode,
			&frequencyRank,
			&notes,
			&word.CreatedAt,
			&word.UpdatedAt,
		); err != nil {
			return nil, err
		}
		word.LemmaNormalized = lemmaNormalized
		word.SearchKey = searchKey
		word.PartOfSpeechID = partOfSpeechID
		word.Romanization = romanization
		word.ScriptCode = scriptCode
		word.FrequencyRank = frequencyRank
		word.Notes = notes
		words = append(words, &word)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// FindWordsByLevelAndLanguages finds words filtered by level and language pair
func (r *wordRepository) FindWordsByLevelAndLanguages(ctx context.Context, levelID int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*model.Word, error) {
	query := `
		SELECT DISTINCT w.id, w.language_id, w.lemma, w.lemma_normalized, w.search_key,
		       w.part_of_speech_id, w.romanization, w.script_code, w.frequency_rank,
		       w.notes, w.created_at, w.updated_at
		FROM words w
		INNER JOIN senses s ON w.id = s.word_id
		WHERE s.level_id = $1
		  AND w.language_id = $2
		  AND EXISTS (
		      SELECT 1
		      FROM sense_translations st
		      INNER JOIN words tw ON st.target_word_id = tw.id
		      WHERE st.source_sense_id = s.id
		        AND tw.language_id = $3
		  )
		ORDER BY w.frequency_rank NULLS LAST, w.id
		LIMIT $4
	`
	rows, err := r.pool.Query(ctx, query, levelID, sourceLanguageID, targetLanguageID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*model.Word
	for rows.Next() {
		var word model.Word
		var lemmaNormalized, searchKey, romanization, scriptCode, notes *string
		var partOfSpeechID *int16
		var frequencyRank *int
		if err := rows.Scan(
			&word.ID,
			&word.LanguageID,
			&word.Lemma,
			&lemmaNormalized,
			&searchKey,
			&partOfSpeechID,
			&romanization,
			&scriptCode,
			&frequencyRank,
			&notes,
			&word.CreatedAt,
			&word.UpdatedAt,
		); err != nil {
			return nil, err
		}
		word.LemmaNormalized = lemmaNormalized
		word.SearchKey = searchKey
		word.PartOfSpeechID = partOfSpeechID
		word.Romanization = romanization
		word.ScriptCode = scriptCode
		word.FrequencyRank = frequencyRank
		word.Notes = notes
		words = append(words, &word)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// FindTranslationsForWord finds translation words for a given source word and target language
func (r *wordRepository) FindTranslationsForWord(ctx context.Context, sourceWordID int64, targetLanguageID int16, limit int) ([]*model.Word, error) {
	query := `
		SELECT DISTINCT tw.id, tw.language_id, tw.lemma, tw.lemma_normalized, tw.search_key,
		       tw.part_of_speech_id, tw.romanization, tw.script_code, tw.frequency_rank,
		       tw.notes, tw.created_at, tw.updated_at
		FROM words sw
		INNER JOIN senses s ON sw.id = s.word_id
		INNER JOIN sense_translations st ON s.id = st.source_sense_id
		INNER JOIN words tw ON st.target_word_id = tw.id
		WHERE sw.id = $1
		  AND tw.language_id = $2
		ORDER BY st.priority, tw.frequency_rank NULLS LAST, tw.id
		LIMIT $3
	`
	rows, err := r.pool.Query(ctx, query, sourceWordID, targetLanguageID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*model.Word
	for rows.Next() {
		var word model.Word
		var lemmaNormalized, searchKey, romanization, scriptCode, notes *string
		var partOfSpeechID *int16
		var frequencyRank *int
		if err := rows.Scan(
			&word.ID,
			&word.LanguageID,
			&word.Lemma,
			&lemmaNormalized,
			&searchKey,
			&partOfSpeechID,
			&romanization,
			&scriptCode,
			&frequencyRank,
			&notes,
			&word.CreatedAt,
			&word.UpdatedAt,
		); err != nil {
			return nil, err
		}
		word.LemmaNormalized = lemmaNormalized
		word.SearchKey = searchKey
		word.PartOfSpeechID = partOfSpeechID
		word.Romanization = romanization
		word.ScriptCode = scriptCode
		word.FrequencyRank = frequencyRank
		word.Notes = notes
		words = append(words, &word)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// SearchWords searches for words using multiple strategies (lemma, normalized, search_key)
func (r *wordRepository) SearchWords(ctx context.Context, query string, languageID *int16, limit, offset int) ([]*model.Word, error) {
	// Build query with multiple search strategies
	// Try exact match on lemma first, then normalized, then search_key, then prefix match
	searchPattern := "%" + query + "%"
	
	var sqlQuery string
	var args []interface{}
	
	if languageID != nil {
		sqlQuery = `
			SELECT DISTINCT w.id, w.language_id, w.lemma, w.lemma_normalized, w.search_key,
			       w.part_of_speech_id, w.romanization, w.script_code, w.frequency_rank,
			       w.notes, w.created_at, w.updated_at
			FROM words w
			WHERE w.language_id = $1
			  AND (
			    w.lemma ILIKE $2
			    OR w.lemma_normalized ILIKE $2
			    OR w.search_key ILIKE $2
			  )
			ORDER BY 
			  CASE 
			    WHEN w.lemma = $3 THEN 1
			    WHEN w.lemma ILIKE $4 THEN 2
			    WHEN w.lemma_normalized ILIKE $4 THEN 3
			    WHEN w.search_key ILIKE $4 THEN 4
			    ELSE 5
			  END,
			  w.frequency_rank NULLS LAST,
			  w.id
			LIMIT $5 OFFSET $6
		`
		exactMatch := query
		args = []interface{}{*languageID, searchPattern, exactMatch, searchPattern, limit, offset}
	} else {
		sqlQuery = `
			SELECT DISTINCT w.id, w.language_id, w.lemma, w.lemma_normalized, w.search_key,
			       w.part_of_speech_id, w.romanization, w.script_code, w.frequency_rank,
			       w.notes, w.created_at, w.updated_at
			FROM words w
			WHERE (
			  w.lemma ILIKE $1
			  OR w.lemma_normalized ILIKE $1
			  OR w.search_key ILIKE $1
			)
			ORDER BY 
			  CASE 
			    WHEN w.lemma = $2 THEN 1
			    WHEN w.lemma ILIKE $3 THEN 2
			    WHEN w.lemma_normalized ILIKE $3 THEN 3
			    WHEN w.search_key ILIKE $3 THEN 4
			    ELSE 5
			  END,
			  w.frequency_rank NULLS LAST,
			  w.id
			LIMIT $4 OFFSET $5
		`
		exactMatch := query
		args = []interface{}{searchPattern, exactMatch, searchPattern, limit, offset}
	}
	
	rows, err := r.pool.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*model.Word
	for rows.Next() {
		var word model.Word
		var lemmaNormalized, searchKey, romanization, scriptCode, notes *string
		var partOfSpeechID *int16
		var frequencyRank *int
		if err := rows.Scan(
			&word.ID,
			&word.LanguageID,
			&word.Lemma,
			&lemmaNormalized,
			&searchKey,
			&partOfSpeechID,
			&romanization,
			&scriptCode,
			&frequencyRank,
			&notes,
			&word.CreatedAt,
			&word.UpdatedAt,
		); err != nil {
			return nil, err
		}
		word.LemmaNormalized = lemmaNormalized
		word.SearchKey = searchKey
		word.PartOfSpeechID = partOfSpeechID
		word.Romanization = romanization
		word.ScriptCode = scriptCode
		word.FrequencyRank = frequencyRank
		word.Notes = notes
		words = append(words, &word)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// CountSearchWords returns the total count of words matching the search query
func (r *wordRepository) CountSearchWords(ctx context.Context, query string, languageID *int16) (int, error) {
	searchPattern := "%" + query + "%"
	
	var sqlQuery string
	var args []interface{}
	
	if languageID != nil {
		sqlQuery = `
			SELECT COUNT(DISTINCT w.id)
			FROM words w
			WHERE w.language_id = $1
			  AND (
			    w.lemma ILIKE $2
			    OR w.lemma_normalized ILIKE $2
			    OR w.search_key ILIKE $2
			  )
		`
		args = []interface{}{*languageID, searchPattern}
	} else {
		sqlQuery = `
			SELECT COUNT(DISTINCT w.id)
			FROM words w
			WHERE (
			  w.lemma ILIKE $1
			  OR w.lemma_normalized ILIKE $1
			  OR w.search_key ILIKE $1
			)
		`
		args = []interface{}{searchPattern}
	}
	
	var count int
	err := r.pool.QueryRow(ctx, sqlQuery, args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

// SenseRepository returns a SenseRepository implementation
func (r *DictionaryRepository) SenseRepository() port.SenseRepository {
	return &senseRepository{r}
}

// senseRepository implements SenseRepository
type senseRepository struct {
	*DictionaryRepository
}

// FindByWordID returns all senses for a word, ordered by sense_order
func (r *senseRepository) FindByWordID(ctx context.Context, wordID int64) ([]*model.Sense, error) {
	query := `
		SELECT id, word_id, sense_order, definition, definition_language_id,
		       usage_label, level_id, note
		FROM senses
		WHERE word_id = $1
		ORDER BY sense_order
	`
	rows, err := r.pool.Query(ctx, query, wordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var senses []*model.Sense
	for rows.Next() {
		var sense model.Sense
		var usageLabel, note *string
		var levelID *int64
		if err := rows.Scan(
			&sense.ID,
			&sense.WordID,
			&sense.SenseOrder,
			&sense.Definition,
			&sense.DefinitionLanguageID,
			&usageLabel,
			&levelID,
			&note,
		); err != nil {
			return nil, err
		}
		sense.UsageLabel = usageLabel
		sense.LevelID = levelID
		sense.Note = note
		senses = append(senses, &sense)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return senses, nil
}

// FindByWordIDs returns senses for multiple words
func (r *senseRepository) FindByWordIDs(ctx context.Context, wordIDs []int64) (map[int64][]*model.Sense, error) {
	if len(wordIDs) == 0 {
		return make(map[int64][]*model.Sense), nil
	}

	query := `
		SELECT id, word_id, sense_order, definition, definition_language_id,
		       usage_label, level_id, note
		FROM senses
		WHERE word_id = ANY($1)
		ORDER BY word_id, sense_order
	`
	rows, err := r.pool.Query(ctx, query, wordIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]*model.Sense)
	for rows.Next() {
		var sense model.Sense
		var usageLabel, note *string
		var levelID *int64
		if err := rows.Scan(
			&sense.ID,
			&sense.WordID,
			&sense.SenseOrder,
			&sense.Definition,
			&sense.DefinitionLanguageID,
			&usageLabel,
			&levelID,
			&note,
		); err != nil {
			return nil, err
		}
		sense.UsageLabel = usageLabel
		sense.LevelID = levelID
		sense.Note = note
		result[sense.WordID] = append(result[sense.WordID], &sense)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// Ensure implementations satisfy interfaces
var (
	_ port.LanguageRepository = (*languageRepository)(nil)
	_ port.TopicRepository    = (*topicRepository)(nil)
	_ port.LevelRepository    = (*levelRepository)(nil)
	_ port.WordRepository     = (*wordRepository)(nil)
	_ port.SenseRepository    = (*senseRepository)(nil)
)
