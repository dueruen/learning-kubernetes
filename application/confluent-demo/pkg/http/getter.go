package http

import (
	"context"
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	backendURL = flag.String("backend-url", os.Getenv("BACKEND_URL"), "BACKEND_URL")
)

func GetFromBackend(logger *zap.Logger, ctx context.Context) {
	// _, span := s.tracer.Start(r.Context(), "Long running task - handler")
	span := trace.SpanFromContext(ctx)
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	if *backendURL != "" {
		ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))
		req, _ := http.NewRequestWithContext(ctx, "GET", *backendURL, nil)

		logger.Debug("Sending request", zap.String("url", *backendURL))
		res, err := client.Do(req)
		if err != nil {
			logger.Error("Failed to do request", zap.String("url", *backendURL), zap.Error(err))
			return
		}

		body, err := ioutil.ReadAll(res.Body)
		logger.Debug("Request returned", zap.String("url", *backendURL), zap.String("body", string(body)))
		span.End()
	}
}
