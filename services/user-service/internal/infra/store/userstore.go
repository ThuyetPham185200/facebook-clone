package store

import (
	dbclient "userservice/internal/infra/postgresclient"
	"userservice/internal/model"
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
