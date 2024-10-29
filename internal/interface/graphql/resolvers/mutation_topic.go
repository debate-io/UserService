package resolvers

import (
	"context"

	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
)

func (m *mutationResolver) SuggestTopic(ctx context.Context, input gen.SuggestTopicInput) (*gen.SuggestTopicOutput, error) {
	output, err := m.useCases.Topics.SuggestTopic(ctx, input)
	if err != nil {
		return nil, NewResolverError("failed to suggest topic", err)
	}
	return output, nil
}
