package main

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	brokerAddr = "localhost:9092" // Kafka broker address
	tEnd       = 100.0            // for how long to run
	dt         float64            // measurement frequency (to be loaded from config)
	volR1i     = 70.0
	topic      = "temperature_data"
)

type DataPacket struct {
	Time         float64
	VolR1        float64
	Humidity     float64
	PeopleInRoom int
}

var rootCmd = &cobra.Command{
	Use:   "producer",
	Short: "Kafka Producer for temperature data",
	Run:   runProducer,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().Float64Var(&dt, "measurement-frequency", 0.04, "Number of times to measure per second")
	viper.BindPFlag("measurement-frequency", rootCmd.PersistentFlags().Lookup("measurement-frequency"))
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file found; using defaults.")
	}
	dt = viper.GetFloat64("measurement-frequency") // Load measurement frequency from config
}

func runProducer(cmd *cobra.Command, args []string) {
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

		// Simulate data packets
		if i < 300 {
			packet.VolR1 = volR1i
			packet.PeopleInRoom = 3
			packet.Humidity = 30.0
		} else if i < 600 {
			packet.VolR1 = 20
			packet.PeopleInRoom = 2
			packet.Humidity = 40.0
		} else if i < 900 {
			packet.VolR1 = 90
			packet.PeopleInRoom = 0
			packet.Humidity = 10.0
		} else if i < 1200 {
			packet.VolR1 = 30
			packet.PeopleInRoom = 4
			packet.Humidity = 60.0
		} else if i < 1500 {
			packet.VolR1 = 80
			packet.PeopleInRoom = 1
			packet.Humidity = 80.0
		} else if i < 1800 {
			packet.VolR1 = 10
			packet.PeopleInRoom = 5
			packet.Humidity = 90.0
		} else if i < 2100 {
			packet.VolR1 = 95
			packet.PeopleInRoom = 0
			packet.Humidity = 20.0
		} else {
			packet.VolR1 = 50
			packet.PeopleInRoom = 2
			packet.Humidity = 50.0
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
