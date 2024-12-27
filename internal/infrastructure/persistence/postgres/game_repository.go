package postgres

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/go-pg/pg/v9"
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
	db      *pg.DB
}

func (g *GameRepository) SetWinnerId(ctx context.Context, roomID string, winnerID int) error {
	game := model.Game{
		WinnerID: int64(winnerID),
	}

	fmt.Printf("%+v\n", game)

	_, err := g.db.ModelContext(ctx, &game).Where("room_uid = ?", roomID).Column("winner_id").Update()
	if err != nil {
		return err
	}

	return nil
}

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

func NewGameRepository(db *pg.DB) *GameRepository {

	results := []string{
		"%s победил в дебатах: его аргументы звучали чётче и убедительнее. Второй участник не смог достойно ответить на его доводы.",
		"В дебатах %s явно переиграл оппонента. Его логика и примеры выглядели сильнее.",
		"%s смог убедить всех своим выступлением. Второй участник не сумел предложить равнозначную аргументацию.",
		"Победа за %s: его речь была яркой и содержательной. Оппонент же допустил много пробелов в логике.",
		"%s уверенно взял верх в дебатах. Его позиция звучала цельно, в отличие от разрозненной речи соперника.",
		"Убедительность %s принесла ему победу. Его оппонент не смог достойно поддержать свою точку зрения.",
		"%s вышел победителем благодаря продуманным доводам. Оппонент выглядел менее подготовленным.",
		"Аргументы %s оказались неоспоримыми. Второй участник явно проиграл в убедительности.",
		"%s выиграл дебаты своей чёткой позицией. Соперник допустил много ошибок в доказательствах.",
		"Победа за %s, чья речь была логичной и убедительной. Его соперник, напротив, выглядел неуверенно.",
		"%s превзошёл своего оппонента в дебатах. Его доводы звучали уверенно и продуманно.",
		"У %s получилось донести свою точку зрения ярче. Второй участник не справился с контраргументами.",
		"Речь %s заслужила победу благодаря логике и структуре. Его соперник запутался в собственных доводах.",
		"%s доминировал в дебатах, приводя сильные доказательства. Второй участник оказался менее убедительным.",
		"Победа за %s: его аргументы были чёткими и последовательными. Второй участник выглядел слабее.",
		"%s уверенно защитил свою точку зрения. Соперник не сумел найти достойных контраргументов.",
		"Аргументы %s оказались более весомыми. Его соперник, напротив, выглядел растерянно.",
		"%s продемонстрировал высокий уровень подготовки. Второй участник, к сожалению, не смог удержаться на таком же уровне.",
		"Победа %s была очевидной: его доводы оставили соперника без шансов.",
	}

	return &GameRepository{
		Games:   make(map[string]model.GameStatus),
		Mu:      sync.Mutex{},
		Results: results,
		db:      db,
	}
}

// FinishGame implements repo.GameRepository.
func (g *GameRepository) FinishGame(ctx context.Context, finishGame model.FinishGame) (model.GameResult, error) {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	game, ok := g.Games[finishGame.RoomID]
	if !ok {
		return model.GameResult{}, tracerr.New("Game not found")
	}

	if game.FirstFinishRequest == nil {
		now := time.Now()
		game.FirstFinishRequest = &now
	}
	if game.FirstPlayerId == finishGame.FromUserID {
		game.FirstPlayerScore = finishGame.SecondsInGame
	} else {
		game.SecondPlayerScore = finishGame.SecondsInGame
	}

	g.Games[finishGame.RoomID] = game

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
		g.Games[finishGame.RoomID] = game

		var chislo rune = 0x00

		for _, v := range game.ID {
			chislo = chislo ^ v
		}

		return model.GameResult{
			RoomID:     game.ID,
			WinnerId:   winnerID,
			ResultText: g.Results[int(chislo)%len(g.Results)],
		}, nil
	}

	return model.GameResult{
		RoomID:     finishGame.RoomID,
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
