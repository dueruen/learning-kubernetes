package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer
var httpClient http.Client
var logger log.Logger

var serviceName = "default"
var traceEndpoint = "default"
var port = "8000"

func getEnvs() {
	serviceNameEnv := os.Getenv("NAME")
	traceEndpointEnv := os.Getenv("TRACE_ENDPOINT")
	portEnv := os.Getenv("PORT")

	if portEnv != "" {
		port = portEnv
	}
	if serviceNameEnv != "" {
		serviceName = serviceNameEnv
	}
	if traceEndpointEnv != "" {
		traceEndpoint = traceEndpointEnv
	}
}

func main() {
	getEnvs()

	shutdown := initProvider(serviceName, traceEndpoint)
	defer shutdown()

	tracer = otel.Tracer("demo-app")
	httpClient = http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	server := instrumentedServer(handler)

	fmt.Println("listening...")
	server.ListenAndServe()
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	longRunningProcess(ctx)
}

func longRunningProcess(ctx context.Context) {
	ctx, sp := tracer.Start(ctx, "Long Running Process")
	defer sp.End()

	time.Sleep(time.Millisecond * 50)
	sp.AddEvent("halfway done!")
	time.Sleep(time.Millisecond * 50)
}

func instrumentedServer(handler http.HandlerFunc) *http.Server {
	// OpenMetrics handler : metrics and exemplars
	omHandleFunc := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		handler.ServeHTTP(w, r)

		ctx := r.Context()
		traceID := (trace.SpanContextFromContext(ctx).TraceID)().String()

		logger.Log("msg", "http request", "traceID", traceID, "path", r.URL.Path, "latency", time.Since(start))
	}

	// OTel handler : traces
	otelHandler := otelhttp.NewHandler(http.HandlerFunc(omHandleFunc), "http")

	r := mux.NewRouter()
	r.Handle("/", otelHandler)

	return &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:" + port,
	}
}
