package api

import (
	"encoding/json"
	"followservice/model"
	"followservice/utils"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type FollowStore interface {
	Follow(follower_id string, followee_id string) (model.Follow, error)
	Unfollow(follower_id string, followee_id string) error
	GetFollowers(userID string) ([]model.Follow, error)
	GetFollowees(userID string) ([]model.Follow, error)
}

type FollowAPI struct {
	followStore FollowStore
}

func NewFollowAPI(followStore_ FollowStore) *FollowAPI {
	return &FollowAPI{
		followStore: followStore_,
	}
}

func (api *FollowAPI) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/follows", api.handleFollow).Methods("POST")
	r.HandleFunc("/follows", api.handleUnFollow).Methods("DELETE")
	r.HandleFunc("/follows/{user_id}/followers", api.handleGetListFollowers).Methods("GET")
	r.HandleFunc("/follows/{user_id}/followees", api.handleGetListFollowees).Methods("GET")
}

type followRequest struct {
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
}

func (api *FollowAPI) handleFollow(w http.ResponseWriter, r *http.Request) {
	var req followRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FollowerID == "" || req.FolloweeID == "" {
		utils.WriteError(w, http.StatusBadRequest, "follower_id and followee_id are required")
		return
	}

	follow, err := api.followStore.Follow(req.FollowerID, req.FolloweeID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := map[string]interface{}{
		"status":      "success",
		"follower_id": follow.FollowerID,
		"followee_id": follow.FolloweeID,
		"created_at":  follow.CreatedAt.Format(time.RFC3339),
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (api *FollowAPI) handleUnFollow(w http.ResponseWriter, r *http.Request) {
	type request struct {
		FollowerID string `json:"follower_id"`
		FolloweeID string `json:"followee_id"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FollowerID == "" || req.FolloweeID == "" {
		utils.WriteError(w, http.StatusBadRequest, "follower_id and followee_id are required")
		return
	}

	err := api.followStore.Unfollow(req.FollowerID, req.FolloweeID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to unfollow: "+err.Error())
		return
	}

	resp := model.UnfollowResponse{
		Status:  "success",
		Message: "unfollowed",
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (api *FollowAPI) handleGetListFollowers(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	if userID == "" {
		utils.WriteError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	followers, err := api.followStore.GetFollowers(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch followers: "+err.Error())
		return
	}

	resp := map[string]interface{}{
		"user_id":   userID,
		"followers": followers,
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

func (api *FollowAPI) handleGetListFollowees(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	if userID == "" {
		utils.WriteError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	followers, err := api.followStore.GetFollowees(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch followers: "+err.Error())
		return
	}

	resp := map[string]interface{}{
		"user_id":   userID,
		"followees": followers,
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}
