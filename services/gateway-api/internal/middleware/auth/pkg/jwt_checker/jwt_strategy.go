package jwt_checker

// JWTStrategy defines a strategy for verifying JWT tokens.
type JWTStrategy interface {
	Verify(tokenString string) (*Claims, error)
}

// HS256Strategy is a concrete strategy using HS256 secret.
type HS256Strategy struct {
	SecretKey string
}

func (h *HS256Strategy) Verify(tokenString string) (*Claims, error) {
	return VerifyHS256(tokenString, h.SecretKey)
}

// RS256Strategy is a concrete strategy using RSA public key.
type RS256Strategy struct {
	PublicKeyPath string
}

func (r *RS256Strategy) Verify(tokenString string) (*Claims, error) {
	return VerifyRS256(tokenString, r.PublicKeyPath)
}
