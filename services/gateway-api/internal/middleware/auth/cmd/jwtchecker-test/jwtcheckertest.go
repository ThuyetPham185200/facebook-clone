package main

import (
	"fmt"
	"net/http"
	"time"

	jwt_checker "auth/pkg/jwt_checker"
)

func main() {
	// Tạo strategy dùng HS256
	strategy := &jwt_checker.HS256Strategy{
		SecretKey: "supersecret",
	}

	// Tạo middleware checker
	jwtCheck := &jwt_checker.JWTChecker{Strategy: strategy}

	// Route public (không check JWT)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to public route"))
	})

	// Route protected (check JWT)
	http.Handle("/protected", jwtCheck.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value("user").(*jwt_checker.Claims)

		fmt.Fprintf(w, "UserID: %s\n", claims.UserID)
		fmt.Fprintf(w, "Email: %s\n", claims.Email)
		fmt.Fprintf(w, "ID (jti): %s\n", claims.ID)
		fmt.Fprintf(w, "IssuedAt: %s\n", claims.IssuedAt.Time.Format(time.RFC3339))
		fmt.Fprintf(w, "ExpiresAt: %s\n", claims.ExpiresAt.Time.Format(time.RFC3339))
		fmt.Fprintf(w, "Server sent data to client!")
	})))

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
