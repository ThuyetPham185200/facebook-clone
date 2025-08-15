package main

import (
	dbclient "postgresclient"
	"postgresclient/tables"
)

func main() {
	client := dbclient.NewPostgresClient(
		"localhost",
		"5432",
		"taopq",
		"123456a@",
		"mydb",
	)
	defer client.Close()

	// Tạo bảng rule
	rulesTable := tables.NewRateLimiterRulesTable(client)
	rulesTable.CreateTable()

	// Insert các rule
	rules := []map[string]interface{}{
		{"action": "max_requests", "limit_per_second": 10000, "description": "Maximum 10,000 requests per second"},
		{"action": "requests_per_ip", "limit_per_second": 10, "description": "10 requests per second per IP"},
		{"action": "post", "limit_per_second": 1, "description": "1 post per second per user"},
		{"action": "like", "limit_per_second": 5, "description": "5 likes per second"},
		{"action": "feed", "limit_per_second": 5, "description": "Get max 5 feeds per second per person"},
		{"action": "comment", "limit_per_second": 5, "description": "Max 5 comments per second"},
		{"action": "follow_unfollow", "limit_per_second": 5, "description": "Follow/unfollow max 5 times per second"},
	}

	for _, rule := range rules {
		rulesTable.Insert(rule)
	}

	// Lấy tất cả rules
	rulesTable.GetAll()
}
