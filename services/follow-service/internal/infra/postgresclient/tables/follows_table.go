package tables

import dbclient "followservice/internal/infra/postgresclient"

// FollowsTable kế thừa BaseTable
type FollowsTable struct {
	dbclient.BaseTable
}

// NewFollowsTable khởi tạo table follows
func NewFollowsTable(client *dbclient.PostgresClient) *FollowsTable {
	return &FollowsTable{
		BaseTable: dbclient.BaseTable{
			Client:    client,
			TableName: "follows",
			Columns: map[string]string{
				"follower_id": "UUID NOT NULL",
				"followee_id": "UUID NOT NULL",
				"created_at":  "TIMESTAMP DEFAULT now()",
			},
			Constraints: []string{
				"PRIMARY KEY (follower_id, followee_id)",
				"CREATE INDEX follower_id_idx ON follows(follower_id)",
				"CREATE INDEX followee_id_idx ON follows(followee_id)",
			},
		},
	}
}
