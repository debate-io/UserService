package app

import (
	"github.com/debate-io/service-auth/internal/infrastructure/smtp"
	"os"
	"time"

	pg "github.com/go-pg/pg/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/debate-io/service-auth/internal/infrastructure/persistence/postgres"
	"github.com/debate-io/service-auth/internal/interface/server"
	"github.com/debate-io/service-auth/internal/registry"
	"github.com/debate-io/service-auth/internal/usecases"
)

type App struct {
	Logger     *zap.Logger
	Server     *server.Server
	DB         *pg.DB
	SmtpSender *smtp.Sender
	Config     *Config
}

func NewApp(config *Config) *App {
	logger := NewLogger(config.IsDebug)

	if err := setOsTimezone("UTC"); err != nil {
		logger.Error("can't set timezone in OS environment", zap.Error(err))
	}

	db, err := postgres.NewPostgresDatabase(config.PostgresDsn, config.ServiceName, logger)
	if err != nil {
		logger.Error("can't connect to postgres database", zap.Error(err))
	}

	smtpClient, err := smtp.NewSender(&smtp.Config{
		Host:     config.Smtp.Host,
		Port:     config.Smtp.Port,
		Username: config.Smtp.Username,
		Password: config.Smtp.Password,
		SSL:      config.Smtp.SSL,
		From:     config.Smtp.From,
	})
	if err != nil {
		logger.Error("can't connect to SMTP server", zap.Error(err))
	}

	return &App{
		Logger:     logger,
		Server:     server.NewServer(logger),
		DB:         db,
		SmtpSender: smtpClient,
		Config:     config,
	}
}

func NewLogger(isDebug bool) *zap.Logger {
	var logger *zap.Logger

	if isDebug {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ = config.Build()
	} else {
		logger, _ = zap.NewProduction()
	}

	return logger
}

func (app *App) CloseConnections() {
	if err := app.DB.Close(); err != nil {
		app.Logger.Fatal("failed to close database connections", zap.Error(err))
	}

	app.Logger.Debug("all connections close")
}

func (app *App) RunApp() {
	if err := app.Server.ListenAndServe(app.Config.Address, app.beforeShutdown); err != nil {
		app.Logger.Fatal("failed to start listen server", zap.Error(err))
	}
}

func (app *App) Initialize() {
	container := app.NewContainer()
	app.Server.InitMiddlewares(app.Config.IsDebug)
	app.Server.InitRoutes(container, app.Config.IsDebug)
}

func (app *App) beforeShutdown() {
	if !app.Config.IsDebug {
		const shutdownIdle = 9 * time.Second

		time.Sleep(shutdownIdle)
	}

	app.Logger.Info("Осторожно двери закрываются, шотдаун, ребзя")
	app.CloseConnections()
}

func (app *App) NewContainer() *registry.Container {
	userRepo := postgres.NewUserRepository(app.DB)
	recoveryCodeRepo := postgres.NewRecoveryCodeRepository(app.DB)
	JwtConfigs := usecases.NewJwtConfigsUseCases(app.Config.JwtSecretAuth, app.Config.JwtSecretMessages, app.Config.DaysAuthExpires, app.Config.DaysRecoveryExpires)

	useCases := &registry.UseCases{
		Users: usecases.NewUserUseCases(userRepo, recoveryCodeRepo, app.SmtpSender, *JwtConfigs),
	}

	return &registry.Container{UseCases: useCases, Logger: app.Logger}
}

func setOsTimezone(tz string) error {
	return os.Setenv("TZ", tz)
}
