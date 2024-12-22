package resolvers

import (
	"context"

	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
)

func (m *queryResolver) GetGameStatus(ctx context.Context, input gen.GameStatusInput) (*gen.GameStatusOutput, error) {
	gameStatus, err := m.useCases.Games.GetGameStatus(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return &gen.GameStatusOutput{
		GameStatus: &gen.GameStatus{
			ID:       gameStatus.ID,
			Status:   string(gameStatus.GameStatusEnum),
			WinnerID: &gameStatus.WinnerId,
			StartAt:  gameStatus.StartAt,
			FinishAt: gameStatus.FinishAt,
		},
	}, nil
}
