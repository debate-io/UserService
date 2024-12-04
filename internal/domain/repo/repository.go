package repo

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindUserByID(ctx context.Context, ID int) (*model.User, error)
	UploadImage(ctx context.Context, userId int, image, hash []byte, contentType string) error
	DownloadImage(ctx context.Context, userId int) ([]byte, string, error)
}

type RecoveryCodeRepository interface {
	CreateRecoveryCode(ctx context.Context, code *model.RecoveryCode) (*model.RecoveryCode, error)
	FindRecoveryCodeByEmailAndCode(ctx context.Context, email string, code string) (*model.RecoveryCode, error)
	ExistsRecoveryCodeByEmailAndCode(ctx context.Context, email string, code string) (bool, error)
}

type GameStatsRepository interface {
	GetTotalGamesStatsByUserId(ctx context.Context, userId int) (*model.UserTotalGamesStats, error)
}

type AchievmentsRepository interface {
	GetAchievmentsByUserId(ctx context.Context, userId int, limit int, offset int) ([]*model.Achievements, error)
}

type TopicRepository interface {
	SuggestTopic(ctx context.Context, topic model.Topic) (*model.Topic, error)
	UpdateTopics(ctx context.Context, topicMetatopics []model.TopicMetatopicIds) ([]model.TopicMetatopics, error)
	GetTopics(ctx context.Context, topicStatuses []model.ApprovingStatusEnum, pageSize, pageNumber int) ([]model.TopicMetatopics, int, error)
	GetTopic(ctx context.Context, topicId int) (*model.TopicMetatopics, error)
	GetMetatopics(ctx context.Context, pageSize, pageNumber int) ([]*model.Metatopic, int, error)
}
