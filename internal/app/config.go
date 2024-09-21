package app

import (
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	ServiceName         string `validate:"required"`
	PostgresDsn         string `validate:"required"`
	Address             string `validate:"required"`
	IsDebug             bool   `validate:"required"`
	JwtSecretAuth       string `validate:"required"`
	JwtSecretMessages   string `validate:"required"`
	DaysAuthExpires     int    `validate:"required"`
	DaysRecoveryExpires int    `validate:"required"`
}

func (c Config) Validate() error {
	return validator.New().Struct(c)
}

func GetAppConfig() (*Config, error) {
	config := &Config{
		ServiceName:       os.Getenv("SERVICE_NAME"),
		PostgresDsn:       os.Getenv("POSTGRES_DSN"),
		Address:           os.Getenv("SERVER_ADDRESS"),
		IsDebug:           os.Getenv("IS_DEBUG") == "true",
		JwtSecretAuth:     os.Getenv("JWT_SECRET_AUTH"),
		JwtSecretMessages: os.Getenv("JWT_SECRET_MESSAGES"),
	}

	daysAuthExpires, err := strconv.Atoi(os.Getenv("DAYS_AUTH_EXPIRES"))
	if err != nil {
		return nil, err
	}

	config.DaysAuthExpires = daysAuthExpires

	daysRecoveryExpires, err := strconv.Atoi(os.Getenv("DAYS_RECOVERY_EXPIRES"))
	if err != nil {
		return nil, err
	}

	config.DaysRecoveryExpires = daysRecoveryExpires

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}
