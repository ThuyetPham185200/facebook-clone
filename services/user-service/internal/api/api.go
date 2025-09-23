package api

import (
	"encoding/json"
	"log"
	"net/http"
	"userservice/internal/model"
	"userservice/utils"

	"github.com/gorilla/mux"
)

// ---- Manager Interfaces ----
type UserStore interface {
	GetOwnProfile(userID string) (*model.User, error)
	CreateUserProfile(username, email string) (*model.User, error)
	UsernameExists(username string) (bool, error)
	EmailExists(username string) (bool, error)
	GetUserByUsername(username string) (*model.User, error)
	SoftDeleteUserProfile(userID string) error
	GetUserByUserID(userID string) (*model.User, error)
}

// ---- API Layer ----
type UserAPI struct {
	userstore UserStore
}

func NewUserAPI(us UserStore) *UserAPI {
	return &UserAPI{us}
}

func (api *UserAPI) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/me", api.handleGetOwnProfile).Methods("GET")
	r.HandleFunc("/users", api.handleCreateUserProfile).Methods("POST")
	r.HandleFunc("/users/exists", api.handleCheckExist).Methods("GET")
	r.HandleFunc("/users/by-username/{username}", api.handleGetUseridByUsername).Methods("GET")
	r.HandleFunc("/users/{user_id}", api.handleDeleteUser).Methods("DELETE")
	r.HandleFunc("/users/{user_id}", api.handleGetUser).Methods("GET")
}

// ---- Handlers ----
func (api *UserAPI) handleGetOwnProfile(w http.ResponseWriter, r *http.Request) {
	// Normally user_id comes from JWT claims or middleware
	userID := r.Context().Value("user_id")
	if userID == nil {
		utils.WriteError(w, http.StatusUnauthorized, "missing user id in context")
		return
	}

	user, err := api.userstore.GetOwnProfile(userID.(string))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

// POST /users
func (api *UserAPI) handleCreateUserProfile(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[UserAPI] invalid request body: %v", err)
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	log.Printf("[UserAPI] handleCreateUserProfile called. username=%s, email=%s", req.Username, req.Email)

	if req.Username == "" || req.Email == "" {
		log.Println("[UserAPI] missing username or email in request body")
		utils.WriteError(w, http.StatusBadRequest, "missing username or email")
		return
	}

	user, err := api.userstore.CreateUserProfile(req.Username, req.Email)
	if err != nil {
		log.Printf("[UserAPI] failed to create user profile (username=%s, email=%s): %v",
			req.Username, req.Email, err)
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[UserAPI] user profile created successfully: %+v", user)
	utils.WriteJSON(w, http.StatusCreated, user) // dùng 201 Created thay vì 200 OK
}

// GET /users/exists?username=abc&email=abc@example.com
func (api *UserAPI) handleCheckExist(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	email := r.URL.Query().Get("email")

	log.Printf("[UserAPI] handleCheckExist called. username=%s, email=%s", username, email)

	existsUsername, err := api.userstore.UsernameExists(username)
	if err != nil {
		log.Printf("[UserAPI] DB error when checking username=%s: %v", username, err)
		utils.WriteError(w, http.StatusInternalServerError, "DB error")
		return
	}

	existsEmail, err := api.userstore.EmailExists(email)
	if err != nil {
		log.Printf("[UserAPI] DB error when checking email=%s: %v", email, err)
		utils.WriteError(w, http.StatusInternalServerError, "DB error")
		return
	}

	log.Printf("[UserAPI] check exist result: username=%s (exists=%v), email=%s (exists=%v)",
		username, existsUsername, email, existsEmail)

	utils.WriteJSON(w, http.StatusOK, model.CheckExistResponse{
		ExistsUsername: existsUsername,
		ExistsEmail:    existsEmail,
	})
}

// GET /users/by-username/{username}
func (api *UserAPI) handleGetUseridByUsername(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	username := vars["username"]

	log.Printf("[UserAPI] handleGetUseridByUsername called. username=%s", username)

	if username == "" {
		utils.WriteError(w, http.StatusBadRequest, "username is required")
		return
	}

	// Query DB
	user, err := api.userstore.GetUserByUsername(username)
	if err != nil {
		log.Printf("[UserAPI] failed to get user_id for username=%s: %v", username, err)
		utils.WriteError(w, http.StatusInternalServerError, "DB error")
		return
	}
	if user == nil {
		utils.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	// Response
	resp := map[string]string{"user_id": user.UserID}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (api *UserAPI) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID == "" {
		utils.WriteError(w, http.StatusBadRequest, "missing userid in path")
		return
	}

	log.Printf("[UserAPI] handleDeleteUser called. userID=%s", userID)

	// gọi xuống UserStore để xóa user
	err := api.userstore.SoftDeleteUserProfile(userID)
	if err != nil {
		log.Printf("[UserAPI] failed to delete user %s: %v", userID, err)
		utils.WriteError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	log.Printf("[UserAPI] user %s deleted successfully", userID)
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user deleted successfully",
	})
}

func (api *UserAPI) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID == "" {
		utils.WriteError(w, http.StatusBadRequest, "missing userid in path")
		return
	}

	log.Printf("[UserAPI] handleDeleteUser called. userID=%s", userID)

	// gọi xuống UserStore để get user
	// Query DB
	user, err := api.userstore.GetUserByUserID(userID)
	if err != nil {
		log.Printf("[UserAPI] failed to get user_id for username=%s: %v", userID, err)
		utils.WriteError(w, http.StatusInternalServerError, "DB error")
		return
	}
	if user == nil {
		utils.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	// Response
	resp := map[string]string{
		"user_id":    user.UserID,
		"username":   user.Username,
		"avatar_url": user.AvatarURL.String,
		"gender":     user.Gender.String,
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}
