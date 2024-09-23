package types

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ztrue/tracerr"
)

const (
	hoursInDay = 24
)

type Claims struct {
	UserID    int              `json:"userId"`
	ExpiredAt *jwt.NumericDate `json:"expiresAt"`
	Role      string           `json:"role"`
	Email     string           `json:"email"`
}

func (c Claims) Valid() error {
	if c.ExpiredAt.Time.Before(time.Now()) {
		return tracerr.New("jwt too older")
	}

	return nil
}

func NewAuthClaims(userID int, email string, roleKey string, daysAuthExpires int) (*Claims, error) {
	return &Claims{
		UserID:    userID,
		ExpiredAt: jwt.NewNumericDate(time.Now().Add(time.Duration(daysAuthExpires*hoursInDay) * time.Hour)),
		Role:      roleKey,
		Email:     email,
	}, nil
}

func NewRecoveryClaims(userID int, daysRecoveryExpires int) (*Claims, error) {
	return &Claims{
		UserID:    userID,
		ExpiredAt: jwt.NewNumericDate(time.Now().Add(time.Duration(daysRecoveryExpires*hoursInDay) * time.Hour)),
	}, nil
}
