package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	brokers                 = flag.String("brokers", os.Getenv("KAFKA_BROKERS"), "KAFKA_BROKERS")
	otelEndpoint            = flag.String("otel-endpoint", os.Getenv("OTEL_ENDPOINT"), "OTEL_ENDPOINT")
	topic                   = flag.String("topic", os.Getenv("TOPIC"), "TOPIC")
	producer                = flag.String("producer", os.Getenv("PRODUCER"), "PRODUCER")
	serviceName             = flag.String("service-name", os.Getenv("SERVICE_NAME"), "SERVICE_NAME")
	instrument              = flag.String("app-instrument", os.Getenv("APP_INSTRUMENT"), "APP_INSTRUMENT")
	messageSize             = flag.String("message-size", os.Getenv("MESSAGE_SIZE"), "MESSAGE_SIZE")
	messageFrequency        = flag.String("message-frequency", os.Getenv("MESSAGE_FREQUENCY"), "MESSAGE_FREQUENCY")
	withConsumerWorkTime    = flag.String("consumer-work-time", os.Getenv("CONSUMER_WORK_TIME"), "CONSUMER_WORK_TIME")
	withConsumerRandomError = flag.String("consumer-random-error", os.Getenv("CONSUMER_RANDOM_ERROR"), "CONSUMER_RANDOM_ERROR")
	appName                 = flag.String("app-name", os.Getenv("APP_NAME"), "APP_NAME")
	debug                   = flag.String("debug", os.Getenv("DEBUG"), "DEBUG")
)

func IsInstrumented() bool {
	return strings.ToLower(*instrument) == "enable"
}

func IsWithWorkTime() bool {
	return strings.ToLower(*withConsumerWorkTime) == "enable"
}

func IsWithRandomError() bool {
	return strings.ToLower(*withConsumerRandomError) == "enable"
}

func IsInDebug() bool {
	return strings.ToLower(*debug) == "enable"
}

func main() {
	flag.Parse()

	validateInputs()

	if IsInstrumented() {
		tp := InitTracer()
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()
	}

	brokerList := strings.Split(*brokers, ",")

	var shutdown chan bool
	if *producer == "" {
		shutdown = StartConsumer(brokerList, *topic)
	} else {
		shutdown = StartProducer(brokerList, *topic)
	}

	select {
	case _ = <-shutdown:
		return
	}
}

func validateInputs() {
	flag.PrintDefaults()

	if *brokers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	fmt.Println("Brokers: ", *brokers)

	if *otelEndpoint == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *topic == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *appName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func InitTracer() *sdktrace.TracerProvider {
	options := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(*otelEndpoint),
	}

	client := otlptracegrpc.NewClient(options...)
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(*serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))
	return tp
}
