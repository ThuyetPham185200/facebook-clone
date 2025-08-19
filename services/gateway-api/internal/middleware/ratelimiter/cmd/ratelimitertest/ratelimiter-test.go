package main

import (
	"fmt"
	ratelimiter "gatewayapi/internal/middleware/ratelimiter" // <--- dùng module name từ go.mod
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// limits := map[string]int{
	// 	"posts":    1,
	// 	"likes":    5,
	// 	"feeds":    5,
	// 	"comments": 5,
	// 	"follow":   5,
	// }

	// algo := algorithm.NewTokenBucketAlgorithm(nil)
	// ipLimiter := factory.CreateLimiter(factory.IP, algo, 10000, nil).(*limiter.IPRateLimiter)
	// featureLimiter := factory.CreateLimiter(factory.Feature, algo, 0, limits).(*limiter.FeatureRateLimiter)

	bo_ratelimiter := ratelimiter.NewRateLimiter()

	ip := "192.168.1.1"
	user := "user123"
	duration := 5 * time.Second
	features := []string{"posts", "likes", "feeds", "comments", "follow"}

	end := time.Now().Add(duration)
	var wg sync.WaitGroup

	for time.Now().Before(end) {
		// Feature requests
		for _, feature := range features {
			reqCount := rand.Intn(10) + 1
			for i := 0; i < reqCount; i++ {
				wg.Add(1)
				go func(f string, idx int) {
					defer wg.Done()
					now := time.Now().Format("15:04:05.000")
					key := user + ":" + f
					if !bo_ratelimiter.FeatureLimiter.Allow(key) {
						fmt.Printf("[%s] Request %d: %s rate limit exceeded\n", now, idx, f)
					} else {
						fmt.Printf("[%s] Request %d: %s allowed\n", now, idx, f)
					}
				}(feature, i)
			}
		}

		// IP requests
		ipReq := rand.Intn(14901) + 100
		for i := 0; i < ipReq; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				now := time.Now().Format("15:04:05.000")
				if !bo_ratelimiter.IPLimiter.Allow(ip) {
					fmt.Printf("[%s] Request %d: IP rate limit exceeded\n", now, idx)
				}
			}(i)
		}

		// Chạy batch 1 giây
		time.Sleep(time.Second)
	}

	wg.Wait()
	fmt.Println("=== Stress test finished ===")
}
