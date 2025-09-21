package tables

import dbclient "feedservice/internal/infra/postgresclient"

// PostsTable kế thừa BaseTable
type PostsTable struct {
	dbclient.BaseTable
}

// NewPostsTable khởi tạo table posts
func NewPostsTable(client *dbclient.PostgresClient) *PostsTable {
	return &PostsTable{
		BaseTable: dbclient.BaseTable{
			Client:    client,
			TableName: "posts",
			Columns: map[string]string{
				"post_id":    "UUID PRIMARY KEY",
				"user_id":    "UUID NOT NULL",
				"content":    "TEXT",
				"media_url":  "TEXT",
				"media_type": "VARCHAR(20)", // 'image' hoặc 'video'
				"created_at": "TIMESTAMP NOT NULL DEFAULT now()",
				"updated_at": "TIMESTAMP",
				"is_deleted": "BOOLEAN NOT NULL DEFAULT FALSE",
			},
			Constraints: []string{
				"FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE",
				"CREATE INDEX idx_posts_user_created ON posts(user_id, created_at DESC)",
				"CREATE INDEX idx_posts_created ON posts(created_at DESC)",
			},
		},
	}
}
