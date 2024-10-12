package usecases

import (
	"context"
	"errors"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/debate-io/service-auth/internal/infrastructure/smtp"
	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
	"github.com/debate-io/service-auth/internal/usecases/mappers"
	"github.com/debate-io/service-auth/internal/usecases/types"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

const (
	CodeLength = 6
	CodeTTL    = 5 // in minute
)

type User struct {
	userRepo         repo.UserRepository
	recoveryCodeRepo repo.RecoveryCodeRepository
	smtpSender       *smtp.Sender
	jwtConfigs       JwtConfigs
}

type JwtConfigs struct {
	jwtSecretAuth       string
	jwtSecretMessages   string
	daysAuthExpires     int
	daysRecoveryExpires int
}

func NewUserUseCases(userRepo repo.UserRepository, recoveryCodeRepo repo.RecoveryCodeRepository, smtpClient *smtp.Sender, jwtConfigs JwtConfigs) *User {
	return &User{userRepo: userRepo, recoveryCodeRepo: recoveryCodeRepo, smtpSender: smtpClient, jwtConfigs: jwtConfigs}
}

func NewJwtConfigsUseCases(jwtSecretAuth string, jwtSecretMessage string, daysAuthExpires int, daysRecoveryExpires int) *JwtConfigs {
	return &JwtConfigs{
		jwtSecretAuth:       jwtSecretAuth,
		jwtSecretMessages:   jwtSecretMessage,
		daysAuthExpires:     daysAuthExpires,
		daysRecoveryExpires: daysRecoveryExpires,
	}
}

func (u *User) CreateUser(
	ctx context.Context,
	input gen.RegisterUserInput,
) (*gen.RegisterUserOutput, error) {
	user := &model.User{
		Email:     input.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Username:  input.Username,
		Password:  input.Password,
		Image:     nil,
		Role:      model.RoleDefaultUser,
	}

	if err := user.Validate(); err != nil {
		return &gen.RegisterUserOutput{
			Error: mappers.NewDTOError(gen.ErrorValidation)}, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	_, err = u.userRepo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, repo.ErrUserAlreadyExist) {
			return &gen.RegisterUserOutput{
				Error: mappers.NewDTOError(gen.ErrorAlreadyExist)}, nil
		}

		return nil, err
	}

	claims, err := types.NewAuthClaims(user.ID, user.Email, string(user.Role), u.jwtConfigs.daysAuthExpires)
	if err != nil {
		return nil, err
	}

	jwt, err := generateTokenByClaims(claims, u.jwtConfigs.jwtSecretAuth)
	if err != nil {
		return nil, err
	}

	return &gen.RegisterUserOutput{User: mappers.MapUserToDTO(user), Jwt: &jwt}, nil
}

func (u *User) AuthenticateUser(
	ctx context.Context,
	input gen.AuthenticateUserInput,
) (*gen.AuthenticateUserOutput, error) {
	user, err := u.userRepo.FindUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return &gen.AuthenticateUserOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound)}, nil
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return &gen.AuthenticateUserOutput{
			Error: mappers.NewDTOError(gen.ErrorInvalidCredentials)}, nil
	}

	claims, err := types.NewAuthClaims(user.ID, user.Email, string(user.Role), u.jwtConfigs.daysAuthExpires)
	if err != nil {
		return nil, err
	}

	token, err := generateTokenByClaims(claims, u.jwtConfigs.jwtSecretAuth)
	if err != nil {
		return nil, err
	}

	return &gen.AuthenticateUserOutput{Jwt: &token}, nil
}

/* func (u *User) GetUser(
	ctx context.Context,
	input gen.GetUserInput,
) (*gen.GetUserOutput, error) {
	user, err := u.userRepo.FindUserByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return &gen.GetUserOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound),
			}, nil
		}

		return nil, err
	}

	if input.GettingAt != nil {
		isUpdated := false
		if input.GettingAt.Before(user.UpdatedAt) {
			isUpdated = true
		}

		return &gen.GetUserOutput{
			IsUpdated: &isUpdated,
		}, nil
	}

	output := &gen.GetUserOutput{}

	if input.GettingAt != nil {
		isUpdated := false
		if user.UpdatedAt.After(*input.GettingAt) {
			isUpdated = true
		}

		output.IsUpdated = &isUpdated
	}

	output.User = mappers.MapUserToDTO(user)

	return output, nil
}
*/

