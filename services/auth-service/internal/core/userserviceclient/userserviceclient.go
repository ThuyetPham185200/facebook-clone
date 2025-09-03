package userserviceclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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
		log.Printf("[UserServiceClient] failed to marshal request body (username=%s, email=%s): %v", username, email, err)
		return "", err
	}

	// set timeout = 3s (có thể config)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("[UserServiceClient] failed to create request: %v", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	log.Printf("[UserServiceClient] sending CreateUserProfile request to %s with username=%s, email=%s", url, username, email)
	resp, err := u.Client.Do(req)
	if err != nil {
		// check if timeout
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("[UserServiceClient] request timeout to %s", url)
		} else {
			log.Printf("[UserServiceClient] error calling UserService: %v", err)
		}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create user, status: %d", resp.StatusCode)
	}

	var res createUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Printf("[UserServiceClient] failed to decode response: %v", err)
		return "", err
	}

	log.Printf("[UserServiceClient] successfully created user with ID=%s", res.UserID)
	return res.UserID, nil
}

// CheckUserExists calls UserService API to check if a username/email already exists
func (u *UserService) CheckUserExists(username, email string) (bool, bool, error) {
	url := fmt.Sprintf("%s/users/exists", u.BaseURL)

	reqBody, err := json.Marshal(map[string]string{
		"username": username,
		"email":    email,
	})
	if err != nil {
		return false, false, err
	}

	resp, err := u.Client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return false, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, false, fmt.Errorf("failed to check user exists, status: %d", resp.StatusCode)
	}

	var res struct {
		ExistsUsername bool `json:"exists_username"`
		ExistsEmail    bool `json:"exists_email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, false, err
	}

	return res.ExistsUsername, res.ExistsEmail, nil
}
