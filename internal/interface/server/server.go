package server

import (
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"net/http"
	"time"

	"github.com/debate-io/service-auth/internal/interface/graphql/gen"
	"github.com/debate-io/service-auth/internal/interface/graphql/resolvers"
	"github.com/debate-io/service-auth/internal/interface/handlers"
	"github.com/debate-io/service-auth/internal/registry"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"gopkg.in/tylerb/graceful.v1"
)

const (
	MaxHTTPAge       = 1000
	shutdownTimeout  = 10 * time.Second
	readWriteTimeout = 90 * time.Second
)

type Server struct {
	router *chi.Mux
	logger *zap.Logger
}

func (s *Server) GetRouter() *chi.Mux {
	return s.router
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{
		router: chi.NewRouter(),
		logger: logger,
	}
}

func (s *Server) InitMiddlewares(isDebug bool) {
	s.router.Use(render.SetContentType(render.ContentTypeJSON))
	s.router.Use(middleware.Recoverer)

	if isDebug {
		s.router.Use(middleware.Logger)
	}

	options := cors.Options{
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           MaxHTTPAge,
		AllowedOrigins:   []string{"*"},
	}

	cors := cors.New(options)
	s.router.Use(cors.Handler)
}

func (s *Server) InitRoutes(container *registry.Container, isDebug bool) {
	graphqlHandler := handlers.NewGraphqlHandler(
		s.logger,
		gen.NewExecutableSchema(
			gen.Config{
				Resolvers: resolvers.NewResolver(
					container.UseCases,
				),
			},
		),
		isDebug,
	)

	graphqlHandler.Use(extension.Introspection{})

	s.router.Handle("/*", graphqlHandler)
}

func (s *Server) ListenAndServe(address string, shutdownInitiated func()) error {
	s.logger.Info("Начинаем внимательно слушать", zap.String("addr", address))

	srv := &graceful.Server{
		Timeout:           shutdownTimeout,
		ShutdownInitiated: shutdownInitiated,
		Server: &http.Server{
			ReadTimeout:  readWriteTimeout,
			WriteTimeout: readWriteTimeout,
			Addr:         address,
			Handler:      s.router,
		},
	}

	return srv.ListenAndServe()
}
