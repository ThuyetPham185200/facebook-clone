package jwt_checker

import (
	"context"
	"net/http"
	"strings"
)

// JWTChecker is a middleware that uses a JWTStrategy to verify tokens.
type JWTChecker struct {
	Strategy JWTStrategy
}

func (j *JWTChecker) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
			return
		}

		claims, err := j.Strategy.Verify(parts[1])
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
