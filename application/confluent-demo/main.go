package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	bootstrapServer = flag.String("bootstrap-server", os.Getenv("BOOTSTRAP_SERVER"), "BOOTSTRAP_SERVER")
	producer        = flag.String("producer", os.Getenv("PRODUCER"), "PRODUCER")
	producerTopic   = flag.String("producer-topic", os.Getenv("PRODUCER_TOPIC"), "PRODUCER_TOPIC")
	consumer        = flag.String("consumer", os.Getenv("CONSUMER"), "CONSUMER")
	consumerTopic   = flag.String("consumer-topic", os.Getenv("CONSUMER_TOPIC"), "CONSUMER_TOPIC")
	consumerGrupId  = flag.String("consumer-group-id", os.Getenv("CONSUMER_GROUP_ID"), "CONSUMER_GROUP_ID")
)

func main() {
	flag.Parse()

	if consumerGrupId == nil {
		*consumerGrupId = "go_example_group_1"
	}

	if consumerTopic == nil {
		*consumerTopic = "default"
	}

	fmt.Println("Starting up")

	waitForProducer := false
	producerShutdown := make(chan bool)
	if producer != nil && *producer != "" {
		waitForProducer = true
		go startProducer(*bootstrapServer, consumerTopic, producerShutdown)
	}

	waitForConsumer := false
	consumerShutdown := make(chan bool)
	if consumer != nil && *consumer != "" {
		waitForConsumer = true
		go startConsumer(*bootstrapServer, *consumerGrupId, consumerTopic, consumerShutdown)
	}

	fmt.Println("Running")
	running := true
	for running {
		select {
		case _ = <-producerShutdown:
			waitForProducer = false
		case _ = <-consumerShutdown:
			waitForConsumer = false
		default:
			if !waitForProducer && !waitForConsumer {
				running = false
			}
		}
	}
	fmt.Println("Shutting down")
}

type KafkaMessage struct {
	Message string
	Count   int
}
