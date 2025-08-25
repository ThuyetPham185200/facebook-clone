package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// ---- Interface ----
type PasswordResetManager interface {
	RequestReset(email string) (string, error)                 // generate reset token & send to user
	VerifyResetToken(token string) (userID int64, err error)  // verify token
	ResetPassword(token string, newPassword string) error     // reset password
}

// ---- Implementation ----
type passwordResetManager struct {
	credStore       CredentialsStore
	resetTokenStore ResetTokenStore
	sessionManager  SessionManager
}

// ---- Constructor ----
func NewPasswordResetManager(
	credStore CredentialsStore,
	resetStore ResetTokenStore,
	sessionMgr SessionManager,
) PasswordResetManager {
	return &passwordResetManager{
		credStore:       credStore,
		resetTokenStore: resetStore,
		sessionManager:  sessionMgr,
	}
}

// ---- RequestReset ----
func (pm *passwordResetManager) RequestReset(email string) (string, error) {
	cred, err := pm.credStore.GetByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// generate secure token
	token, err := generateSecureToken(32)
	if err != nil {
		return "", err
	}

	// save token with TTL (15m)
	if err := pm.resetTokenStore.Save(cred.ID, token, 15*time.Minute); err != nil {
		return "", err
	}

	// In thực tế: send email với link reset
	// vd: https://example.com/reset-password?token=xxxx
	return token, nil
}

// ---- VerifyResetToken ----
func (pm *passwordResetManager) VerifyResetToken(token string) (int64, error) {
	return pm.resetTokenStore.GetUserID(token)
}

// ---- ResetPassword ----
func (pm *passwordResetManager) ResetPassword(token string, newPassword string) error {
	userID, err := pm.resetTokenStore.GetUserID(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	// hash mật khẩu mới
	hashed, salt, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	// update credential
	if err := pm.credStore.UpdatePassword(userID, hashed, salt); err != nil {
		return err
	}

	// revoke all sessions for security
	if err := pm.sessionManager.LogoutAll(userID); err != nil {
		return err
	}

	// delete token (one-time use)
	if err := pm.resetTokenStore.Delete(token); err != nil {
		return err
	}

	return nil
}

// ---- Helpers ----
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

//
// ---- In-Memory ResetTokenStore (demo)
// In production: use Redis or Postgres
//
type ResetTokenStore interface {
	Save(userID int64, token string, ttl time.Duration) error
	GetUserID(token string) (int64, error)
	Delete(token string) error
}

type inMemoryResetTokenStore struct {
	mu     sync.Mutex
	tokens map[string]struct {
		userID int64
		exp    time.Time
	}
}

func NewInMemoryResetTokenStore() ResetTokenStore {
	return &inMemoryResetTokenStore{
		tokens: make(map[string]struct {
			userID int64
			exp    time.Time
		}),
	}
}

func (s *inMemoryResetTokenStore) Save(userID int64, token string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tokens[token] = struc
