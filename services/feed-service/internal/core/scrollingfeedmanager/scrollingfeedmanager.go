package scrollingfeedmanager

import (
	"feedservice/internal/core/userserviceclient"
	"feedservice/internal/infra/redisclient"
	"feedservice/internal/infra/store"
)

type SrollingFeedManager struct {
	PostStore         *store.PostStore
	userserviceclient *userserviceclient.UserService
	redisclient       *redisclient.RedisClient
}
