// internal/limiter/ip_limiter.go
package limiter

import "rate-limiter/algorithm"

type IPRateLimiter struct {
	algorithm algorithm.RateLimitAlgorithm
	limit     int
}

func NewIPRateLimiter(algo algorithm.RateLimitAlgorithm, limit int) *IPRateLimiter {
	return &IPRateLimiter{
		algorithm: algo,
		limit:     limit,
	}
}

func (l *IPRateLimiter) Allow(ip string) bool {
	return l.algorithm.Allow(ip, l.limit)
}
