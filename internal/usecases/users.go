package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
	"github.com/debate-io/service-auth/internal/usecases/mappers"
	"github.com/debate-io/service-auth/internal/usecases/types"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	RoleUserID  = 1
	RoleUserKey = "USER"
)

type User struct {
	userRepo   repo.UserRepository
	roleRepo   repo.RoleRepository
	jwtConfigs JwtConfigs
}

type JwtConfigs struct {
	jwtSecretAuth       string
	jwtSecretMessages   string
	daysAuthExpires     int
	daysRecoveryExpires int
}

func NewUserUseCases(userRepo repo.UserRepository, roleRepo repo.RoleRepository, jwtConfigs JwtConfigs) *User {
	return &User{userRepo: userRepo, roleRepo: roleRepo, jwtConfigs: jwtConfigs}
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
	input gen.CreateUserInput,
) (*gen.CreateUserOutput, error) {
	user := &model.User{
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Email:         input.Email,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Username:      input.Username,
		AvatarImageID: input.AvatarImageID,
		BirthDate:     input.BirthDate,
		Gender:        input.Gender,
		RoleID:        RoleUserID,
		Status:        gen.StatusNotConfirmed,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	if err := user.Validate(); err != nil {
		return &gen.CreateUserOutput{
			Error: mappers.NewDTOError(gen.ErrorValidation)}, nil
	}

	_, err = u.userRepo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, repo.ErrUserAlreadyExist) {
			return &gen.CreateUserOutput{
				Error: mappers.NewDTOError(gen.ErrorAlreadyExist)}, nil
		}

		return nil, err
	}

	claims, err := types.NewAuthClaims(user.ID, user.Email, RoleUserKey, u.jwtConfigs.daysAuthExpires)
	if err != nil {
		return nil, err
	}

	jwt, err := generateTokenByClaims(claims, u.jwtConfigs.jwtSecretAuth)
	if err != nil {
		return nil, err
	}

	return &gen.CreateUserOutput{User: mappers.MapUserToDTO(user), Jwt: &jwt}, nil
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

	role, err := u.roleRepo.FindUserByID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}

	claims, err := types.NewAuthClaims(user.ID, user.Email, role.HeaderKey, u.jwtConfigs.daysAuthExpires)
	if err != nil {
		return nil, err
	}

	token, err := generateTokenByClaims(claims, u.jwtConfigs.jwtSecretAuth)
	if err != nil {
		return nil, err
	}

	return &gen.AuthenticateUserOutput{Jwt: &token}, nil
}

func (u *User) GetClaims(
	ctx context.Context,
	input gen.GetClaimsInput,
) (*gen.GetClaimsOutput, error) {
	token, err := jwt.ParseWithClaims(input.Jwt, &types.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtConfigs.jwtSecretAuth), nil
	})
	if err != nil {
		return &gen.GetClaimsOutput{
			Error: mappers.NewDTOError(gen.ErrorValidation),
		}, nil
	}

	claims, ok := token.Claims.(*types.Claims)
	if err := claims.Valid(); err != nil || !ok {
		return &gen.GetClaimsOutput{
			Error: mappers.NewDTOError(gen.ErrorValidation)}, nil
	}

	return &gen.GetClaimsOutput{Claims: mappers.MapClaimsToDTO(claims)}, nil
}

func (u *User) GetUser(
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

func (u *User) FindUsers(ctx context.Context, input gen.FindUsersInput) (*gen.FindUsersOutput, error) {
	res, err := u.userRepo.FindUsers(ctx, &repo.FindUsersQuery{IDAnyOf: input.IDAnyOf})
	if err != nil {
		return nil, err
	}

	return &gen.FindUsersOutput{Users: mappers.MapUsersToDTO(res)}, nil
}

func (u *User) GetMessageToken(
	ctx context.Context,
	input gen.GetMessageTokenInput,
) (*gen.GetMessageTokenOutput, error) {
	user, err := u.userRepo.FindUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return &gen.GetMessageTokenOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound)}, nil
		}

		return nil, err
	}

	claims, err := types.NewRecoveryClaims(user.ID, u.jwtConfigs.daysRecoveryExpires)
	if err != nil {
		return nil, err
	}

	token, err := generateTokenByClaims(claims, u.jwtConfigs.jwtSecretMessages)
	if err != nil {
		return nil, err
	}

	return &gen.GetMessageTokenOutput{Jwt: &token}, nil
}

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

func generateTokenByClaims(claims *types.Claims, secret string) (string, error) {
	signBytes := []byte(secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(signBytes)
	if err != nil {
		return "", err
	}

	return ss, nil
}
