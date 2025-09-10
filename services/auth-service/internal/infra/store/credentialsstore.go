package store

import (
	dbclient "authservice/internal/infra/postgresclient"
	"authservice/internal/infra/redisclient"
	"authservice/internal/model"
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
			if err := c.RedisClient.SetKey("user:"+username, userID, 24*time.Hour); err != nil {
				fmt.Printf("[CredentialsStore] failed to cache username: %v\n", err)
			} else {
				_, err := c.RedisClient.KeyExists(username)
				if err != nil {
					fmt.Printf("%s not stored\n", username)
				}
			}

			if err := c.RedisClient.SetKey("user:"+userID, username, 24*time.Hour); err != nil {
				fmt.Printf("[CredentialsStore] failed to cache username: %v\n", err)
			} else {
				_, err := c.RedisClient.KeyExists(userID)
				if err != nil {
					fmt.Printf("%s not stored\n", userID)
				}
			}
		}

		// Cache email
		if email != "" {
			if err := c.RedisClient.SetKey("email:"+email, userID, 24*time.Hour); err != nil {
				fmt.Printf("[CredentialsStore] failed to cache email: %v\n", err)
			}
		}
	}

	return nil
}

func (c *CredentialsStore) GetUserIdByName(username string) (string, error) {
	//

	if c.RedisClient.GetClient() == nil {
		return "", fmt.Errorf("[CredentialsStore] redis client not initialized")
	}
	userid, err := c.RedisClient.GetKey("user:" + username)
	if err != nil {
		fmt.Printf("[CredentialsStore] failed to get username from cache: %v\n", err)
		return "", err
	}
	return userid, nil
}

func (c *CredentialsStore) GetCredentialByUserID(userID string) (*model.Credential, error) {
	query := `
		SELECT id, user_id, password_hash, mfa_secret, status, created_at, updated_at
		FROM credentials WHERE user_id = $1
	`
	row := c.DBclient.DB.QueryRow(query, userID)

	var cre model.Credential
	err := row.Scan(
		&cre.ID,
		&cre.UserID,
		&cre.PasswordHash,
		&cre.MFASecret,
		&cre.Status,
		&cre.CreatedAt,
		&cre.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("[CredentialsStore] failed to scan credential: %w", err)
	}

	return &cre, nil
}

func (c *CredentialsStore) UpdatePassword(userID, newHash string) error {
	query := `
		UPDATE credentials
		SET password_hash = $1,
		    updated_at = now()
		WHERE user_id = $2
	`
	res, err := c.DBclient.DB.Exec(query, newHash, userID)
	if err != nil {
		return fmt.Errorf("[CredentialsStore] failed to update password for user_id=%s: %w", userID, err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("[CredentialsStore] no rows updated for user_id=%s", userID)
	}

	return nil
}

func (c *CredentialsStore) MarkDeleted(userID string) error {
	query := `
		UPDATE credentials
		SET status = 'disabled',
		    updated_at = now()
		WHERE user_id = $1
	`
	result, err := c.DBclient.DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("[CredentialsStore] failed to disable credentials for user_id=%s: %w", userID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("[CredentialsStore] failed to check rows affected for user_id=%s: %w", userID, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("[CredentialsStore] no credentials found for user_id=%s", userID)
	}
	fmt.Printf("[CredentialsStore] Disabled credentials for %s in Database\n", userID)

	// --- Cache cleanup ---
	existed, err := c.RedisClient.KeyExists("user:" + userID)
	if err != nil {
		return fmt.Errorf("[CredentialsStore] failed to check cache for user %s: %w", userID, err)
	}
	if existed {
		fmt.Printf("[CredentialsStore] Found %s in cache, deleting...\n", userID)

		// Lấy username từ cache bằng userID
		username, err := c.RedisClient.GetKey("user:" + userID)
		if err != nil {
			return fmt.Errorf("[CredentialsStore] failed to get username for user_id=%s from cache: %w", userID, err)
		}
		if username != "" {
			if err := c.RedisClient.DeleteKey("user:" + username); err != nil {
				return fmt.Errorf("[CredentialsStore] failed to delete username %s from cache: %w", username, err)
			}
			fmt.Printf("[CredentialsStore] Deleted username %s from cache\n", username)
		}

		// Xóa luôn cache theo userID
		if err := c.RedisClient.DeleteKey("user:" + userID); err != nil {
			return fmt.Errorf("[CredentialsStore] failed to delete userID %s from cache: %w", userID, err)
		}
		fmt.Printf("[CredentialsStore] Deleted userID %s from cache\n", userID)
	}

	return nil
}
