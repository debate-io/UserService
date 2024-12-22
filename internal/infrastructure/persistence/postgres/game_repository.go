package postgres

import (
	"context"
	"time"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
)

var (
	_ repo.GameRepository = (*GameRepository)(nil)
)

const (
	gameDuration    = time.Minute * 10
	waitingDuration = time.Second * 20
)

type GameRepository struct {
	Games map[int]model.GameStatus
}

func NewGameRepository() *GameRepository {
	return &GameRepository{
		Games: make(map[int]model.GameStatus),
	}
}

// FinishGame implements repo.GameRepository.
func (g *GameRepository) FinishGame(ctx context.Context, startGame model.FinishGame) (model.GameResult, error) {
	panic("unimplemented")
}

func (g *GameRepository) FinishGameByDeadline(ctx context.Context, startGameRequest model.StartGame, currentGameStatus model.GameStatus) (model.GameStatus, error) {
	if startGameRequest.FromUserID == currentGameStatus.FirstPlayerId {
		currentGameStatus.WinnerId = currentGameStatus.SecondPlayerId
	} else {
		currentGameStatus.WinnerId = currentGameStatus.FirstPlayerId
	}

	currentGameStatus.FinishAt = time.Now().UTC()
	currentGameStatus.GameStatusEnum = model.GameStatusDeclined
	g.Games[currentGameStatus.ID] = currentGameStatus
	return g.Games[currentGameStatus.ID], nil
}

func (g *GameRepository) GetGameById(ctx context.Context, id int) (model.GameStatus, error) {
	return g.Games[id], nil
}

func (g *GameRepository) StartGame(ctx context.Context, startGame model.StartGame) (model.GameStatus, error) {
	game, exist := g.Games[startGame.ID]

	if !exist {
		// Пришёл первый игрок
		newGame := model.GameStatus{
			ID:             startGame.ID,
			FirstPlayerId:  startGame.FromUserID,
			FirstRequest:   time.Now().UTC(),
			GameStatusEnum: model.GameStatusPending,
			StartAt:        time.Now().UTC(),
			FinishAt:       time.Now().UTC(),
		}
		g.Games[newGame.ID] = newGame
	} else {
		// Ретрай от первого игрока
		if game.FirstPlayerId == startGame.FromUserID {
			return game, nil
		}
		// Завершение игры по дедлайну ожидания подтверждения начала игры от второго игрока
		deadline := game.FirstRequest.UTC().Add(waitingDuration)
		if time.Now().UTC().After(deadline) {
			return g.FinishGameByDeadline(ctx, startGame, game)
		}

		// Пришёл второй игрок
		newGame := model.GameStatus{
			ID:             game.ID,
			FirstPlayerId:  game.FirstPlayerId,
			SecondPlayerId: startGame.FromUserID,
			FirstRequest:   time.Now().UTC(),
			GameStatusEnum: model.GameStatusStarted,
			StartAt:        time.Now().UTC(),
			FinishAt:       time.Now().Add(gameDuration).UTC(),
		}
		g.Games[newGame.ID] = newGame
	}

	return g.Games[startGame.ID], nil
}
