module github.com/dueruen/learning-kubernetes/application/confluent-demo

go 1.18

require (
	github.com/confluentinc/confluent-kafka-go v1.8.2
	github.com/gorilla/mux v1.8.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.31.0
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.31.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.31.0
	go.opentelemetry.io/contrib/propagators/b3 v1.6.0
	go.opentelemetry.io/contrib/propagators/jaeger v1.6.0
	go.opentelemetry.io/contrib/propagators/ot v1.6.0
	go.opentelemetry.io/otel v1.6.3
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.6.3
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.6.3
	go.opentelemetry.io/otel/sdk v1.6.3
	go.opentelemetry.io/otel/trace v1.6.3
	go.uber.org/zap v1.21.0
)

require (
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.6.3 // indirect
	go.opentelemetry.io/otel/metric v0.28.0 // indirect
	go.opentelemetry.io/proto/otlp v0.15.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20211118181313-81c1377c94b1 // indirect
	google.golang.org/grpc v1.45.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
