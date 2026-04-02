// Package redisdb is used to connect and write to the Redis stream
package redisdb

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func ConnectToDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	return rdb
}

func WriteToDB(client *redis.Client, ctx context.Context, temp float64, humidity float64) {
	client.XAdd(ctx, &redis.XAddArgs{
		Stream: "stream",
		Values: map[string]interface{}{"temp": temp, "humidity": humidity},
		ID:     "*",
	})
}
