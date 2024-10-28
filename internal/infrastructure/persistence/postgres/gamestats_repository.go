package postgres

import (
	"context"
	"fmt"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"
)

var (
	_ repo.GameStatsRepository = (*GameStatsRepository)(nil)
)

type GameStatsRepository struct {
	db *pg.DB
}

func NewGameStatsRepository(
	db *pg.DB,
) *GameStatsRepository {
	return &GameStatsRepository{
		db: db,
	}
}

type UserGameStatsResult struct {
	UserGameStats model.UserGameStats
	Error         error
}

func (u *GameStatsRepository) GetTotalGamesStatsByUserId(ctx context.Context, userId int) (*model.UserTotalGamesStats, error) {
	totalStatsChan := u.getTotalUserGameStats(ctx, userId)
	metaTopicStatsChan := u.getUserMetaTopicStats(ctx, userId)

	totalStatsResult := <-totalStatsChan
	if err := totalStatsResult.Error; err != nil {
		return &model.UserTotalGamesStats{}, err
	}

	metaTopicStatsResult := <-metaTopicStatsChan
	if err := metaTopicStatsResult.Error; err != nil {
		return &model.UserTotalGamesStats{}, err
	}

	return &model.UserTotalGamesStats{
		TotalGamesStats: totalStatsResult.UserGameStats,
		MetaTopicStats:  metaTopicStatsResult.MetaTopicStats,
	}, nil
}

func (u *GameStatsRepository) getTotalUserGameStats(ctx context.Context, userId int) <-chan UserGameStatsResult {
	res := make(chan UserGameStatsResult, 1)
	go func() {
		defer close(res)

		var totalStats model.UserGameStats
		_, err := u.db.QueryOneContext(ctx, &totalStats, `
		SELECT ? AS user_id,
		       COUNT(*) AS games_amount,
		       COUNT(*) FILTER (WHERE winner_id = ?) AS wins_amount
		FROM games
		WHERE first_player_id = ? OR second_player_id = ?
	`, userId, userId, userId, userId)

		if err != nil {
			res <- UserGameStatsResult{
				Error: fmt.Errorf("cannot select user games stats: %v", zap.Error(err)),
			}
			return
		}

		res <- UserGameStatsResult{
			UserGameStats: totalStats,
			Error:         err,
		}
	}()

	return res
}

type UserMetaTopicStatsResult struct {
	MetaTopicStats map[string]model.UserGameStats
	Error          error
}

func (u *GameStatsRepository) getUserMetaTopicStats(ctx context.Context, userId int) <-chan UserMetaTopicStatsResult {
	res := make(chan UserMetaTopicStatsResult, 1)
	go func() {
		defer close(res)

		var metaTopicStats []struct {
			MetaTopicName string
			GamesAmount   int
			WinsAmount    int
		}

		_, err := u.db.QueryContext(ctx, &metaTopicStats, `
			SELECT m.name AS meta_topic_name,
			       COUNT(*) AS games_amount,
			       COUNT(*) FILTER (WHERE winner_id = ?) AS wins_amount
			FROM games g
			JOIN metatopics m ON g.metatopic_id = m.id
			WHERE (g.first_player_id = ? OR g.second_player_id = ?)
			GROUP BY m.name
		`, userId, userId, userId)

		if err != nil {
			res <- UserMetaTopicStatsResult{
				Error: fmt.Errorf("cannot select user meta topic stats: %v", zap.Error(err)),
			}
			return
		}

		metaTopicMap := make(map[string]model.UserGameStats)
		for _, stat := range metaTopicStats {
			metaTopicMap[stat.MetaTopicName] = model.UserGameStats{
				UserId:      userId,
				GamesAmount: stat.GamesAmount,
				WinsAmount:  stat.WinsAmount,
			}
		}

		res <- UserMetaTopicStatsResult{
			MetaTopicStats: metaTopicMap,
			Error:          nil,
		}
	}()

	return res
}
