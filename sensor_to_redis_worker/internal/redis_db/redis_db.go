// Package redisdb is used to connect and write to the Redis stream
package redisdb

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func ConnectToDB(connStr string, password string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     connStr,
		Password: password,
		DB:       0,
		Protocol: 2,
	})
	return rdb
}

func CreateConsumerGroups(ctx context.Context, client *redis.Client) {
	_, err := client.XGroupCreate(ctx, "stream", "server-consumer", "$").Result()
	if err != nil {
		panic(err)
	}
	_, err = client.XGroupCreate(ctx, "stream", "histo-db-consumer", "$").Result()
	if err != nil {
		panic(err)
	}
}

func WriteToDB(client *redis.Client, ctx context.Context, temp float64, humidity float64) {
	client.XAdd(ctx, &redis.XAddArgs{
		Stream: "stream",
		Values: map[string]interface{}{"temp": temp, "humidity": humidity},
		ID:     "*",
	})
}
