package store

import (
	dbclient "authservice/internal/infra/postgresclient"
	"time"
)

// RefreshTokenStore interface cho việc quản lý refresh token
type RefreshTokenStore interface {
	Save(userID string, refreshToken string, ttl time.Duration) error
	Exists(userID string, refreshToken string) (bool, error)
	Delete(userID string, refreshToken string) error
	DeleteAll(userID string) error
}

// ===================== Postgres Implementation (Production) =====================

type postgresTokenStore struct {
	DB *dbclient.PostgresClient
}

func NewPostgresTokenStore(posgresconfig *PostGresConfig) RefreshTokenStore {
	return &postgresTokenStore{
		DB: dbclient.NewPostgresClient(posgresconfig.Host, posgresconfig.Port, posgresconfig.User, posgresconfig.Password, posgresconfig.DBname),
	}
}

func (s *postgresTokenStore) Save(userID string, refreshToken string, ttl time.Duration) error {
	query := `
		INSERT INTO sessions (user_id, refresh_token, refresh_expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
		ON CONFLICT (refresh_token) DO UPDATE
		SET refresh_expires_at = EXCLUDED.refresh_expires_at,
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
			WHERE user_id = $1 AND refresh_token = $2 AND refresh_expires_at > now()
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
