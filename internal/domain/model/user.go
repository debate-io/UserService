package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RoleEnum string

const (
	RoleAdmin          RoleEnum = "ADMIN"
	RoleContentManager RoleEnum = "CONTENT_MANAGER"
	RoleDefaultUser    RoleEnum = "USER"
)

type Image struct {
	tableName struct{} `pg:"images"`
	ID        int64    `pg:"id,pk"`
	Format    string   `pg:"format"`
	File      []byte   `pg:"file"`
}

type User struct {
	tableName struct{}  `pg:"users"`
	ID        int64     `pg:"id,pk"`
	Role      RoleEnum  `pg:"role`
	Username  string    `pg:"username"`
	Email     string    `pg:"email"`
	Password  string    `pg:"password"`
	CreatedAt time.Time `pg:"created_at`
	UpdatedAt time.Time `pg:"updated_at"`
	ImageId   *Image    `pg:"rel:has-one"`
}

// Валидация полей структуры User
func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email), // Проверка, что Email валиден
		validation.Field(&u.Role, validation.Required),            // Поле обязательно
		validation.Field(&u.Username, validation.Required),        // Поле обязательно
		validation.Field(&u.Password, validation.Required),        // Поле обязательно
		validation.Field(&u.CreatedAt, validation.Required),       // Поле обязательно
		validation.Field(&u.UpdatedAt, validation.Required),       // Поле обязательно

	)
}

// Валидация полей структуры Image
func (i *Image) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.Format, validation.Required), // Поле обязательно
		validation.Field(&i.File, validation.Required),   // Поле обязательно
	)
}
