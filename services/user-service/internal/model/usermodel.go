package model

import "time"

// ---- DTOs ----
type User struct {
	UserID      string     `json:"user_id"` // UUID
	Username    string     `json:"username"`
	Email       string     `json:"email,omitempty"`
	Bio         string     `json:"bio,omitempty"`
	Gender      string     `json:"gender,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"` // NULLABLE
	AvatarURL   string     `json:"avatar_url,omitempty"`
	IsDeleted   bool       `json:"is_deleted,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
