// internal/algorithm/token_bucket_algo.go
package algorithm

import (
	"rate-limiter/repository"
	"sync"
	"time"
)

type TokenBucketAlgorithm struct {
	buckets map[string]*TokenBucket
	mu      sync.Mutex
	repo    repository.RedisStateRepository
}

type TokenBucket struct {
	Capacity   int
	Tokens     int
	RefillRate int
	LastRefill time.Time
}

func NewTokenBucketAlgorithm(repo repository.RedisStateRepository) *TokenBucketAlgorithm {
	return &TokenBucketAlgorithm{
		buckets: make(map[string]*TokenBucket),
		repo:    repo,
	}
}

func (t *TokenBucketAlgorithm) Allow(key string, limit int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	bucket, ok := t.buckets[key]
	if !ok {
		bucket = &TokenBucket{
			Capacity:   limit,
			Tokens:     limit,
			RefillRate: limit,
			LastRefill: time.Now(),
		}
		t.buckets[key] = bucket
	}

	now := time.Now()
	elapsed := now.Sub(bucket.LastRefill).Seconds()
	bucket.Tokens += int(elapsed * float64(bucket.RefillRate))
	if bucket.Tokens > bucket.Capacity {
		bucket.Tokens = bucket.Capacity
	}
	bucket.LastRefill = now

	if bucket.Tokens > 0 {
		bucket.Tokens--
		return true
	}
	return false
}
