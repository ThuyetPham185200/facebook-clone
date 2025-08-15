package dbclient

import (
	"fmt"
	"log"
	"strings"
)

type BaseTable struct {
	Client    *PostgresClient
	TableName string
	Columns   map[string]string // column_name -> type (VD: "id": "SERIAL PRIMARY KEY")
}

// CreateTable tạo bảng dựa trên metadata
func (bt *BaseTable) CreateTable() {
	var cols []string
	for col, typ := range bt.Columns {
		cols = append(cols, fmt.Sprintf("%s %s", col, typ))
	}
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (%s)`,
		bt.TableName,
		strings.Join(cols, ", "),
	)
	_, err := bt.Client.DB.Exec(query)
	if err != nil {
		log.Fatalf("❌ Lỗi tạo bảng %s: %v", bt.TableName, err)
	}
	log.Printf("✅ Bảng %s sẵn sàng.", bt.TableName)
}

// Insert thêm dữ liệu vào bảng
func (bt *BaseTable) Insert(values map[string]interface{}) {
	cols := []string{}
	vals := []interface{}{}
	placeholders := []string{}

	i := 1
	for col, val := range values {
		cols = append(cols, col)
		vals = append(vals, val)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		i++
	}

	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`,
		bt.TableName,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)
	_, err := bt.Client.DB.Exec(query, vals...)
	if err != nil {
		log.Printf("❌ Lỗi insert vào %s: %v", bt.TableName, err)
	} else {
		log.Printf("✅ Insert thành công vào %s", bt.TableName)
	}
}

// GetAll lấy tất cả dữ liệu
func (bt *BaseTable) GetAll() {
	query := fmt.Sprintf(`SELECT * FROM %s`, bt.TableName)
	rows, err := bt.Client.DB.Query(query)
	if err != nil {
		log.Printf("❌ Lỗi query %s: %v", bt.TableName, err)
		return
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))

	for rows.Next() {
		for i := range cols {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		rowData := map[string]interface{}{}
		for i, col := range cols {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				rowData[col] = string(b)
			} else {
				rowData[col] = val
			}
		}
		log.Println(rowData)
	}
}
