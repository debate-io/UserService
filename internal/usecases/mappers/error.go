package mappers

import "github.com/debate-io/service-auth/internal/interface/graphql/gen"

func NewDTOError(val gen.Error) *gen.Error {
	return &val
}
