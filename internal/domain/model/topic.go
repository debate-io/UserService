package model

import "time"

type ApprovingStatusEnum string

const (
	StatusPending  ApprovingStatusEnum = "PENDING"
	StatusApproved ApprovingStatusEnum = "APPROVED"
	StatusDeclined ApprovingStatusEnum = "DECLINED"
)

type Metatopic struct {
	tableName struct{}            `pg:"public.metatopic"`
	ID        int                 `pg:"id,pk`
	Name      string              `pg:"name"`
	Status    ApprovingStatusEnum `pg:"is_approved"` // always use approved
	CreatedAt time.Time           `pg:"created_at"`
}

type Topic struct {
	tableName struct{}            `pg:"public.topics"`
	ID        int                 `pg:"id,pk"`
	Name      string              `pg:"name"`
	Status    ApprovingStatusEnum `pg:"is_approved"`
	CreatedAt time.Time           `pg:"created_at"`
}

type MetatopicsTopics struct {
	tableName   struct{} `pg:"public.metatopics_topics"`
	MetatopicID int      `pg:"metatopics_id"`
	TopicID     int      `pg:"topics_id"`
}

type UserMetatopic struct {
	tableName   struct{} `pg:"public.users_metatopics"`
	UserID      int      `pg:"user_id"`
	MetatopicID int      `pg:"metatopic_id"`
}
