package store

import (
	dbclient "feedservice/internal/infra/postgresclient"
	"fmt"
	"time"
)

type PostStore struct {
	DBClient *dbclient.PostgresClient
}

func NewPostStore(postgrescfg *PostGresConfig) *PostStore {
	mediaStore := &PostStore{}
	mediaStore.DBClient = dbclient.NewPostgresClient(postgrescfg.Host, postgrescfg.Port, postgrescfg.User, postgrescfg.Password, postgrescfg.DBname)
	return mediaStore
}

type Post struct {
	ID        string
	UserID    string
	Content   string
	CreatedAt time.Time
}

func (p *PostStore) InsertPost(postID, userID, content string) error {
	query := `INSERT INTO posts (post_id, user_id, content, created_at) VALUES ($1, $2, $3, now())`
	_, err := p.DBClient.DB.Exec(query, postID, userID, content)
	if err != nil {
		return fmt.Errorf("failed to insert post: %w", err)
	}
	return nil
}

func (p *PostStore) GetPostByID(postID string) (*Post, error) {
	query := `SELECT post_id, user_id, content, created_at FROM posts WHERE post_id = $1`

	row := p.DBClient.DB.QueryRow(query, postID)

	post := &Post{}
	err := row.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post %s: %w", postID, err)
	}

	return post, nil
}
