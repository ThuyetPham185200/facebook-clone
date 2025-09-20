package tables

import dbclient "feedservice/internal/infra/postgresclient"

// PostsTable kế thừa BaseTable
type PostsTable struct {
	dbclient.BaseTable
}

// NewSessionsTable khởi tạo table sessions
func NewPostsTable(client *dbclient.PostgresClient) *PostsTable {
	return &PostsTable{
		BaseTable: dbclient.BaseTable{
			Client:    client,
			TableName: "posts",
			Columns: map[string]string{
				"id":                 "UUID PRIMARY KEY",
				"user_id":            "UUID NOT NULL",
				"refresh_token":      "TEXT UNIQUE NOT NULL",
				"status":             "VARCHAR(16) NOT NULL DEFAULT 'active'",
				"refresh_expires_at": "TIMESTAMP NOT NULL",
				"created_at":         "TIMESTAMP DEFAULT now()",
				"updated_at":         "TIMESTAMP DEFAULT now()",
			},
			Constraints: []string{
				// Mỗi session phải gắn với user tồn tại
				"FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE",

				// status chỉ có thể nhận một trong các giá trị hợp lệ
				"CHECK (status IN ('active','revoked','expired'))",
			},
		},
	}
}
