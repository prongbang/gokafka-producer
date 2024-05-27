package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	// Create a new Kafka producer instance
	server := os.Getenv("KAFKA_SERVER")
	fmt.Println("KAFKA_SERVER:", server)
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": server,
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %s", err)
	}
	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Topic)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition.Topic)
				}
			}
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
			err = p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte(message),
			}, nil)
			if err != nil {
				log.Printf("Failed to produce message: %v", err)
			}
		}
	}
}
