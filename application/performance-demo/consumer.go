package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func StartConsumer(brokerList []string, kafkaTopic string) chan bool {
	shutdown := make(chan bool)

	consumerGroupHandler := Consumer{}
	var handler sarama.ConsumerGroupHandler
	if IsInstrumented() {
		handler = otelsarama.WrapConsumerGroupHandler(&consumerGroupHandler)
	}

	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	var consumerGroup sarama.ConsumerGroup
	retries := 5
	retry := true
	for retry {
		con, err := sarama.NewConsumerGroup(brokerList, "consumer", config)
		if err != nil {
			if retries != 0 {
				retries--
				log.Println("RETRY: Failed to start consumer group:", err)
				time.Sleep(2 * time.Second)
				continue
			}
			log.Fatalln("Failed to start consumer group:", err)
		}
		log.Println("Started consumer group")
		consumerGroup = con
		retry = false
	}

	retries = 5
	retry = true
	for retry {
		var err error
		if IsInstrumented() {
			err = consumerGroup.Consume(context.Background(), []string{kafkaTopic}, handler)
		} else {
			err = consumerGroup.Consume(context.Background(), []string{kafkaTopic}, &consumerGroupHandler)
		}

		if err != nil {
			if retries != 0 {
				retries--
				log.Println("RETRY: Failed to consume via handler:", err)
				time.Sleep(2 * time.Second)
				continue
			}
			log.Fatalln("Failed to consume via handler:", err)
		}
		log.Println("Started consumer via handler")
		retry = false
	}

	go func() {
		select {
		case _ = <-sigchan:
			shutdown <- true
		}
	}()

	return shutdown
}

func printMessage(msg *sarama.ConsumerMessage) {
	go func() {
		var span trace.Span
		if IsInstrumented() {
			// Extract tracing info from message
			ctx := otel.GetTextMapPropagator().Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(msg))

			tr := otel.Tracer("consumer")
			_, span = tr.Start(ctx, "consume message", trace.WithAttributes(
				semconv.MessagingOperationProcess,
			))
			defer span.End()
			fmt.Println("traceId: ", span.SpanContext().TraceID())
		}

		if IsWithWorkTime() {
			// Emulate Work loads
			time.Sleep(1 * time.Second)
		}

		if IsWithRandomError() {
			randNum := rand.Intn(100)
			if randNum < 7 {
				time.Sleep(3 * time.Second)
				if IsInstrumented() {
					span.SetStatus(codes.Error, "This is some error")
				}

				log.Println("Failed to close producer:", "This is some error")
			}
		}

		traceId := ""
		for _, d := range msg.Headers {
			if bytes.Compare(d.Key, []byte("id")) == 0 {
				traceId = string(d.Value)
			}
		}

		t := time.Now().UnixNano()
		fmt.Println(t, "consumer.consumed", "app=", appName, "id=", traceId)
	}()
}

type Consumer struct {
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		printMessage(message)
		session.MarkMessage(message, "")
	}

	return nil
}
