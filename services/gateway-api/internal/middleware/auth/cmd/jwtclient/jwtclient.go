package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func main() {
	secret := []byte("supersecret")
	claims := MyClaims{
		UserID: "12345",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "12345",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(secret)

	fmt.Println("Generated Token:", tokenString)

	// Chuẩn bị lệnh curl
	url := "http://localhost:8080/protected"
	cmd := exec.Command("curl", "-H", "Authorization: Bearer "+tokenString, url)

	// Lấy output từ curl
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running curl:", err)
		return
	}

	fmt.Println("Response from server:")
	fmt.Println(string(output))
}
