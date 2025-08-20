// internal/algorithm/token_bucket_algo.go
package algorithm

import (
	"sync"
	"time"
)

type TokenBucketAlgorithm struct {
	buckets map[string]*TokenBucket
	mu      sync.Mutex
}

type TokenBucket struct {
	Capacity   int
	Tokens     int
	RefillRate int
	LastRefill time.Time
}

func NewTokenBucketAlgorithm() *TokenBucketAlgorithm {
	return &TokenBucketAlgorithm{
		buckets: make(map[string]*TokenBucket),
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

	// // In toàn bộ giá trị trong bucket
	// fmt.Printf("🚀 New Bucket Created:\n"+
	// 	"  Capacity: %d\n"+
	// 	"  Tokens: %d\n"+
	// 	"  RefillRate: %d\n"+
	// 	"  LastRefill: %v\n",
	// 	bucket.Capacity,
	// 	bucket.Tokens,
	// 	bucket.RefillRate,
	// 	bucket.LastRefill,
	// )

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
