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

func (t *Topic) UpdateTopics(
	ctx context.Context,
	input []model.TopicMetatopicIds,
) ([]model.TopicMetatopics, error) {
	for _, v := range input {
		ids := v.MetatopicIds
		if (v.Topic.Status == model.StatusDeclined && len(ids) != 0) ||
			(v.Topic.Status == model.StatusApproved && len(ids) == 0) {
			return nil, repo.ErrValidation
		}
	}

	return t.topicRepo.UpdateTopics(ctx, input)
}

func (t *Topic) GetTopics(
	ctx context.Context,
	topicStatuses []model.ApprovingStatusEnum,
	pageSize, pageNumber int,
) ([]model.TopicMetatopics, int, error) {
	if len(topicStatuses) == 0 || pageSize <= 0 || pageNumber < 0 {
		return nil, 0, repo.ErrValidation
	}

	return t.topicRepo.GetTopics(ctx, topicStatuses, pageSize, pageNumber)
}

func (t *Topic) GetMetatopics(
	ctx context.Context,
	pageSize, pageNumber int,
) ([]*model.Metatopic, int, error) {
	if pageSize <= 0 || pageNumber < 0 {
		return nil, 0, repo.ErrValidation
	}

	return t.topicRepo.GetMetatopics(ctx, pageSize, pageNumber)
}
