package model

import "time"

type Role struct {
	tableName   struct{}  `pg:"public.roles,alias:roles"` // nolint
	ID          int       `pg:"id"`
	Title       string    `pg:"title"`
	CreatedAt time.Time `pg:"createdAt"`
	HeaderKey   string    `pg:"headerKey"`
}
