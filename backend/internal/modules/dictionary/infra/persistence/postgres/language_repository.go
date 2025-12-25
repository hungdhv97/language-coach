package dictionary

import (
	"context"

	"github.com/english-coach/backend/internal/modules/dictionary/domain"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
)

// languageRepository implements LanguageRepository using sqlc
type languageRepository struct {
	*DictionaryRepository
}

// FindAll returns all languages
func (r *languageRepository) FindAll(ctx context.Context) ([]*domain.Language, error) {
	rows, err := r.queries.FindAllLanguages(ctx)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindAll")
	}

	languages := make([]*domain.Language, 0, len(rows))
	for _, row := range rows {
		languages = append(languages, &domain.Language{
			ID:   row.ID,
			Code: row.Code,
			Name: row.Name,
		})
	}

	return languages, nil
}

// FindByID returns a language by ID
func (r *languageRepository) FindByID(ctx context.Context, id int16) (*domain.Language, error) {
	row, err := r.queries.FindLanguageByID(ctx, id)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindByID")
	}

	return &domain.Language{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}, nil
}

// FindByCode returns a language by code
func (r *languageRepository) FindByCode(ctx context.Context, code string) (*domain.Language, error) {
	row, err := r.queries.FindLanguageByCode(ctx, code)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindByCode")
	}

	return &domain.Language{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}, nil
}
