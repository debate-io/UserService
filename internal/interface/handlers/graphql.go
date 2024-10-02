package handlers

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/inhies/go-bytesize"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/zap"
)

const (
	maxUploadSize = 100 * bytesize.MB
	maxMemorySize = 1 * bytesize.MB
	cacheQuery    = 1000
)

func NewGraphqlHandler(
	logger *zap.Logger,
	schema graphql.ExecutableSchema,
	isDebug bool,
) *handler.Server {
	srv := handler.New(schema)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{
		MaxUploadSize: int64(maxUploadSize),
		MaxMemory:     int64(maxMemorySize),
	})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](cacheQuery))
	srv.Use(extension.Introspection{})

	if isDebug {
		srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			if rc := graphql.GetOperationContext(ctx); rc != nil {
				logger.Debug("graphql query information",
					zap.String("query", rc.RawQuery),
					zap.Reflect("variables", rc.Variables),
				)
			}

			return next(ctx)
		})
	}

	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) (userMessage error) {
		logger.Error("recover graphql", zap.Reflect("err:", err))

		return fmt.Errorf("internal server error")
	})

	return srv
}
