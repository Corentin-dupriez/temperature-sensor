package main

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"go.bug.st/serial"
)

func main() {
	port, err := OpenPort()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	client := ConnectToDB()
	ReadPort(port, ctx, client)

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
		if port == "/dev/tty.usbserial-130" {
			port, err := serial.Open(port, mode)
			if err != nil {
				panic(err)
			}
			return port, nil
		}
	}
	return nil, errors.New("arduino is not connected")
}

func WriteToDB(client *redis.Client, ctx context.Context, temp float64, humidity float64) {
	client.XAdd(ctx, &redis.XAddArgs{
		Stream: "stream",
		Values: map[string]interface{}{"temp": temp, "humidity": humidity},
		ID:     "*",
	})
}

func ReadPort(p serial.Port, ctx context.Context, client *redis.Client) {
	buf := make([]byte, 100)
	for {
		n, err := p.Read(buf)
		if err != nil {
			panic(err)
		}
		if n == 0 {
			break
		}
		temperature, humidity := ProcessBuffer(buf)
		WriteToDB(client, ctx, temperature, humidity)
	}
}

func SplitString(s string, sep string, index int) (string, error) {
	splitString := strings.Split(s, sep)
	if len(splitString) < index {
		return "", errors.New("unable to split the string correctly")
	}
	return splitString[index], nil
}

func ProcessBuffer(b []byte) (float64, float64) {
	str := string(b)
	string := strings.Split(str, " ")

	tempString, err := SplitString(string[0], ":", 1)
	if err != nil {
		panic(err)
	}

	humidityString, err := SplitString(string[1], ":", 1)
	if err != nil {
		panic(err)
	}
	humidityString, err = SplitString(humidityString, "%", 0)
	if err != nil {
		panic(err)
	}

	tempFloat, err := strconv.ParseFloat(tempString, 64)
	if err != nil {
		slog.Warn("Unable to parse the temperature")
	}

	humidityFloat, err := strconv.ParseFloat(humidityString, 64)
	if err != nil {
		slog.Warn("Unable to parse the humidity")
	}

	return tempFloat, humidityFloat
}
