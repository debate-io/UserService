package usecases

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
)

type Topic struct {
	topicRepo repo.TopicRepository
}

func NewTopicUseCase(topicRepo repo.TopicRepository) *Topic {
	return &Topic{
		topicRepo: topicRepo,
	}
}

func (t *Topic) SuggestTopic(
	ctx context.Context,
	input *model.Topic,
) (*model.Topic, error) {
	return t.topicRepo.SuggestTopic(ctx, *input)
}
