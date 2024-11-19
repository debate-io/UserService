package resolvers

import (
	"context"

	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
)

func (m *mutationResolver) RegisterUser(ctx context.Context, input gen.RegisterUserInput) (*gen.RegisterUserOutput, error) {
	output, err := m.useCases.Users.CreateUser(ctx, input)
	if err != nil {
		return nil, NewResolverError("failed to create user", err)
	}

	return output, nil
}

func (m *mutationResolver) UpdateUser(ctx context.Context, input gen.UpdateUserInput) (*gen.UpdateUserOutput, error) {
	output, err := m.useCases.Users.UpdateUser(ctx, input)
	if err != nil {
		return nil, NewResolverError("failed to update user", err)
	}

	return output, err
}

func (m *mutationResolver) UpdatePassword(ctx context.Context, input gen.UpdatePasswordInput) (*gen.UpdatePasswordOutput, error) {
	output, err := m.useCases.Users.UpdatePassword(ctx, input)
	if err != nil {
		return nil, NewResolverError("failed to update password", err)
	}

	return output, err
}

func (m *mutationResolver) UpdateEmail(ctx context.Context, input gen.UpdateEmailInput) (*gen.UpdateEmailOutput, error) {
	output, err := m.useCases.Users.UpdateEmail(ctx, input)
	if err != nil {
		return nil, NewResolverError("failed to update email", err)
	}

	return output, err
}

func (m mutationResolver) RecoveryPassword(
	ctx context.Context,
	input gen.RecoveryPasswordInput,
) (*gen.RecoveryPasswordOutput, error) {

	output, err := m.useCases.Users.RecoveryPassword(ctx, input)
	if err != nil {
		return nil, NewResolverError("can't recovery password", err)
	}

	return output, nil
}

func (m mutationResolver) ResetPassword(
	ctx context.Context,
	input gen.ResetPasswordInput,
) (*gen.ResetPasswordOutput, error) {

	output, err := m.useCases.Users.ResetPassword(ctx, input)
	if err != nil {
		return nil, NewResolverError("can't reset password", err)
	}

	return output, nil
}
