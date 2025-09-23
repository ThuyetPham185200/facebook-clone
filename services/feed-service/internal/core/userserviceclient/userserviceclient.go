package userserviceclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserService struct {
	BaseURL string
	Client  *http.Client
}

func NewUserServiceClient(baseURL string) *UserService {
	return &UserService{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

type UserProfileResponse struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Gender    string `json:"gender"`
}

func (u *UserService) GetUserProfile(userID string) (UserProfileResponse, error) {
	url := fmt.Sprintf("%s/users/%s", u.BaseURL, userID)

	req, err := http.NewRequest(http.MethodGet, url, nil) // ✅ Use GET
	if err != nil {
		return UserProfileResponse{}, fmt.Errorf("[UserServiceClient] failed to build get request: %w", err)
	}

	resp, err := u.Client.Do(req)
	if err != nil {
		return UserProfileResponse{}, fmt.Errorf("[UserServiceClient] failed to get user profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UserProfileResponse{}, fmt.Errorf("[UserServiceClient] unexpected status code: %d", resp.StatusCode)
	}

	var res UserProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return UserProfileResponse{}, fmt.Errorf("[UserServiceClient] failed to decode response: %w", err)
	}

	return res, nil
}
