package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	redisdb "github.com/Corentin-dupriez/temperature-sensor/internal/redis_db"
	sensorworker "github.com/Corentin-dupriez/temperature-sensor/internal/sensor_worker"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		slog.Error("error reading the .env file")
		fmt.Println(os.Getenv("ARDUINO_PORT"))
	}
	port, err := sensorworker.OpenPort(os.Getenv("ARDUINO_PORT"))
	if err != nil {
		slog.Error("Error reading the port: ", "error", err)
		os.Exit(1)
	}
	ctx := context.Background()
	client := redisdb.ConnectToDB(os.Getenv("REDIS_CONN_STR"), os.Getenv("REDIS_PASSWORD"))
	slog.Info("Connected to Redis DB")
	redisdb.CreateConsumerGroups(ctx, client)
	sensorworker.ReadPort(port, ctx, client)

	defer client.Close()
}
