package tables

import dbclient "feedservice/internal/infra/postgresclient"

// PostMediaTable kế thừa BaseTable
type PostMediaTable struct {
	dbclient.BaseTable
}

// NewPostMediaTable khởi tạo table post_media
func NewPostMediaTable(client *dbclient.PostgresClient) *PostMediaTable {
	return &PostMediaTable{
		BaseTable: dbclient.BaseTable{
			Client:    client,
			TableName: "post_media",
			Columns: map[string]string{
				"post_id":  "UUID NOT NULL",
				"media_id": "UUID NOT NULL",
			},
			Constraints: []string{
				"PRIMARY KEY (post_id, media_id)",
				"FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE",
				"FOREIGN KEY (media_id) REFERENCES medias(media_id) ON DELETE CASCADE",
			},
		},
	}
}
