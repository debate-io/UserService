package model

import "time"

type ApprovingStatusEnum string

const (
	StatusPending  ApprovingStatusEnum = "PENDING"
	StatusApproved ApprovingStatusEnum = "APPROVED"
	StatusDeclined ApprovingStatusEnum = "DECLINED"
)

type Metatopic struct {
	tableName struct{}            `pg:"mrtatopics"`
	ID        int64               `pg:"id,pk`
	Name      string              `pg:"name"`
	Status    ApprovingStatusEnum `pg:"is_approved"` // always use approved
	CreatedAt time.Time           `pg:"created_at"`
}

type Topic struct {
	tableName struct{}            `pg:"topics"`
	ID        int64               `pg:"id,pk"`
	Name      string              `pg:"name"`
	Status    ApprovingStatusEnum `pg:"is_approved"`
	CreatedAt time.Time           `pg:"created_at"`
}

type MetatopicsTopics struct {
	tableName   struct{} `pg:"metatopics_topics"`
	MetatopicID int64    `db:"metatopics_id"`
	TopicID     int64    `db:"topics_id"`
}

type UserMetatopic struct {
	tableName   struct{} `pg:"users_metatopics"`
	UserID      int64    `db:"user_id"`
	MetatopicID int64    `db:"metatopic_id"`
}
