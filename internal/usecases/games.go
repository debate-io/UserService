package usecases

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
)

type Game struct {
	gameRepo repo.GameRepository
}

func NewGameUseCase(gameRepo repo.GameRepository) *Game {
	return &Game{
		gameRepo: gameRepo,
	}
}

func (g *Game) StartGame(ctx context.Context, startGameRequest model.StartGame) (model.GameStatus, error) {

	game, err := g.gameRepo.StartGame(ctx, startGameRequest)
	if err != nil {
		return model.GameStatus{}, err
	}

	return game, nil
}

func (g *Game) FinishGame(ctx context.Context, finishGameRequest model.FinishGame) (model.GameResult, error) {
	return g.gameRepo.FinishGame(ctx, finishGameRequest)
}

func (g *Game) GetGameStatus(ctx context.Context, gameID int) (model.GameStatus, error) {
	return g.gameRepo.GetGameById(ctx, gameID)
}
