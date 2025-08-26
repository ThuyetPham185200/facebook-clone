package session

import (
	"authservice/internal/infra/store"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ---- Interface ----
type SessionManager interface {
	CreateSession(userID int64) (accessToken, refreshToken string, err error)
	RefreshToken(refreshToken string) (newAccess, newRefresh string, err error)
	Logout(userID int64, refreshToken string) error
	LogoutAll(userID int64) error
}

// ---- JWT Config ----
type JwtConfig struct {
	AccessSecret  []byte
	RefreshSecret []byte
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

// ---- Implementation ----
type sessionManager struct {
	cfg   *JwtConfig
	store store.RefreshTokenStore // backend to persist refresh tokens
}

// ---- Constructor ----
func NewSessionManager(cfg *JwtConfig, store store.RefreshTokenStore) SessionManager {
	return &sessionManager{cfg: cfg, store: store}
}

// ---- CreateSession ----
func (sm *sessionManager) CreateSession(userID int64) (string, string, error) {
	// create access token
	accessToken, err := sm.generateToken(userID, sm.cfg.AccessSecret, sm.cfg.AccessTTL)
	if err != nil {
		return "", "", err
	}

	// create refresh token
	refreshToken, err := sm.generateToken(userID, sm.cfg.RefreshSecret, sm.cfg.RefreshTTL)
	if err != nil {
		return "", "", err
	}

	// persist refresh token
	if err := sm.store.Save(string(userID), refreshToken, sm.cfg.RefreshTTL); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ---- RefreshToken ----
func (sm *sessionManager) RefreshToken(refreshToken string) (string, string, error) {
	// validate refresh token
	_, userID, err := sm.parseToken(refreshToken, sm.cfg.RefreshSecret)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// check if token exists in store
	ok, err := sm.store.Exists(string(userID), refreshToken)
	if err != nil || !ok {
		return "", "", errors.New("refresh token revoked or not found")
	}

	// issue new tokens
	newAccess, newRefresh, err := sm.CreateSession(userID)
	if err != nil {
		return "", "", err
	}

	// delete old refresh token
	if err := sm.store.Delete(string(userID), refreshToken); err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}

// ---- Logout single session ----
func (sm *sessionManager) Logout(userID int64, refreshToken string) error {
	return sm.store.Delete(string(userID), refreshToken)
}

// ---- Logout all sessions ----
func (sm *sessionManager) LogoutAll(userID int64) error {
	return sm.store.DeleteAll(string(userID))
}

// ---- Helpers ----
func (sm *sessionManager) generateToken(userID int64, secret []byte, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (sm *sessionManager) parseToken(tokenStr string, secret []byte) (jwt.MapClaims, int64, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, 0, errors.New("invalid claims")
	}

	uid, ok := claims["user_id"].(float64)
	if !ok {
		return nil, 0, errors.New("invalid user_id claim")
	}

	return claims, int64(uid), nil
}

// ---- Simple In-Memory RefreshTokenStore (for demo/testing)
