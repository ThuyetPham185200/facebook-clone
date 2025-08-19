// dbclient/tables/rate_limiter_rules.go
package tables

import dbclient "gatewayapi/internal/repository/postgresclient"

type RateLimiterRulesTable struct {
	dbclient.BaseTable
}

func NewRateLimiterRulesTable(client *dbclient.PostgresClient) *RateLimiterRulesTable {
	return &RateLimiterRulesTable{
		BaseTable: dbclient.BaseTable{
			Client:    client,
			TableName: "rate_limiter_rules",
			Columns: map[string]string{
				"id":               "SERIAL PRIMARY KEY",
				"action":           "VARCHAR(50) NOT NULL",
				"limit_per_second": "INT NOT NULL",
				"description":      "TEXT",
			},
		},
	}
}
