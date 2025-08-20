package repository

import (
	dbclient "gatewayapi/internal/repository/postgresclient"
	"gatewayapi/internal/repository/postgresclient/tables"
	redisclient "gatewayapi/internal/repository/redisclient"
	gmodel "gatewayapi/model"
)

type GateWayRepository struct {
	Redisrepo        *redisclient.RedisClient
	Postgresqlrepo   *dbclient.PostgresClient
	RateLimiterModel *gmodel.RateLimterModel
}

func NewGateWayRepository() *GateWayRepository {
	repo := &GateWayRepository{}
	repo.Redisrepo = redisclient.InitSingleton("127.0.0.1:6379", "", 0)
	repo.Postgresqlrepo = dbclient.NewPostgresClient(
		"localhost", // IP
		"5432",      // Port
		"taopq",     // user_name
		"123456a@",  // password
		"mydb",      // db
	)
	rulesTable := tables.NewRateLimiterRulesTable(repo.Postgresqlrepo)
	ret := rulesTable.GetRateLimitMap()
	repo.RateLimiterModel = gmodel.NewRateLimterModel(ret, ret["requests_per_ip"], ret["max_requests"])

	return repo
}

func (g *GateWayRepository) Close() {
	g.Postgresqlrepo.Close()
	g.Redisrepo.Close()
}
