package resolvers

import (
	"context"
	"errors"

	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
	"github.com/debate-io/service-auth/internal/usecases/mappers"
)

func (m *mutationResolver) SuggestTopic(ctx context.Context, input gen.SuggestTopicInput) (*gen.SuggestTopicOutput, error) {
	output, err := m.useCases.Topics.SuggestTopic(ctx, mappers.MapSuggestInputToTopic(&input))

	if err != nil {
		if errors.Is(err, repo.ErrAlreadyExist) {
			return &gen.SuggestTopicOutput{
				Topic: nil,
				Error: mappers.NewDTOError(gen.ErrorAlreadyExist),
			}, nil
		}
		return nil, NewResolverError("failed to suggest topic", err)
	}

	return &gen.SuggestTopicOutput{
		Topic: mappers.MapTopicToTopicDTO(output),
		Error: nil,
	}, nil
}
