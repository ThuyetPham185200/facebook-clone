package jwt_checker

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func VerifyHS256(tokenString, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("token expired")
		}
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}

// Placeholder for RS256 verification (load public key, etc.)
func VerifyRS256(tokenString, publicKeyPath string) (*Claims, error) {
	return nil, errors.New("RS256 verification not implemented yet")
}
