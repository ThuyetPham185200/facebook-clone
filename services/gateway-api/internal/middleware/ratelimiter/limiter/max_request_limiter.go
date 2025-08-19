// internal/limiter/max_request_limiter.go
package limiter

import (
	"gatewayapi/internal/middleware/ratelimiter/algorithm"
)

type MaxRequestLimiter struct {
	algorithm algorithm.RateLimitAlgorithm
	limit     int
}

func NewMaxRequestLimiter(algo algorithm.RateLimitAlgorithm, limits int) *MaxRequestLimiter {
	return &MaxRequestLimiter{
		algorithm: algo,
		limit:     limits,
	}
}

func (m *MaxRequestLimiter) Allow(key string) bool {
	return m.algorithm.Allow(key, m.limit)
}
