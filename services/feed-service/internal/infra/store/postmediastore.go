package store

import (
	dbclient "feedservice/internal/infra/postgresclient"
	"fmt"
)

type PostMediaStore struct {
	DBClient *dbclient.PostgresClient
}

func NewPostMediaStore(postgrescfg *PostGresConfig) *PostStore {
	mediaStore := &PostStore{}
	mediaStore.DBClient = dbclient.NewPostgresClient(postgrescfg.Host, postgrescfg.Port, postgrescfg.User, postgrescfg.Password, postgrescfg.DBname)
	return mediaStore
}

func (pm *PostMediaStore) LinkMediaToPost(postID string, mediaIDs []string) error {
	if len(mediaIDs) == 0 {
		return nil
	}

	query := `INSERT INTO post_media (post_id, media_id) VALUES ($1, $2)`
	for _, mediaID := range mediaIDs {
		if _, err := pm.DBClient.DB.Exec(query, postID, mediaID); err != nil {
			return fmt.Errorf("failed to link media %s to post %s: %w", mediaID, postID, err)
		}
	}
	return nil
}

func (pm *PostMediaStore) GetMediaByPostID(postID string) ([]string, error) {
	query := `SELECT media_id FROM post_media WHERE post_id = $1`

	rows, err := pm.DBClient.DB.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch media for post %s: %w", postID, err)
	}
	defer rows.Close()

	var mediaIDs []string
	for rows.Next() {
		var mediaID string
		if err := rows.Scan(&mediaID); err != nil {
			return nil, fmt.Errorf("failed to scan media_id: %w", err)
		}
		mediaIDs = append(mediaIDs, mediaID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return mediaIDs, nil
}
