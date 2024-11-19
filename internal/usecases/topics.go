package usecases

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/debate-io/service-auth/internal/interface/server/middleware"
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
	claims := ctx.Value(middleware.JwtClaimsKey).(*model.Claims)
	if claims == nil {
		return nil, repo.ErrUnauthorized
	}

	role := claims.Role
	if role != model.RoleAdmin && role != model.RoleContentManager && role != model.RoleDefaultUser {
		return nil, repo.ErrUnauthorized
	}

	return t.topicRepo.SuggestTopic(ctx, *input)
}

func (t *Topic) UpdateTopics(
	ctx context.Context,
	input []model.TopicMetatopicIds,
) ([]model.TopicMetatopics, error) {
	claims := ctx.Value(middleware.JwtClaimsKey).(*model.Claims)
	if claims == nil {
		return nil, repo.ErrUnauthorized
	}

	role := claims.Role
	if role != model.RoleAdmin && role != model.RoleContentManager {
		return nil, repo.ErrUnauthorized
	}

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
	claims := ctx.Value(middleware.JwtClaimsKey).(*model.Claims)
	if claims == nil {
		return nil, 0, repo.ErrUnauthorized
	}

	role := claims.Role
	if role != model.RoleAdmin && role != model.RoleContentManager && role != model.RoleDefaultUser {
		return nil, 0, repo.ErrUnauthorized
	}

	if len(topicStatuses) == 0 || pageSize <= 0 || pageNumber < 0 {
		return nil, 0, repo.ErrValidation
	}

	return t.topicRepo.GetTopics(ctx, topicStatuses, pageSize, pageNumber)
}

func (t *Topic) GetMetatopics(
	ctx context.Context,
	pageSize, pageNumber int,
) ([]*model.Metatopic, int, error) {
	claims := ctx.Value(middleware.JwtClaimsKey).(*model.Claims)
	if claims == nil || claims.Role != model.RoleAdmin && claims.Role != model.RoleContentManager && claims.Role != model.RoleDefaultUser {
		return nil, 0, repo.ErrUnauthorized
	}

	if pageSize <= 0 || pageNumber < 0 {
		return nil, 0, repo.ErrValidation
	}

	return t.topicRepo.GetMetatopics(ctx, pageSize, pageNumber)
}
