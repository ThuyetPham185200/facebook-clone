package model

import (
	"database/sql"
	"time"
)

// ---- DTOs ----
type Media struct {
	MediaID       string         `json:"media_id"`
	UserID        string         `json:"user_id"`
	MediaFileName string         `json:"file_name"`
	MediaType     string         `json:"media_type"`
	Url           sql.NullString `json:"url"`
	Status        string         `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
}

type NewPostEvent struct {
	PostID string
	UserID string
}
