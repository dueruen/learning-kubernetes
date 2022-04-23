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
	getter "github.com/dueruen/learning-kubernetes/application/confluent-demo/pkg/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func (srv *KafkaServer) startConsumer(bootstrapServer string, groupId string, topic *string, consumerShutdown chan bool, topicProduce *string) {
	srv.logger.Sugar().Infof("Starting consumer - server: %s - groupId: %s - topic: %s\n", bootstrapServer, groupId, *topic)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               bootstrapServer,
		"group.id":                        groupId,
		"auto.offset.reset":               "earliest",
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"enable.partition.eof":            true,
		"session.timeout.ms":              6000,
		// "fetch.max.bytes":                 2147483135,
		// "receive.message.max.bytes":       2147483647,
	})

	if err != nil {
		srv.logger.Sugar().Errorf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	err = c.SubscribeTopics([]string{*topic}, nil)
	if err != nil {
		srv.logger.Sugar().Errorf("Failed to subsribe to topic: %s  - Error: %s", *topic, err)
		os.Exit(1)
	}

	// totalCount := 0
	run := true
	srv.logger.Sugar().Infof("Consumer running")
	for run == true {
		select {
		case sig := <-sigchan:
			srv.logger.Sugar().Infof("Caught signal %v: terminating consumer\n", sig)
			consumerShutdown <- true
			run = false
		case ev := <-c.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Unassign()
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				srv.handleMessage(e, topic, bootstrapServer, topicProduce)
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			}
		}
	}

	srv.logger.Sugar().Infof("END of consumer")
	c.Close()
}

func (srv *KafkaServer) handleMessage(msg *kafka.Message, topic *string, bootstrapServer string, topicProduce *string) {
	carrier := NewConsumerMessageCarrier(msg)
	parentSpanContext := otel.GetTextMapPropagator().Extract(context.Background(), carrier)

	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKey.String("kafka"),
		semconv.MessagingDestinationKindTopic,
		semconv.MessagingDestinationKey.String(*msg.TopicPartition.Topic),
		semconv.MessagingOperationReceive,
		semconv.MessagingMessageIDKey.String(strconv.FormatInt(int64(msg.TopicPartition.Offset), 10)),
		attribute.Key("messaging.kafka.partition").Int64(int64(msg.TopicPartition.Partition)),
	}
	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindConsumer),
	}
	ctx, span := srv.tracer.Start(parentSpanContext, "kafka.consume", opts...)

	otel.GetTextMapPropagator().Inject(ctx, carrier)

	// ctx := otel.GetTextMapPropagator().Extract(context.Background(), NewConsumerMessageCarrier(msg))
	// ctx, span := srv.tracer.Start(ctx, "consume message", trace.WithAttributes(
	// 	semconv.MessagingOperationProcess,
	// 	semconv.MessagingSystemKey.String("kafka"),
	// 	semconv.MessagingDestinationKindTopic,
	// 	semconv.MessagingDestinationKey.String(*msg.TopicPartition.Topic),
	// 	semconv.MessagingMessageIDKey.String(strconv.FormatInt(int64(msg.TopicPartition.Offset), 10)),
	// 	attribute.Key("messaging.kafka.partition").Int64(int64(msg.TopicPartition.Partition)),
	// ), trace.WithSpanKind(trace.SpanKindConsumer))

	if msg.Headers != nil {
		srv.logger.Sugar().Warnf("%% Headers: %v\n", msg.Headers)
	} else {
		srv.logger.Sugar().Warnf("NO HEADERS\n")
	}

	if topicProduce != nil && *topicProduce != "" {
		go func() {
			srv.testPro(ctx, bootstrapServer, topicProduce)
		}()
	}

	go func() {
		getter.GetFromBackend(srv.logger, ctx, srv.tracer)
		span.End()
	}()
	span.End()
}

func (srv *KafkaServer) testPro(ctx context.Context, bootstrapServer string, topic *string) {
	ctxRoot, spanRoot := srv.tracer.Start(ctx, "kafka.produce.root")
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
	})

	if err != nil {
		srv.logger.Sugar().Errorf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	select {
	default:
		number := rand.Intn(60)
		time.Sleep(time.Duration(number/10) * time.Second)

		//ctx, span := srv.tracer.Start(context.Background(), "produce message", trace.WithSpanKind(trace.SpanKindProducer))
		//defer span.End()

		recordKey := "bo"
		data := &KafkaMessage{
			Message: strconv.Itoa(number),
			Count:   42,
		}
		recordValue, _ := json.Marshal(&data)
		srv.logger.Sugar().Debugf("Preparing to produce record: %s\t%s\n", recordKey, recordValue)

		msg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
			Key:            []byte(recordKey),
			Value:          []byte(recordValue),
		}

		carrier := NewProducerMessageCarrier(msg)
		//	ctx := otel.GetTextMapPropagator().Extract(context.Background(), carrier)
		otel.GetTextMapPropagator().Inject(ctxRoot, NewProducerMessageCarrier(msg))

		attrs := []attribute.KeyValue{
			semconv.MessagingSystemKey.String("kafka"),
			semconv.MessagingDestinationKindTopic,
			semconv.MessagingDestinationKey.String(*msg.TopicPartition.Topic),
		}
		opts := []trace.SpanStartOption{
			trace.WithAttributes(attrs...),
			trace.WithSpanKind(trace.SpanKindProducer),
		}
		ctx, span := srv.tracer.Start(ctxRoot, "kafka.produce", opts...)
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
	}
}
