package store

import (
	dbclient "authservice/internal/infra/postgresclient"
	"database/sql"
	"fmt"
	"time"
)

type CredentialsStore struct {
	DBclient *dbclient.PostgresClient
}

func NewCredentialsStore(host, port, user, password, dbname string) *CredentialsStore {
	return &CredentialsStore{
		DBclient: dbclient.NewPostgresClient(host, port, user, password, dbname),
	}
}

// Exists kiểm tra username hoặc email có tồn tại không
func (c *CredentialsStore) Exists(username, email string) (bool, error) {
	if c.DBclient.DB == nil {
		return false, fmt.Errorf("[CredentialsStore] database client not initialized")
	}

	query := `
		SELECT EXISTS (
			SELECT 1 FROM credentials
			WHERE username = $1 OR email = $2
		)
	`

	var exists bool
	err := c.DBclient.DB.QueryRow(query, username, email).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("[CredentialsStore] query failed: %w", err)
	}

	return exists, nil
}

// Save lưu thông tin user mới vào bảng credentials
func (c *CredentialsStore) Save(userID, username, email, hashed string) error {
	if c.DBclient.DB == nil {
		return fmt.Errorf("[CredentialsStore] database client not initialized")
	}

	query := `
		INSERT INTO credentials (
			id, username, email, password_hash, created_at, updated_at
		) VALUES ($1, $2, $3, $4, 'active', $5, $6)
	`

	now := time.Now()

	_, err := c.DBclient.DB.Exec(query,
		userID, username, email, hashed, now, now,
	)
	if err != nil {
		return fmt.Errorf("[CredentialsStore] failed to insert credential: %w", err)
	}

	return nil
}
