package registry

import (
	"github.com/debate-io/service-auth/internal/usecases"
	"go.uber.org/zap"
)

type UseCases struct {
	Users  *usecases.User
	Topics *usecases.Topic
}

type Container struct {
	UseCases *UseCases
	Logger   *zap.Logger
}
