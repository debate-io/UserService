// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gen

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

// ###################################################
//
//	input GetUserInput {
//	    id: Int!
//	    gettingAt: Time
//	}
//
//	type GetUserOutput {
//	    user: User
//	    isUpdated: Boolean
//	    error: Error
//	}
type AuthenticateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateUserOutput struct {
	Jwt   *string `json:"jwt,omitempty"`
	Error *Error  `json:"error,omitempty"`
}

type Mutation struct {
}

type Query struct {
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

type User struct {
	ID        int       `json:"id"`
	Role      Role      `json:"role"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ImageURL  string    `json:"imageUrl"`
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
