package store

import (
	"followservice/model"
	"time"

	dbclient "followservice/internal/infra/postgresclient"
)

type FollowStore struct {
	DBClient *dbclient.PostgresClient
}

type PostGresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
}

func NewFollowStore(postgrescfg *PostGresConfig) *FollowStore {
	return &FollowStore{
		DBClient: dbclient.NewPostgresClient(
			postgrescfg.Host,
			postgrescfg.Port,
			postgrescfg.User,
			postgrescfg.Password,
			postgrescfg.DBname,
		),
	}
}

// Follow inserts a new follow relationship
func (f *FollowStore) Follow(follower_id string, followee_id string) (model.Follow, error) {
	query := `
		INSERT INTO follows (follower_id, followee_id, created_at)
		VALUES ($1, $2, $3)
		RETURNING follower_id, followee_id, created_at
	`
	row := f.DBClient.DB.QueryRow(query, follower_id, followee_id, time.Now())

	var follow model.Follow
	if err := row.Scan(&follow.FollowerID, &follow.FolloweeID, &follow.CreatedAt); err != nil {
		return model.Follow{}, err
	}

	return follow, nil
}

// Unfollow deletes an existing follow relationship
func (f *FollowStore) Unfollow(follower_id string, followee_id string) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2`
	_, err := f.DBClient.DB.Exec(query, follower_id, followee_id)
	return err
}

// GetFollowers fetches all followers for a given user (followee_id)
func (f *FollowStore) GetFollowers(userID string) ([]model.Follow, error) {
	query := `SELECT follower_id, followee_id, created_at FROM follows WHERE followee_id = $1`
	rows, err := f.DBClient.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []model.Follow
	for rows.Next() {
		var follow model.Follow
		if err := rows.Scan(&follow.FollowerID, &follow.FolloweeID, &follow.CreatedAt); err != nil {
			return nil, err
		}
		followers = append(followers, follow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

// GetFollowees fetches all followers for a given user (followee_id)
func (f *FollowStore) GetFollowees(userID string) ([]model.Follow, error) {
	query := `SELECT follower_id, followee_id, created_at FROM follows WHERE follower_id = $1`
	rows, err := f.DBClient.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followees []model.Follow
	for rows.Next() {
		var follow model.Follow
		if err := rows.Scan(&follow.FollowerID, &follow.FolloweeID, &follow.CreatedAt); err != nil {
			return nil, err
		}
		followees = append(followees, follow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followees, nil
}
