package model

import "time"

type Follow struct {
	FollowerID string    `json:"follower_id"`
	FolloweeID string    `json:"followee_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// Response when unfollow is successful
type UnfollowResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
