package main

import (
	"context"
	"log/slog"
	"os"

	redisdb "github.com/Corentin-dupriez/temperature-sensor/internal/redis_db"
	sensorworker "github.com/Corentin-dupriez/temperature-sensor/internal/sensor_worker"
)

func main() {
	port, err := sensorworker.OpenPort()
	if err != nil {
		slog.Error("Error reading the port: ", "error", err)
		os.Exit(1)
	}
	ctx := context.Background()
	client := redisdb.ConnectToDB()
	sensorworker.ReadPort(port, ctx, client)

	defer client.Close()
}
