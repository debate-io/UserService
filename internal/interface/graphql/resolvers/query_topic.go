package resolvers

import (
	"context"
	"errors"
	"math"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
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

func (q *queryResolver) GetTopic(ctx context.Context, input gen.GetTopicInput) (*gen.GetTopicOutput, error) {
	topicMetatopics, err := q.useCases.Topics.GetTopic(ctx, input.ID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return &gen.GetTopicOutput{
				Topic: nil,
				Error: mappers.NewDTOError(gen.ErrorNotFound),
			}, nil
		}
		if errors.Is(err, repo.ErrValidation) {
			return &gen.GetTopicOutput{
				Topic: nil,
				Error: mappers.NewDTOError(gen.ErrorValidation),
			}, nil
		}
		if errors.Is(err, repo.ErrUnauthorized) {
			return &gen.GetTopicOutput{
				Topic: nil,
				Error: mappers.NewDTOError(gen.ErrorUnauthorized),
			}, nil
		}
		return nil, NewResolverError("failed get topics", err)
	}
	return &gen.GetTopicOutput{
		Topic: mappers.MapTopicMetatopicToTopicMetatopicsDTO([]model.TopicMetatopics{*topicMetatopics})[0],
		Error: nil,
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
