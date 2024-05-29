package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	// Kafka producer configuration :9093
	server := os.Getenv("KAFKA_SERVER")
	fmt.Println("KAFKA_SERVER:", server)

	// Create a new Kafka producer configuration
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true // Enable delivery report for successful messages

	// Create a new Kafka producer instance
	producer, err := sarama.NewSyncProducer([]string{server}, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %s", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Failed to close Kafka producer: %s", err)
		}
	}()

	// Define the topic and message payload
	topic := "my-topic"
	message := "Hello, Kafka!"

	// Create a ticker to send messages at a specified interval
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send the message to Kafka
			msgBytes := []byte(message)
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.ByteEncoder(msgBytes),
			}

			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				log.Printf("Failed to produce message: %v", err)
			} else {
				fmt.Printf("Message delivered to partition %d at offset %d\n", partition, offset)
			}
		}
	}
}
