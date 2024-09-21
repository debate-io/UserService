package resolvers

import (
	"fmt"

	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
	"github.com/debate-io/service-auth/internal/registry"
	"github.com/ztrue/tracerr"
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

// Mutation returns gen.MutationResolver implementation.
func (r *Resolver) Mutation() gen.MutationResolver { return &mutationResolver{r} }

// Query returns gen.QueryResolver implementation.
func (r *Resolver) Query() gen.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

func NewResolverError(
	responseError string,
	originalError error,
) *ResolverError {
	fmt.Println(originalError)

	return &ResolverError{
		tracerr.Errorf(responseError),
		originalError,
	}
}

type ResolverError struct {
	ResponseError error
	OriginalError error
}

func (e *ResolverError) Error() string {
	return e.ResponseError.Error()
}
