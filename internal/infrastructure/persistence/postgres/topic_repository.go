package postgres

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
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

func (t *TopicRepository) UpdateTopics(ctx context.Context, topicMetatopicIds []model.TopicMetatopicIds) ([]model.TopicMetatopics, error) {
	var updateTopics []model.Topic
	var updateMetatopicsTopics []model.MetatopicsTopics

	for _, v := range topicMetatopicIds {
		updateTopics = append(updateTopics, v.Topic)
		for _, id := range v.MetatopicIds {
			updateMetatopicsTopics = append(updateMetatopicsTopics, model.MetatopicsTopics{
				MetatopicID: id,
				TopicID:     v.Topic.ID,
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

	result, err := t.selectTopicMetatopics(ctx, transactional, t.topicsToTopicIds(updateTopics))
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

func (t *TopicRepository) GetTopics(ctx context.Context, topicStatuses []model.ApprovingStatusEnum, pageSize, pageNumber int) ([]model.TopicMetatopics, int, error) {
	rows, err := t.db.ModelContext(ctx, &model.Topic{}).
		WhereIn("status IN (?)", topicStatuses).
		Count()
	if err != nil {
		return nil, 0, tracerr.New("failed get topics")
	}

	var resultTopics []model.Topic
	err = t.db.ModelContext(ctx, &resultTopics).
		WhereIn("status IN (?)", topicStatuses).
		Offset(pageNumber * pageSize).
		Limit(pageSize).
		Select()
	if err != nil {
		return nil, 0, tracerr.New("failed get topics")
	}

	var result []model.TopicMetatopics
	for _, topic := range resultTopics {
		var metatopics []model.Metatopic
		err = t.db.ModelContext(ctx, &metatopics).
			Join("JOIN metatopics_topics mt ON mt.metatopics_id = metatopic.id").
			Where("mt.topics_id = ?", topic.ID).
			Select()
		if err != nil {
			return nil, 0, tracerr.New("failed get topics")
		}

		result = append(result, model.TopicMetatopics{
			Topic:      topic,
			Metatopics: metatopics,
		})
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

func (t *TopicRepository) updateTopics(ctx context.Context, db *pg.Tx, topics []model.Topic, metatopicTopic []model.MetatopicsTopics) (err error) {
	_, err = db.ModelContext(ctx, &topics).
		Column("name", "status").
		Update()
	if err != nil {
		if isNoRowsError(err) {
			return repo.ErrNotFound
		}
		return tracerr.Errorf("failed update topics: %w", err)
	}

	_, err = db.ModelContext(ctx, &model.MetatopicsTopics{}).
		WhereIn("topics_id IN (?)", t.topicsToTopicIds(topics)).
		Delete()
	if err != nil {
		return tracerr.Errorf("failed update topics: %w", err)
	}

	if len(metatopicTopic) == 0 {
		return
	}

	err = db.Insert(&metatopicTopic)
	if err != nil {
		return tracerr.Errorf("failed update topics: %w", err)
	}

	return
}

func (t *TopicRepository) selectTopicMetatopics(ctx context.Context, transactional *pg.Tx, topicIds []int) ([]model.TopicMetatopics, error) {
	var result []model.TopicMetatopics
	var topics []model.Topic

	var db orm.DB = t.db
	if transactional != nil {
		db = transactional
	}

	err := db.ModelContext(ctx, &topics).
		WhereIn("id IN (?)", topicIds).
		Select()
	if err != nil {
		return nil, tracerr.Errorf("failed select topics: %w", err)
	}

	for _, topic := range topics {
		var metatopics []model.Metatopic
		err = db.ModelContext(ctx, &metatopics).
			Join("JOIN metatopics_topics mt ON mt.metatopics_id = metatopic.id").
			Where("mt.topics_id = ?", topic.ID).
			Select()
		if err != nil {
			return nil, tracerr.Errorf("failed select topics: %w", err)
		}

		result = append(result, model.TopicMetatopics{
			Topic:      topic,
			Metatopics: metatopics,
		})
	}
	return result, err
}

func (t *TopicRepository) topicsToTopicIds(topics []model.Topic) (topicIds []int) {
	for _, topic := range topics {
		topicIds = append(topicIds, topic.ID)
	}
	return
}
