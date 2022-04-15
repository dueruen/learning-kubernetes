package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	url          = flag.String("server", os.Getenv("SERVER"), "URL of server to call")
	otelEndpoint = flag.String("otel-endpoint", os.Getenv("OTEL_ENDPOINT"), "OTEL_ENDPOINT")
)

func initTracer(ctx context.Context) *sdktrace.TracerProvider {
	// Create stdout exporter to be able to retrieve
	// the collected spans.
	options := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(*otelEndpoint),
	}

	client := otlptracegrpc.NewClient(options...)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func main() {
	flag.Parse()
	bag, _ := baggage.Parse("username=donuts")
	ctx := baggage.ContextWithBaggage(context.Background(), bag)

	tp := initTracer(ctx)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	for {
		client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

		var body []byte

		tr := otel.Tracer("example/client")
		err := func(ctx context.Context) error {
			ctx, span := tr.Start(ctx, "say hello", trace.WithAttributes(semconv.PeerServiceKey.String("ExampleService")))
			defer span.End()
			req, _ := http.NewRequestWithContext(ctx, "GET", *url, nil)

			fmt.Printf("Sending request...\n")
			res, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return err
			}
			body, err = ioutil.ReadAll(res.Body)
			_ = res.Body.Close()

			return err
		}(ctx)

		if err != nil {
			log.Println("ERROR waiting a few seconds")
			time.Sleep(10 * time.Second)
			continue
		}

		fmt.Printf("Response Received: %s\n\n\n", body)
		fmt.Printf("Waiting for few seconds to export spans ...\n\n")
		time.Sleep(10 * time.Second)
		fmt.Printf("Inspect traces on stdout\n")
	}
}
