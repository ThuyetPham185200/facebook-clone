package main

import (
	"fmt"

	redisclient "redis-client"
)

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

	rdb.Close()
}
