package dictionary

import (
	"context"

	"github.com/english-coach/backend/internal/domain/dictionary/model"
	"github.com/english-coach/backend/internal/infrastructure/repository/common"
)

// languageRepository implements LanguageRepository using sqlc
type languageRepository struct {
	*DictionaryRepository
}

// FindAll returns all languages
func (r *languageRepository) FindAll(ctx context.Context) ([]*model.Language, error) {
	rows, err := r.queries.FindAllLanguages(ctx)
	if err != nil {
		return nil, common.MapPgError(err)
	}

	languages := make([]*model.Language, 0, len(rows))
	for _, row := range rows {
		languages = append(languages, &model.Language{
			ID:   row.ID,
			Code: row.Code,
			Name: row.Name,
		})
	}

	return languages, nil
}

// FindByID returns a language by ID
func (r *languageRepository) FindByID(ctx context.Context, id int16) (*model.Language, error) {
	row, err := r.queries.FindLanguageByID(ctx, id)
	if err != nil {
		return nil, common.MapPgError(err)
	}

	return &model.Language{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}, nil
}

// FindByCode returns a language by code
func (r *languageRepository) FindByCode(ctx context.Context, code string) (*model.Language, error) {
	row, err := r.queries.FindLanguageByCode(ctx, code)
	if err != nil {
		return nil, common.MapPgError(err)
	}

	return &model.Language{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}, nil
}
