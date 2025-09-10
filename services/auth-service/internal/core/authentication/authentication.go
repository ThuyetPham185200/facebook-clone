package auth

import (
	auth "authservice/internal/core/session"
	"authservice/internal/core/userserviceclient"
	"authservice/internal/infra/store"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// ---- Interface ----
type AuthenticationManager interface {
	Register(username, email, password string) (userID string, accessToken, refreshToken string, err error)
	Login(login, password string) (accessToken, refreshToken string, err error)
	ChangePassword(string string, oldPassword, newPassword string) error
	DeleteAccount(userID string) error
	// Logout(userID string, refreshToken string) error
	// LogoutAll(userID string) error
	// RefreshToken(refreshToken string) (newAccess, newRefresh string, err error)
}

// ---- Implementation ----
type authenticationManager struct {
	credStore      *store.CredentialsStore // abstract interface to Credentials DB
	userService    *userserviceclient.UserService
	sessionManager auth.SessionManager
}

func NewAuthenticationManager(cs *store.CredentialsStore, us *userserviceclient.UserService, sm auth.SessionManager) AuthenticationManager {
	return &authenticationManager{
		credStore:      cs,
		userService:    us,
		sessionManager: sm,
	}
}

// ---- Register ----
func (am *authenticationManager) Register(username, email, password string) (string, string, string, error) {
	// check in credential store first (cache or DB)
	exists, errc := am.credStore.ExistsUser(username, email)
	if errc != nil {
		return "", "", "", errc
	}

	if exists {
		return "", "", "", errors.New("[authenticationManager] user already exists in credentials")
	}

	// check duplicate username/email in UserService
	usernameexists, emailexists, err := am.userService.CheckUserExists(username, email)
	if err != nil {
		return "", "", "", err
	}
	if usernameexists {
		return "", "", "", errors.New("[authenticationManager] username already exists")
	}
	if emailexists {
		return "", "", "", errors.New("[authenticationManager] email already exists")
	}

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", "", err
	}

	// create user profile in UserService
	userID, err := am.userService.CreateUserProfile(username, email)
	if err != nil {
		return "", "", "", err
	}

	// save credentials in Credentials DB
	err = am.credStore.Save(userID, username, email, string(hashed))
	if err != nil {
		return "", "", "", err
	}

	// create session
	access, refresh, err := am.sessionManager.CreateSession(userID)
	if err != nil {
		return "", "", "", err
	}

	return userID, access, refresh, nil
}

// ---- Login ----
func (am *authenticationManager) Login(login, password string) (string, string, error) {
	// step 1: get userid
	var userid string
	var err error
	userid, err = am.credStore.GetUserIdByName(login)
	if err != nil {
		userid, err = am.userService.GetUserIdByName(login)
		if err != nil {
			return "", "", err
		} else {
			fmt.Printf("[authenticationManager - Login] got user_id on Userservice %s", userid)
		}
	} else {
		fmt.Printf("[authenticationManager - Login] got user_id on cache %s", userid)
	}

	// step 2: get hasedpassword
	cred, err := am.credStore.GetCredentialByUserID(userid)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cred.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	access, refresh, err := am.sessionManager.CreateSession(cred.UserID)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

// ---- Change Password ----
func (am *authenticationManager) ChangePassword(userID string, oldPassword, newPassword string) error {
	cred, err := am.credStore.GetCredentialByUserID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cred.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := am.credStore.UpdatePassword(userID, string(newHash)); err != nil {
		return err
	}

	// revoke all sessions after password change
	return am.sessionManager.LogoutAll(userID)
}

// // ---- Delete Account ----
func (am *authenticationManager) DeleteAccount(userID string) error {
	// soft delete user in User Service
	if err := am.userService.DeleteUserProfile(userID); err != nil {
		return err
	}

	// mark credentials deleted
	if err := am.credStore.MarkDeleted(userID); err != nil {
		return err
	}

	// revoke sessions
	return am.sessionManager.LogoutAll(userID)
}

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
