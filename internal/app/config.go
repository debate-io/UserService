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
	IsDebug             bool   `validate:"omitempty"`
	JwtSecretAuth       string `validate:"required"`
	JwtSecretMessages   string `validate:"required"`
	DaysAuthExpires     int    `validate:"required"`
	DaysRecoveryExpires int    `validate:"required"`
	Smtp                SmtpConfig
}

type SmtpConfig struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required"`
	Username string `validate:"required"`
	Password string `validate:"required"`
	SSL      bool   `validate:"omitempty"`
	From     string `validate:"required"`
	SSL      bool
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
		Smtp: SmtpConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			SSL:      os.Getenv("SMTP_SSL") == "true",
			From:     os.Getenv("SMTP_FROM"),
		},
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

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return nil, err
	}

	config.Smtp.Port = smtpPort

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}
