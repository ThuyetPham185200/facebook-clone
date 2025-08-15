// internal/factory/limiter_factory.go
package factory

import (
	"rate-limiter/algorithm"
	"rate-limiter/limiter"
)

type LimiterType string

const (
	IP      LimiterType = "ip"
	Feature LimiterType = "feature"
)

func CreateLimiter(t LimiterType, algo algorithm.RateLimitAlgorithm, limit int, featureLimits map[string]int) limiter.RateLimiter {
	switch t {
	case IP:
		return limiter.NewIPRateLimiter(algo, limit)
	case Feature:
		return limiter.NewFeatureRateLimiter(algo, featureLimits)
	default:
		return nil
	}
}
