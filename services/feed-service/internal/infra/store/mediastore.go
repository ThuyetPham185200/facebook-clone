package store

import (
	"database/sql"
	dbclient "feedservice/internal/infra/postgresclient"
	s3 "feedservice/internal/infra/s3client"
	"feedservice/internal/model"
	"fmt"
	"time"
)

type MediaStore struct {
	DBClient *dbclient.PostgresClient
	S3client *s3.S3Client
}

type PostGresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
}

type S3Config struct {
	Endpoint  string
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
}

func NewMediaStore(postgrescfg *PostGresConfig, s3cfg *S3Config) *MediaStore {
	mediaStore := &MediaStore{}
	mediaStore.DBClient = dbclient.NewPostgresClient(postgrescfg.Host, postgrescfg.Port, postgrescfg.User, postgrescfg.Password, postgrescfg.DBname)
	mediaStore.S3client = s3.NewS3Client(s3cfg.Endpoint, s3cfg.Region, s3cfg.AccessKey, s3cfg.SecretKey, s3cfg.Bucket)
	return mediaStore
}

func (m *MediaStore) ValidateUserMedia(userID string, mediaIDs []string) error {
	if len(mediaIDs) == 0 {
		return nil
	}

	query := `SELECT media_id FROM medias WHERE user_id=$1 AND media_id = ANY($2)`
	rows, err := m.DBClient.DB.Query(query, userID, mediaIDs)
	if err != nil {
		return fmt.Errorf("DB query failed: %w", err)
	}
	defer rows.Close()

	validIDs := map[string]struct{}{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan media_id: %w", err)
		}
		validIDs[id] = struct{}{}
	}

	if len(validIDs) != len(mediaIDs) {
		return fmt.Errorf("one or more media IDs are invalid or not owned by user")
	}

	return nil
}

func (m *MediaStore) UpdateMediaStatus(mediaIDs []string, status string) error {
	if len(mediaIDs) == 0 {
		return nil
	}

	query := `UPDATE medias SET status=$1 WHERE media_id = ANY($2)`
	_, err := m.DBClient.DB.Exec(query, status, mediaIDs)
	if err != nil {
		return fmt.Errorf("failed to update media status: %w", err)
	}
	return nil
}

func (s *MediaStore) CreateMediaRecord(mediaID string, userID string, mediaType string, mediaStatus string, objectkeys3 string) (model.Media, error) {

	query := `
		INSERT INTO medias (media_id, user_id, media_type, objectkeys3, status, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING media_id, user_id, media_type, objectkeys3, status, created_at
	`

	var media model.Media
	var objectkey sql.NullString
	if objectkeys3 != "" {
		objectkey = sql.NullString{String: objectkeys3, Valid: true}
	} else {
		objectkey = sql.NullString{Valid: false}
	}

	err := s.DBClient.DB.QueryRow(query,
		mediaID, userID, mediaType, objectkey, mediaStatus,
	).Scan(&media.MediaID, &media.UserID, &media.MediaType, &media.Objectkeys3, &media.Status, &media.CreatedAt)

	if err != nil {
		return model.Media{}, fmt.Errorf("failed to insert media: %w", err)
	}

	// Fill extra fields (if any missing from RETURNING)
	media.MediaFileName = "" // you can later update this field when binding presigned URL file name

	return media, nil
}

func (s *MediaStore) GetMediaByID(mID string) (model.Media, error) {
	query := `
		SELECT media_id, user_id, media_type, objectkeys3, status, created_at
		FROM medias
		WHERE media_id = $1
	`

	var media model.Media
	var objectkey sql.NullString

	err := s.DBClient.DB.QueryRow(query, mID).Scan(
		&media.MediaID,
		&media.UserID,
		&media.MediaType,
		&objectkey,
		&media.Status,
		&media.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Media{}, nil // not found
		}
		return model.Media{}, fmt.Errorf("failed to get media by id %s: %w", mID, err)
	}

	if objectkey.Valid {
		// 🔑 Generate a temporary pre-signed URL (e.g., 15m)
		media.URL = s.S3client.GeneratePreSignedGetURL(objectkey.String, 15*time.Minute)
	}

	// Optional: fill defaults for non-persistent fields
	media.MediaFileName = ""

	return media, nil
}
