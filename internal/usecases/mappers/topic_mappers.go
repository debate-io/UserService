package mappers

import (
	"regexp"
	"strings"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
	"github.com/ztrue/tracerr"
)

func MapSuggestInputToTopic(input *gen.SuggestTopicInput) *model.Topic {
	return &model.Topic{
		Name:   cleanString(input.Name),
		Status: model.StatusPending,
	}
}

func MapUpdateTopicInputToTopicMetatopicIds(input *gen.UpdateTopicInput) (output []model.TopicMetatopicIds) {
	for _, topicInput := range input.Topics {
		output = append(output, model.TopicMetatopicIds{
			Topic: model.Topic{
				ID:     topicInput.ID,
				Name:   cleanString(topicInput.Name),
				Status: MapTopicStatusesToApprovingStatus(topicInput.Status)[0],
			},
			MetatopicIds: topicInput.MetatopicIds,
		})
	}

	return output
}

func MapTopicMetatopicToTopicMetatopicsDTO(topicMetatopics []model.TopicMetatopics) (output []*gen.TopicMetatopics) {
	for _, v := range topicMetatopics {
		var metatopicsDto []*gen.Metatopic
		for _, metatopic := range v.Metatopics {
			metatopicsDto = append(metatopicsDto, MapMetatopicToMetatopicDTO(&metatopic))
		}

		output = append(output, &gen.TopicMetatopics{
			Topic:      MapTopicToTopicDTO(&v.Topic),
			Metatopics: metatopicsDto,
		})
	}
	return
}

func MapTopicStatusesToApprovingStatus(topicStatuses ...gen.TopicStatus) (output []model.ApprovingStatusEnum) {
	for _, status := range topicStatuses {
		switch status {
		case gen.TopicStatusApproved:
			output = append(output, model.StatusApproved)
		case gen.TopicStatusDeclined:
			output = append(output, model.StatusDeclined)
		case gen.TopicStatusPending:
			output = append(output, model.StatusPending)
		default:
			panic(tracerr.New("error enum mapping"))
		}
	}
	return
}

func MapTopicToTopicDTO(topic *model.Topic) *gen.Topic {
	var genStatus gen.TopicStatus
	switch topic.Status {
	case model.StatusApproved:
		genStatus = gen.TopicStatusApproved
	case model.StatusPending:
		genStatus = gen.TopicStatusPending
	case model.StatusDeclined:
		genStatus = gen.TopicStatusDeclined
	default:
		panic(tracerr.New("error enum mapping"))
	}

	return &gen.Topic{
		ID:        int(topic.ID),
		Name:      topic.Name,
		Status:    genStatus,
		CreatedAt: topic.CreatedAt,
	}
}

func MapMetatopicToMetatopicDTO(metatopic *model.Metatopic) *gen.Metatopic {
	return &gen.Metatopic{
		ID:        metatopic.ID,
		Name:      metatopic.Name,
		CreatedAt: metatopic.CreatedAt,
	}
}

func cleanString(input string) string {
	lowerString := strings.ToLower(input)
	cleanSpaceString := strings.TrimSpace(lowerString)

	re := regexp.MustCompile(`[\s]+`)
	cleanString := re.ReplaceAllString(cleanSpaceString, " ")

	return cleanString
}
