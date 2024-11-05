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

func MapUpdateTopicInputToTopicMetatopicIds(input *gen.UpdateTopicInput) map[*model.Topic][]int {
	output := make((map[*model.Topic][]int))
	for _, topicInput := range input.Topics {
		topic := &model.Topic{
			ID:     topicInput.ID,
			Name:   cleanString(topicInput.Name),
			Status: MapTopicStatusesToApprovingStatus(topicInput.Status)[0],
		}
		output[topic] = topicInput.MetatopicIds
	}

	return output
}

func MapTopicMetatopicToTopicMetatopicsDTO(topicMetatopics map[*model.Topic][]*model.Metatopic) (output []*gen.TopicMetatopics) {
	for topic, metatopics := range topicMetatopics {
		var metatopicsDto []*gen.Metatopic
		for _, metatopic := range metatopics {
			metatopicsDto = append(metatopicsDto, MapMetatopicToMetatopicDTO(metatopic))
		}

		output = append(output, &gen.TopicMetatopics{
			Topic:      MapTopicToTopicDTO(topic),
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
