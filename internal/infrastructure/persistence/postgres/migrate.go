package postgres

import (
	"errors"

	goMigrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Импорт драйвера для PostgreSQL
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Импорт драйвера для работы с файлами
	"go.uber.org/zap"
)

type migrateLogger struct {
	Log *zap.SugaredLogger
}

func (l *migrateLogger) Printf(format string, v ...interface{}) {
	l.Log.Info(v)
}

func (l *migrateLogger) Verbose() bool {
	return false
}

func startMigrate(dsn string, logger *zap.Logger) error {

	migrate, err := goMigrate.New("file://migrations", dsn)
	if err != nil {
		return errors.Join(errors.New("migrate failed: "), err)
	}
	defer migrate.Close()

	migrate.Log = &migrateLogger{
		Log: logger.Sugar(),
	}

	if err = migrate.Up(); err != nil && err != goMigrate.ErrNoChange {
		return errors.Join(errors.New("migrate failed: "), err)
	}

	if err == goMigrate.ErrNoChange {
		logger.Info("База в актуальном состоянии")
	} else {
		logger.Info("Миграции успешно установлены")
	}

	return nil
}
