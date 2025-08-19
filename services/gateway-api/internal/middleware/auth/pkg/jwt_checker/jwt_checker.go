package jwt_checker

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// JWTChecker is a middleware that uses a JWTStrategy to verify tokens.
type JWTChecker struct {
	Strategy JWTStrategy
}

func NewJWTChecker() *JWTChecker {
	// Tạo strategy dùng HS256
	strategy := &HS256Strategy{
		SecretKey: "supersecret",
	}
	// Tạo middleware checker
	jwtCheck := &JWTChecker{Strategy: strategy}
	return jwtCheck
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

func (j *JWTChecker) TokenCheck(authHeader string) bool {
	if authHeader == "" {
		fmt.Println("Missing Authorization header")
		return false
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		fmt.Println("Invalid Authorization format")
		return false
	}

	_, err := j.Strategy.Verify(parts[1])
	if err != nil {
		fmt.Println("Invalid token")
		return false
	}
	return true
}
