package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func startConsumer(bootstrapServer string, groupId string, topic *string, consumerShutdown chan bool) {
	fmt.Println("Starting consumer - topic: " + *topic)

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

	run := true

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	// err = c.SubscribeTopics([]string{*topic}, nil)
	// if err != nil {
	// 	fmt.Printf("Failed to subsribe to topic: %s  - Error: %s", *topic, err)
	// 	os.Exit(1)
	// }

	// totalCount := 0
	// run := true
	// fmt.Println("Consumer running")
	// for run == true {
	// 	select {
	// 	case sig := <-sigchan:
	// 		fmt.Printf("Caught signal %v: terminating consumer\n", sig)
	// 		consumerShutdown <- true
	// 		run = false
	// 	default:
	// 		msg, err := c.ReadMessage(100 * time.Millisecond)
	// 		if err != nil {
	// 			// Errors are informational and automatically handled by the consumer
	// 			continue
	// 		}
	// 		recordKey := string(msg.Key)
	// 		message := msg.Value
	// 		data := KafkaMessage{}
	// 		err = json.Unmarshal(message, &data)
	// 		if err != nil {
	// 			fmt.Printf("Failed to decode JSON at offset %d: %v", msg.TopicPartition.Offset, err)
	// 			continue
	// 		}
	// 		count := data.Count
	// 		totalCount += count
	// 		fmt.Printf("Consumed record with key %s and value %d, and updated total count to %d -- Message was: %s\n", recordKey, data.Count, totalCount, data.Message)
	// 	}
	// }

	c.Close()
}
