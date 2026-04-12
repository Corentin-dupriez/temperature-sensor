package redisdb

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	redis "github.com/redis/go-redis/v9"
)

type TempReading struct {
	humidity    float64
	temperature float64
}

func ConnectToRedis() *redis.Client {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	redisConnString := os.Getenv("REDIS_CONN_STR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConnString,
		Password: redisPassword,
		DB:       0,
		Protocol: 2,
	})
	slog.Info("Connected to Redis DB")
	return rdb
}

func ReadFromRedis(ctx context.Context, rdb *redis.Client) []redis.XStream {
	res, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Streams: []string{"stream", ">"},
		Group:   "histo-db-consumer",
		Count:   1,
	}).Result()
	if err != nil {
		panic(err)
	}
	return res
}

func ParseFromStreamResult(res []redis.XStream) {
	for _, stream := range res {
		for _, msg := range stream.Messages {
			reading := msg.Values
			humidity, err := strconv.ParseFloat(reading["humidity"].(string), 64)
			if err != nil {
				panic(err)
			}
			temp, err := strconv.ParseFloat(reading["temp"].(string), 64)
			if err != nil {
				panic(err)
			}

			tempreading := TempReading{
				humidity:    humidity,
				temperature: temp,
			}
			fmt.Println(tempreading.humidity)
			fmt.Println(tempreading.temperature)
		}
	}
}
