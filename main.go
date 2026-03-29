package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.bug.st/serial"
)

func main() {
	port, err := OpenPort()
	if err != nil {
		panic(err)
	}
	fmt.Println(port)
	client := ConnectToDB()
	ctx := context.Background()

	client.XAdd(ctx, &redis.XAddArgs{
		Stream: "stream",
		Values: map[string]interface{}{"temp": "20", "humidity": "40"},
		ID:     "*",
	})
	defer client.Close()
}

func OpenPort() (serial.Port, error) {
	// dev/cu.usbserial-130
	mode := &serial.Mode{
		BaudRate: 9600,
	}
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, errors.New("not able to read ports")
	}
	for _, port := range ports {
		if port == "dev/cu.usbserial-130" {
			port, err := serial.Open(port, mode)
			if err != nil {
				panic(err)
			}
			return port, nil
		}
	}
	return nil, errors.New("arduino is not connected")
}
