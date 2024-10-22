package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	brokerAddr = "localhost:9092" // Kafka broker address
	topic      = "temperature_data"
	dt         = 0.04 // measurement frequency (to be loaded from config)
	densityAir = 1000.0

	// PID controller parameters
	Kp float64 // Proportional gain
	Ki float64 // Integral gain
	Kd float64 // Derivative gain

	// Initial conditions
	volO1i           = 30.0
	TemperatureRoom1 = volO1i

	// PID state variables
	integral  float64
	lastError float64
)

var rootCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Kafka Consumer for temperature data",
	Run:   runConsumer,
}

func init() {
	cobra.OnInitialize(initConfig)

	dt = 0.04 // Hardcoded measurement frequency

	// Define persistent flags
	rootCmd.PersistentFlags().Float64Var(&Kp, "Kp", 1000.0, "Proportional gain for PID controller")
	rootCmd.PersistentFlags().Float64Var(&Ki, "Ki", 0.1, "Integral gain for PID controller")
	rootCmd.PersistentFlags().Float64Var(&Kd, "Kd", 0.01, "Derivative gain for PID controller")

	// Bind flags to viper
	viper.BindPFlag("Kp", rootCmd.PersistentFlags().Lookup("Kp"))
	viper.BindPFlag("Ki", rootCmd.PersistentFlags().Lookup("Ki"))
	viper.BindPFlag("Kd", rootCmd.PersistentFlags().Lookup("Kd"))
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file found; using defaults.")
	}
	Kp = viper.GetFloat64("Kp") // Load proportional gain from config
	Ki = viper.GetFloat64("Ki") // Load integral gain from config
	Kd = viper.GetFloat64("Kd") // Load derivative gain from config
}
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
		error1 = volR1 - TemperatureRoom1               // Calculate the error
		integral += error1 * dt                         // Integral component
		derivative := (error1 - lastError) / dt         // Derivative component
		mDot1 = Kp*error1 + Ki*integral + Kd*derivative // PID control input
		TemperatureRoom1 += (mDot1 / densityAir) * dt   // Update temperature based on control input
		lastError = error1                              // Save the error for the next derivative calculation
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

func runConsumer(cmd *cobra.Command, args []string) {
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
