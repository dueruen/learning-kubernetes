package kafka

import (
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel"
)

func (srv *KafkaServer) startProducer(bootstrapServer string, topic *string, producerShutdown chan bool) {
	srv.logger.Sugar().Infof("Starting producer - server: %s - topic: %s\n", bootstrapServer, *topic)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServer})
	if err != nil {
		srv.logger.Sugar().Errorf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	srv.logger.Sugar().Infof("Created Producer %v\n", p)

	go func() {
		srv.logger.Sugar().Debugf("Loop events")
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					srv.logger.Sugar().Errorf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					srv.logger.Sugar().Debugf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}

				// mtx.Lock()
				// if ev.TopicPartition.Metadata == nil {
				// 	srv.logger.Warn("METADATA IS NIL")
				// } else {
				// 	srv.logger.Sugar().Debugf("CAN GET ID", zap.String("data", *ev.TopicPartition.Metadata))
				// }

				// for _, d := range producerMessageContexts {
				// 	srv.logger.Sugar().Debugf("header", zap.String("some id: ", d.SpanContext().SpanID().String()))
				// }
				// mtx.Unlock()

				// carrier := NewProducerMessageCarrier(ev)
				// srv.logger.Sugar().Debugf("After success or failed %d", len(carrier.msg.Headers), zap.Int("header length", len(carrier.msg.Headers)))
				// for _, f := range carrier.msg.Headers {
				// 	srv.logger.Sugar().Debugf("header", zap.String(f.Key, string(f.Value)))
				// 	if f.Key == "traceparent" {
				// 		srv.logger.Sugar().Infof("HEADER HAS traceparent", zap.String("traceparent", string(f.Value)))
				// 		// mtx.Lock()
				// 		// if span, ok := producerMessageContexts[string(f.Value)]; ok {
				// 		// 	srv.logger.Sugar().Infof("FOUND traceparent", zap.String("traceparent", string(f.Value)))
				// 		// 	delete(producerMessageContexts, string(f.Value))
				// 		// 	//finishProducerSpan(mc.span, msg.Partition, msg.Offset, nil)
				// 		// 	//msg.Metadata = mc.metadataBackup // Restore message metadata
				// 		// 	span.SetAttributes(
				// 		// 		semconv.MessagingSystemKey.String("kafka"),
				// 		// 		semconv.MessagingDestinationKindTopic,
				// 		// 		semconv.MessagingDestinationKey.String(*ev.TopicPartition.Topic),
				// 		// 		semconv.MessagingMessageIDKey.String(strconv.FormatInt(int64(ev.TopicPartition.Offset), 10)),
				// 		// 		attribute.Key("messaging.kafka.partition").Int64(int64(ev.TopicPartition.Partition)),
				// 		// 	)

				// 		// 	span.End()
				// 		// }
				// 		// mtx.Unlock()
				// 	}
				// }

				// carrier := NewProducerMessageCarrier(ev)
				// ctx := otel.GetTextMapPropagator().Extract(context.Background(), carrier)

				// fmt.Println(ctx)

				// c := ctx.Value("b3KeyType")
				// if c != nil {
				// 	fmt.Println(c)
				// }

				// c = ctx.Value("b3.b3KeyType")
				// if c != nil {
				// 	fmt.Println(c)
				// }

				// attrs := []attribute.KeyValue{
				// 	semconv.MessagingSystemKey.String("kafka"),
				// 	semconv.MessagingDestinationKindTopic,
				// 	semconv.MessagingDestinationKey.String(*ev.TopicPartition.Topic),
				// }

				// opts := []trace.SpanStartOption{
				// 	trace.WithAttributes(attrs...),
				// 	trace.WithSpanKind(trace.SpanKindProducer),
				// }
				// ctx, span := srv.tracer.Start(ctx, "kafka.produce", opts...)

				// span.SetAttributes(
				// 	semconv.MessagingMessageIDKey.String(strconv.FormatInt(int64(ev.TopicPartition.Offset), 10)),
				// 	attribute.Key("messaging.kafka.partition").Int64(int64(ev.TopicPartition.Partition)),
				// )

				// span.End()
			}
		}
		srv.logger.Sugar().Infof("Loop events OVER")
	}()

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
			ctx, span := srv.tracer.Start(context.Background(), "produce message")
			//defer span.End()

			recordKey := "alice"
			data := &KafkaMessage{
				Message: "Hello from a consumer random number: " + strconv.Itoa(number),
				Count:   count,
			}
			recordValue, _ := json.Marshal(&data)
			srv.logger.Sugar().Debugf("Preparing to produce record: %s\t%s\n", recordKey, recordValue)

			msg := &kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
				Key:            []byte(recordKey),
				Value:          []byte(recordValue),
				Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
			}

			otel.GetTextMapPropagator().Inject(ctx, NewProducerMessageCarrier(msg))

			//fmt.Printf("%s", *msg)

			// srv.logger.Sugar().Infof("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			// for _, f := range msg.Headers {
			// 	srv.logger.Sugar().Infof("key; %s - value: %s", f.Key, string(f.Value))
			// }

			// srv.logger.Sugar().Infof("spanID", zap.String("spanid", *msg.TopicPartition.Metadata))
			//msg.TopicPartition.Metadata = &spanId
			// mtx.Lock()
			// // for _, f := range msg.Headers {
			// // 	if f.Key == "traceparent" {
			// // 		srv.logger.Sugar().Infof("Added traceparent", zap.String("traceparent", string(f.Value)))
			// // 		producerMessageContexts[string(f.Value)] = span
			// // 	}
			// // }

			// mtx.Unlock()

			p.Produce(msg, nil)

			count++
			span.End()
			// run = false
		}
	}
	srv.logger.Sugar().Infof("Producer down")
}
