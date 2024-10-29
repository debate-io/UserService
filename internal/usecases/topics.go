package usecases

import (
	"context"
	"errors"

	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
	"github.com/debate-io/service-auth/internal/usecases/mappers"
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
	input gen.SuggestTopicInput,
) (*gen.SuggestTopicOutput, error) {
	topic, err := t.topicRepo.SuggestTopic(ctx, *mappers.MapSuggestInputToTopic(&input))
	if err != nil {
		if errors.Is(err, repo.ErrTopicAlreadyExist) {
			return &gen.SuggestTopicOutput{
				Topic: nil,
				Error: mappers.NewDTOError(gen.ErrorAlreadyExist),
			}, nil
		}
		return nil, err
	}
	return &gen.SuggestTopicOutput{
		Topic: mappers.MapTopicToTopicDTO(topic),
		Error: nil,
	}, nil
}
