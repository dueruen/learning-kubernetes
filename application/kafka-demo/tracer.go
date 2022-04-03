package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"

	// "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func initProvider(serviceName string, traceEndpoint string) func() {
	ctx := context.Background()

	resource := createResource(ctx, serviceName)

	exporter := createExporter(ctx, traceEndpoint)

	provider := createProvider(exporter, resource)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		// Shutdown will flush any remaining spans and shut down the exporter.
		handleErr(provider.Shutdown(ctx), "failed to shutdown TracerProvider")
	}
}

func createResource(ctx context.Context, serviceName string) *resource.Resource {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	handleErr(err, "failed to create resource")

	return res
}

func createExporter(ctx context.Context, traceEndpoint string) *otlptrace.Exporter {
	// conn, err := grpc.DialContext(ctx,
	// 	traceEndpoint,
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 	grpc.WithBlock())

	// handleErr(err, "failed to create gRPC connection to collector")

	options := []otlptracehttp.Option{
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(traceEndpoint),
		otlptracehttp.WithURLPath("/")}

	traceExporter, err := otlptracehttp.New(ctx, options...)
	handleErr(err, "failed to create trace exporter")

	return traceExporter
}

// func createExporter(ctx context.Context, traceEndpoint string) *otlptrace.Exporter {
// 	conn, err := grpc.DialContext(ctx,
// 		traceEndpoint,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithBlock())

// 	handleErr(err, "failed to create gRPC connection to collector")

// 	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
// 	handleErr(err, "failed to create trace exporter")

// 	return traceExporter
// }

func createProvider(exporter *otlptrace.Exporter, resource *resource.Resource) *sdktrace.TracerProvider {
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource),
		sdktrace.WithSpanProcessor(bsp),
	)

	return tracerProvider
}

func handleErr(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", err, message))
	}
}
