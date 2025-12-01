package dictionary

import (
	"context"

	"github.com/english-coach/backend/internal/domain/dictionary/model"
	"github.com/english-coach/backend/internal/infrastructure/repository/common"
)

// senseRepository implements SenseRepository using sqlc
type senseRepository struct {
	*DictionaryRepository
}

// FindByWordID returns all senses for a word, ordered by sense_order
func (r *senseRepository) FindByWordID(ctx context.Context, wordID int64) ([]*model.Sense, error) {
	rows, err := r.queries.FindSensesByWordID(ctx, wordID)
	if err != nil {
		return nil, common.MapPgError(err)
	}

	senses := make([]*model.Sense, 0, len(rows))
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

		senses = append(senses, &model.Sense{
			ID:                   row.ID,
			WordID:               row.WordID,
			SenseOrder:           row.SenseOrder,
			Definition:           row.Definition,
			DefinitionLanguageID: row.DefinitionLanguageID,
			UsageLabel:           usageLabel,
			LevelID:              levelID,
			Note:                 note,
		})
	}

	return senses, nil
}

// FindByWordIDs returns senses for multiple words
func (r *senseRepository) FindByWordIDs(ctx context.Context, wordIDs []int64) (map[int64][]*model.Sense, error) {
	if len(wordIDs) == 0 {
		return make(map[int64][]*model.Sense), nil
	}

	rows, err := r.queries.FindSensesByWordIDs(ctx, wordIDs)
	if err != nil {
		return nil, common.MapPgError(err)
	}

	result := make(map[int64][]*model.Sense)
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

		sense := &model.Sense{
			ID:                   row.ID,
			WordID:               row.WordID,
			SenseOrder:           row.SenseOrder,
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
