package main

import (
	"context"
	histodb "histo-db/internal/historical-db"
	redisdb "histo-db/internal/redis_db"
)

func main() {
	rdb := redisdb.ConnectToRedis()
	ctx := context.Background()
	res := redisdb.ReadFromRedis(ctx, rdb)
	redisdb.ParseFromStreamResult(res)
	db := histodb.ConnectToHistoricalDB()

	defer db.Close()
}
