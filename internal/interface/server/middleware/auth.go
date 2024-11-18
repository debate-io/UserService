package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/debate-io/service-auth/internal/domain/model"
	"github.com/debate-io/service-auth/internal/infrastructure/auth"
)

type authKey string

const (
	JwtClaimsKey authKey = "jwtClaims"
)

func AuthMiddleware(auth *auth.AuthService) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			jwt := r.Header.Get("Authorization")
			log.Println("JWT", jwt)
			if jwt == "" {
				r = r.WithContext(context.WithValue(r.Context(), JwtClaimsKey, (*model.Claims)(nil)))
				h.ServeHTTP(w, r)
				return
			}
			tokenString := strings.Split(jwt, "Bearer ")[1]
			if tokenString == "" {
				http.Error(w, "Bearer token not found in Authorization header", http.StatusUnauthorized)
				return
			}

			claims, err := auth.ParseToken(tokenString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), JwtClaimsKey, claims))
			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
