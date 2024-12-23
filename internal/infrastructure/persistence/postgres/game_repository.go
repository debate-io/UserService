package postgres

import (
	"context"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/ztrue/tracerr"
)

var (
	_ repo.GameRepository = (*GameRepository)(nil)
)

const (
	gameDuration    = time.Second * 40
	waitingDuration = time.Second * 20
)

type GameRepository struct {
	Games   map[string]model.GameStatus
	Mu      sync.Mutex
	Results []string
}

// IsGameOverByDeadline implements repo.GameRepository.
func (g *GameRepository) IsGameOverByDeadline(ctx context.Context, roomId string) bool {
	g.Mu.Lock()
	game, ok := g.Games[roomId]
	g.Mu.Unlock()

	if !ok {
		return true
	}

	deadline := game.FirstRequest.UTC().Add(waitingDuration)
	return time.Now().UTC().After(deadline)
}

func NewGameRepository() *GameRepository {
	results := []string{
		"Затычка %s",
	}
	return &GameRepository{
		Games:   make(map[string]model.GameStatus),
		Mu:      sync.Mutex{},
		Results: results,
	}
}

// FinishGame implements repo.GameRepository.
func (g *GameRepository) FinishGame(ctx context.Context, startGame model.FinishGame) (model.GameResult, error) {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	game, ok := g.Games[startGame.RoomID]
	if !ok {
		return model.GameResult{}, tracerr.New("Game not found")
	}

	if game.FirstFinishRequest == nil {
		now := time.Now()
		game.FirstFinishRequest = &now
	}
	if game.FirstPlayerId == startGame.FromUserID {
		game.FirstPlayerScore = startGame.SecondsInGame
	} else {
		game.SecondPlayerScore = startGame.SecondsInGame
	}

	g.Games[startGame.RoomID] = game

	gameFinished := game.FirstPlayerScore != 0 &&
		game.SecondPlayerScore != 0 ||
		time.Now().UTC().After(game.FirstFinishRequest.UTC().Add(waitingDuration))

	if gameFinished {
		var winnerID int
		delta := game.FirstPlayerScore - game.SecondPlayerScore
		if delta > int(gameDuration/10) {
			winnerID = game.FirstPlayerId
		} else if delta < -int(gameDuration/10) {
			winnerID = game.SecondPlayerId
		} else if game.FirstFinishRequest.Second()%2 == 0 {
			winnerID = game.FirstPlayerId
		} else {
			winnerID = game.SecondPlayerId
		}
		g.Games[startGame.RoomID] = game

		return model.GameResult{
			RoomID:     game.ID,
			WinnerId:   winnerID,
			ResultText: g.Results[rand.Int()%len(g.Results)],
		}, nil
	}

	return model.GameResult{
		RoomID:     startGame.RoomID,
		WinnerId:   0,
		ResultText: "",
	}, nil
}

func (g *GameRepository) FinishGameByDeadline(ctx context.Context, fromUserId int, currentGameStatus model.GameStatus) (model.GameStatus, error) {
	if fromUserId == currentGameStatus.FirstPlayerId {
		currentGameStatus.WinnerId = currentGameStatus.SecondPlayerId
	} else {
		currentGameStatus.WinnerId = currentGameStatus.FirstPlayerId
	}

	currentGameStatus.FinishAt = time.Now().UTC()
	currentGameStatus.GameStatusEnum = model.GameStatusDeclined

	g.Mu.Lock()
	defer g.Mu.Unlock()
	g.Games[currentGameStatus.ID] = currentGameStatus
	return g.Games[currentGameStatus.ID], nil
}

func (g *GameRepository) GetGameById(ctx context.Context, roomId string) (model.GameStatus, error) {
	g.Mu.Lock()
	defer g.Mu.Unlock()
	return g.Games[roomId], nil
}

func (g *GameRepository) StartGame(ctx context.Context, startGame model.StartGame) (model.GameStatus, error) {
	g.Mu.Lock()
	game, exist := g.Games[startGame.RoomID]
	g.Mu.Unlock()

	if !exist {
		// Пришёл первый игрок
		newGame := model.GameStatus{
			ID:             startGame.RoomID,
			FirstPlayerId:  startGame.FromUserID,
			FirstRequest:   time.Now().UTC(),
			GameStatusEnum: model.GameStatusPending,
			StartAt:        time.Now().UTC(),
			FinishAt:       time.Now().UTC(),
		}
		g.Mu.Lock()
		g.Games[newGame.ID] = newGame
		g.Mu.Unlock()
	} else {
		// Ретрай от первого игрока
		if game.FirstPlayerId == startGame.FromUserID {
			return game, nil
		}
		// Завершение игры по дедлайну ожидания подтверждения начала игры от второго игрока
		if g.IsGameOverByDeadline(ctx, game.ID) {
			return g.FinishGameByDeadline(ctx, startGame.FromUserID, game)
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
		g.Mu.Lock()
		g.Games[newGame.ID] = newGame
		g.Mu.Unlock()
	}

	g.Mu.Lock()
	defer g.Mu.Unlock()
	return g.Games[startGame.RoomID], nil
}
