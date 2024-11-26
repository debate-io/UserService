package usecases

import (
	"context"
	"errors"

	"math/rand"
	"time"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/domain/repo"
	"github.com/debate-io/service-auth/internal/infrastructure/auth"
	"github.com/debate-io/service-auth/internal/infrastructure/smtp"
	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
	"github.com/debate-io/service-auth/internal/interface/server/middleware"
	"github.com/debate-io/service-auth/internal/usecases/mappers"
	"golang.org/x/crypto/bcrypt"
)

const (
	CodeLength = 6
	CodeTTL    = 5 // in minute
)

type User struct {
	userRepo         repo.UserRepository
	recoveryCodeRepo repo.RecoveryCodeRepository
	gameStatsRepo    repo.GameStatsRepository
	achievementRepo  repo.AchievmentsRepository
	smtpSender       *smtp.Sender
	authService      *auth.AuthService
}

func NewUserUseCases(userRepo repo.UserRepository, recoveryCodeRepo repo.RecoveryCodeRepository, gameStatsRepo repo.GameStatsRepository, achievementRepo repo.AchievmentsRepository, smtpClient *smtp.Sender, authService *auth.AuthService) *User {
	return &User{
		userRepo:         userRepo,
		recoveryCodeRepo: recoveryCodeRepo,
		gameStatsRepo:    gameStatsRepo,
		achievementRepo:  achievementRepo,
		smtpSender:       smtpClient,
		authService:      authService,
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
		if errors.Is(err, repo.ErrAlreadyExist) {
			return &gen.RegisterUserOutput{
				Error: mappers.NewDTOError(gen.ErrorAlreadyExist)}, nil
		}

		return nil, err
	}

	claims, err := model.NewAuthClaims(user.ID, user.Email, user.Role, u.authService.GetDaysAuthExpires())
	if err != nil {
		return nil, err
	}

	jwt, err := u.authService.GenerateTokenByClaims(claims)
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
		if errors.Is(err, repo.ErrNotFound) {
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

	claims, err := model.NewAuthClaims(user.ID, user.Email, user.Role, u.authService.GetDaysAuthExpires())
	if err != nil {
		return nil, err
	}

	jwt, err := u.authService.GenerateTokenByClaims(claims)
	if err != nil {
		return nil, err
	}

	return &gen.AuthenticateUserOutput{Jwt: &jwt}, nil
}

func (u *User) GetUser(
	ctx context.Context,
	input gen.GetUserInput,
) (*gen.GetUserOutput, error) {
	claims := ctx.Value(middleware.JwtClaimsKey).(*model.Claims)
	if claims == nil {
		return &gen.GetUserOutput{
			Error: mappers.NewDTOError(gen.ErrorUnauthorized),
		}, nil
	}

	role := claims.Role
	if role != model.RoleAdmin && role != model.RoleContentManager && role != model.RoleDefaultUser {
		return &gen.GetUserOutput{
			Error: mappers.NewDTOError(gen.ErrorUnauthorized),
		}, nil
	}

	user, err := u.userRepo.FindUserByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return &gen.GetUserOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound),
			}, nil
		}

		return nil, err
	}

	return &gen.GetUserOutput{User: mappers.MapUserToDTO(user)}, nil
}

func (u *User) GetGamesStats(
	ctx context.Context,
	input gen.GetGamesStatsInput,
) (*gen.GetGamesStatsOutput, error) {
	claims := ctx.Value(middleware.JwtClaimsKey).(*model.Claims)
	if claims == nil {
		return &gen.GetGamesStatsOutput{
			Error: mappers.NewDTOError(gen.ErrorUnauthorized),
		}, nil
	}

	role := claims.Role
	if role != model.RoleAdmin && role != model.RoleContentManager && role != model.RoleDefaultUser {
		return &gen.GetGamesStatsOutput{
			Error: mappers.NewDTOError(gen.ErrorUnauthorized),
		}, nil
	}

	stat, err := u.gameStatsRepo.GetTotalGamesStatsByUserId(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	result := &gen.GetGamesStatsOutput{
		GamesAmount: stat.TotalGamesStats.GamesAmount,
		WinsAmount:  stat.TotalGamesStats.WinsAmount,
	}
	if stat.TotalGamesStats.GamesAmount != 0 {
		result.WinsPercents = 100. * float64(stat.TotalGamesStats.WinsAmount) / float64(stat.TotalGamesStats.GamesAmount)
	}

	for metatopic, stat := range stat.MetaTopicStats {
		element := &gen.MetaTopicsStats{
			MetaTopic:   metatopic,
			GamesAmount: stat.GamesAmount,
			WinsAmount:  stat.WinsAmount,
		}

		if stat.GamesAmount != 0 {
			element.WinsPercents = 100. * float64(stat.WinsAmount) / float64(stat.GamesAmount)
		}

		result.MetaTopicsStats = append(result.MetaTopicsStats, element)
	}

	return result, nil
}

func (u *User) UpdateUser(
	ctx context.Context,
	input gen.UpdateUserInput,
) (output *gen.UpdateUserOutput, err error) {
	claims := ctx.Value(middleware.JwtClaimsKey).(*model.Claims)
	if claims == nil || claims.Role != model.RoleAdmin && claims.UserID != input.ID {
		return &gen.UpdateUserOutput{
			Error: mappers.NewDTOError(gen.ErrorUnauthorized),
		}, nil
	}

	user, err := u.userRepo.FindUserByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return &gen.UpdateUserOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound)}, nil
		}

		return nil, err
	}

	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.ImageID != nil {
		user.Image.ID = *input.ImageID
	}
	if input.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}
	if input.Email != nil {
		user.Email = *input.Email
	}

	if err := user.Validate(); err != nil {
		return &gen.UpdateUserOutput{
			Error: mappers.NewDTOError(gen.ErrorValidation)}, nil
	}

	if _, err = u.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return &gen.UpdateUserOutput{User: mappers.MapUserToDTO(user)}, nil
}

