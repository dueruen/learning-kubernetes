package kafka

import (
	"encoding/json"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

func startProducer(bootstrapServer string, topic *string, producerShutdown chan bool, logger *zap.Logger) {
	logger.Sugar().Infof("Starting producer - server: %s - topic: %s\n", bootstrapServer, *topic)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServer})
	if err != nil {
		logger.Sugar().Errorf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	logger.Sugar().Infof("Created Producer %v\n", p)

	go func() {
		logger.Sugar().Debugf("Loop events")
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.Sugar().Errorf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					logger.Sugar().Debugf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
		logger.Sugar().Infof("Loop events OVER")
	}()

	count := 0
	run := true
	logger.Sugar().Infof("Producer running")
	for run == true {
		select {
		case sig := <-sigchan:
			p.Flush(10)
			logger.Sugar().Infof("Caught signal %v: terminating producer\n", sig)
			producerShutdown <- true
			run = false
		default:
			number := rand.Intn(60)

			recordKey := "alice"
			data := &KafkaMessage{
				Message: "Hello from a consumer random number: " + strconv.Itoa(number),
				Count:   count,
			}
			recordValue, _ := json.Marshal(&data)
			logger.Sugar().Debugf("Preparing to produce record: %s\t%s\n", recordKey, recordValue)
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
				Key:            []byte(recordKey),
				Value:          []byte(recordValue),
				Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
			}, nil)

			time.Sleep(time.Duration(number/10) * time.Second)
			count++
		}
	}
	logger.Sugar().Infof("Producer down")
}
