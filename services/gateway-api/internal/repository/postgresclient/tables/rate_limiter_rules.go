// dbclient/tables/rate_limiter_rules.go
package tables

import (
	"fmt"
	dbclient "gatewayapi/internal/repository/postgresclient"
	"strconv"
)

type RateLimiterRulesTable struct {
	dbclient.BaseTable
}

func NewRateLimiterRulesTable(client *dbclient.PostgresClient) *RateLimiterRulesTable {
	table := &RateLimiterRulesTable{
		BaseTable: dbclient.BaseTable{
			Client:    client,
			TableName: "rate_limiter_rules",
			Columns: map[string]string{
				"id":               "SERIAL PRIMARY KEY",
				"action":           "VARCHAR(50) NOT NULL",
				"limit_per_second": "INT NOT NULL",
			},
		},
	}
	return table
}

func (r *RateLimiterRulesTable) GetRateLimitMap() map[string]int {
	result := make(map[string]int)

	rows, err := r.GetAll()
	if err != nil {
		return result
	}

	for _, rule := range rows {
		action, ok1 := rule["action"].(string)
		if !ok1 {
			continue
		}

		var limitInt int
		switch v := rule["limit_per_second"].(type) {
		case int:
			limitInt = v
		case int64:
			limitInt = int(v)
		case float64:
			limitInt = int(v)
		case []byte: // phòng khi nó về []uint8
			parsed, err := strconv.Atoi(string(v))
			if err == nil {
				limitInt = parsed
			}
		default:
			fmt.Printf("⚠️ Unknown type for limit_per_second: %T\n", v)
		}

		if limitInt > 0 {
			result[action] = limitInt
		}
	}

	return result
}
