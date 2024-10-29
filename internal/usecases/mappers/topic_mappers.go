package mappers

import (
	"time"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
)

func MapSuggestInputToTopic(input *gen.SuggestTopicInput) *model.Topic {
	return &model.Topic{
		Name:      input.Name,
		Status:    model.StatusPending,
		CreatedAt: time.Now(),
	}
}

func MapTopicToTopicDTO(topic *model.Topic) *gen.Topic {
	return &gen.Topic{
		ID:        int(topic.ID),
		Name:      topic.Name,
		Status:    string(topic.Status),
		CreatedAt: topic.CreatedAt,
	}
}
