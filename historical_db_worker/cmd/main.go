package main

import (
	"context"
	"fmt"
	histodb "histo-db/internal/historical-db"
	redisdb "histo-db/internal/redis_db"
	"time"
)

func main() {
	rdb := redisdb.ConnectToRedis()
	ctx := context.Background()
	res := redisdb.ReadFromRedis(ctx, rdb)
	reading := redisdb.ParseFromStreamResult(res, ctx, rdb)
	db := histodb.ConnectToHistoricalDB()

	fails := 0

	for {
		id, err := histodb.AddTempReading(db, reading)
		fmt.Printf("Inserted row #%d\n", id)
		if err != nil {
			panic(err)
		}
		if id == 0 {
			fails++
		}
		if fails == 10 {
			break
		}
		time.Sleep(2 * time.Second)
	}

	defer db.Close()
}
