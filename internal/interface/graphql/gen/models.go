// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gen

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Achievement struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type AuthenticateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateUserOutput struct {
	Jwt   *string `json:"jwt,omitempty"`
	Error *Error  `json:"error,omitempty"`
}

type GetGamesStatsInput struct {
	UserID int `json:"userId"`
}

type GetGamesStatsOutput struct {
	GamesAmount     int                `json:"gamesAmount"`
	WinsAmout       int                `json:"winsAmout"`
	WinsPercents    float64            `json:"WinsPercents"`
	MetatopicsStats []*MetatopicsStats `json:"metatopicsStats,omitempty"`
	Error           *Error             `json:"error,omitempty"`
}

type GetUserInput struct {
	ID int `json:"id"`
}

type GetUserOutput struct {
	User  *User  `json:"user,omitempty"`
	Error *Error `json:"error,omitempty"`
}

type MetatopicsStats struct {
	Matatpoic    string  `json:"matatpoic"`
	GamesAmount  int     `json:"gamesAmount"`
	WinsAmout    int     `json:"winsAmout"`
	WinsPercents float64 `json:"WinsPercents"`
}

type Mutation struct {
}

type Query struct {
}

type RecoveryPasswordInput struct {
	Email string `json:"email"`
}

type RecoveryPasswordOutput struct {
	Error *Error `json:"error,omitempty"`
}

type RegisterUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserOutput struct {
	User  *User   `json:"user,omitempty"`
	Jwt   *string `json:"jwt,omitempty"`
	Error *Error  `json:"error,omitempty"`
}

type ResetPasswordInput struct {
	Code     string `json:"code"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPasswordOutput struct {
	Error *Error `json:"error,omitempty"`
}

type SuggestTopicInput struct {
	Name string `json:"name"`
}

type SuggestTopicOutput struct {
	Topic *Topic `json:"topic,omitempty"`
	Error *Error `json:"error,omitempty"`
}

type Topic struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateEmailInput struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateEmailOutput struct {
	Error *Error `json:"error,omitempty"`
}

type UpdatePasswordInput struct {
	ID          int    `json:"id"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UpdatePasswordOutput struct {
	Error *Error `json:"error,omitempty"`
}

type UpdateUserInput struct {
	ID       int     `json:"id"`
	Username *string `json:"username,omitempty"`
	ImageID  *int    `json:"imageId,omitempty"`
}

type UpdateUserOutput struct {
	User  *User  `json:"user"`
	Error *Error `json:"error,omitempty"`
}

type User struct {
	ID        int       `json:"id"`
	Role      Role      `json:"role"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ImageURL  string    `json:"imageUrl"`
}

type VerifyRecoveryCodeInput struct {
	Code  string `json:"code"`
	Email string `json:"email"`
}

type VerifyRecoveryCodeOutput struct {
	Error *Error `json:"error,omitempty"`
}

// Чтобы понять какая придёт, смотри описание метода API
type Error string

const (
	ErrorNotFound           Error = "NOT_FOUND"
	ErrorValidation         Error = "VALIDATION"
	ErrorInvalidCredentials Error = "INVALID_CREDENTIALS"
	ErrorAlreadyExist       Error = "ALREADY_EXIST"
)

var AllError = []Error{
	ErrorNotFound,
	ErrorValidation,
	ErrorInvalidCredentials,
	ErrorAlreadyExist,
}

func (e Error) IsValid() bool {
	switch e {
	case ErrorNotFound, ErrorValidation, ErrorInvalidCredentials, ErrorAlreadyExist:
		return true
	}
	return false
}

func (e Error) String() string {
	return string(e)
}

func (e *Error) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Error(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Error", str)
	}
	return nil
}

func (e Error) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Role string

const (
	RoleUser           Role = "USER"
	RoleContentManager Role = "CONTENT_MANAGER"
	RoleAdmin          Role = "ADMIN"
)

var AllRole = []Role{
	RoleUser,
	RoleContentManager,
	RoleAdmin,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleUser, RoleContentManager, RoleAdmin:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
