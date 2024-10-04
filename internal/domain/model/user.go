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
	RoleDefaultUser    RoleEnum = "DEFAULT_USER"
)

type User struct {
	tableName struct{}  `pg:"public.users,alias:users"`
	ID        int64     `pg:"id,pk"`      // BIGINT и PRIMARY KEY
	Role      RoleEnum  `pg:"role`        // role_enum с умолчанием
	Username  string    `pg:"username"`   // Текстовое поле, обязательное
	Email     string    `pg:"email"`      // Email с уникальностью
	Password  string    `pg:"password"`   // Пароль, обязательное
	CreatedAt time.Time `pg:"created_at`  // Время создания с умолчанием
	UpdatedAt time.Time `pg:"updated_at"` // Время обновления с умолчанием
	Image     *int      `pg:"image"`      // Поле для OID (ссылки на файл), может быть NULL
}

// Валидация полей структуры User
func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email), // Проверка, что Email валиден
		validation.Field(&u.Username, validation.Required),        // Поле username обязательно
		validation.Field(&u.Password, validation.Required),        // Поле password обязательно
	)
}
