package store

import (
	dbclient "authservice/internal/infra/postgresclient"
	"authservice/internal/infra/redisclient"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type CredentialsStore struct {
	DBclient    *dbclient.PostgresClient
	RedisClient *redisclient.RedisClient
}

type PostGresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DBNumber int
}

func NewCredentialsStore(posgresconfig *PostGresConfig, redisconfig *RedisConfig) *CredentialsStore {
	return &CredentialsStore{
		DBclient:    dbclient.NewPostgresClient(posgresconfig.Host, posgresconfig.Port, posgresconfig.User, posgresconfig.Password, posgresconfig.DBname),
		RedisClient: redisclient.InitSingleton(redisconfig.Host+":"+redisconfig.Port, redisconfig.Password, redisconfig.DBNumber),
	}
}

// ExistsUser kiểm tra username hoặc email có tồn tại không
func (c *CredentialsStore) ExistsUser(username, email string) (bool, error) {
	if c.RedisClient.GetClient() == nil {
		return false, fmt.Errorf("[CredentialsStore] redis client not initialized")
	}

	// Check username
	if username != "" {
		key := "user:" + username
		exists, err := c.RedisClient.KeyExists(key)
		if err != nil {
			log.Printf("[CredentialsStore] error checking username=%s: %v", username, err)
			return false, err
		}
		log.Printf("[CredentialsStore] checked username=%s, exists=%v", username, exists)
		if exists {
			return true, nil
		}
	}

	// Check email
	if email != "" {
		key := "email:" + email
		exists, err := c.RedisClient.KeyExists(key)
		if err != nil {
			log.Printf("[CredentialsStore] error checking email=%s: %v", email, err)
			return false, err
		}
		log.Printf("[CredentialsStore] checked email=%s, exists=%v", email, exists)
		if exists {
			return true, nil
		}
	}

	log.Printf("[CredentialsStore] username=%s and email=%s not found in cache", username, email)
	return false, nil
}

// Save lưu thông tin credential (auth info) cho user
func (c *CredentialsStore) Save(userID, username, email, hashed string) error {
	if c.DBclient.DB == nil {
		return fmt.Errorf("[CredentialsStore] database client not initialized")
	}
	if c.RedisClient.GetClient() == nil {
		return fmt.Errorf("[CredentialsStore] redis client not initialized")
	}

	newUUID := uuid.New().String() // generate id for credentials

	query := `
		INSERT INTO credentials (
			id, user_id, password_hash, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()
	status := "active"

	_, err := c.DBclient.DB.Exec(query, newUUID, userID, hashed, status, now, now)
	if err != nil {
		return fmt.Errorf("[CredentialsStore] failed to insert credential: %w", err)
	}

	// Update Redis cache
	if c.RedisClient != nil && c.RedisClient.GetClient() != nil {
		// Cache username
		if username != "" {
			if err := c.RedisClient.SetString("user:"+username, userID, 24*time.Hour); err != nil {
				fmt.Printf("[CredentialsStore] failed to cache username: %v\n", err)
			}
		}

		// Cache email
		if email != "" {
			if err := c.RedisClient.SetString("email:"+email, userID, 24*time.Hour); err != nil {
				fmt.Printf("[CredentialsStore] failed to cache email: %v\n", err)
			}
		}
	}

	return nil
}
