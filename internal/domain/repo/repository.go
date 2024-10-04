package repo

import (
	"context"

	"github.com/debate-io/service-auth/internal/domain/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	//UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	//FindUserByID(ctx context.Context, ID int) (*model.User, error)
}
