// internal/algorithm/token_bucket_algo.go
package algorithm

import (
	redisclient "gatewayapi/internal/repository/redisclient"
	"strconv"
	"sync"
	"time"
)

type TokenBucketAlgorithm struct {
	buckets     map[string]*TokenBucket
	mu          sync.Mutex
	redisClient *redisclient.RedisClient
}

type TokenBucket struct {
	Capacity   int
	Tokens     int
	RefillRate int
	LastRefill time.Time
}

func NewTokenBucketAlgorithm(client *redisclient.RedisClient) *TokenBucketAlgorithm {
	return &TokenBucketAlgorithm{
		buckets:     make(map[string]*TokenBucket),
		redisClient: client,
	}
}

func (t *TokenBucketAlgorithm) Allow(key string, limit int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	redisKey := key
	var bucket *TokenBucket

	// 1. Thử lấy bucket từ Redis
	data, err := t.redisClient.HGetAll(redisKey)
	if err == nil && len(data) > 0 {
		capacity, _ := strconv.Atoi(data["capacity"])
		tokens, _ := strconv.Atoi(data["tokens"])
		refillRate, _ := strconv.Atoi(data["refill_rate"])
		lastRefillUnix, _ := strconv.ParseInt(data["last_refill"], 10, 64)

		//fmt.Printf("📤 Data from Redis [%s]: %+v\n", redisKey, data)

		bucket = &TokenBucket{
			Capacity:   capacity,
			Tokens:     tokens,
			RefillRate: refillRate,
			LastRefill: time.Unix(lastRefillUnix, 0),
		}
	} else {
		// 2. Nếu Redis chưa có thì tạo mới
		bucket = &TokenBucket{
			Capacity:   limit,
			Tokens:     limit,
			RefillRate: limit,
			LastRefill: time.Now(),
		}
		// lưu bucket mới vào Redis
		// fmt.Printf("📥 Save bucket to Redis [%s]: capacity=%d, tokens=%d, refill_rate=%d, last_refill=%d\n",
		// 	redisKey,
		// 	bucket.Capacity,
		// 	bucket.Tokens,
		// 	bucket.RefillRate,
		// 	bucket.LastRefill.Unix(),
		// )
		t.redisClient.HSet(redisKey, "capacity", bucket.Capacity)
		t.redisClient.HSet(redisKey, "tokens", bucket.Tokens)
		t.redisClient.HSet(redisKey, "refill_rate", bucket.RefillRate)
		t.redisClient.HSet(redisKey, "last_refill", bucket.LastRefill.Unix())
	}

	// 3. Tính refill
	now := time.Now()
	elapsed := now.Sub(bucket.LastRefill).Seconds()

	// lazy refill: tính lại số token dựa trên thời gian trôi qua
	refilled := bucket.Tokens + int(elapsed*float64(bucket.RefillRate))
	if refilled > bucket.Capacity {
		refilled = bucket.Capacity
	}
	bucket.Tokens = refilled
	bucket.LastRefill = now

	// 4. Kiểm tra và trừ token
	allowed := false
	if bucket.Tokens > 0 {
		bucket.Tokens--
		allowed = true
	}

	// 5. Update lại Redis
	t.redisClient.HSet(redisKey, "tokens", bucket.Tokens)
	t.redisClient.HSet(redisKey, "last_refill", bucket.LastRefill.Unix())

	// Lưu lại vào local cache (optional)
	t.buckets[key] = bucket

	return allowed
}
