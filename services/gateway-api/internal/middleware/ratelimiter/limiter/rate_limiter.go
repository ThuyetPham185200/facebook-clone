// internal/limiter/rate_limiter.go
package limiter

// RateLimiter là interface chung cho tất cả limiter
type RateLimiter interface {
	Allow(key string) bool
}
