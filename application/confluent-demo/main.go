package main

import (
	"flag"

	"github.com/dueruen/learning-kubernetes/application/confluent-demo/pkg/http"
	"github.com/dueruen/learning-kubernetes/application/confluent-demo/pkg/kafka"
	"github.com/dueruen/learning-kubernetes/application/confluent-demo/pkg/middleware"
	"go.uber.org/zap"
)

func main() {
	flag.Parse()

	logger, _ := middleware.InitZap()
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	srv, _ := http.NewServer(logger)
	httpDone := false
	httpShutdown := srv.ListenAndServe()

	kafkaShutdown := kafka.InitKafka(logger)
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
