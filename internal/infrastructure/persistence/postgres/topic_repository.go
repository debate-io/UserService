package postgres

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/go-pg/pg/v9"
	"github.com/ztrue/tracerr"
)

var (
	_ repo.TopicRepository = (*TopicRepository)(nil)
)

type TopicRepository struct {
	db *pg.DB
}

func NewTopicRepository(db *pg.DB) *TopicRepository {
	return &TopicRepository{
		db: db,
	}
}

func (t *TopicRepository) SuggestTopic(ctx context.Context, topic model.Topic) (*model.Topic, error) {
	count, err := t.db.ModelContext(ctx, &model.Topic{}).
		Where("name = ?", topic.Name).
		Count()

	if err == nil && count == 0 {
		_, err = t.db.ModelContext(ctx, &topic).Insert()
	}

	if err != nil {
		return nil, tracerr.Errorf("failed suggest topic: %w", err)
	}
	if count > 0 {
		return nil, repo.ErrTopicAlreadyExist
	}
	return &topic, nil
}