/*
func (u *User) UpdateUserCredentials(
	ctx context.Context,
	input gen.UpdateUserCredentialsInput,
) (*gen.UpdateUserCredentialsOutput, error) {
	token, err := jwt.ParseWithClaims(input.Jwt, &types.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtConfigs.jwtSecretMessages), nil
	})
	if err != nil {
		return &gen.UpdateUserCredentialsOutput{Error: mappers.NewDTOError(gen.ErrorValidation)}, nil
	}

	claims, ok := token.Claims.(*types.Claims)
	if err := claims.Valid(); err != nil || !ok {
		return &gen.UpdateUserCredentialsOutput{
			Error: mappers.NewDTOError(gen.ErrorValidation)}, nil
	}

	user, err := u.userRepo.FindUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	_, err = u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &gen.UpdateUserCredentialsOutput{Ok: true}, nil
}
*/
/*
func (u *User) ConfirmUser(
	ctx context.Context,
	input gen.ConfirmUserInput,
) (*gen.ConfirmUserOutput, error) {
	token, err := jwt.ParseWithClaims(input.Jwt, &types.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtConfigs.jwtSecretMessages), nil
	})
	if err != nil {
		return &gen.ConfirmUserOutput{
			Error: mappers.NewDTOError(gen.ErrorValidation)}, nil
	}

	claims, ok := token.Claims.(*types.Claims)
	if err := claims.Valid(); err != nil || !ok {
		return &gen.ConfirmUserOutput{
			Error: mappers.NewDTOError(gen.ErrorValidation)}, nil
	}

	user, err := u.userRepo.FindUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	user.Status = gen.StatusConfirmed

	_, err = u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &gen.ConfirmUserOutput{Ok: true}, nil
}
*/
/*
func (u *User) UpdateUser(
	ctx context.Context,
	input gen.UpdateUserInput,
) (output *gen.UpdateUserOutput, err error) {
	user, err := u.userRepo.FindUserByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return &gen.UpdateUserOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound)}, nil
		}

		return nil, err
	}

	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.AvatarImageID = input.AvatarImageID

	if err := user.Validate(); err != nil {
		return &gen.UpdateUserOutput{
			Error: mappers.NewDTOError(gen.ErrorValidation)}, nil
	}

	if _, err = u.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return &gen.UpdateUserOutput{User: mappers.MapUserToDTO(user)}, nil
}
*/
/*
func (u *User) DeleteUser(
	ctx context.Context,
	input gen.DeleteUserInput,
) (output *gen.DeleteUserOutput, err error) {
	user, err := u.userRepo.FindUserByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return &gen.DeleteUserOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound)}, nil
		}

		return nil, err
	}

	if user.Status != gen.StatusConfirmed {
		return &gen.DeleteUserOutput{
			Error: mappers.NewDTOError(gen.ErrorValidation)}, nil
	}

	user.Status = gen.StatusDeleted

	if _, err = u.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return &gen.DeleteUserOutput{Ok: true}, nil
}
*/

func (u *User) RecoveryPassword(ctx context.Context, input gen.RecoveryPasswordInput) (*gen.RecoveryPasswordOutput, error) {
	user, err := u.userRepo.FindUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return &gen.RecoveryPasswordOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound)}, nil
		}

		return nil, err
	}

	code := &model.RecoveryCode{
		UserEmail: user.Email,
		User:      user,
		Code:      generateCode(CodeLength),
		ExpiredAt: time.Now().Add(CodeTTL * time.Minute),
	}
	_, err = u.recoveryCodeRepo.CreateRecoveryCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// TODO: Add rendering template HTML message
	err = u.smtpSender.SendPlainMessage("Код для восстановления пароля", code.Code, input.Email)
	if err != nil {
		return nil, err
	}

	return &gen.RecoveryPasswordOutput{}, nil
}

func (u *User) VerifyRecoveryCode(ctx context.Context, input gen.VerifyRecoveryCodeInput) (*gen.VerifyRecoveryCodeOutput, error) {
	exists, err := u.recoveryCodeRepo.ExistsRecoveryCodeByEmailAndCode(ctx, input.Email, input.Code)
	if err != nil {
		return nil, err
	}

	if !exists {
		return &gen.VerifyRecoveryCodeOutput{
			Error: mappers.NewDTOError(gen.ErrorNotFound)}, nil
	}

	return &gen.VerifyRecoveryCodeOutput{}, nil
}

func (u *User) ResetPassword(ctx context.Context, input gen.ResetPasswordInput) (*gen.ResetPasswordOutput, error) {
	code, err := u.recoveryCodeRepo.FindRecoveryCodeByEmailAndCode(ctx, input.Email, input.Code)
	if err != nil {
		if errors.Is(err, repo.ErrRecoveryCodeNotFound) {
			return &gen.ResetPasswordOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound)}, nil
		}

		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	code.User.Password = string(hashedPassword)

	_, err = u.userRepo.UpdateUser(ctx, code.User)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return &gen.ResetPasswordOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound)}, nil
		}

		return nil, err
	}

	return &gen.ResetPasswordOutput{}, nil
}

// дважды перепроверить функцию, но вроде она не затронута
func generateTokenByClaims(claims *types.Claims, secret string) (string, error) {
	signBytes := []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(signBytes)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func generateCode(length int) string {
	code := make([]byte, length)

	for i := range code {
		code[i] = byte(rand.Intn(10) + '0')
	}

	return string(code)
}
