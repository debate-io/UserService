package resolvers

import (
	"fmt"

	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
	"github.com/debate-io/service-auth/internal/registry"
	"github.com/ztrue/tracerr"
)

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

type queryResolver struct{ *Resolver }

type mutationResolver struct{ *Resolver }

func (r *Resolver) Query() gen.QueryResolver { return &queryResolver{r} }

func (r *Resolver) Mutation() gen.MutationResolver { return &mutationResolver{r} }

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
