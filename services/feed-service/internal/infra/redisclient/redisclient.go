package redisclient

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

var (
	instance *RedisClient
	once     sync.Once
	ctx      = context.Background()
)

// InitSingleton - khởi tạo 1 lần duy nhất
func InitSingleton(addr, password string, db int) *RedisClient {
	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password, // "" nếu không có password
			DB:       db,
		})

		// Test kết nối
		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			panic(fmt.Sprintf("❌ Không kết nối được Redis: %v", err))
		}

		fmt.Println("✅ Redis connected:", addr)

		instance = &RedisClient{
			client: rdb,
		}
	})
	return instance
}

// GetInstance - lấy instance Redis
func GetInstance() *RedisClient {
	if instance == nil {
		panic("⚠ Redis chưa được init! Gọi InitSingleton trước.")
	}
	return instance
}

// Close - đóng kết nối Redis
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// GetClient - lấy raw *redis.Client nếu cần
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

// SetKey - set key với TTL
func (r *RedisClient) SetKey(key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// GetKey - lấy value
func (r *RedisClient) GetKey(key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// IncrKey - tăng giá trị integer
func (r *RedisClient) IncrKey(key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// KeyExists - kiểm tra key có tồn tại trong Redis
func (r *RedisClient) KeyExists(key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DeleteKey - xóa 1 key
func (r *RedisClient) DeleteKey(key string) error {
	return r.client.Del(ctx, key).Err()
}

// ExpireKey - đặt lại TTL cho 1 key
func (r *RedisClient) ExpireKey(key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}

// GetTTL - lấy TTL còn lại của 1 key
func (r *RedisClient) GetTTL(key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}
