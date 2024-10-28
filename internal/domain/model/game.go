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
