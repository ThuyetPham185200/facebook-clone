package userserviceclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserService struct {
	BaseURL string
	Client  *http.Client
}

func NewUserService(baseURL string) *UserService {
	return &UserService{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

type createUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type createUserResponse struct {
	UserID string `json:"user_id"`
}

// CreateUserProfile calls UserService API to create user
func (u *UserService) CreateUserProfile(username, email string) (string, error) {
	url := fmt.Sprintf("%s/users", u.BaseURL)

	reqBody, err := json.Marshal(&createUserRequest{
		Username: username,
		Email:    email,
	})
	if err != nil {
		return "", err
	}

	resp, err := u.Client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create user, status: %d", resp.StatusCode)
	}

	var res createUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.UserID, nil
}
