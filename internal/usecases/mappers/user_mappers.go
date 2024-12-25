package mappers

import (
	"fmt"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
)

func MapUserToDTO(user *model.User) *gen.User {
	url := "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQlwq2pr4enZ_frUAdm0vcxieKI3E1ZYxA-8Q&s"
	if user.Image != nil {
		url = fmt.Sprintf("http://185.84.163.166:9090/user/%d/image/", user.ID)
	}

	return &gen.User{
		ID:        int(user.ID),
		Role:      gen.Role(user.Role),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		ImageURL:  url,
	}
}

func MapUsersToDTO(users []*model.User) []*gen.User {
	var genUsers []*gen.User
	for i := range users {
		genUsers = append(genUsers, MapUserToDTO(users[i]))
	}

	return genUsers
}
