package resolvers

import (
	"context"

	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
)

func (q queryResolver) AuthenticateUser(
	ctx context.Context,
	input gen.AuthenticateUserInput,
) (*gen.AuthenticateUserOutput, error) {
	panic("not implemented")
	/*
		 	output, err := q.useCases.Users.AuthenticateUser(ctx, input)
			if err != nil {
				return nil, NewResolverError("can't authenticate user", err)
			}

			return output, nil
	*/
}

/* func (q queryResolver) GetUser(ctx context.Context, input gen.GetUserInput) (*gen.GetUserOutput, error) {
	output, err := q.useCases.Users.GetUser(ctx, input)
	if err != nil {
		return nil, NewResolverError("can't get user", err)
	}

	return output, nil
} */
