package dbclient

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresClient struct {
	DB *sql.DB
}

func NewPostgresClient(host, port, user, password, dbname string) *PostgresClient {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Không thể mở kết nối DB: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Không thể ping DB: %v", err)
	}

	log.Println("✅ Kết nối PostgreSQL thành công!")
	return &PostgresClient{DB: db}
}

func (pc *PostgresClient) Close() {
	if pc.DB != nil {
		pc.DB.Close()
	}
}
