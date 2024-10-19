package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	brokerAddr = "localhost:9092" // Kafka broker address
	topic      = "temperature_data"
	dt         = 0.04
	densityAir = 1000.0
	Kp1        = 1000.0
	volO1i     = 30.0
)

var (
	TemperatureRoom1 = volO1i
)

func CreateCSVFile(filename string) (*os.File, *csv.Writer, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}

	writer := csv.NewWriter(file)
	writer.Write([]string{
		"Time", "Reference Temperature", "Actual Temperature", "Error", "Control Input", "People In Room", "Humidity", "Annotation",
	})
	writer.Flush()

	return file, writer, nil
}

func handleMessage(payload string, writer *csv.Writer) {
	fields := strings.Split(payload, ",")
	if len(fields) < 4 {
		fmt.Println("Invalid payload:", payload)
		return
	}

	time, _ := strconv.ParseFloat(fields[0], 64)
	volR1, _ := strconv.ParseFloat(fields[1], 64)
	humidity, _ := strconv.ParseFloat(fields[2], 64)
	peopleInRoom, _ := strconv.Atoi(fields[3])

	var error1, mDot1 float64
	var annotation string

	if peopleInRoom > 0 {
		error1 = volR1 - TemperatureRoom1
		mDot1 = Kp1 * error1
		TemperatureRoom1 += (mDot1 / densityAir) * dt
	} else {
		annotation = "No people in room"
	}

	writer.Write([]string{
		fmt.Sprintf("%.2f", time),
		fmt.Sprintf("%.2f", volR1),
		fmt.Sprintf("%.2f", TemperatureRoom1),
		fmt.Sprintf("%.2f", error1),
		fmt.Sprintf("%.2f", mDot1),
		fmt.Sprintf("%d", peopleInRoom),
		fmt.Sprintf("%.2f", humidity),
		annotation,
	})
	writer.Flush()
}

func main() {
	// Set up the CSV file
	file, writer, err := CreateCSVFile("temperature_data.csv")
	if err != nil {
		log.Fatalf("Could not open CSV file: %v", err)
	}
	defer file.Close()
	defer writer.Flush()

	// Set up Kafka consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokerAddr,
		"group.id":          "group1",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{topic}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			handleMessage(string(msg.Value), writer)
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
