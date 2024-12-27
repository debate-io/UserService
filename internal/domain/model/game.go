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

type Game struct {
	ID             int64     `pg:"id, pk"`
	FirstPlayerID  int64     `pg:"first_player_id"`
	SecondPlayerID int64     `pg:"second_player_id"`
	RoomID         string    `pg:"room_id"`
	WinnerID       int64     `pg:"winner_id"`
	MetatopicID    int64     `pg:"metatopic_id"`
	TopicID        int64     `pg:"topic_id"`
	CreatedAt      time.Time `pg:"created_at"`
}

type StartGame struct {
	RoomID     string `pg:"id,pk"`
	FromUserID int    `pg:"from_user_id"` // references
}

type FinishGame struct {
	RoomID        string `pg:"id,pk"`
	FromUserID    int    `pg:"from_user_id"` // references
	SecondsInGame int    `pg:"seconds_in_game"`
}

type GameResult struct {
	RoomID     string `pg:"id,pk"`
	WinnerId   int    `pg:"winner_id"` // references
	ResultText string `pg:"result_text"`
}

type GameStatusEnum string

const (
	GameStatusPending  GameStatusEnum = "PENDING"
	GameStatusStarted  GameStatusEnum = "STARTED"
	GameStatusDeclined GameStatusEnum = "DECLINED"
	GameStatusFinished GameStatusEnum = "FINISHED"
)

type GameStatus struct {
	ID             string `pg:"id,pk"`
	FirstPlayerId  int    `pg:"first_player_id"`  // references
	SecondPlayerId int    `pg:"second_player_id"` // references

	FirstPlayerScore  int
	SecondPlayerScore int

	FirstRequest       time.Time
	FirstFinishRequest *time.Time

	GameStatusEnum GameStatusEnum `pg:"status"`
	WinnerId       int            `pg:"winner_id"` // references
	StartAt        time.Time      `pg:"start_at"`
	FinishAt       time.Time      `pg:"finish_at"`
}
