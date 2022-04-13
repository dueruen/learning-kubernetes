package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func startConsumer(bootstrapServer string, groupId string, topic *string, consumerShutdown chan bool) {
	fmt.Println("Startin consumer")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	err = c.SubscribeTopics([]string{*topic}, nil)

	totalCount := 0
	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating consumer\n", sig)
			consumerShutdown <- true
			run = false
		default:
			msg, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			recordKey := string(msg.Key)
			message := msg.Value
			data := KafkaMessage{}
			err = json.Unmarshal(message, &data)
			if err != nil {
				fmt.Printf("Failed to decode JSON at offset %d: %v", msg.TopicPartition.Offset, err)
				continue
			}
			count := data.Count
			totalCount += count
			fmt.Printf("Consumed record with key %s and value %d, and updated total count to %d -- Message was: %s\n", recordKey, data.Count, totalCount, data.Message)
		}
	}

	c.Close()
}
