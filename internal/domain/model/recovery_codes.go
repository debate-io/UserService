package model

import "time"

type RecoveryCode struct {
	tableName struct{}  `pg:"recovery_codes,alias:rc"`
	UserEmail string    `pg:"email,pk"`
	User      *User     `pg:"fk:email,rel:has-one"`
	Code      string    `pg:"code"`
	ExpiredAt time.Time `pg:"expired_at"`
}
