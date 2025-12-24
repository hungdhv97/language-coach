package dictionary

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/english-coach/backend/internal/modules/dictionary/domain"
	db "github.com/english-coach/backend/internal/platform/db/sqlc/gen/dictionary"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
)

// wordRepository implements WordRepository using sqlc
type wordRepository struct {
	*DictionaryRepository
}

// FindByID returns a word by ID
func (r *wordRepository) FindByID(ctx context.Context, id int64) (*domain.Word, error) {
	row, err := r.queries.FindWordByID(ctx, id)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindWordByID")
	}

	return r.mapWordRow(row), nil
}

// FindByIDs returns multiple words by their IDs
func (r *wordRepository) FindByIDs(ctx context.Context, ids []int64) ([]*domain.Word, error) {
	if len(ids) == 0 {
		return []*domain.Word{}, nil
	}

	rows, err := r.queries.FindWordsByIDs(ctx, ids)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindByIDs")
	}

	words := make([]*domain.Word, 0, len(rows))
	for _, row := range rows {
		words = append(words, r.mapWordRow(row))
	}

	return words, nil
}

// FindWordsByTopicAndLanguages finds words filtered by topic and language pair
func (r *wordRepository) FindWordsByTopicAndLanguages(ctx context.Context, topicID int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*domain.Word, error) {
	rows, err := r.queries.FindWordsByTopicAndLanguages(ctx, db.FindWordsByTopicAndLanguagesParams{
		TopicID:          topicID,
		SourceLanguageID: sourceLanguageID,
		TargetLanguageID: targetLanguageID,
		Limit:            int32(limit),
	})
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindWordsByTopicAndLanguages")
	}

	words := make([]*domain.Word, 0, len(rows))
	for _, row := range rows {
		words = append(words, r.mapWordRow(row))
	}

	return words, nil
}

// FindWordsByLevelAndLanguages finds words filtered by level and language pair
func (r *wordRepository) FindWordsByLevelAndLanguages(ctx context.Context, levelID int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*domain.Word, error) {
	levelIDPg := pgtype.Int8{Int64: levelID, Valid: true}
	rows, err := r.queries.FindWordsByLevelAndLanguages(ctx, db.FindWordsByLevelAndLanguagesParams{
		LevelID:          levelIDPg,
		SourceLanguageID: sourceLanguageID,
		TargetLanguageID: targetLanguageID,
		Limit:            int32(limit),
	})
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindWordsByLevelAndLanguages")
	}

	words := make([]*domain.Word, 0, len(rows))
	for _, row := range rows {
		words = append(words, r.mapWordRow(row))
	}

	return words, nil
}

// FindWordsByLevelAndTopicsAndLanguages finds words filtered by level, optional topics, and language pair
// If topicIDs is nil or empty, returns all words for the level (no topic filter)
func (r *wordRepository) FindWordsByLevelAndTopicsAndLanguages(ctx context.Context, levelID int64, topicIDs []int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*domain.Word, error) {
	// First, fetch words by level (without topic filter)
	// We'll fetch more than needed to account for topic filtering
	fetchLimit := limit
	if len(topicIDs) > 0 {
		fetchLimit = limit * 3 // Fetch more to have enough after filtering
	}

	words, err := r.FindWordsByLevelAndLanguages(ctx, levelID, sourceLanguageID, targetLanguageID, fetchLimit)
	if err != nil {
		return nil, err
	}

	// If no topic filter, return all words (up to limit)
	if len(topicIDs) == 0 || topicIDs == nil {
		if len(words) > limit {
			return words[:limit], nil
		}
		return words, nil
	}

	// Filter by topics: get word IDs that have any of the specified topics
	// We need to query word_topics table to filter
	wordIDs := make([]int64, 0, len(words))
	for _, word := range words {
		wordIDs = append(wordIDs, word.ID)
	}

	if len(wordIDs) == 0 {
		return []*domain.Word{}, nil
	}

	// Query to find which words have the specified topics
	// Using a raw query since we need to check word_topics
	query := `
		SELECT DISTINCT wt.word_id
		FROM word_topics wt
		WHERE wt.word_id = ANY($1::bigint[])
		  AND wt.topic_id = ANY($2::bigint[])
	`

	rows, err := r.DictionaryRepository.pool.Query(ctx, query, wordIDs, topicIDs)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindWordsByLevelAndTopicsAndLanguages")
	}
	defer rows.Close()

	validWordIDs := make(map[int64]bool)
	for rows.Next() {
		var wordID int64
		if err := rows.Scan(&wordID); err != nil {
			return nil, err
		}
		validWordIDs[wordID] = true
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Filter words to only include those with valid topics
	filteredWords := make([]*domain.Word, 0, limit)
	for _, word := range words {
		if validWordIDs[word.ID] {
			filteredWords = append(filteredWords, word)
			if len(filteredWords) >= limit {
				break
			}
		}
	}

	return filteredWords, nil
}

