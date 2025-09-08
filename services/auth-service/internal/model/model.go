package model

import (
	"database/sql"
	"time"
)

// ---- DTOs ----

type Credential struct {
	ID           string         `json:"id"`            // UUID, generated in app
	UserID       string         `json:"user_id"`       // UUID, FK -> users(user_id)
	PasswordHash string         `json:"password_hash"` // hashed password
	MFASecret    sql.NullString `json:"mfa_secret"`    // nullable
	Status       string         `json:"status"`        // active/locked/disabled
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type ResetRequest struct {
	Email string `json:"email"`
}
type ResetConfirmRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}
