package main

import (
	"context"
	"flag"

	"github.com/dueruen/learning-kubernetes/application/confluent-demo/pkg/http"
	"github.com/dueruen/learning-kubernetes/application/confluent-demo/pkg/kafka"
	"github.com/dueruen/learning-kubernetes/application/confluent-demo/pkg/middleware"
	"go.uber.org/zap"
)

func main() {
	flag.Parse()
	ctx := context.Background()

	logger, _ := middleware.InitZap()
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	tracerProvider, tracer := middleware.InitTracer(logger, ctx)
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			logger.Warn("stopping tracer provider", zap.Error(err))
		}
	}()

	srv, _ := http.NewServer(logger, tracerProvider, tracer)
	httpDone := false
	httpShutdown := srv.ListenAndServe()

	kafkaServer, _ := kafka.NewKafkaServer(logger, tracerProvider, tracer, ctx)
	kafkaShutdown := kafkaServer.InitKafka()
	kafkaDone := false
	run := true
	for run {
		select {
		case _ = <-kafkaShutdown:
			kafkaDone = true
		case _ = <-httpShutdown:
			httpDone = true
		default:
			if kafkaDone && httpDone {
				run = false
			}
		}
	}

	logger.Sugar().Infof("Shutting down")
}
