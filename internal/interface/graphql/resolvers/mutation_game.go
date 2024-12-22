package resolvers

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
)

func (m *mutationResolver) StartGame(ctx context.Context, input gen.StartGameInput) (*gen.StartGameOutput, error) {
	startGameRequest := model.StartGame{
		ID:         input.ID,
		FromUserID: input.FromUserID,
	}

	gameStatus, err := m.useCases.Games.StartGame(ctx, startGameRequest)
	if err != nil {
		return &gen.StartGameOutput{}, err
	}
	return &gen.StartGameOutput{
		GameStatus: &gen.GameStatus{
			ID:       gameStatus.ID,
			Status:   string(gameStatus.GameStatusEnum),
			WinnerID: &gameStatus.WinnerId,
			StartAt:  gameStatus.StartAt,
			FinishAt: gameStatus.FinishAt,
		},
	}, err
}

func (m *mutationResolver) FinishGame(ctx context.Context, input gen.FinishGameInput) (*gen.FinishGameOutput, error) {
	return &gen.FinishGameOutput{}, nil
}
