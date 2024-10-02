package resolvers

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"

	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
	"github.com/debate-io/service-auth/internal/registry"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	useCases *registry.UseCases
}

func NewResolver(
	useCases *registry.UseCases,
) *Resolver {
	return &Resolver{
		useCases,
	}
}

// AuthenticateUser is the resolver for the authenticateUser field.
// это нахуй вынести в 1 из 2х файлов выше
func (r *queryResolver) AuthenticateUser(ctx context.Context, input gen.AuthenticateUserInput) (*gen.AuthenticateUserOutput, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

// Query returns gen.QueryResolver implementation.
func (r *Resolver) Query() gen.QueryResolver { return &queryResolver{r} }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
type Resolver struct{}
*/

//старое ниже

/*

type mutationResolver struct{ *Resolver }
func (r *Resolver) Mutation() gen.MutationResolver { return &mutationResolver{r} }


func NewResolverError(
	responseError string,
	originalError error,
) *ResolverError {
	fmt.Println(originalError)
	return &ResolverError{
	}
}

type ResolverError struct {
	ResponseError error
	OriginalError error
}

func (e *ResolverError) Error() string {
	return e.ResponseError.Error()
} */
