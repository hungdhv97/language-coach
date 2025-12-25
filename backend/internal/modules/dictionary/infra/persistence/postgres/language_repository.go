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

// FindAllLanguages returns all languages
func (r *languageRepository) FindAllLanguages(ctx context.Context) ([]*domain.Language, error) {
	rows, err := r.queries.FindAllLanguages(ctx)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindAllLanguages")
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

// FindLanguageByID returns a language by ID
func (r *languageRepository) FindLanguageByID(ctx context.Context, id int16) (*domain.Language, error) {
	row, err := r.queries.FindLanguageByID(ctx, id)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindLanguageByID")
	}

	return &domain.Language{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}, nil
}

// FindLanguageByCode returns a language by code
func (r *languageRepository) FindLanguageByCode(ctx context.Context, code string) (*domain.Language, error) {
	row, err := r.queries.FindLanguageByCode(ctx, code)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindLanguageByCode")
	}

	return &domain.Language{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}, nil
}
