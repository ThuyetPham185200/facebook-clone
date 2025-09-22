package store

import (
	dbclient "feedservice/internal/infra/postgresclient"
	"fmt"
)

type PostStore struct {
	DBClient *dbclient.PostgresClient
}

func NewPostStore(postgrescfg *PostGresConfig) *PostStore {
	mediaStore := &PostStore{}
	mediaStore.DBClient = dbclient.NewPostgresClient(postgrescfg.Host, postgrescfg.Port, postgrescfg.User, postgrescfg.Password, postgrescfg.DBname)
	return mediaStore
}

func (p *PostStore) InsertPost(postID, userID, content string) error {
	query := `INSERT INTO posts (post_id, user_id, content, created_at) VALUES ($1, $2, $3, now())`
	_, err := p.DBClient.DB.Exec(query, postID, userID, content)
	if err != nil {
		return fmt.Errorf("failed to insert post: %w", err)
	}
	return nil
}
