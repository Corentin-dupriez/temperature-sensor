package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go.bug.st/serial"
)

func main() {
	port, err := OpenPort()
	if err != nil {
		panic(err)
	}
	fmt.Println(port)
	ReadPort(port)
	client := ConnectToDB()
	// ctx := context.Background()

	// client.XAdd(ctx, &redis.XAddArgs{
	// 	Stream: "stream",
	// 	Values: map[string]interface{}{"temp": "20", "humidity": "40"},
	// 	ID:     "*",
	// })
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

func ReadPort(p serial.Port) {
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
		fmt.Println(temperature, humidity)
	}
}

func ProcessBuffer(b []byte) (float64, float64) {
	str := string(b)
	string := strings.Split(str, " ")

	tempString := strings.Split(string[0], ":")[1]
	humidityString := strings.Split(string[1], ":")[1]

	tempFloat, err := strconv.ParseFloat(tempString, 64)
	if err != nil {
		panic(err)
	}

	humidityFloat, err := strconv.ParseFloat(humidityString, 64)
	if err != nil {
		panic(err)
	}

	return tempFloat, humidityFloat
}
