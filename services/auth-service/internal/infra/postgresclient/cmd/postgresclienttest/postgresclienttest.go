package main

import (
	dbclient "authservice/internal/infra/postgresclient"
	"authservice/internal/infra/postgresclient/tables"
	"fmt"
	"log"
)

func main() {
	client := dbclient.NewPostgresClient(
		"localhost", // IP
		"5432",      // Port
		"taopq",     // user_name
		"123456a@",  // password
		"mydb",      // db
	)
	defer client.Close()

	// Tạo bảng rules
	cedentialsTable := tables.NewCredentialsTable(client)

	if !client.SearchTable(cedentialsTable.TableName) {
		fmt.Printf("%s NOT EXIST - CREATION PROCESS STARTING\n", cedentialsTable.TableName)
		cedentialsTable.CreateTable()
	} else {
		fmt.Printf("%s EXISTED\n", cedentialsTable.TableName)
	}

	// Lấy tất cả rules
	rows, err := cedentialsTable.GetAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range rows {
		fmt.Println(row)
	}
}
