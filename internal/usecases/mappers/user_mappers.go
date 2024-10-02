package mappers

import (
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

/* func MapClaimsToDTO(claims *types.Claims) *gen.Claims {
	return &gen.Claims{
		UserID:    claims.UserID,
		Role:      gen.Role(claims.Role),
		ExpiredAt: claims.ExpiredAt.Time,
		Email:     claims.Email,
	}
} */
