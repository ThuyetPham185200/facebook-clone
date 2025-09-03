package model

import (
	"database/sql"
	"time"
)

// ---- DTOs ----
type User struct {
	UserID      string         `json:"user_id"`
	Username    string         `json:"username"`
	Email       string         `json:"email,omitempty"`
	Bio         sql.NullString `json:"bio,omitempty"`
	Gender      sql.NullString `json:"gender,omitempty"`
	DateOfBirth sql.NullTime   `json:"date_of_birth,omitempty"`
	AvatarURL   sql.NullString `json:"avatar_url,omitempty"`
	IsDeleted   bool           `json:"is_deleted,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type CheckExistRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type CheckExistResponse struct {
	ExistsUsername bool `json:"exists_username"`
	ExistsEmail    bool `json:"exists_email"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
