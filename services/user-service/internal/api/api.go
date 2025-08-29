package api

import (
	"net/http"
	"userservice/internal/model"
	"userservice/utils"

	"github.com/gorilla/mux"
)

// ---- Manager Interfaces ----
type UserStore interface {
	GetOwnProfile(userID string) (*model.User, error)
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
