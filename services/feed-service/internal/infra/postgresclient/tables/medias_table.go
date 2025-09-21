package tables

import dbclient "feedservice/internal/infra/postgresclient"

// MediasTable kế thừa BaseTable
type MediasTable struct {
	dbclient.BaseTable
}

// NewMediasTable khởi tạo table medias
func NewMediasTable(client *dbclient.PostgresClient) *MediasTable {
	return &MediasTable{
		BaseTable: dbclient.BaseTable{
			Client:    client,
			TableName: "medias",
			Columns: map[string]string{
				"media_id":   "UUID PRIMARY KEY",
				"user_id":    "UUID NOT NULL",
				"media_type": "VARCHAR(10) NOT NULL CHECK (media_type IN ('image','video'))",
				"url":        "TEXT NOT NULL", // S3/CDN URL
				"status":     "VARCHAR(10) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','uploaded','failed'))",
				"created_at": "TIMESTAMP NOT NULL DEFAULT now()",
			},
			Constraints: []string{
				"FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE",
				"CREATE INDEX idx_medias_user_status ON medias(user_id, status)",
			},
		},
	}
}
