package dictionary

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/english-coach/backend/internal/modules/dictionary/domain"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
)

// levelRepository implements LevelRepository using sqlc
type levelRepository struct {
	*DictionaryRepository
}

// FindAllLevels returns all levels
func (r *levelRepository) FindAllLevels(ctx context.Context) ([]*domain.Level, error) {
	rows, err := r.queries.FindAllLevels(ctx)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindAllLevels")
	}

	levels := make([]*domain.Level, 0, len(rows))
	for _, row := range rows {
		var description *string
		var languageID *int16
		var difficultyOrder *int16

		if row.Description.Valid {
			description = &row.Description.String
		}
		if row.LanguageID.Valid {
			val := int16(row.LanguageID.Int16)
			languageID = &val
		}
		if row.DifficultyOrder.Valid {
			val := int16(row.DifficultyOrder.Int16)
			difficultyOrder = &val
		}

		levels = append(levels, &domain.Level{
			ID:              row.ID,
			Code:            row.Code,
			Name:            row.Name,
			Description:     description,
			LanguageID:      languageID,
			DifficultyOrder: difficultyOrder,
		})
	}

	return levels, nil
}

// FindLevelByID returns a level by ID
func (r *levelRepository) FindLevelByID(ctx context.Context, id int64) (*domain.Level, error) {
	row, err := r.queries.FindLevelByID(ctx, id)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindLevelByID")
	}

	var description *string
	var languageID *int16
	var difficultyOrder *int16

	if row.Description.Valid {
		description = &row.Description.String
	}
	if row.LanguageID.Valid {
		val := int16(row.LanguageID.Int16)
		languageID = &val
	}
	if row.DifficultyOrder.Valid {
		val := int16(row.DifficultyOrder.Int16)
		difficultyOrder = &val
	}

	return &domain.Level{
		ID:              row.ID,
		Code:            row.Code,
		Name:            row.Name,
		Description:     description,
		LanguageID:      languageID,
		DifficultyOrder: difficultyOrder,
	}, nil
}

// FindLevelByCode returns a level by code
func (r *levelRepository) FindLevelByCode(ctx context.Context, code string) (*domain.Level, error) {
	row, err := r.queries.FindLevelByCode(ctx, code)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindLevelByCode")
	}

	var description *string
	var languageID *int16
	var difficultyOrder *int16

	if row.Description.Valid {
		description = &row.Description.String
	}
	if row.LanguageID.Valid {
		val := int16(row.LanguageID.Int16)
		languageID = &val
	}
	if row.DifficultyOrder.Valid {
		val := int16(row.DifficultyOrder.Int16)
		difficultyOrder = &val
	}

	return &domain.Level{
		ID:              row.ID,
		Code:            row.Code,
		Name:            row.Name,
		Description:     description,
		LanguageID:      languageID,
		DifficultyOrder: difficultyOrder,
	}, nil
}

// FindLevelsByLanguageID returns all levels for a specific language
func (r *levelRepository) FindLevelsByLanguageID(ctx context.Context, languageID int16) ([]*domain.Level, error) {
	langIDPg := pgtype.Int2{Int16: languageID, Valid: true}
	rows, err := r.queries.FindLevelsByLanguageID(ctx, langIDPg)
	if err != nil {
		return nil, sharederrors.MapDictionaryRepositoryError(err, "FindLevelsByLanguageID")
	}

	levels := make([]*domain.Level, 0, len(rows))
	for _, row := range rows {
		var description *string
		var languageID *int16
		var difficultyOrder *int16

		if row.Description.Valid {
			description = &row.Description.String
		}
		if row.LanguageID.Valid {
			val := int16(row.LanguageID.Int16)
			languageID = &val
		}
		if row.DifficultyOrder.Valid {
			val := int16(row.DifficultyOrder.Int16)
			difficultyOrder = &val
		}

		levels = append(levels, &domain.Level{
			ID:              row.ID,
			Code:            row.Code,
			Name:            row.Name,
			Description:     description,
			LanguageID:      languageID,
			DifficultyOrder: difficultyOrder,
		})
	}

	return levels, nil
}
