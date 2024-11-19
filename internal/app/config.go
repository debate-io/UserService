package app

import (
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	ServiceName string `validate:"required"`
	PostgresDsn string `validate:"required"`
	Address     string `validate:"required"`
	IsDebug     bool   `validate:"omitempty"`
	Smtp        SmtpConfig
	Jwt         jwtConfig
}

type SmtpConfig struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required"`
	Username string `validate:"required"`
	Password string `validate:"required"`
	From     string `validate:"required"`
	SSL      bool
}

type jwtConfig struct {
	JwtSecretAuth       string `validate:"required"`
	JwtSecretMessages   string `validate:"required"`
	DaysAuthExpires     int    `validate:"required"`
	DaysRecoveryExpires int    `validate:"required"`
}

func (c Config) Validate() error {
	return validator.New().Struct(c)
}

func GetAppConfig() (*Config, error) {
	daysAuthExpires, err := strconv.Atoi(os.Getenv("DAYS_AUTH_EXPIRES"))
	if err != nil {
		return nil, err
	}

	daysRecoveryExpires, err := strconv.Atoi(os.Getenv("DAYS_RECOVERY_EXPIRES"))
	if err != nil {
		return nil, err
	}

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return nil, err
	}

	config := &Config{
		ServiceName: os.Getenv("SERVICE_NAME"),
		PostgresDsn: os.Getenv("POSTGRES_DSN"),
		Address:     os.Getenv("SERVER_ADDRESS"),
		IsDebug:     os.Getenv("IS_DEBUG") == "true",
		Smtp: SmtpConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     smtpPort,
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     os.Getenv("SMTP_FROM"),
			SSL:      os.Getenv("SMTP_SSL") == "true",
		},
		Jwt: jwtConfig{
			JwtSecretAuth:       os.Getenv("JWT_SECRET_AUTH"),
			JwtSecretMessages:   os.Getenv("JWT_SECRET_MESSAGES"),
			DaysAuthExpires:     daysAuthExpires,
			DaysRecoveryExpires: daysRecoveryExpires,
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}
