// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gen

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type AuthenticateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateUserOutput struct {
	Jwt   *string `json:"jwt"`
	Error *Error  `json:"error"`
}

type Claims struct {
	UserID    int       `json:"userId"`
	Role      Role      `json:"role"`
	ExpiredAt time.Time `json:"expiredAt"`
	Email     string    `json:"email"`
}

type ConfirmUserInput struct {
	Jwt string `json:"jwt"`
}

type ConfirmUserOutput struct {
	Ok    bool   `json:"ok"`
	Error *Error `json:"error"`
}

type CreateUserInput struct {
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Password      string    `json:"password"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	BirthDate     time.Time `json:"birthDate"`
	Gender        Gender    `json:"gender"`
	AvatarImageID *int      `json:"avatarImageId"`
}

type CreateUserOutput struct {
	User  *User   `json:"user"`
	Jwt   *string `json:"jwt"`
	Error *Error  `json:"error"`
}

type DeleteUserInput struct {
	ID int `json:"id"`
}

type DeleteUserOutput struct {
	Ok    bool   `json:"ok"`
	Error *Error `json:"error"`
}

type FindUsersInput struct {
	IDAnyOf []int `json:"idAnyOf"`
}

type FindUsersOutput struct {
	Users []*User `json:"users"`
}

type GetClaimsInput struct {
	Jwt string `json:"jwt"`
}

type GetClaimsOutput struct {
	Claims *Claims `json:"claims"`
	Error  *Error  `json:"error"`
}

type GetMessageTokenInput struct {
	Email string `json:"email"`
}

type GetMessageTokenOutput struct {
	Jwt   *string `json:"jwt"`
	Error *Error  `json:"error"`
}

type GetUserInput struct {
	ID        int        `json:"id"`
	GettingAt *time.Time `json:"gettingAt"`
}

type GetUserOutput struct {
	User      *User  `json:"user"`
	IsUpdated *bool  `json:"isUpdated"`
	Error     *Error `json:"error"`
}

type UpdateUserCredentialsInput struct {
	Jwt      string `json:"jwt"`
	Password string `json:"password"`
}

type UpdateUserCredentialsOutput struct {
	Ok    bool   `json:"ok"`
	Error *Error `json:"error"`
}

type UpdateUserInput struct {
	ID            int    `json:"id"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	AvatarImageID *int   `json:"avatarImageId"`
}

type UpdateUserOutput struct {
	User  *User  `json:"user"`
	Error *Error `json:"error"`
}

type User struct {
	ID            int       `json:"id"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	AvatarImageID *int      `json:"avatarImageId"`
	BirthDate     time.Time `json:"birthDate"`
	Gender        Gender    `json:"gender"`
	Status        Status    `json:"status"`
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

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)

var AllGender = []Gender{
	GenderMale,
	GenderFemale,
}

func (e Gender) IsValid() bool {
	switch e {
	case GenderMale, GenderFemale:
		return true
	}
	return false
}

func (e Gender) String() string {
	return string(e)
}

func (e *Gender) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Gender(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Gender", str)
	}
	return nil
}

func (e Gender) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Role string

const (
	RoleUser    Role = "USER"
	RoleService Role = "SERVICE"
)

var AllRole = []Role{
	RoleUser,
	RoleService,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleUser, RoleService:
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

type Status string

const (
	StatusDeleted      Status = "DELETED"
	StatusBlocked      Status = "BLOCKED"
	StatusConfirmed    Status = "CONFIRMED"
	StatusNotConfirmed Status = "NOT_CONFIRMED"
)

var AllStatus = []Status{
	StatusDeleted,
	StatusBlocked,
	StatusConfirmed,
	StatusNotConfirmed,
}

func (e Status) IsValid() bool {
	switch e {
	case StatusDeleted, StatusBlocked, StatusConfirmed, StatusNotConfirmed:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

func (e *Status) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

func (e Status) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
