package model

import "time"

type UserGameStats struct {
	UserId      int
	GamesAmount int
	WinsAmount  int
}

type UserTotalGamesStats struct {
	TotalGamesStats UserGameStats
	MetaTopicStats  map[string]UserGameStats
}

type Achievements struct {
	ID          int       `pg:"id,pk"`
	Name        string    `pg:"name"`
	Description string    `pg:"description"`
	CreateAt    time.Time `pg:"created_at"`
}

type StartGame struct {
	ID         int `pg:"id,pk"`
	FromUserID int `pg:"from_user_id"` // references
}

type FinishGame struct {
	ID            int `pg:"id,pk"`
	FromUserID    int `pg:"from_user_id"` // references
	SecondsInGame int `pg:"seconds_in_game"`
}

type GameResult struct {
	ID         int    `pg:"id,pk"`
	WinnerId   int    `pg:"winner_id"` // references
	ResultText string `pg:"result_text"`
}

type GameStatusEnum string

const (
	GameStatusPending  GameStatusEnum = "PENDING"
	GameStatusStarted  GameStatusEnum = "STARTED"
	GameStatusDeclined GameStatusEnum = "DECLINED"
)

type GameStatus struct {
	ID             int `pg:"id,pk"`
	FirstPlayerId  int `pg:"first_player_id"`  // references
	SecondPlayerId int `pg:"second_player_id"` // references

	FirstRequest time.Time

	GameStatusEnum GameStatusEnum `pg:"status"`
	WinnerId       int            `pg:"winner_id"` // references
	StartAt        time.Time      `pg:"start_at"`
	FinishAt       time.Time      `pg:"finish_at"`
}
