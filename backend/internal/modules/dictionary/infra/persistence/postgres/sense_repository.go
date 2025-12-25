package dictionary

import (
	"context"

	"github.com/english-coach/backend/internal/modules/dictionary/domain"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
)

// senseRepository implements SenseRepository using sqlc
type senseRepository struct {
	*DictionaryRepository
}

// FindSensesByWordID returns all senses for a word, ordered by sense_order
func (r *senseRepository) FindSensesByWordID(ctx context.Context, wordID int64) ([]*domain.Sense, error) {
	rows, err := r.queries.FindSensesByWordID(ctx, wordID)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindSensesByWordID")
	}

	senses := make([]*domain.Sense, 0, len(rows))
	for _, row := range rows {
		var usageLabel, note *string
		var levelID *int64

		if row.UsageLabel.Valid {
			usageLabel = &row.UsageLabel.String
		}
		if row.Note.Valid {
			note = &row.Note.String
		}
		if row.LevelID.Valid {
			val := row.LevelID.Int64
			levelID = &val
		}

		senses = append(senses, &domain.Sense{
			ID:                   row.ID,
			WordID:               row.WordID,
			SenseOrder:           row.SenseOrder,
			PartOfSpeechID:       row.PartOfSpeechID,
			Definition:           row.Definition,
			DefinitionLanguageID: row.DefinitionLanguageID,
			UsageLabel:           usageLabel,
			LevelID:              levelID,
			Note:                 note,
		})
	}

	return senses, nil
}

// FindSensesByWordIDs returns senses for multiple words
func (r *senseRepository) FindSensesByWordIDs(ctx context.Context, wordIDs []int64) (map[int64][]*domain.Sense, error) {
	if len(wordIDs) == 0 {
		return make(map[int64][]*domain.Sense), nil
	}

	rows, err := r.queries.FindSensesByWordIDs(ctx, wordIDs)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindSensesByWordIDs")
	}

	result := make(map[int64][]*domain.Sense)
	for _, row := range rows {
		var usageLabel, note *string
		var levelID *int64

		if row.UsageLabel.Valid {
			usageLabel = &row.UsageLabel.String
		}
		if row.Note.Valid {
			note = &row.Note.String
		}
		if row.LevelID.Valid {
			val := row.LevelID.Int64
			levelID = &val
		}

		sense := &domain.Sense{
			ID:                   row.ID,
			WordID:               row.WordID,
			SenseOrder:           row.SenseOrder,
			PartOfSpeechID:       row.PartOfSpeechID,
			Definition:           row.Definition,
			DefinitionLanguageID: row.DefinitionLanguageID,
			UsageLabel:           usageLabel,
			LevelID:              levelID,
			Note:                 note,
		}
		result[sense.WordID] = append(result[sense.WordID], sense)
	}

	return result, nil
}
