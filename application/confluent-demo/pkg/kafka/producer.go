package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func (srv *KafkaServer) startProducer(bootstrapServer string, topic *string, producerShutdown chan bool) {
	srv.logger.Sugar().Infof("Starting producer - server: %s - topic: %s\n", bootstrapServer, *topic)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		// "fetch.max.bytes":           2147483135,
		// "receive.message.max.bytes": 2147483647,
	})
	// "fetch.message.max.bytes":   18000000,
	// "receive.message.max.bytes": 2147483647,
	// "security.protocol":         "plaintext",
	// "sasl.mechanism":            "PLAIN",
	// "sasl.username":             "*",
	// "sasl.password":             "*",

	if err != nil {
		srv.logger.Sugar().Errorf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	srv.logger.Sugar().Infof("Created Producer %v\n", p)

	count := 0
	run := true
	srv.logger.Sugar().Infof("Producer running")
	for run == true {
		select {
		case sig := <-sigchan:
			p.Flush(10)
			srv.logger.Sugar().Infof("Caught signal %v: terminating producer\n", sig)
			producerShutdown <- true
			run = false
		default:
			number := rand.Intn(60)
			time.Sleep(time.Duration(number/10) * time.Second)

			go func(count int) {
				ctx, spanRoot := srv.tracer.Start(context.Background(), "kafka.test")
				//ctx, span := srv.tracer.Start(context.Background(), "produce message", trace.WithSpanKind(trace.SpanKindProducer))
				//defer span.End()

				recordKey := "alice"
				data := &KafkaMessage{
					Message: strconv.Itoa(number),
					Count:   count,
				}
				recordValue, _ := json.Marshal(&data)
				srv.logger.Sugar().Debugf("Preparing to produce record: %s\t%s\n", recordKey, recordValue)

				msg := &kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
					Key:            []byte(recordKey),
					Value:          []byte(recordValue),
				}

				carrier := NewProducerMessageCarrier(msg)
				//ctx := otel.GetTextMapPropagator().Extract(ctxRoot, carrier)
				otel.GetTextMapPropagator().Inject(ctx, carrier)

				attrs := []attribute.KeyValue{
					semconv.MessagingSystemKey.String("kafka"),
					semconv.MessagingDestinationKindTopic,
					semconv.MessagingDestinationKey.String(*msg.TopicPartition.Topic),
				}
				opts := []trace.SpanStartOption{
					trace.WithAttributes(attrs...),
					trace.WithSpanKind(trace.SpanKindProducer),
				}
				ctx, span := srv.tracer.Start(ctx, "kafka.produce", opts...)
				otel.GetTextMapPropagator().Inject(ctx, carrier)

				deliveryChan := make(chan kafka.Event)

				p.Produce(msg, deliveryChan)

				e := <-deliveryChan
				m := e.(*kafka.Message)

				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}

				close(deliveryChan)

				span.SetAttributes(
					semconv.MessagingMessageIDKey.String(strconv.FormatInt(int64(m.TopicPartition.Offset), 10)),
					attribute.Key("messaging.kafka.partition").Int64(int64(m.TopicPartition.Partition)),
				)

				span.End()
				spanRoot.End()
			}(count)
			count++
		}
	}
	srv.logger.Sugar().Infof("Producer down")
}
