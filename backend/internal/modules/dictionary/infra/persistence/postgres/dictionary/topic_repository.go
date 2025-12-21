package dictionary

import (
	"context"

	"github.com/english-coach/backend/internal/modules/dictionary/domain"
	"github.com/english-coach/backend/internal/shared/errors"
)

// topicRepository implements TopicRepository using sqlc
type topicRepository struct {
	*DictionaryRepository
}

// FindAll returns all topics
func (r *topicRepository) FindAll(ctx context.Context) ([]*domain.Topic, error) {
	rows, err := r.queries.FindAllTopics(ctx)
	if err != nil {
		return nil, errors.MapPgError(err)
	}

	topics := make([]*domain.Topic, 0, len(rows))
	for _, row := range rows {
		topics = append(topics, &domain.Topic{
			ID:   row.ID,
			Code: row.Code,
			Name: row.Name,
		})
	}

	return topics, nil
}

// FindByID returns a topic by ID
func (r *topicRepository) FindByID(ctx context.Context, id int64) (*domain.Topic, error) {
	row, err := r.queries.FindTopicByID(ctx, id)
	if err != nil {
		return nil, errors.MapPgError(err)
	}

	return &domain.Topic{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}, nil
}

// FindByCode returns a topic by code
func (r *topicRepository) FindByCode(ctx context.Context, code string) (*domain.Topic, error) {
	row, err := r.queries.FindTopicByCode(ctx, code)
	if err != nil {
		return nil, errors.MapPgError(err)
	}

	return &domain.Topic{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}, nil
}

