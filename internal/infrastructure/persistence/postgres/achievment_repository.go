package postgres

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/go-pg/pg/v9"
)

type AchievementRepository struct {
	db *pg.DB
}

var (
	_ repo.AchievmentsRepository = (*AchievementRepository)(nil)
)

func NewAchievementRepository(db *pg.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

func (r *AchievementRepository) GetAchievmentsByUserId(ctx context.Context, userId int, limit int, offset int) ([]*model.Achievements, error) {
	var achievements []*model.Achievements

	err := r.db.ModelContext(ctx, &achievements).
		Join("JOIN users_achievements ua ON ua.achievement_id = achievements.id").
		Where("ua.user_id = ?", userId).
		Limit(limit).
		Offset(offset).
		Select()

	if err != nil {
		return nil, err
	}

	return achievements, nil
}
