package api

import (
	"encoding/json"
	"feedservice/internal/model"
	"feedservice/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type FanoutManager interface {
}

type PostManager interface {
	CreateMedia(userID string, mediatype string, mediafilename string) (model.Media, error)
	CreatePost(userID string, content string, mediaIDs []string) (string, error)
}

type FeedAPI struct {
	FanoutInterface FanoutManager
	PostInteface    PostManager
}

func NewFeedAPI(fanoutInterface FanoutManager, postInteface PostManager) *FeedAPI {
	return &FeedAPI{
		FanoutInterface: fanoutInterface,
		PostInteface:    postInteface,
	}
}

func (api *FeedAPI) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/media", api.handleCreateMedia).Methods("POST")
	r.HandleFunc("/posts", api.handleCreatePost).Methods("POST")
}

func (api *FeedAPI) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Content  string   `json:"content"`
		MediaIDs []string `json:"media_ids"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Content == "" {
		utils.WriteError(w, http.StatusBadRequest, "content is required")
		return
	}

	// Get user ID from header (set by auth middleware / gateway)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	postID, err := api.PostInteface.CreatePost(userID, req.Content, req.MediaIDs)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to create post: "+err.Error())
		return
	}

	resp := map[string]interface{}{
		"post_id": postID,
		"message": "Post created",
	}
	utils.WriteJSON(w, http.StatusCreated, resp)
}

func (api *FeedAPI) handleCreateMedia(w http.ResponseWriter, r *http.Request) {
	type request struct {
		MediaTypes []string `json:"media_type"`
		FileNames  []string `json:"file_name"`
	}
	type response struct {
		MediaIDs   []string `json:"media_id"`
		UploadURLs []string `json:"upload_url"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.MediaTypes) == 0 || len(req.FileNames) == 0 || len(req.MediaTypes) != len(req.FileNames) {
		utils.WriteError(w, http.StatusBadRequest, "media_type and file_name must be non-empty arrays of equal length")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	if len(req.MediaTypes) != len(req.FileNames) {
		utils.WriteError(w, http.StatusBadRequest, "media_types and file_names must have the same length")
		return
	}

	var resp response
	for i := range req.MediaTypes {
		media, err := api.PostInteface.CreateMedia(userID, req.MediaTypes[i], req.FileNames[i])
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "failed to create media: "+err.Error())
			return
		}
		resp.MediaIDs = append(resp.MediaIDs, media.MediaID)
		resp.UploadURLs = append(resp.UploadURLs, media.Url.String) // presigned URL
	}

	utils.WriteJSON(w, http.StatusCreated, resp)
}
