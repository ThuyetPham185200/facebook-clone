package store

import (
	dbclient "authservice/internal/infra/postgresclient"
	"sync"
	"time"
)

// RefreshTokenStore interface cho việc quản lý refresh token
type RefreshTokenStore interface {
	Save(userID string, refreshToken string, ttl time.Duration) error
	Exists(userID string, refreshToken string) (bool, error)
	Delete(userID string, refreshToken string) error
	DeleteAll(userID string) error
}

// ===================== In-Memory Implementation =====================
type inMemoryTokenStore struct {
	mu     sync.Mutex
	tokens map[string]map[string]time.Time // userID -> token -> expiry
}

func NewInMemoryTokenStore() RefreshTokenStore {
	return &inMemoryTokenStore{
		tokens: make(map[string]map[string]time.Time),
	}
}

func (s *inMemoryTokenStore) Save(userID string, refreshToken string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tokens[userID] == nil {
		s.tokens[userID] = make(map[string]time.Time)
	}
	s.tokens[userID][refreshToken] = time.Now().Add(ttl)
	return nil
}

func (s *inMemoryTokenStore) Exists(userID string, refreshToken string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userTokens, ok := s.tokens[userID]
	if !ok {
		return false, nil
	}

	exp, ok := userTokens[refreshToken]
	if !ok {
		return false, nil
	}

	if time.Now().After(exp) {
		delete(userTokens, refreshToken)
		return false, nil
	}

	return true, nil
}

func (s *inMemoryTokenStore) Delete(userID string, refreshToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tokens[userID]; ok {
		delete(s.tokens[userID], refreshToken)
	}
	return nil
}

func (s *inMemoryTokenStore) DeleteAll(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tokens, userID)
	return nil
}

// ===================== Postgres Implementation (Production) =====================

type postgresTokenStore struct {
	DB *dbclient.PostgresClient
}

func NewPostgresTokenStore(db *dbclient.PostgresClient) RefreshTokenStore {
	return &postgresTokenStore{DB: db}
}

func (s *postgresTokenStore) Save(userID string, refreshToken string, ttl time.Duration) error {
	query := `
		INSERT INTO sessions (user_id, refresh_token, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
		ON CONFLICT (user_id, refresh_token) DO UPDATE
		SET expires_at = EXCLUDED.expires_at,
		    updated_at = now()
	`
	expiry := time.Now().Add(ttl)
	_, err := s.DB.DB.Exec(query, userID, refreshToken, expiry)
	return err
}

func (s *postgresTokenStore) Exists(userID string, refreshToken string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM sessions
			WHERE user_id = $1 AND refresh_token = $2 AND expires_at > now()
		)
	`
	var exists bool
	err := s.DB.DB.QueryRow(query, userID, refreshToken).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *postgresTokenStore) Delete(userID string, refreshToken string) error {
	query := `
		DELETE FROM sessions
		WHERE user_id = $1 AND refresh_token = $2
	`
	_, err := s.DB.DB.Exec(query, userID, refreshToken)
	return err
}

func (s *postgresTokenStore) DeleteAll(userID string) error {
	query := `
		DELETE FROM sessions
		WHERE user_id = $1
	`
	_, err := s.DB.DB.Exec(query, userID)
	return err
}
