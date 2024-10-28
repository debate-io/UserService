package resolvers

import (
	"context"

	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
)

func (q queryResolver) AuthenticateUser(
	ctx context.Context,
	input gen.AuthenticateUserInput,
) (*gen.AuthenticateUserOutput, error) {

	output, err := q.useCases.Users.AuthenticateUser(ctx, input)
	if err != nil {
		return nil, NewResolverError("can't authenticate user", err)
	}

	return output, nil
}

func (q queryResolver) VerifyRecoveryCode(
	ctx context.Context,
	input gen.VerifyRecoveryCodeInput,
) (*gen.VerifyRecoveryCodeOutput, error) {

	output, err := q.useCases.Users.VerifyRecoveryCode(ctx, input)
	if err != nil {
		return nil, NewResolverError("can't verify recovery code", err)
	}

	return output, nil
}

func (q queryResolver) GetUser(ctx context.Context, input gen.GetUserInput) (*gen.GetUserOutput, error) {
	output, err := q.useCases.Users.GetUser(ctx, input)
	if err != nil {
		return nil, NewResolverError("can't get user", err)
	}

	return output, nil
}

func (q queryResolver) GetGamesStats(ctx context.Context, input gen.GetGamesStatsInput) (*gen.GetGamesStatsOutput, error) {
	output, err := q.useCases.Users.GetGamesStats(ctx, input)
	if err != nil {
		return nil, NewResolverError("can't get user's game stat", err)
	}

	return output, nil
}

func (q queryResolver) GetUserAchievements(ctx context.Context, userId int, limit int, offset int) ([]*gen.Achievement, error) {
	output, err := q.useCases.Users.GetAchievmentsByUserId(ctx, userId, limit, offset)
	if err != nil {
		return nil, NewResolverError("can't get user's achievs", err)
	}

	return output, nil
}
