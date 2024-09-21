package repo

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/ztrue/tracerr"
)

var (
	ErrUserAlreadyExist = tracerr.New("user already register")
	ErrUserNotFound     = tracerr.New("user not found")
)

type FindUsersQuery struct {
	IDAnyOf []int
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindUsers(ctx context.Context, query *FindUsersQuery) ([]*model.User, error)
	FindUserByID(ctx context.Context, ID int) (*model.User, error)
}
