package tables

import dbclient "postgresclient"

// UsersTable kế thừa BaseTable
type UsersTable struct {
	dbclient.BaseTable
}

func NewUsersTable(client *dbclient.PostgresClient) *UsersTable {
	return &UsersTable{
		BaseTable: dbclient.BaseTable{
			Client:    client,
			TableName: "users",
			Columns: map[string]string{
				"id":    "SERIAL PRIMARY KEY",
				"name":  "TEXT NOT NULL",
				"email": "TEXT UNIQUE NOT NULL",
			},
		},
	}
}
