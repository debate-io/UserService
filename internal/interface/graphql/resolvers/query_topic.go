package resolvers

import (
	"context"
	"math"

	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
	"github.com/debate-io/service-auth/internal/usecases/mappers"
)

func (q *queryResolver) GetTopics(ctx context.Context, input gen.GetTopicsInput) (*gen.GetTopicsOutput, error) {
	topicMetatopics, rows, err := q.useCases.Topics.GetTopics(ctx, mappers.MapTopicStatusesToApprovingStatus(input.TopicStatus...), input.PageSize, input.PageNumber)
	if err != nil {
		return nil, NewResolverError("failed get topics", err)
	}

	return &gen.GetTopicsOutput{
		PageSize:   input.PageSize,
		PageNumber: input.PageNumber,
		PageCount:  int(math.Ceil(float64(rows) / float64(input.PageSize))),
		Topics:     mappers.MapTopicMetatopicToTopicMetatopicsDTO(topicMetatopics),
	}, nil

}

func (q *queryResolver) GetMetatopics(ctx context.Context, input gen.GetMetatopicsInput) (*gen.GetMetatopicsOutput, error) {
	metatopics, rows, err := q.useCases.Topics.GetMetatopics(ctx, input.PageSize, input.PageNumber)
	if err != nil {
		return nil, NewResolverError("failed get metatopics", err)
	}

	var metatopicDtos []*gen.Metatopic
	for _, metatopic := range metatopics {
		metatopicDtos = append(metatopicDtos, mappers.MapMetatopicToMetatopicDTO(metatopic))
	}

	return &gen.GetMetatopicsOutput{
		PageSize:   input.PageSize,
		PageNumber: input.PageNumber,
		PageCount:  int(math.Ceil(float64(rows) / float64(input.PageSize))),
		Metatopics: metatopicDtos,
	}, nil
}
