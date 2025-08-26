package auth

import (
	auth "authservice/internal/core/session"
	"authservice/internal/infra/store"
	"encoding/binary"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}
func (u *UserService) CreateUserProfile(username, email string) (int64, error) {
	newUUID := uuid.New()
	// Lấy 8 byte đầu của UUID để convert sang int64
	userID := int64(binary.BigEndian.Uint64(newUUID[:8]))

	// TODO: insert vào bảng users ở DB, có cột id kiểu BIGINT
	return userID, nil
}

// ---- Interface ----
type AuthenticationManager interface {
	Register(username, email, password string) (userID int64, accessToken, refreshToken string, err error)
	// Login(login, password string) (accessToken, refreshToken string, err error)
	// ChangePassword(userID int64, oldPassword, newPassword string) error
	// DeleteAccount(userID int64) error
	// Logout(userID int64, refreshToken string) error
	// LogoutAll(userID int64) error
	// RefreshToken(refreshToken string) (newAccess, newRefresh string, err error)
}

// ---- Implementation ----
type authenticationManager struct {
	credStore      *store.CredentialsStore // abstract interface to Credentials DB
	userService    *UserService
	sessionManager auth.SessionManager
}

func NewAuthenticationManager(cs *store.CredentialsStore, us *UserService, sm auth.SessionManager) AuthenticationManager {
	return &authenticationManager{
		credStore:      cs,
		userService:    us,
		sessionManager: sm,
	}
}

// ---- Register ----
func (am *authenticationManager) Register(username, email, password string) (int64, string, string, error) {
	// check duplicate username/email
	exists, err := am.credStore.Exists(username, email)

	if err != nil {
		return 0, "", "", err
	}
	if exists {
		return 0, "", "", errors.New("username or email already exists")
	}

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, "", "", err
	}

	// create user profile in UserService
	userID, err := am.userService.CreateUserProfile(username, email)

	if err != nil {
		return 0, "", "", err
	}

	// save credentials in Credentials DB
	err = am.credStore.Save(strconv.FormatInt(userID, 10), username, email, string(hashed))
	if err != nil {
		return 0, "", "", err
	}

	//create session
	access, refresh, err := am.sessionManager.CreateSession(int64(userID))
	if err != nil {
		return 0, "", "", err
	}

	return userID, access, refresh, nil
}

// // ---- Login ----
// func (am *authenticationManager) Login(login, password string) (string, string, error) {
// 	cred, err := am.credStore.GetByLogin(login)
// 	if err != nil {
// 		return "", "", errors.New("invalid credentials")
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(cred.PasswordHash), []byte(password)); err != nil {
// 		return "", "", errors.New("invalid credentials")
// 	}

// 	access, refresh, err := am.sessionManager.CreateSession(cred.UserID)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	return access, refresh, nil
// }

// // ---- Change Password ----
// func (am *authenticationManager) ChangePassword(userID int64, oldPassword, newPassword string) error {
// 	cred, err := am.credStore.GetByUserID(userID)
// 	if err != nil {
// 		return errors.New("user not found")
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(cred.PasswordHash), []byte(oldPassword)); err != nil {
// 		return errors.New("invalid old password")
// 	}

// 	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}

// 	if err := am.credStore.UpdatePassword(userID, string(newHash)); err != nil {
// 		return err
// 	}

// 	// revoke all sessions after password change
// 	return am.sessionManager.LogoutAll(userID)
// }

// // ---- Delete Account ----
// func (am *authenticationManager) DeleteAccount(userID int64) error {
// 	// soft delete user in User Service
// 	if err := am.userService.DeleteUserProfile(userID); err != nil {
// 		return err
// 	}

// 	// mark credentials deleted
// 	if err := am.credStore.MarkDeleted(userID); err != nil {
// 		return err
// 	}

// 	// revoke sessions
// 	return am.sessionManager.LogoutAll(userID)
// }

// // ---- Logout single session ----
// func (am *authenticationManager) Logout(userID int64, refreshToken string) error {
// 	return am.sessionManager.Logout(userID, refreshToken)
// }

// // ---- Logout all sessions ----
// func (am *authenticationManager) LogoutAll(userID int64) error {
// 	return am.sessionManager.LogoutAll(userID)
// }

// // ---- Refresh token ----
// func (am *authenticationManager) RefreshToken(refreshToken string) (string, string, error) {
// 	return am.sessionManager.RefreshToken(refreshToken)
// }
