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
		// TODO: change mock to real image URL
		ImageURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQlwq2pr4enZ_frUAdm0vcxieKI3E1ZYxA-8Q&s",
	}
}
