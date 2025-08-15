// internal/algorithm/algorithm.go
package algorithm

// RateLimitAlgorithm Strategy Pattern interface
type RateLimitAlgorithm interface {
	Allow(key string, limit int) bool
}
