package dictionary

import (
	"context"

	"github.com/english-coach/backend/internal/domain/dictionary/model"
	"github.com/english-coach/backend/internal/infrastructure/repository/common"
)

// partOfSpeechRepository implements PartOfSpeechRepository using sqlc
type partOfSpeechRepository struct {
	*DictionaryRepository
}

// FindAll returns all parts of speech
func (r *partOfSpeechRepository) FindAll(ctx context.Context) ([]*model.PartOfSpeech, error) {
	rows, err := r.queries.FindAllPartsOfSpeech(ctx)
	if err != nil {
		return nil, common.MapPgError(err)
	}

	partsOfSpeech := make([]*model.PartOfSpeech, 0, len(rows))
	for _, row := range rows {
		partsOfSpeech = append(partsOfSpeech, &model.PartOfSpeech{
			ID:   row.ID,
			Code: row.Code,
			Name: row.Name,
		})
	}

	return partsOfSpeech, nil
}

// FindByID returns a part of speech by ID
func (r *partOfSpeechRepository) FindByID(ctx context.Context, id int16) (*model.PartOfSpeech, error) {
	row, err := r.queries.FindPartOfSpeechByID(ctx, id)
	if err != nil {
		return nil, common.MapPgError(err)
	}

	return &model.PartOfSpeech{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}, nil
}

// FindByCode returns a part of speech by code
func (r *partOfSpeechRepository) FindByCode(ctx context.Context, code string) (*model.PartOfSpeech, error) {
	row, err := r.queries.FindPartOfSpeechByCode(ctx, code)
	if err != nil {
		return nil, common.MapPgError(err)
	}

	return &model.PartOfSpeech{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}, nil
}

// FindByIDs returns parts of speech by their IDs
func (r *partOfSpeechRepository) FindByIDs(ctx context.Context, ids []int16) (map[int16]*model.PartOfSpeech, error) {
	if len(ids) == 0 {
		return make(map[int16]*model.PartOfSpeech), nil
	}

	rows, err := r.queries.FindPartsOfSpeechByIDs(ctx, ids)
	if err != nil {
		return nil, common.MapPgError(err)
	}

	result := make(map[int16]*model.PartOfSpeech)
	for _, row := range rows {
		result[row.ID] = &model.PartOfSpeech{
			ID:   row.ID,
			Code: row.Code,
			Name: row.Name,
		}
	}

	return result, nil
}
