package tables

import dbclient "userservice/internal/infra/postgresclient"

// UsersTable kế thừa BaseTable
type UserTable struct {
	dbclient.BaseTable
}

// NewUserTable khởi tạo table users với schema chuẩn (User Profile)
func NewUserTable(client *dbclient.PostgresClient) *UserTable {
	return &UserTable{
		BaseTable: dbclient.BaseTable{
			Client:    client,
			TableName: "users",
			Columns: map[string]string{
				"user_id":       "UUID PRIMARY KEY",
				"username":      "VARCHAR(50) NOT NULL",
				"email":         "VARCHAR(255) NOT NULL",
				"bio":           "TEXT",
				"gender":        "VARCHAR(16)",
				"date_of_birth": "DATE",
				"avatar_url":    "VARCHAR(255)",
				"is_deleted":    "BOOLEAN NOT NULL DEFAULT FALSE",
				"created_at":    "TIMESTAMP NOT NULL DEFAULT now()",
				"updated_at":    "TIMESTAMP NOT NULL DEFAULT now()",
			},
			Constraints: []string{
				"UNIQUE (username)",
				"UNIQUE (email)",
				"CHECK (gender IN ('male','female','other'))",
			},
		},
	}
}
