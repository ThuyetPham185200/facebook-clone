package main

import (
	"fmt"
	"log"
	dbclient "userservice/internal/infra/postgresclient"
	"userservice/internal/infra/postgresclient/tables"
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
	usersTable := tables.NewUserTable(client)

	if !client.SearchTable(usersTable.TableName) {
		fmt.Printf("%s NOT EXIST - CREATION PROCESS STARTING\n", usersTable.TableName)
		usersTable.CreateTable()
	} else {
		fmt.Printf("%s EXISTED\n", usersTable.TableName)
	}

	// Lấy tất cả rules
	rows, err := usersTable.GetAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range rows {
		fmt.Println(row)
	}
}
