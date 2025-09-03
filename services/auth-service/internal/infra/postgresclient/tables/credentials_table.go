package tables

import dbclient "authservice/internal/infra/postgresclient"

// UsersTable kế thừa BaseTable
type CredentialsTable struct {
	dbclient.BaseTable
}

// NewCredentialsTable khởi tạo table credentials với schema chuẩn
func NewCredentialsTable(client *dbclient.PostgresClient) *CredentialsTable {
	return &CredentialsTable{
		BaseTable: dbclient.BaseTable{
			Client:    client,
			TableName: "credentials",
			Columns: map[string]string{
				"id":            "UUID PRIMARY KEY",     // tự generate trong app
				"user_id":       "UUID NOT NULL UNIQUE", // 1-1 với users
				"password_hash": "VARCHAR(255) NOT NULL",
				"mfa_secret":    "BYTEA",
				"status":        "VARCHAR(16) NOT NULL DEFAULT 'active'",
				"created_at":    "TIMESTAMP DEFAULT now()",
				"updated_at":    "TIMESTAMP DEFAULT now()",
			},
			Constraints: []string{
				"FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE",
				"CHECK (status IN ('active','locked','disabled'))",
			},
		},
	}
}
