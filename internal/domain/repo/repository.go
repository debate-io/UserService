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
