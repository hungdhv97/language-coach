package dictionary

import (
	"context"

	"github.com/english-coach/backend/internal/domain/dictionary/model"
	"github.com/english-coach/backend/internal/infrastructure/repository/common"
)

// topicRepository implements TopicRepository using sqlc
type topicRepository struct {
	*DictionaryRepository
}

// FindAll returns all topics
func (r *topicRepository) FindAll(ctx context.Context) ([]*model.Topic, error) {
	rows, err := r.queries.FindAllTopics(ctx)
	if err != nil {
		return nil, common.MapPgError(err)
	}

	topics := make([]*model.Topic, 0, len(rows))
	for _, row := range rows {
		topics = append(topics, &model.Topic{
			ID:   row.ID,
			Code: row.Code,
			Name: row.Name,
		})
	}

	return topics, nil
}

// FindByID returns a topic by ID
func (r *topicRepository) FindByID(ctx context.Context, id int64) (*model.Topic, error) {
	row, err := r.queries.FindTopicByID(ctx, id)
	if err != nil {
		return nil, common.MapPgError(err)
	}

	return &model.Topic{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}, nil
}

// FindByCode returns a topic by code
func (r *topicRepository) FindByCode(ctx context.Context, code string) (*model.Topic, error) {
	row, err := r.queries.FindTopicByCode(ctx, code)
	if err != nil {
		return nil, common.MapPgError(err)
	}

	return &model.Topic{
		ID:   row.ID,
		Code: row.Code,
		Name: row.Name,
	}, nil
}

