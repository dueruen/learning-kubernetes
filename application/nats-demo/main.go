package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

var serviceName = "default"
var httpPort = "8080"
var natsURL = nats.DefaultURL
var natsTopic = "default"
var natsProducer = ""

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		publicToNats()
		w.WriteHeader(http.StatusOK)
	case "POST":
		publicToNats()
		w.WriteHeader(http.StatusOK)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getEnvs() {
	serviceNameEnv := os.Getenv("NAME")
	httpPortEnv := os.Getenv("HTTP_PORT")
	natsURLEnv := os.Getenv("NATS_URL")
	natsTopicEnv := os.Getenv("NATS_TOPIC")
	natsProducerEnv := os.Getenv("NATS_PRODUCER")

	if serviceNameEnv != "" {
		serviceName = serviceNameEnv
	}
	if httpPortEnv != "" {
		httpPort = httpPortEnv
	}
	if natsURLEnv != "" {
		natsURL = natsURLEnv
	}
	if natsTopicEnv != "" {
		natsTopic = natsTopicEnv
	}
	if natsProducerEnv != "" {
		natsProducer = natsProducerEnv
	}
}

func publicToNats() {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer ec.Close()

	type data struct {
		Text   string
		Number int
	}

	// Publish the message
	if err := ec.Publish(natsTopic, &data{Text: serviceName, Number: rand.Intn(10000)}); err != nil {
		log.Fatal(err)
	}

	log.Println("Nats done")
}

func main() {
	getEnvs()

	if natsProducer == "" {
		nc, err := nats.Connect(natsURL)
		if err != nil {
			log.Fatal(err)
		}
		defer nc.Close()
		ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
		if err != nil {
			log.Fatal(err)
		}
		defer ec.Close()

		// Define the object
		type data struct {
			Text   string
			Number int
		}

		if _, err := ec.Subscribe(natsTopic, func(s *data) {
			log.Printf("Text: %s - Number: %v", s.Text, s.Number)
			//wg.Done()
		}); err != nil {
			log.Fatal(err)
			return
		}
	} else {
		go func() {
			for {
				publicToNats()

				sleepTime := time.Duration(rand.Intn(5)) * time.Second
				time.Sleep(sleepTime)
			}
		}()
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/readiness", readinessHandler)

	srv := &http.Server{
		Addr:         ":" + httpPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server: " + serviceName + "  on port: " + httpPort)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
