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
	Objectkeys3   sql.NullString `json:"Objectkeys3"`
	URL           string         `json:"url"`
	Status        string         `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
}

type NewPostEvent struct {
	PostID string
	UserID string
}
