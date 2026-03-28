package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := ConnectToDB()
	ctx := context.Background()

	client.XAdd(ctx, &redis.XAddArgs{
		Stream: "stream",
		Values: map[string]interface{}{"temp": "20", "humidity": "40"},
		ID:     "*",
	})
	val, err := client.XRange(ctx, "stream", "-", "+").Result()
	if err != nil {
		panic(err)
	}

	if len(val) > 0 {
		for i := 0; i < len(val); i++ {
			first := val[i].Values
			fmt.Println(first)
		}
	}

	defer client.Close()
}
