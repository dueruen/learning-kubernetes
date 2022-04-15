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
			_, span := srv.tracer.Start(context.Background(), "produce message")

			number := rand.Intn(60)

			recordKey := "alice"
			data := &KafkaMessage{
				Message: "Hello from a consumer random number: " + strconv.Itoa(number),
				Count:   count,
			}
			recordValue, _ := json.Marshal(&data)
			srv.logger.Sugar().Debugf("Preparing to produce record: %s\t%s\n", recordKey, recordValue)
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
				Key:            []byte(recordKey),
				Value:          []byte(recordValue),
				Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
			}, nil)

			span.End()
			time.Sleep(time.Duration(number/10) * time.Second)
			count++
		}
	}
	srv.logger.Sugar().Infof("Producer down")
}
