package mappers

import (
	"fmt"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
)

func MapUserToDTO(user *model.User) *gen.User {
	return &gen.User{
		ID:        int(user.ID),
		Role:      gen.Role(user.Role),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		ImageURL:  fmt.Sprintf("/user_image/%s", string(user.ID)),
	}
}

func MapUsersToDTO(values []*model.User) []*gen.User {
	res := []*gen.User{}

	if values == nil || len(values) == 0 {
		return res
	}

	for _, val := range values {
		res = append(res, MapUserToDTO(val))
	}

	return res
}
