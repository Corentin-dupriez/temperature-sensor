package redisdb

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	redis "github.com/redis/go-redis/v9"
)

type TempReading struct {
	Humidity    float64
	Temperature float64
	TimeReading time.Time
}

func ConnectToRedis() *redis.Client {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("env file not found, gathering secrets from environment")
	}
	redisConnString := os.Getenv("REDIS_CONN_STR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	fmt.Printf("Connection string: %v, password: %v", redisConnString, redisPassword)
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

func ParseFromStreamResult(res []redis.XStream, ctx context.Context, rdb *redis.Client) TempReading {
	var tempreading TempReading
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

			timeReading, err := time.Parse(time.RFC3339Nano, reading["time"].(string))
			if err != nil {
				panic(err)
			}

			tempreading = TempReading{
				Humidity:    humidity,
				Temperature: temp,
				TimeReading: timeReading,
			}
			rdb.XAck(ctx, "stream", "histo-db-consumer", msg.ID)
		}
	}
	return tempreading
}
