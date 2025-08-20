package main

import (
	"fmt"
	dbclient "gatewayapi/internal/repository/postgresclient"
	"gatewayapi/internal/repository/postgresclient/tables"
	gmodel "gatewayapi/model"
	"log"
)

type RateLimitAPIModel struct {
	Gmodel gmodel.GatewayModel
}

func NewRateLimitAPIModel() *RateLimitAPIModel {
	r := &RateLimitAPIModel{}
	r.Gmodel = *gmodel.NewGatewayModel()
	return r
}

func (r *RateLimitAPIModel) RateLimitAdapt() []map[string]interface{} {
	var rules []map[string]interface{}
	for action, limit := range r.Gmodel.RateLimitMap {
		rules = append(rules, map[string]interface{}{
			"action":           action,
			"limit_per_second": limit,
		})
	}
	return rules
}

func main() {
	client := dbclient.NewPostgresClient(
		"localhost", // IP
		"5432",      // Port
		"taopq",     // user_name
		"123456a@",  // password
		"mydb",      // db
	)
	defer client.Close()

	r := NewRateLimitAPIModel()

	// Tạo bảng rules
	rulesTable := tables.NewRateLimiterRulesTable(client)

	if !client.SearchTable(rulesTable.TableName) {
		fmt.Printf("%s NOT EXIST - CREATION PROCESS STARTING\n", rulesTable.TableName)
		rulesTable.CreateTable()
		rules := r.RateLimitAdapt()
		rules = append(rules, map[string]interface{}{
			"action":           "max_requests",
			"limit_per_second": 10000,
		})
		rules = append(rules, map[string]interface{}{
			"action":           "requests_per_ip",
			"limit_per_second": 10,
		})
		for _, rule := range rules {
			rulesTable.Insert(rule)
		}
	} else {
		fmt.Printf("%s EXISTED\n", rulesTable.TableName)
	}

	// Lấy tất cả rules
	rows, err := rulesTable.GetAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range rows {
		fmt.Println(row)
	}
}
