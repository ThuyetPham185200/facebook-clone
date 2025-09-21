package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type FanoutManager interface {
}

type PostManager interface {
}

type FeedAPI struct {
	FanoutInterface *FanoutManager
	PostInteface    *PostManager
}

func NewFeedAPI(fanoutInterface *FanoutManager, postInteface *PostManager) *FeedAPI {
	return &FeedAPI{
		FanoutInterface: fanoutInterface,
		PostInteface:    postInteface,
	}
}

func (api *FeedAPI) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/media", api.handleCreateMedia).Methods("POST")
}

func (api *FeedAPI) handleCreateMedia(w http.ResponseWriter, r *http.Request) {

}
