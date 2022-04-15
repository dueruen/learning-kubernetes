package kafka

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	getter "github.com/dueruen/learning-kubernetes/application/confluent-demo/pkg/http"
	"go.opentelemetry.io/otel/baggage"
	"go.uber.org/zap"
)

func startConsumer(bootstrapServer string, groupId string, topic *string, consumerShutdown chan bool, logger *zap.Logger) {
	logger.Sugar().Infof("Starting consumer - server: %s - groupId: %s - topic: %s\n", bootstrapServer, groupId, *topic)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.Sugar().Errorf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	err = c.SubscribeTopics([]string{*topic}, nil)
	if err != nil {
		logger.Sugar().Errorf("Failed to subsribe to topic: %s  - Error: %s", *topic, err)
		os.Exit(1)
	}

	totalCount := 0
	run := true
	logger.Sugar().Infof("Consumer running")
	for run == true {
		select {
		case sig := <-sigchan:
			logger.Sugar().Infof("Caught signal %v: terminating consumer\n", sig)
			consumerShutdown <- true
			run = false
		default:
			msg, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}

			if msg.Headers != nil {
				logger.Sugar().Debugf("%% Headers: %v\n", msg.Headers)
			}

			recordKey := string(msg.Key)
			message := msg.Value
			data := KafkaMessage{}
			err = json.Unmarshal(message, &data)
			if err != nil {
				logger.Sugar().Errorf("Failed to decode JSON at offset %d: %v", msg.TopicPartition.Offset, err)
				continue
			}
			count := data.Count
			totalCount = count
			logger.Sugar().Debugf("Consumed record with key %s and value %d, and updated total count to %d -- Message was: %s\n", recordKey, data.Count, totalCount, data.Message)

			go func() {
				ctx := baggage.ContextWithoutBaggage(context.Background())
				getter.GetFromBackend(logger, ctx)
			}()
		}
	}

	logger.Sugar().Infof("END of consumer")
	c.Close()
}
