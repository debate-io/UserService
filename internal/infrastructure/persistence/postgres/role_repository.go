package postgres

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/go-pg/pg/v9"
	"github.com/ztrue/tracerr"
)

var (
	_ repo.RoleRepository = (*RoleRepository)(nil)
)

type RoleRepository struct {
	db *pg.DB
}

func NewRoleRepository(
	db *pg.DB,
) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) FindUserByID(ctx context.Context, id int) (*model.Role, error) {
	result := &model.Role{}
	query := r.db.ModelContext(ctx, result).
		Where("id in (?)", id)

	if err := query.Select(); err != nil {
		if isNoRowsError(err) {
			return nil, repo.ErrRoleNotFound
		}

		return nil, tracerr.Errorf("failed to find role: %w", err)
	}

	return result, nil
}
