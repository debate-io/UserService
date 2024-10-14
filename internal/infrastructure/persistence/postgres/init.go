package postgres

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/ztrue/tracerr"
	"go.uber.org/zap"
)

func NewPostgresDatabase(dsn string, appName string, logger *zap.Logger) (*pg.DB, error) {
	options, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, tracerr.Errorf("can't parse postgres dsn: %w", err)
	}

	options.ApplicationName = fmt.Sprintf("[%s]", appName)
	options.TLSConfig = nil

	db := pg.Connect(options)
	db.AddQueryHook(QueryLogger{Logger: logger})

	if err := startMigrate(dsn, logger); err != nil {
		return nil, err
	}

	return db, nil
}

type QueryLogger struct {
	Logger *zap.Logger
}

func (l QueryLogger) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (l QueryLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	sql, err := q.FormattedQuery()
	if err != nil {
		l.Logger.Error("SQL error", zap.String("sql", sql), zap.Error(err))
	} else {
		l.Logger.Debug(fmt.Sprintf("SQL: %s", sql))
	}

	return nil
}
