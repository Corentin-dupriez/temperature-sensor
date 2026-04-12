package main

import (
	"context"
	redisdb "histo-db/internal/redis_db"
)

func main() {
	rdb := redisdb.ConnectToRedis()
	ctx := context.Background()
	res := redisdb.ReadFromRedis(ctx, rdb)
	redisdb.ParseFromStreamResult(res)
}
