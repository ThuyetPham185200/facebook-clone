package apis

import "net/http"

// ===== Struct cho 1 Endpoint =====
type Endpoint struct {
	Name        string
	Method      string
	Path        string
	RequireAuth bool
	RateLimit   int // per second
}

// ===== Struct cho Group Endpoint / Internal Service =====
type ServiceGroup struct {
	Name      string
	IP        string
	Port      int
	Endpoints []Endpoint
}

// ===== Auth Service =====
var AuthService = ServiceGroup{
	Name: "AuthService",
	IP:   "localhost",
	Port: 9090,
	Endpoints: []Endpoint{
		{
			Name:        "Register",
			Method:      http.MethodPost,
			Path:        "/register",
			RequireAuth: false,
			RateLimit:   1,
		},
		{
			Name:        "Login",
			Method:      http.MethodPost,
			Path:        "/login",
			RequireAuth: false,
			RateLimit:   1,
		},
		{
			Name:        "ChangePassword",
			Method:      http.MethodPut,
			Path:        "/me/password",
			RequireAuth: true,
			RateLimit:   1,
		},
		{
			Name:        "DeleteAccount",
			Method:      http.MethodDelete,
			Path:        "/me",
			RequireAuth: true,
			RateLimit:   1,
		},
	},
}

// ===== User Service =====
var UserService = ServiceGroup{
	Name: "UserService",
	IP:   "localhost",
	Port: 9091,
	Endpoints: []Endpoint{
		{
			Name:        "GetUserProfile",
			Method:      http.MethodGet,
			Path:        "/users/{user_id}",
			RequireAuth: true, // tùy chọn
			RateLimit:   5,
		},
		{
			Name:        "UpdateOwnProfile",
			Method:      http.MethodPatch,
			Path:        "/me",
			RequireAuth: true,
			RateLimit:   1,
		},
		{
			Name:        "SearchUsers",
			Method:      http.MethodGet,
			Path:        "/users",
			RequireAuth: true, // tùy hệ thống
			RateLimit:   3,
		},
	},
}

// ===== Posts Service =====
var PostsService = ServiceGroup{
	Name: "PostsService",
	IP:   "localhost",
	Port: 9092,
	Endpoints: []Endpoint{
		{
			Name:        "GetPost",
			Method:      http.MethodGet,
			Path:        "/posts/{post_id}",
			RequireAuth: true, // tùy chọn
			RateLimit:   5,
		},
		{
			Name:        "GetUserPosts",
			Method:      http.MethodGet,
			Path:        "/users/{user_id}/posts",
			RequireAuth: true,
			RateLimit:   5,
		},
		{
			Name:        "GetOwnPosts",
			Method:      http.MethodGet,
			Path:        "/me/posts",
			RequireAuth: true,
			RateLimit:   5,
		},
		{
			Name:        "CreatePost",
			Method:      http.MethodPost,
			Path:        "/posts",
			RequireAuth: true,
			RateLimit:   1,
		},
		{
			Name:        "UpdatePost",
			Method:      http.MethodPatch,
			Path:        "/posts/{post_id}",
			RequireAuth: true,
			RateLimit:   1,
		},
		{
			Name:        "DeletePost",
			Method:      http.MethodDelete,
			Path:        "/posts/{post_id}",
			RequireAuth: true,
			RateLimit:   1,
		},
	},
}

// ===== Reactions Service =====
var ReactionsService = ServiceGroup{
	Name: "ReactionsService",
	IP:   "localhost",
	Port: 9093,
	Endpoints: []Endpoint{
		{
			Name:        "GetReactions",
			Method:      http.MethodGet,
			Path:        "/posts/{post_id}/reactions",
			RequireAuth: true, // tùy chọn
			RateLimit:   5,
		},
		{
			Name:        "ReactToPost",
			Method:      http.MethodPost,
			Path:        "/posts/{post_id}/reactions",
			RequireAuth: true,
			RateLimit:   5,
		},
		{
			Name:        "RemoveReaction",
			Method:      http.MethodDelete,
			Path:        "/posts/{post_id}/reactions",
			RequireAuth: true,
			RateLimit:   1,
		},
	},
}
