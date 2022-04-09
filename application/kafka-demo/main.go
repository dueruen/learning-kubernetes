package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	otelServiceName = flag.String("otel-service-name", os.Getenv("OTEL_SERVICE_NAME"), "Service name for reporting to open telemetry address, when not set tracing is disabled")
)

func main() {
	flag.Parse()

	logger, _ := initZap("debug")
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	srvCfg := Config{
		Unhealthy:                 false,
		Unready:                   false,
		HttpServerTimeout:         30 * time.Second,
		HttpServerShutdownTimeout: 30 * time.Second,
		Host:                      "",
		Port:                      "9898",
	}

	srv, _ := NewServer(&srvCfg, logger)
	stopCh := SetupSignalHandler()
	srv.ListenAndServe(stopCh)
}

func initZap(logLevel string) (*zap.Logger, error) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	}

	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	zapConfig := zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zapEncoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return zapConfig.Build()
}

var onlyOneSignalHandler = make(chan struct{})

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/go-kit/log"
// 	"github.com/gorilla/mux"
// 	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/trace"
// )

// var tracer trace.Tracer
// var httpClient http.Client
// var logger log.Logger

// var serviceName = "default"
// var traceEndpoint = "default"
// var traceURL = "/v1/traces"
// var port = "8000"

// func getEnvs() {
// 	serviceNameEnv := os.Getenv("NAME")
// 	traceEndpointEnv := os.Getenv("TRACE_ENDPOINT")
// 	traceURLEnv := os.Getenv("TRACE_URL")
// 	portEnv := os.Getenv("PORT")

// 	if portEnv != "" {
// 		port = portEnv
// 	}
// 	if serviceNameEnv != "" {
// 		serviceName = serviceNameEnv
// 	}
// 	if traceEndpointEnv != "" {
// 		traceEndpoint = traceEndpointEnv
// 	}
// 	if traceURLEnv != "" {
// 		traceURL = traceURLEnv
// 	}
// }

// func main() {
// 	getEnvs()

// 	shutdown := initProvider(serviceName, traceEndpoint, traceURL)
// 	defer shutdown()

// 	tracer = otel.Tracer("demo-app")
// 	httpClient = http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
// 	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
// 	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

// 	server := instrumentedServer(handler)

// 	fmt.Println("listening...")
// 	server.ListenAndServe()
// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	longRunningProcess(ctx)
// }

// func longRunningProcess(ctx context.Context) {
// 	ctx, sp := tracer.Start(ctx, "Long Running Process")
// 	defer sp.End()

// 	time.Sleep(time.Millisecond * 50)
// 	sp.AddEvent("halfway done!")
// 	time.Sleep(time.Millisecond * 50)
// }

// func instrumentedServer(handler http.HandlerFunc) *http.Server {
// 	// OpenMetrics handler : metrics and exemplars
// 	omHandleFunc := func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()

// 		handler.ServeHTTP(w, r)

// 		ctx := r.Context()
// 		traceID := (trace.SpanContextFromContext(ctx).TraceID)().String()

// 		logger.Log("msg", "http request", "traceID", traceID, "path", r.URL.Path, "latency", time.Since(start))
// 	}

// 	// OTel handler : traces
// 	otelHandler := otelhttp.NewHandler(http.HandlerFunc(omHandleFunc), "http")

// 	r := mux.NewRouter()
// 	r.Handle("/", otelHandler)

// 	return &http.Server{
// 		Handler: r,
// 		Addr:    "0.0.0.0:" + port,
// 	}
// }
