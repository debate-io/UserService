package usecases

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/debate-io/service-auth/internal/interface/server/middleware"
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
	claims := ctx.Value(middleware.JwtClaimsKey).(*model.Claims)
	if claims == nil {
		return model.GameStatus{}, repo.ErrUnauthorized
	}

	game, err := g.gameRepo.GetGameById(ctx, gameID)
	if err != nil {
		return model.GameStatus{}, err
	}

	// Ретрай от первого игрока
	if game.FirstPlayerId == claims.UserID {
		return game, nil
	}

	// Завершение игры по дедлайну ожидания подтверждения начала игры от второго игрока
	// deadline := game.FirstRequest.UTC().Add(waitingDuration)
	// if time.Now().UTC().After(deadline) {
	return g.gameRepo.FinishGameByDeadline(ctx, claims.UserID, game)
	// }
}
