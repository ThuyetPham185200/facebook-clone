package ratelimiter

import (
	"rate-limiter/algorithm"
	"rate-limiter/factory"
	"rate-limiter/limiter"
)

type RateLimiter struct {
	Limits         map[string]int
	Algo           *algorithm.TokenBucketAlgorithm
	IPLimiter      *limiter.IPRateLimiter
	FeatureLimiter *limiter.FeatureRateLimiter
}

func NewRateLimiter() *RateLimiter {
	limits := map[string]int{
		"posts":    1,
		"likes":    5,
		"feeds":    5,
		"comments": 5,
		"follow":   5,
	}

	algo := algorithm.NewTokenBucketAlgorithm(nil)

	return &RateLimiter{
		Limits:         limits,
		Algo:           algo,
		IPLimiter:      factory.CreateLimiter(factory.IP, algo, 10000, nil).(*limiter.IPRateLimiter),
		FeatureLimiter: factory.CreateLimiter(factory.Feature, algo, 0, limits).(*limiter.FeatureRateLimiter),
	}
}

func NewRateLimiterWithConfig(limitsInput map[string]int, algoInput *algorithm.TokenBucketAlgorithm, ipLimiterInput *limiter.IPRateLimiter, featureLimiterInput *limiter.FeatureRateLimiter) *RateLimiter {
	return &RateLimiter{
		Limits:         limitsInput,
		Algo:           algoInput,
		IPLimiter:      ipLimiterInput,
		FeatureLimiter: featureLimiterInput,
	}
}
