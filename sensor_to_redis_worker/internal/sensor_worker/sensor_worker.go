// Package sensorworker contains functions used to read and parse the serial port of the computer to get details from the Arduino board
package sensorworker

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"strings"

	redisdb "github.com/Corentin-dupriez/temperature-sensor/internal/redis_db"
	"github.com/redis/go-redis/v9"
	"go.bug.st/serial"
)

func OpenPort(portName string) (serial.Port, error) {
	// dev/cu.usbserial-130
	mode := &serial.Mode{
		BaudRate: 9600,
	}
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, errors.New("not able to read ports")
	}
	for _, port := range ports {
		if port == portName {
			port, err := serial.Open(port, mode)
			if err != nil {
				panic(err)
			}
			return port, nil
		}
	}
	return nil, errors.New("arduino is not connected")
}

func ReadPort(p serial.Port, ctx context.Context, client *redis.Client) {
	buf := make([]byte, 100)
	for {
		n, err := p.Read(buf)
		if err != nil {
			slog.Error("impossible to read the port", "error", err)
		}
		if n == 0 {
			break
		}
		temperature, humidity := processBuffer(buf)
		redisdb.WriteToDB(client, ctx, temperature, humidity)
	}
}

func splitString(s string, sep string, index int) (string, error) {
	splitString := strings.Split(s, sep)
	if len(splitString) < index {
		return "", errors.New("unable to split the string correctly")
	}
	return splitString[index], nil
}

func processBuffer(b []byte) (float64, float64) {
	str := string(b)
	string := strings.Split(str, " ")

	tempString, err := splitString(string[0], ":", 1)
	if err != nil {
		panic(err)
	}

	humidityString, err := splitString(string[1], ":", 1)
	if err != nil {
		panic(err)
	}
	humidityString, err = splitString(humidityString, "%", 0)
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
