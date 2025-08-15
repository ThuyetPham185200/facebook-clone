// internal/limiter/feature_limiter.go
package limiter

import (
	"rate-limiter/algorithm"
	"strings"
)

type FeatureRateLimiter struct {
	algorithm     algorithm.RateLimitAlgorithm
	featureLimits map[string]int
}

func NewFeatureRateLimiter(algo algorithm.RateLimitAlgorithm, limits map[string]int) *FeatureRateLimiter {
	return &FeatureRateLimiter{
		algorithm:     algo,
		featureLimits: limits,
	}
}

func (f *FeatureRateLimiter) Allow(key string) bool {
	// key = "userID:feature", tách ra nếu cần
	parts := strings.Split(key, ":")
	if len(parts) != 2 {
		return false
	}
	userID, feature := parts[0], parts[1]

	limit, ok := f.featureLimits[feature]
	if !ok {
		limit = 1
	}

	return f.algorithm.Allow(userID+":"+feature, limit)
}
