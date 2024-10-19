package main

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	brokerAddr = "localhost:9092" // Kafka broker address
	tEnd       = 100.0            // for how long to run
	dt         = 0.04             // how many times a second to measure
	volR1i     = 70.0
)

var topic = "temperature_data"

type DataPacket struct {
	Time         float64
	VolR1        float64
	Humidity     float64
	PeopleInRoom int
}

func main() {
	// Set up Kafka producer configuration
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokerAddr})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	steps := int(tEnd / dt)
	lastVolR1 := volR1i // Initialize with the initial temperature reference

	for i := 1; i < steps; i++ {
		packet := DataPacket{
			Time: float64(i) * dt,
		}

		if i < 300 {
			packet.VolR1 = volR1i
			packet.PeopleInRoom = 3 // Random number between 0 and 5
			packet.Humidity = 30.0  // Placeholder value
		} else if i < 600 {
			packet.VolR1 = 20
			packet.PeopleInRoom = 2 // Random number between 0 and 5
			packet.Humidity = 40.0  // Placeholder value
		} else if i < 900 {
			packet.VolR1 = 90
			packet.PeopleInRoom = 0 // Random number between 0 and 5
			packet.Humidity = 10.0  // Placeholder value
		} else if i < 1200 {
			packet.VolR1 = 30
			packet.PeopleInRoom = 4 // Random number between 0 and 5
			packet.Humidity = 60.0  // Placeholder value
		} else if i < 1500 {
			packet.VolR1 = 80
			packet.PeopleInRoom = 1 // Random number between 0 and 5
			packet.Humidity = 80.0  // Placeholder value
		} else if i < 1800 {
			packet.VolR1 = 10
			packet.PeopleInRoom = 5 // Random number between 0 and 5
			packet.Humidity = 90.0  // Placeholder value
		} else if i < 2100 {
			packet.VolR1 = 95
			packet.PeopleInRoom = 0 // Random number between 0 and 5
			packet.Humidity = 20.0  // Placeholder value
		} else {
			packet.VolR1 = 50
			packet.PeopleInRoom = 2 // Random number between 0 and 5
			packet.Humidity = 50.0  // Placeholder value
		}

		// Create a payload to send
		payload := fmt.Sprintf("%f,%f,%f,%d", packet.Time, packet.VolR1, packet.Humidity, packet.PeopleInRoom)

		// Send the payload to Kafka
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(payload),
		}, nil)
		if err != nil {
			fmt.Printf("Failed to send data to Kafka: %v\n", err)
		}

		// Wait for the reference temperature change
		if packet.VolR1 != lastVolR1 {
			fmt.Printf("Reference temperature changed to: %.2f after %.2f seconds\n", packet.VolR1, packet.Time)
			lastVolR1 = packet.VolR1
		}

		time.Sleep(time.Duration(dt*1000) * time.Millisecond)
	}
}
