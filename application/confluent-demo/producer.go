package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func startProducer(bootstrapServer string, topic *string, producerShutdown chan bool) {
	fmt.Println("Starting producer - topic: " + *topic)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServer})
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)

	deliveryChan := make(chan kafka.Event)

	value := "Hello Go!"
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
		Value:          []byte(value),
		Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, deliveryChan)

	fmt.Printf("Waiting for delivery of message on topic %s\n", *topic)
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)

	// go func() {
	// 	fmt.Println("Loop events")
	// 	for e := range p.Events() {
	// 		switch ev := e.(type) {
	// 		case *kafka.Message:
	// 			if ev.TopicPartition.Error != nil {
	// 				fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
	// 			} else {
	// 				fmt.Printf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
	// 					*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
	// 			}
	// 		}
	// 	}
	// 	fmt.Println("Loop events OVER")
	// }()

	// count := 0
	// run := true
	// fmt.Println("Producer running")
	// for run == true {
	// 	select {
	// 	case sig := <-sigchan:
	// 		p.Flush(10)
	// 		fmt.Printf("Caught signal %v: terminating producer\n", sig)
	// 		producerShutdown <- true
	// 		run = false
	// 	default:
	// 		number := rand.Intn(60)

	// 		recordKey := "alice"
	// 		data := &KafkaMessage{
	// 			Message: "Hello from a consumer random number: " + strconv.Itoa(number),
	// 			Count:   count,
	// 		}
	// 		recordValue, _ := json.Marshal(&data)
	// 		fmt.Printf("Preparing to produce record: %s\t%s\n", recordKey, recordValue)
	// 		p.Produce(&kafka.Message{
	// 			TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
	// 			Key:            []byte(recordKey),
	// 			Value:          []byte(recordValue),
	// 		}, nil)

	// 		p.Flush(10)
	// 		time.Sleep(time.Duration(number/10) * time.Second)
	// 		count++
	// 	}
	// }
}
