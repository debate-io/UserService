package auth

import (
	"time"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ztrue/tracerr"
)

const (
	hoursInDay = 24
)

type Config struct {
	JwtSecretAuth       string
	JwtSecretMessages   string
	DaysAuthExpires     int
	DaysRecoveryExpires int
}

type AuthService struct {
	cfg Config
}

func NewAuthService(cfg Config) *AuthService {
	return &AuthService{
		cfg: cfg,
	}
}

func (a *AuthService) GenerateTokenByClaims(claims *model.Claims) (string, error) {
	signBytes := []byte(a.cfg.JwtSecretAuth)
	claimsWithExpiredAt := &model.Claims{
		UserID:    claims.UserID,
		ExpiredAt: jwt.NewNumericDate(time.Now().Add(time.Duration(a.cfg.DaysAuthExpires*hoursInDay) * time.Hour)),
		Role:      claims.Role,
		Email:     claims.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsWithExpiredAt)

	ss, err := token.SignedString(signBytes)
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	return ss, nil
}

func (a *AuthService) ParseToken(jwtStr string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(jwtStr, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.cfg.JwtSecretAuth), nil
	})
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	claims := token.Claims.(model.Claims)
	if claims.Valid() != nil {
		return nil, tracerr.Wrap(err)
	}

	return &claims, nil
}

func (a *AuthService) GetDaysAuthExpires() int {
	return a.cfg.DaysAuthExpires
}

func (a *AuthService) GetDaysRecoveryExpires() int {
	return a.cfg.DaysRecoveryExpires
}
