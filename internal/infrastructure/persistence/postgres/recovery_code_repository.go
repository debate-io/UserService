package postgres

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/go-pg/pg/v9"
	"github.com/ztrue/tracerr"
)

var (
	_ repo.RecoveryCodeRepository = (*RecoveryCodeRepository)(nil)
)

type RecoveryCodeRepository struct {
	db *pg.DB
}

func NewRecoveryCodeRepository(
	db *pg.DB,
) *RecoveryCodeRepository {
	return &RecoveryCodeRepository{
		db: db,
	}
}

func (c *RecoveryCodeRepository) CreateRecoveryCode(ctx context.Context, code *model.RecoveryCode) (*model.RecoveryCode, error) {
	_, err := c.db.ModelContext(ctx, code).
		OnConflict("(email) DO UPDATE").
		Set("code = ?", code.Code).
		Set("expired_at = ?", code.ExpiredAt).
		Insert()

	if err != nil {
		return nil, tracerr.Errorf("failed insert user code: %w", err)
	}

	return code, nil
}

func (c *RecoveryCodeRepository) FindRecoveryCodeByEmailAndCode(ctx context.Context, email, code string) (*model.RecoveryCode, error) {
	result := &model.RecoveryCode{}
	q := c.db.ModelContext(ctx, result).
		Relation("User").
		Where("rc.email = ?", email).
		Where("rc.code = ?", code).
		Where("rc.expired_at > now()")

	if err := q.Select(); err != nil {
		if isNoRowsError(err) {
			return nil, repo.ErrNotFound
		}

		return nil, tracerr.Errorf("failed to find recovery code by code: %w", err)
	}

	return result, nil
}

func (c *RecoveryCodeRepository) ExistsRecoveryCodeByEmailAndCode(ctx context.Context, email, code string) (bool, error) {
	result := &model.RecoveryCode{}
	q := c.db.ModelContext(ctx, result).
		Where("email = ?", email).
		Where("code = ?", code).
		Where("expired_at > now()")

	count, err := q.Count()
	if err != nil {
		return false, tracerr.Errorf("failed to find recovery code by code: %w", err)
	}

	return count > 0, nil
}
