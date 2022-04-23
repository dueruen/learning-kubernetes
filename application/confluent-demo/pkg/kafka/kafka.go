package kafka

import (
	"context"
	"flag"
	"os"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	bootstrapServer = flag.String("bootstrap-server", os.Getenv("BOOTSTRAP_SERVER"), "BOOTSTRAP_SERVER")
	producer        = flag.String("producer", os.Getenv("PRODUCER"), "PRODUCER")
	producerTopic   = flag.String("producer-topic", os.Getenv("PRODUCER_TOPIC"), "PRODUCER_TOPIC")
	consumer        = flag.String("consumer", os.Getenv("CONSUMER"), "CONSUMER")
	consumerTopic   = flag.String("consumer-topic", os.Getenv("CONSUMER_TOPIC"), "CONSUMER_TOPIC")
	consumerGrupId  = flag.String("consumer-group-id", os.Getenv("CONSUMER_GROUP_ID"), "CONSUMER_GROUP_ID")
)

type KafkaServer struct {
	logger         *zap.Logger
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
	context        context.Context
}

func NewKafkaServer(logger *zap.Logger, tracerProvider *sdktrace.TracerProvider, tracer trace.Tracer, context context.Context) (*KafkaServer, error) {
	srv := &KafkaServer{
		logger:         logger,
		tracer:         tracer,
		tracerProvider: tracerProvider,
		context:        context,
	}

	return srv, nil
}

func (srv *KafkaServer) InitKafka() chan bool {
	shutdown := make(chan bool)

	if consumerGrupId == nil {
		*consumerGrupId = "go_example_group_1"
	}

	if producerTopic == nil {
		*producerTopic = "default"
	}

	if consumerTopic == nil {
		*consumerTopic = "default"
	}

	waitForProducer := false
	producerShutdown := make(chan bool)
	if producer != nil && *producer != "" {
		waitForProducer = true
		go srv.startProducer(*bootstrapServer, producerTopic, producerShutdown)
	}

	waitForConsumer := false
	consumerShutdown := make(chan bool)
	if consumer != nil && *consumer != "" {
		waitForConsumer = true
		go srv.startConsumer(*bootstrapServer, *consumerGrupId, consumerTopic, consumerShutdown, producerTopic)
	}

	go func() {
		srv.logger.Debug("Kafka running")
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
		shutdown <- true
	}()

	return shutdown
}
