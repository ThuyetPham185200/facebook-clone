package tables

import dbclient "authservice/internal/infra/postgresclient"

// SessionsTable kế thừa BaseTable
type SessionsTable struct {
	dbclient.BaseTable
}

// NewSessionsTable khởi tạo table sessions
func NewSessionsTable(client *dbclient.PostgresClient) *SessionsTable {
	return &SessionsTable{
		BaseTable: dbclient.BaseTable{
			Client:    client,
			TableName: "sessions",
			Columns: map[string]string{
				"id":            "UUID PRIMARY KEY",
				"user_id":       "UUID NOT NULL REFERENCES credentials(id) ON DELETE CASCADE",
				"access_token":  "TEXT UNIQUE NOT NULL",
				"refresh_token": "TEXT UNIQUE NOT NULL",
				"expires_at":    "TIMESTAMP NOT NULL",
				"created_at":    "TIMESTAMP DEFAULT now()",
				"updated_at":    "TIMESTAMP DEFAULT now()",
			},
		},
	}
}
