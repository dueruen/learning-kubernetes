package main

import (
	"context"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	// "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

const (
	instrumentationName = "kafka-demo"
)

func (s *Server) initTracer(ctx context.Context) {
	s.logger.Info("otelServiceName", zap.String("otelServiceName", *otelServiceName))
	if *otelServiceName == "" {
		nop := trace.NewNoopTracerProvider()
		s.tracer = nop.Tracer(*otelServiceName)
		return
	}

	options := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(),
	}

	client := otlptracegrpc.NewClient(options...)
	//client := otlptracegrpc.NewHTTPConfig()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		s.logger.Error("creating OTLP trace exporter", zap.Error(err))
	}

	s.tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(*otelServiceName),
			semconv.ServiceVersionKey.String("1.0.0"),
		)),
	)

	otel.SetTracerProvider(s.tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		b3.New(),
		&jaeger.Jaeger{},
		&ot.OT{},
	))

	s.tracer = s.tracerProvider.Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion("1.0.0"),
		trace.WithSchemaURL(semconv.SchemaURL),
	)
}

func NewOpenTelemetryMiddleware() mux.MiddlewareFunc {
	return otelmux.Middleware(*otelServiceName)
}

// func initProvider(serviceName string, traceEndpoint string, traceURL string) func() {
// 	ctx := context.Background()

// 	resource := createResource(ctx, serviceName)

// 	exporter := createExporter(ctx, traceEndpoint, traceURL)

// 	provider := createProvider(exporter, resource)

// 	otel.SetTracerProvider(provider)
// 	otel.SetTextMapPropagator(propagation.TraceContext{})

// 	return func() {
// 		// Shutdown will flush any remaining spans and shut down the exporter.
// 		handleErr(provider.Shutdown(ctx), "failed to shutdown TracerProvider")
// 	}
// }

// func createResource(ctx context.Context, serviceName string) *resource.Resource {
// 	res, err := resource.New(ctx,
// 		resource.WithAttributes(
// 			semconv.ServiceNameKey.String(serviceName),
// 		),
// 	)
// 	handleErr(err, "failed to create resource")

// 	return res
// }

// func createExporter(ctx context.Context, traceEndpoint string, traceURL string) *otlptrace.Exporter {
// 	// conn, err := grpc.DialContext(ctx,
// 	// 	traceEndpoint,
// 	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	// 	grpc.WithBlock())

// 	// handleErr(err, "failed to create gRPC connection to collector")

// 	options := []otlptracehttp.Option{
// 		otlptracehttp.WithInsecure(),
// 		otlptracehttp.WithEndpoint(traceEndpoint),
// 		otlptracehttp.WithURLPath(traceURL)}

// 	traceExporter, err := otlptracehttp.New(ctx, options...)
// 	handleErr(err, "failed to create trace exporter")

// 	return traceExporter
// }

// // func createExporter(ctx context.Context, traceEndpoint string) *otlptrace.Exporter {
// // 	conn, err := grpc.DialContext(ctx,
// // 		traceEndpoint,
// // 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// // 		grpc.WithBlock())

// // 	handleErr(err, "failed to create gRPC connection to collector")

// // 	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
// // 	handleErr(err, "failed to create trace exporter")

// // 	return traceExporter
// // }

// func createProvider(exporter *otlptrace.Exporter, resource *resource.Resource) *sdktrace.TracerProvider {
// 	bsp := sdktrace.NewBatchSpanProcessor(exporter)
// 	tracerProvider := sdktrace.NewTracerProvider(
// 		sdktrace.WithSampler(sdktrace.AlwaysSample()),
// 		sdktrace.WithResource(resource),
// 		sdktrace.WithSpanProcessor(bsp),
// 	)

// 	return tracerProvider
// }

// func handleErr(err error, message string) {
// 	if err != nil {
// 		panic(fmt.Sprintf("%s: %s", err, message))
// 	}
// }
