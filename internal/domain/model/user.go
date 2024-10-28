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
	tableName   struct{}  `pg:"images"`
	ID          int       `pg:"id,pk"`
	ContentType string    `pg:"content_type"`
	File        []byte    `pg:"file"`
	CreatedAt   time.Time `pg:"created_at"`
	UpdatedAt   time.Time `pg:"updated_at"`
}

type User struct {
	tableName struct{}  `pg:"users"`
	ID        int       `pg:"id,pk"`
	Role      RoleEnum  `pg:"role"`
	Username  string    `pg:"username"`
	Email     string    `pg:"email"`
	Password  string    `pg:"password"`
	CreatedAt time.Time `pg:"created_at"`
	UpdatedAt time.Time `pg:"updated_at"`
	Image     *Image    `pg:"image_id, rel:has-one"`
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
		validation.Field(&i.ContentType, validation.Required), // Поле обязательно
		validation.Field(&i.File, validation.Required),        // Поле обязательно
	)
}