// FindTranslationsForWord finds translation words for a given source word and target language
func (r *wordRepository) FindTranslationsForWord(ctx context.Context, sourceWordID int64, targetLanguageID int16, limit int) ([]*domain.Word, error) {
	rows, err := r.queries.FindTranslationsForWord(ctx, db.FindTranslationsForWordParams{
		SourceWordID:     sourceWordID,
		TargetLanguageID: targetLanguageID,
		Limit:            int32(limit),
	})
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindTranslationsForWord")
	}

	words := make([]*domain.Word, 0, len(rows))
	for _, row := range rows {
		words = append(words, r.mapWordRow(row))
	}

	return words, nil
}

// SearchWords searches for words using multiple strategies (lemma, normalized, search_key)
func (r *wordRepository) SearchWords(ctx context.Context, query string, languageID int16, limit, offset int) ([]*domain.Word, error) {
	searchPattern := "%" + query + "%"

	wordRows, err := r.queries.SearchWords(ctx, db.SearchWordsParams{
		LanguageID:    languageID,
		SearchPattern: searchPattern,
		ExactMatch:    query,
		Limit:         int32(limit),
		Offset:        int32(offset),
	})

	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "SearchWords")
	}

	words := make([]*domain.Word, 0, len(wordRows))
	for _, row := range wordRows {
		words = append(words, r.mapWordRow(row))
	}

	return words, nil
}

// CountSearchWords returns the total count of words matching the search query
func (r *wordRepository) CountSearchWords(ctx context.Context, query string, languageID int16) (int, error) {
	searchPattern := "%" + query + "%"

	count, err := r.queries.CountSearchWords(ctx, db.CountSearchWordsParams{
		LanguageID:    languageID,
		SearchPattern: searchPattern,
	})

	if err != nil {
		return 0, sharederrors.MapDictionaryRepositoryError(err, "CountSearchWords")
	}

	return int(count), nil
}

// mapWordRow maps sqlc generated row to domain model
func (r *wordRepository) mapWordRow(row db.Word) *domain.Word {
	var lemmaNormalized, searchKey, romanization, scriptCode, note *string
	var frequencyRank *int

	if row.LemmaNormalized.Valid {
		lemmaNormalized = &row.LemmaNormalized.String
	}
	if row.SearchKey.Valid {
		searchKey = &row.SearchKey.String
	}
	if row.Romanization.Valid {
		romanization = &row.Romanization.String
	}
	if row.ScriptCode.Valid {
		scriptCode = &row.ScriptCode.String
	}
	if row.Note.Valid {
		note = &row.Note.String
	}
	if row.FrequencyRank.Valid {
		val := int(row.FrequencyRank.Int32)
		frequencyRank = &val
	}

	return &domain.Word{
		ID:              row.ID,
		LanguageID:      row.LanguageID,
		Lemma:           row.Lemma,
		LemmaNormalized: lemmaNormalized,
		SearchKey:       searchKey,
		Romanization:    romanization,
		ScriptCode:      scriptCode,
		FrequencyRank:   frequencyRank,
		Note:            note,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}
}
