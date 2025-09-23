package store

import (
	"database/sql"
	"fmt"
	"log"
	dbclient "userservice/internal/infra/postgresclient"
	"userservice/internal/model"

	"github.com/google/uuid"
)

type UserStore struct {
	DBclient *dbclient.PostgresClient
}

func NewUserStore(host, port, user, password, dbname string) *UserStore {
	return &UserStore{
		DBclient: dbclient.NewPostgresClient(host, port, user, password, dbname),
	}
}

func (us *UserStore) GetOwnProfile(userID string) (*model.User, error) {
	query := `
		SELECT user_id, username, email, bio, gender, date_of_birth, avatar_url, 
		       is_deleted, created_at, updated_at
		FROM users
		WHERE user_id = $1 AND is_deleted = FALSE
	`
	row := us.DBclient.DB.QueryRow(query, userID)

	var u model.User
	err := row.Scan(
		&u.UserID,
		&u.Username,
		&u.Email,
		&u.Bio,
		&u.Gender,
		&u.DateOfBirth,
		&u.AvatarURL,
		&u.IsDeleted,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateUserProfile inserts a new user into "users" table and returns the created User
func (us *UserStore) CreateUserProfile(username, email string) (*model.User, error) {
	newUUID := uuid.New().String()

	log.Printf("[UserStore] CreateUserProfile called. user_id=%s, username=%s, email=%s",
		newUUID, username, email)

	query := `
		INSERT INTO users (user_id, username, email, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
		RETURNING user_id, username, email, bio, gender, date_of_birth, avatar_url, created_at, updated_at
	`

	row := us.DBclient.DB.QueryRow(query, newUUID, username, email)

	var user model.User
	err := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Bio,
		&user.Gender,
		&user.DateOfBirth,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Printf("[UserStore] failed to insert user (user_id=%s, username=%s, email=%s): %v",
			newUUID, username, email, err)
		return nil, err
	}

	log.Printf("[UserStore] user created successfully: %+v", user)
	return &user, nil
}

func (us *UserStore) UsernameExists(username string) (bool, error) {
	query := `SELECT 1 FROM users WHERE username = $1 AND is_deleted = FALSE LIMIT 1`
	var exists int
	err := us.DBclient.DB.QueryRow(query, username).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("[UserStore] failed to check username=%s: %w", username, err)
	}
	return true, nil
}

func (us *UserStore) EmailExists(email string) (bool, error) {
	query := `SELECT 1 FROM users WHERE email = $1 AND is_deleted = FALSE LIMIT 1`
	var exists int
	err := us.DBclient.DB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (us *UserStore) GetUserByUserID(userID string) (*model.User, error) {
	query := `
		SELECT user_id, username, email, bio, gender, date_of_birth, avatar_url, is_deleted, created_at, updated_at
		FROM users
		WHERE user_id = $1 AND is_deleted = FALSE
	`

	row := us.DBclient.DB.QueryRow(query, userID)

	var u model.User
	err := row.Scan(
		&u.UserID,
		&u.Username,
		&u.Email,
		&u.Bio,
		&u.Gender,
		&u.DateOfBirth,
		&u.AvatarURL,
		&u.IsDeleted,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (us *UserStore) GetUserByUsername(username string) (*model.User, error) {
	query := `
		SELECT user_id, username, email, bio, gender, date_of_birth, avatar_url, is_deleted, created_at, updated_at
		FROM users
		WHERE username = $1 AND is_deleted = FALSE
	`

	row := us.DBclient.DB.QueryRow(query, username)

	var u model.User
	err := row.Scan(
		&u.UserID,
		&u.Username,
		&u.Email,
		&u.Bio,
		&u.Gender,
		&u.DateOfBirth,
		&u.AvatarURL,
		&u.IsDeleted,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (us *UserStore) HardDeleteUserProfile(userID string) error {
	query := `
		delete from users where user_id = $1
	`
	result, err := us.DBclient.DB.Exec(query, userID)

	if err != nil {
		return fmt.Errorf("[UserStore] failed to delete user %s: %w", userID, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("[UserStore] failed to check rows affected for user %s: %w", userID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("[UserStore] no user found with id %s", userID)
	}

	return nil
}

func (us *UserStore) SoftDeleteUserProfile(userID string) error {
	query := `
		update users
		set is_deleted = TRUE
		where user_id = $1
	`
	result, err := us.DBclient.DB.Exec(query, userID)

	if err != nil {
		return fmt.Errorf("[UserStore] failed to soft delete user %s: %w", userID, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("[UserStore] failed to check rows affected for user %s: %w", userID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("[UserStore] no user found with id %s", userID)
	}

	return nil
}
