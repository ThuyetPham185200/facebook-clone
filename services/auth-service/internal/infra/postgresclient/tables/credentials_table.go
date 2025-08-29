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
				"id":            "UUID PRIMARY KEY", // bỏ DEFAULT gen_random_uuid()
				"password_hash": "VARCHAR(255) NOT NULL",
				"mfa_secret":    "BYTEA", // VARBINARY trong Postgres là BYTEA
				"status":        "VARCHAR(16) NOT NULL DEFAULT 'active' CHECK (status IN ('active','locked','disabled'))",
				"created_at":    "TIMESTAMP DEFAULT now()",
				"updated_at":    "TIMESTAMP DEFAULT now()",
			},
		},
	}
}
