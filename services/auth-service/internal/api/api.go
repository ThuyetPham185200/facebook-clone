package api

import (
	auth "authservice/internal/core/authentication"
	"authservice/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// // ---- Manager Interfaces ----
type AuthenticationManager interface {
	Register(username, email, password string) (userID int64, accessToken, refreshToken string, err error)
	// Login(login, password string) (accessToken, refreshToken string, err error)
	// ChangePassword(userID int64, oldPwd, newPwd string) error
	// DeleteAccount(userID int64) error
}

type SessionManager interface {
	RefreshToken(refreshToken string) (newAccess, newRefresh string, err error)
	Logout(userID int64, refreshToken string) error
	LogoutAll(userID int64) error
}

type PasswordResetManager interface {
	RequestReset(email string) error
	ConfirmReset(token, newPassword string) error
}

// ---- API Layer ----
type AuthAPI struct {
	authManager auth.AuthenticationManager
	// sessionManager SessionManager
	// resetManager   PasswordResetManager
}

// func NewAuthAPI(am AuthenticationManager, sm SessionManager, prm PasswordResetManager) *AuthAPI {
// 	return &AuthAPI{am, sm, prm}
// }

func NewAuthAPI(am auth.AuthenticationManager) *AuthAPI {
	//return &AuthAPI{am, sm, prm}
	return &AuthAPI{am}
}

func (api *AuthAPI) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/register", api.handleRegister).Methods("POST")
	// r.HandleFunc("/login", api.handleLogin).Methods("POST")
	// r.HandleFunc("/me/password", api.handleChangePassword).Methods("PUT")
	// r.HandleFunc("/me", api.handleDeleteAccount).Methods("DELETE")
	// r.HandleFunc("/refresh", api.handleRefreshToken).Methods("POST")
	// r.HandleFunc("/logout", api.handleLogout).Methods("POST")
	// r.HandleFunc("/logout_all", api.handleLogoutAll).Methods("POST")
	// r.HandleFunc("/password/reset/request", api.handleRequestPasswordReset).Methods("POST")
	// r.HandleFunc("/password/reset/confirm", api.handleConfirmPasswordReset).Methods("POST")
}

// ---- DTOs ----

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// type loginRequest struct {
// 	Login    string `json:"login"`
// 	Password string `json:"password"`
// }
// type changePasswordRequest struct {
// 	OldPassword string `json:"old_password"`
// 	NewPassword string `json:"new_password"`
// }
// type refreshRequest struct {
// 	RefreshToken string `json:"refresh_token"`
// }
// type logoutRequest struct {
// 	RefreshToken string `json:"refresh_token"`
// }
// type resetRequest struct {
// 	Email string `json:"email"`
// }
// type resetConfirmRequest struct {
// 	Token       string `json:"token"`
// 	NewPassword string `json:"new_password"`
// }

// ---- Handlers ----

func (api *AuthAPI) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid data")
		return
	}
	userID, access, refresh, err := api.authManager.Register(req.Username, req.Email, req.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":       userID,
		"access_token":  access,
		"refresh_token": refresh,
	})
}

// func (api *AuthAPI) handleLogin(w http.ResponseWriter, r *http.Request) {
// 	var req loginRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.WriteJSON(w, http.StatusBadRequest, "Invalid data")
// 		return
// 	}
// 	access, refresh, err := api.authManager.Login(req.Login, req.Password)
// 	if err != nil {
// 		utils.WriteJSON(w, http.StatusUnauthorized, "Invalid credentials")
// 		return
// 	}
// 	utils.WriteJSON(w, http.StatusOK, map[string]string{
// 		"access_token":  access,
// 		"refresh_token": refresh,
// 	})
// }

// func (api *AuthAPI) handleChangePassword(w http.ResponseWriter, r *http.Request) {
// 	var req changePasswordRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.WriteJSON(w, http.StatusBadRequest, "Invalid data")
// 		return
// 	}
// 	userID := r.Context().Value("user_id").(int64)
// 	if err := api.authManager.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
// 		utils.WriteJSON(w, http.StatusForbidden, err.Error())
// 		return
// 	}
// 	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Password updated"})
// }

// func (api *AuthAPI) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
// 	userID := r.Context().Value("user_id").(int64)
// 	if err := api.authManager.DeleteAccount(userID); err != nil {
// 		utils.WriteJSON(w, http.StatusForbidden, "Unauthorized")
// 		return
// 	}
// 	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Account deleted"})
// }

// func (api *AuthAPI) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
// 	var req refreshRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.WriteJSON(w, http.StatusBadRequest, "Invalid data")
// 		return
// 	}
// 	access, refresh, err := api.sessionManager.RefreshToken(req.RefreshToken)
// 	if err != nil {
// 		utils.WriteJSON(w, http.StatusUnauthorized, "Invalid refresh token")
// 		return
// 	}
// 	utils.WriteJSON(w, http.StatusOK, map[string]string{
// 		"access_token":  access,
// 		"refresh_token": refresh,
// 	})
// }

// func (api *AuthAPI) handleLogout(w http.ResponseWriter, r *http.Request) {
// 	userID := r.Context().Value("user_id").(int64)
// 	var req logoutRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.WriteJSON(w, http.StatusBadRequest, "Invalid data")
// 		return
// 	}
// 	if err := api.sessionManager.Logout(userID, req.RefreshToken); err != nil {
// 		utils.WriteJSON(w, http.StatusUnauthorized, "Invalid token")
// 		return
// 	}
// 	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Logged out from current session"})
// }

// func (api *AuthAPI) handleLogoutAll(w http.ResponseWriter, r *http.Request) {
// 	userID := r.Context().Value("user_id").(int64)
// 	if err := api.sessionManager.LogoutAll(userID); err != nil {
// 		utils.WriteJSON(w, http.StatusUnauthorized, "Invalid token")
// 		return
// 	}
// 	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Logged out from all sessions"})
// }

// func (api *AuthAPI) handleRequestPasswordReset(w http.ResponseWriter, r *http.Request) {
// 	var req resetRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.WriteJSON(w, http.StatusBadRequest, "Invalid data")
// 		return
// 	}
// 	_ = api.resetManager.RequestReset(req.Email) // không leak email tồn tại
// 	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "If this email exists, a reset link has been sent"})
// }

// func (api *AuthAPI) handleConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
// 	var req resetConfirmRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		utils.WriteJSON(w, http.StatusBadRequest, "Invalid data")
// 		return
// 	}
// 	if err := api.resetManager.ConfirmReset(req.Token, req.NewPassword); err != nil {
// 		utils.WriteJSON(w, http.StatusBadRequest, "Invalid or expired token")
// 		return
// 	}
// 	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Password has been reset"})
// }
