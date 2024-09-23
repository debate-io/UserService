package repo

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/ztrue/tracerr"
)

var (
	ErrRoleNotFound = tracerr.New("role not found")
)

type RoleRepository interface {
	FindUserByID(ctx context.Context, id int) (*model.Role, error)
}