func (u *User) RecoveryPassword(ctx context.Context, input gen.RecoveryPasswordInput) (*gen.RecoveryPasswordOutput, error) {
	user, err := u.userRepo.FindUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
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

func (u *User) UpdatePassword(ctx context.Context, input gen.UpdatePasswordInput) (*gen.UpdatePasswordOutput, error) {
	claims := ctx.Value(middleware.JwtClaimsKey).(*model.Claims)
	if claims == nil || claims.Role != model.RoleAdmin && claims.UserID != input.ID {
		return &gen.UpdatePasswordOutput{
			Error: mappers.NewDTOError(gen.ErrorUnauthorized),
		}, nil
	}

	user, err := u.userRepo.FindUserByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return &gen.UpdatePasswordOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound)}, nil
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		return &gen.UpdatePasswordOutput{
			Error: mappers.NewDTOError(gen.ErrorInvalidCredentials)}, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	if _, err := u.userRepo.UpdateUser(ctx, user); err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return &gen.UpdatePasswordOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound)}, nil
		}
		return &gen.UpdatePasswordOutput{}, err
	}

	return &gen.UpdatePasswordOutput{}, nil
}

func (u *User) UpdateEmail(ctx context.Context, input gen.UpdateEmailInput) (*gen.UpdateEmailOutput, error) {
	claims := ctx.Value(middleware.JwtClaimsKey).(*model.Claims)
	if claims == nil || claims.Role != model.RoleAdmin && claims.UserID != input.ID {
		return &gen.UpdateEmailOutput{
			Error: mappers.NewDTOError(gen.ErrorUnauthorized),
		}, nil
	}

	user, err := u.userRepo.FindUserByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return &gen.UpdateEmailOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound)}, nil
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return &gen.UpdateEmailOutput{
			Error: mappers.NewDTOError(gen.ErrorInvalidCredentials)}, nil
	}

	user.Email = input.Email
	if _, err := u.userRepo.UpdateUser(ctx, user); err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return &gen.UpdateEmailOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound)}, nil
		}
		return &gen.UpdateEmailOutput{}, err
	}

	return &gen.UpdateEmailOutput{}, nil
}

func (u *User) GetAchievmentsByUserId(ctx context.Context, userId int, limit int, offset int) ([]*gen.Achievement, error) {
	achievs, err := u.achievementRepo.GetAchievmentsByUserId(ctx, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	var result []*gen.Achievement
	for _, achiev := range achievs {
		element := &gen.Achievement{
			ID:          achiev.ID,
			Name:        achiev.Name,
			Description: achiev.Description,
			CreatedAt:   achiev.CreateAt,
		}
		result = append(result, element)
	}
	return result, nil
}

func (u *User) ResetPassword(ctx context.Context, input gen.ResetPasswordInput) (*gen.ResetPasswordOutput, error) {
	code, err := u.recoveryCodeRepo.FindRecoveryCodeByEmailAndCode(ctx, input.Email, input.Code)
	if err != nil {
		if errors.Is(err, repo.ErrAlreadyExist) {
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
		if errors.Is(err, repo.ErrNotFound) {
			return &gen.ResetPasswordOutput{
				Error: mappers.NewDTOError(gen.ErrorNotFound)}, nil
		}

		return nil, err
	}

	return &gen.ResetPasswordOutput{}, nil
}

func generateCode(length int) string {
	code := make([]byte, length)

	for i := range code {
		code[i] = byte(rand.Intn(10) + '0')
	}

	return string(code)
}
