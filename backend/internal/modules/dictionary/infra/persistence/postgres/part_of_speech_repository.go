package dictionary

import (
	"context"

	"github.com/english-coach/backend/internal/modules/dictionary/domain"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
)

// partOfSpeechRepository implements PartOfSpeechRepository using sqlc
type partOfSpeechRepository struct {
	*DictionaryRepository
}

// FindAllPartsOfSpeech returns all parts of speech
func (r *partOfSpeechRepository) FindAllPartsOfSpeech(ctx context.Context) ([]*domain.PartOfSpeech, error) {
	rows, err := r.queries.FindAllPartsOfSpeech(ctx)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindAllPartsOfSpeech")
	}

	partsOfSpeech := make([]*domain.PartOfSpeech, 0, len(rows))
	for _, row := range rows {
		partsOfSpeech = append(partsOfSpeech, &domain.PartOfSpeech{
			ID:   row.ID,
			Code: row.Code,
			Name: row.Name,
		})
	}

	return partsOfSpeech, nil
}

// FindPartOfSpeechByID returns a part of speech by ID
func (r *partOfSpeechRepository) FindPartOfSpeechByID(ctx context.Context, id int16) (*domain.PartOfSpeech, error) {
	row, err := r.queries.FindPartOfSpeechByID(ctx, id)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindPartOfSpeechByID")
	}

	return &domain.PartOfSpeech{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}, nil
}

// FindPartOfSpeechByCode returns a part of speech by code
func (r *partOfSpeechRepository) FindPartOfSpeechByCode(ctx context.Context, code string) (*domain.PartOfSpeech, error) {
	row, err := r.queries.FindPartOfSpeechByCode(ctx, code)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindPartOfSpeechByCode")
	}

	return &domain.PartOfSpeech{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}, nil
}

// FindPartsOfSpeechByIDs returns parts of speech by their IDs
func (r *partOfSpeechRepository) FindPartsOfSpeechByIDs(ctx context.Context, ids []int16) (map[int16]*domain.PartOfSpeech, error) {
	if len(ids) == 0 {
		return make(map[int16]*domain.PartOfSpeech), nil
	}

	rows, err := r.queries.FindPartsOfSpeechByIDs(ctx, ids)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindPartsOfSpeechByIDs")
	}

	result := make(map[int16]*domain.PartOfSpeech)
	for _, row := range rows {
		result[row.ID] = &domain.PartOfSpeech{
			ID:   row.ID,
			Code: row.Code,
			Name: row.Name,
		}
	}

	return result, nil
}
