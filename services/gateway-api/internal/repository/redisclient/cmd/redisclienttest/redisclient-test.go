package main

import (
	"fmt"
	"gatewayapi/internal/repository/redisclient"
	"time"
)

type TokenBucket struct {
	Tokens     int
	LastRefill int64
}

func main() {
	rdb := redisclient.InitSingleton("127.0.0.1:6379", "", 0)

	// Test SetInt & GetInt
	rdb.SetInt("views", 100, 0)
	val, _ := rdb.GetInt("views")
	fmt.Println("📊 views =", val)

	// Test IncrBy
	newVal, _ := rdb.IncrBy("views", 5)
	fmt.Println("📊 views after +5 =", newVal)

	// Test DecrBy
	newVal, _ = rdb.DecrBy("views", 3)
	fmt.Println("📊 views after -3 =", newVal)

	// -------------------------------------
	key := "ratelimit:user:123:post"

	// save bucket to Redis
	bucket := TokenBucket{
		Tokens:     5,
		LastRefill: time.Now().Unix(),
	}

	err := rdb.HSet(key, "tokens", bucket.Tokens)
	if err != nil {
		panic(err)
	}
	err = rdb.HSet(key, "last_refill", bucket.LastRefill)
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Saved bucket to Redis")

	// read all fields back
	data, err := rdb.HGetAll(key)
	if err != nil {
		panic(err)
	}

	fmt.Println("📦 Read from Redis:", data)

	// parse back into TokenBucket
	tokens := data["tokens"]
	lastRefill := data["last_refill"]
	fmt.Printf("👉 tokens=%s, last_refill=%s\n", tokens, lastRefill)
	rdb.Close()

}
