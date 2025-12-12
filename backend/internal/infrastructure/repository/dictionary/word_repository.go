package dictionary

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/english-coach/backend/internal/domain/dictionary/model"
	db "github.com/english-coach/backend/internal/infrastructure/db/sqlc/gen/dictionary"
	"github.com/english-coach/backend/internal/infrastructure/repository/common"
)

// wordRepository implements WordRepository using sqlc
type wordRepository struct {
	*DictionaryRepository
}

// FindByID returns a word by ID
func (r *wordRepository) FindByID(ctx context.Context, id int64) (*model.Word, error) {
	row, err := r.queries.FindWordByID(ctx, id)
	if err != nil {
		return nil, common.MapPgError(err)
	}

	return r.mapWordRow(row), nil
}

// FindByIDs returns multiple words by their IDs
func (r *wordRepository) FindByIDs(ctx context.Context, ids []int64) ([]*model.Word, error) {
	if len(ids) == 0 {
		return []*model.Word{}, nil
	}

	rows, err := r.queries.FindWordsByIDs(ctx, ids)
	if err != nil {
		return nil, common.MapPgError(err)
	}

	words := make([]*model.Word, 0, len(rows))
	for _, row := range rows {
		words = append(words, r.mapWordRow(row))
	}

	return words, nil
}

// FindWordsByTopicAndLanguages finds words filtered by topic and language pair
func (r *wordRepository) FindWordsByTopicAndLanguages(ctx context.Context, topicID int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*model.Word, error) {
	rows, err := r.queries.FindWordsByTopicAndLanguages(ctx, db.FindWordsByTopicAndLanguagesParams{
		TopicID:          topicID,
		SourceLanguageID: sourceLanguageID,
		TargetLanguageID: targetLanguageID,
		Limit:            int32(limit),
	})
	if err != nil {
		return nil, common.MapPgError(err)
	}

	words := make([]*model.Word, 0, len(rows))
	for _, row := range rows {
		words = append(words, r.mapWordRow(row))
	}

	return words, nil
}

// FindWordsByLevelAndLanguages finds words filtered by level and language pair
func (r *wordRepository) FindWordsByLevelAndLanguages(ctx context.Context, levelID int64, sourceLanguageID, targetLanguageID int16, limit int) ([]*model.Word, error) {
	levelIDPg := pgtype.Int8{Int64: levelID, Valid: true}
	rows, err := r.queries.FindWordsByLevelAndLanguages(ctx, db.FindWordsByLevelAndLanguagesParams{
		LevelID:          levelIDPg,
		SourceLanguageID: sourceLanguageID,
		TargetLanguageID: targetLanguageID,
		Limit:            int32(limit),
	})
	if err != nil {
		return nil, common.MapPgError(err)
	}

	words := make([]*model.Word, 0, len(rows))
	for _, row := range rows {
		words = append(words, r.mapWordRow(row))
	}

	return words, nil
}

// FindTranslationsForWord finds translation words for a given source word and target language
func (r *wordRepository) FindTranslationsForWord(ctx context.Context, sourceWordID int64, targetLanguageID int16, limit int) ([]*model.Word, error) {
	rows, err := r.queries.FindTranslationsForWord(ctx, db.FindTranslationsForWordParams{
		SourceWordID:     sourceWordID,
		TargetLanguageID: targetLanguageID,
		Limit:            int32(limit),
	})
	if err != nil {
		return nil, common.MapPgError(err)
	}

	words := make([]*model.Word, 0, len(rows))
	for _, row := range rows {
		words = append(words, r.mapWordRow(row))
	}

	return words, nil
}

// SearchWords searches for words using multiple strategies (lemma, normalized, search_key)
func (r *wordRepository) SearchWords(ctx context.Context, query string, languageID int16, limit, offset int) ([]*model.Word, error) {
	searchPattern := "%" + query + "%"

	wordRows, err := r.queries.SearchWords(ctx, db.SearchWordsParams{
		LanguageID:    languageID,
		SearchPattern: searchPattern,
		ExactMatch:    query,
		Limit:         int32(limit),
		Offset:        int32(offset),
	})

	if err != nil {
		return nil, common.MapPgError(err)
	}

	words := make([]*model.Word, 0, len(wordRows))
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
		return 0, common.MapPgError(err)
	}

	return int(count), nil
}

// mapWordRow maps sqlc generated row to domain model
func (r *wordRepository) mapWordRow(row db.Word) *model.Word {
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

	return &model.Word{
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
