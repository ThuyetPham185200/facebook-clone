package followserviceclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FollowServiceClient struct {
	BaseURL string
	Client  *http.Client
}

func NewFollowServiceClient(baseURL string) *FollowServiceClient {
	return &FollowServiceClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (c *FollowServiceClient) GetFollowers(userID string) ([]string, error) {
	url := fmt.Sprintf("%s/follows/%s/followers", c.BaseURL, userID)

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call follow service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("follow service returned %d", resp.StatusCode)
	}

	var result struct {
		UserID    string `json:"user_id"`
		Followers []struct {
			FollowerID string    `json:"follower_id"`
			CreatedAt  time.Time `json:"created_at"`
		} `json:"followers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	// Extract only follower IDs
	followerIDs := make([]string, 0, len(result.Followers))
	for _, f := range result.Followers {
		followerIDs = append(followerIDs, f.FollowerID)
	}

	return followerIDs, nil
}
