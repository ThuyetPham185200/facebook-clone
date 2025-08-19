package ratelimiter

import (
	"gatewayapi/internal/middleware/ratelimiter/algorithm"
	"gatewayapi/internal/middleware/ratelimiter/factory"
	"gatewayapi/internal/middleware/ratelimiter/limiter"
)

type RateLimiter struct {
	Limits         map[string]int
	Algo           *algorithm.TokenBucketAlgorithm
	IPLimiter      *limiter.IPRateLimiter
	FeatureLimiter *limiter.FeatureRateLimiter
	MaxReqLimiter  *limiter.MaxRequestLimiter
}

func NewRateLimiter() *RateLimiter {
	limits := map[string]int{
		"posts":    1,
		"likes":    5,
		"feeds":    5,
		"comments": 5,
		"follow":   5,
	}

	algo := algorithm.NewTokenBucketAlgorithm()

	return &RateLimiter{
		Limits:         limits,
		Algo:           algo,
		IPLimiter:      factory.CreateLimiter(factory.IP, algo, 10, nil).(*limiter.IPRateLimiter),
		FeatureLimiter: factory.CreateLimiter(factory.Feature, algo, 0, limits).(*limiter.FeatureRateLimiter),
		MaxReqLimiter:  factory.CreateLimiter(factory.MaxRequest, algo, 2000, nil).(*limiter.MaxRequestLimiter),
	}
}

func NewRateLimiterWithConfig(limitsInput map[string]int, algoInput *algorithm.TokenBucketAlgorithm, ipLimiterInput *limiter.IPRateLimiter,
	featureLimiterInput *limiter.FeatureRateLimiter, maxReqLimiterInput *limiter.MaxRequestLimiter) *RateLimiter {
	return &RateLimiter{
		Limits:         limitsInput,
		Algo:           algoInput,
		IPLimiter:      ipLimiterInput,
		FeatureLimiter: featureLimiterInput,
		MaxReqLimiter:  maxReqLimiterInput,
	}
}
