// Package redisdb is used to connect and write to the Redis stream
package redisdb

import (
	"context"
	"log/slog"

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
		if redis.HasErrorPrefix(err, "BUSYGROUP") {
			slog.Info("Group already exists, skipping", "group", "server-consumer")
		} else {
			panic(err)
		}
	}
	_, err = client.XGroupCreate(ctx, "stream", "histo-db-consumer", "$").Result()
	if err != nil {
		if redis.HasErrorPrefix(err, "BUSYGROUP") {
			slog.Info("Group already exists, skipping", "group", "server-consumer")
		} else {
			panic(err)
		}
	}
}

func WriteToDB(client *redis.Client, ctx context.Context, temp float64, humidity float64) {
	client.XAdd(ctx, &redis.XAddArgs{
		Stream: "stream",
		Values: map[string]interface{}{"temp": temp, "humidity": humidity},
		ID:     "*",
	})
}
