package ratelimiter

import (
	"gatewayapi/internal/middleware/ratelimiter/algorithm"
	"gatewayapi/internal/middleware/ratelimiter/factory"
	"gatewayapi/internal/middleware/ratelimiter/limiter"
	"gatewayapi/model"
)

type RateLimiter struct {
	Limits         map[string]int
	Algo           *algorithm.TokenBucketAlgorithm
	IPLimiter      *limiter.IPRateLimiter
	FeatureLimiter *limiter.FeatureRateLimiter
	MaxReqLimiter  *limiter.MaxRequestLimiter
}

func NewRateLimiter(model model.RateLimterModel) *RateLimiter {
	limits := model.FeatureLimits

	algo := algorithm.NewTokenBucketAlgorithm()

	return &RateLimiter{
		Limits:         limits,
		Algo:           algo,
		IPLimiter:      factory.CreateLimiter(factory.IP, algo, model.IPLimit, nil).(*limiter.IPRateLimiter),
		FeatureLimiter: factory.CreateLimiter(factory.Feature, algo, 0, limits).(*limiter.FeatureRateLimiter),
		MaxReqLimiter:  factory.CreateLimiter(factory.MaxRequest, algo, model.MaxRequestLimit, nil).(*limiter.MaxRequestLimiter),
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
