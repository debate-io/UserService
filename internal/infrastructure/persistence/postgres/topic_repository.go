package postgres

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/go-pg/pg/v9"
	"github.com/ztrue/tracerr"
)

var (
	_ repo.TopicRepository = (*TopicRepository)(nil)
)

type TopicRepository struct {
	db *pg.DB
}

func NewTopicRepository(db *pg.DB) *TopicRepository {
	return &TopicRepository{
		db: db,
	}
}

func (t *TopicRepository) SuggestTopic(ctx context.Context, topic model.Topic) (*model.Topic, error) {
	_, err := t.db.ModelContext(ctx, &topic).Insert()
	if err != nil {
		if getConstraint(err) != "" {
			return nil, repo.ErrAlreadyExist
		}
		return nil, tracerr.Errorf("failed suggest topic: %w", err)
	}

	return &topic, nil
}

func (t *TopicRepository) UpdateTopics(ctx context.Context, topicMetatopics map[*model.Topic][]int) (map[*model.Topic][]*model.Metatopic, error) {
	var updateTopics []*model.Topic
	var updateMetatopicsTopics []*model.MetatopicsTopics

	for topic, ids := range topicMetatopics {
		updateTopics = append(updateTopics, topic)
		for _, id := range ids {
			updateMetatopicsTopics = append(updateMetatopicsTopics, &model.MetatopicsTopics{
				MetatopicID: id,
				TopicID:     topic.ID,
			})
		}
	}

	transactional, err := t.db.Begin()
	if err != nil {
		return nil, tracerr.Errorf("failed update topics: %w", err)
	}

	err = t.updateTopics(ctx, transactional, updateTopics, updateMetatopicsTopics)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	result, err := t.selectTopicMetatopicsMap(ctx, t.topicsToTopicIds(updateTopics))
	if err != nil {
		transactional.Rollback()
		return nil, tracerr.Wrap(err)

	}

	err = transactional.Commit()
	if err != nil {
		return nil, tracerr.Errorf("failed update topics: %w", err)
	}
	return result, nil
}

func (t *TopicRepository) GetTopics(ctx context.Context, topicStatuses []model.ApprovingStatusEnum, pageSize, pageNumber int) (map[*model.Topic][]*model.Metatopic, int, error) {
	var resultTopics []*model.Topic

	rows, err := t.db.ModelContext(ctx, &resultTopics).
		WhereIn("status IN (?)", topicStatuses).
		Count()
	if err != nil {
		return nil, 0, tracerr.New("failed get topics")
	}

	err = t.db.ModelContext(ctx, &resultTopics).
		WhereIn("status IN (?)", topicStatuses).
		Offset(pageNumber * pageSize).
		Limit(pageSize).
		Select()
	if err != nil {
		return nil, 0, tracerr.New("failed get topics")
	}

	result := make(map[*model.Topic][]*model.Metatopic)
	for _, topic := range resultTopics {
		var metatopics []*model.Metatopic
		err = t.db.ModelContext(ctx, &metatopics).
			Join("JOIN metatopics_topics mt ON mt.metatopics_id = metatopic.id").
			WhereIn("mt.topic_id IN (?)", topic.ID).
			Select()
		if err != nil {
			return nil, 0, tracerr.New("failed get topics")
		}

		result[topic] = metatopics
	}

	return result, rows, nil
}

func (t *TopicRepository) GetMetatopics(ctx context.Context, pageSize, pageNumber int) ([]*model.Metatopic, int, error) {
	var metatopics []*model.Metatopic

	rows, err := t.db.ModelContext(ctx, &metatopics).
		Count()
	if err != nil {
		return nil, 0, tracerr.New("failed get metatopics")
	}

	err = t.db.ModelContext(ctx, &metatopics).
		Offset(pageNumber * pageSize).
		Limit(pageSize).
		Select()
	if err != nil {
		return nil, 0, tracerr.New("failed get metatopics")
	}

	return metatopics, rows, nil
}

func (t *TopicRepository) updateTopics(ctx context.Context, db *pg.Tx, topics []*model.Topic, metatopicTopic []*model.MetatopicsTopics) (err error) {
	err = db.Update(&topics)
	if err != nil {
		if isNoRowsError(err) {
			return repo.ErrNotFound
		}
		return tracerr.Errorf("failed update topics: %w", err)
	}

	_, err = db.ModelContext(ctx, &model.MetatopicsTopics{}).
		WhereIn("topic_id IN (?)", t.topicsToTopicIds(topics)).
		Delete()
	if err != nil {
		return tracerr.Errorf("failed update topics: %w", err)
	}
	err = db.Insert(metatopicTopic)
	if err != nil {
		return tracerr.Errorf("failed update topics: %w", err)
	}

	return
}

func (t *TopicRepository) selectTopicMetatopicsMap(ctx context.Context, topicIds []int) (map[*model.Topic][]*model.Metatopic, error) {
	var result = make(map[*model.Topic][]*model.Metatopic)
	var topics []*model.Topic

	err := t.db.ModelContext(ctx, &topics).
		WhereIn("id IN (?)", topicIds).
		Select()
	if err != nil {
		return nil, tracerr.Errorf("failed select topics: %w", err)
	}

	for _, topic := range topics {
		var metatopics []*model.Metatopic
		err = t.db.ModelContext(ctx, &metatopics).
			Join("JOIN metatopics_topics mt ON mt.metatopics_id = metatopic.id").
			WhereIn("mt.topic_id IN (?)", topic.ID).
			Select()
		if err != nil {
			return nil, tracerr.Errorf("failed select topics: %w", err)
		}

		result[topic] = metatopics
	}
	return result, err
}

func (t *TopicRepository) topicsToTopicIds(topics []*model.Topic) (topicIds []int) {
	for _, topic := range topics {
		topicIds = append(topicIds, topic.ID)
	}
	return
}
