package model

import (
	"time"

	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	tableName     struct{}   `pg:"public.users,alias:users"` // nolint
	ID            int        `pg:"id"`
	FirstName     string     `pg:"firstName"`
	LastName      string     `pg:"lastName"`
	Email         string     `pg:"email"`
	Username      string     `pg:"username"`
	CreatedAt     time.Time  `pg:"createdAt"`
	UpdatedAt     time.Time  `pg:"updatedAt"`
	AvatarImageID *int       `pg:"avatarImageId,,use_zero"`
	BirthDate     time.Time  `pg:"birthDate"`
	Gender        gen.Gender `pg:"gender"`
	RoleID        int        `pg:"roleId"`
	Status        gen.Status `pg:"status"`
	Password      string     `pg:"password"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.BirthDate, validation.Required),
	)
}

func (u *User) SetAvatarImageID(val int) {
	u.AvatarImageID = &val
}
