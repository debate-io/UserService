package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ztrue/tracerr"
)

const (
	hoursInDay = 24
)

var (
	_ jwt.Claims = (*Claims)(nil)
)

type Claims struct {
	UserID    int              `json:"userId"`
	ExpiredAt *jwt.NumericDate `json:"expiresAt"`
	Role      RoleEnum         `json:"role"`
	Email     string           `json:"email"`
}

func (c Claims) Valid() error {
	if c.ExpiredAt.Time.Before(time.Now()) {
		return tracerr.New("jwt too older")
	}

	return nil
}

func NewAuthClaims(userID int, email string, role RoleEnum, daysAuthExpires int) (*Claims, error) {
	return &Claims{
		UserID:    userID,
		ExpiredAt: jwt.NewNumericDate(time.Now().Add(time.Duration(daysAuthExpires*hoursInDay) * time.Hour)),
		Role:      role,
		Email:     email,
	}, nil
}
