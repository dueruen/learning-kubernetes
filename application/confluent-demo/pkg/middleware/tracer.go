package middleware

import (
	"context"
	"flag"
	"os"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	otelServiceName     = flag.String("otel-service-name", os.Getenv("OTEL_SERVICE_NAME"), "OTEL_SERVICE_NAME")
	version             = flag.String("version", os.Getenv("VERSION"), "VERSION")
	instrumentationName = flag.String("instrumentation-name", os.Getenv("INSTRUMENTATION_NAME"), "INSTRUMENTATION_NAME")
	otelEndpoint        = flag.String("otel-endpoint", os.Getenv("OTEL_ENDPOINT"), "OTEL_ENDPOINT")
)

func InitTracer(logger *zap.Logger, ctx context.Context) (tracerProvider *sdktrace.TracerProvider, tracer trace.Tracer) {
	logger.Debug("Init tracer", zap.String("otelServiceName", *otelServiceName))

	if *otelServiceName == "" {
		logger.Warn("NO OTEL SERVICE NAME - using noop provider")
		nop := trace.NewNoopTracerProvider()
		tracer = nop.Tracer(*otelServiceName)
		return
	}

	options := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(),
	}

	if *otelEndpoint != "" {
		options = append(options, otlptracegrpc.WithEndpoint(*otelEndpoint))
	}

	client := otlptracegrpc.NewClient(options...)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		logger.Error("creating OTLP trace exporter", zap.Error(err))
	}

	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(*otelServiceName),
			semconv.ServiceVersionKey.String(*version),
		)),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		b3.New(),
		&jaeger.Jaeger{},
		&ot.OT{},
	))

	// exporter, err := stdout.New(stdout.WithPrettyPrint())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// tracerProvider = sdktrace.NewTracerProvider(
	// 	sdktrace.WithSampler(sdktrace.AlwaysSample()),
	// 	sdktrace.WithBatcher(exporter),
	// )
	// otel.SetTracerProvider(tracerProvider)
	// otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))

	tracer = tracerProvider.Tracer(
		*instrumentationName,
		trace.WithInstrumentationVersion(*version),
		trace.WithSchemaURL(semconv.SchemaURL),
	)

	return
}

func NewOpenTelemetryMiddleware() mux.MiddlewareFunc {
	return otelmux.Middleware(*otelServiceName)
}
