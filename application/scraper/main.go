package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	port            = flag.String("port", os.Getenv("PORT"), "The port to serve http requests")
	brokers         = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma separated list")
	scrapeEndpoints = flag.String("scrapeEndpoints", os.Getenv("SCRAPE_ENDPOINTS"), "The http metric endpoints to scrape, as a comma separated list")
)

var (
	logTopic    = "logs"
	metricTopic = "metrics"
	traceTopic  = "traces"
)

func main() {
	flag.Parse()

	if *port == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Println("Service scraper is starting...")

	go startMetricScraping(*scrapeEndpoints)

	log.Println("Service scraper api running listering at port: " + *port)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/logs", handleLogs)
	r.Post("/traces", handleLogs)
	http.ListenAndServe(":"+*port, r)
}

func handleLogs(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		_, err := ioutil.ReadAll(r.Body) //body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		//log.Print("Method:: ", r.Method+" - READY TO SEND LOGS TO KAFKA   body size: "+strconv.Itoa(len(body)))
	} else {
		log.Print("Method:: ", r.Method+" - No body")
	}
}

func handleTraces(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		_, err := ioutil.ReadAll(r.Body) //body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		//log.Print("Method:: ", r.Method+" - READY TO SEND TRACES TO KAFKA   body size: "+strconv.Itoa(len(body)))
	} else {
		log.Print("Method:: ", r.Method+" - No body")
	}
}

func startMetricScraping(scrapeEndpoints string) {
	endpointList := strings.Split(scrapeEndpoints, ",")

	if len(endpointList) == 0 || endpointList[0] == "" {
		log.Printf("No scrape endpoints provided, no scraping")
		return
	}

	for _, endpoint := range endpointList {
		go scrapeMetrics(endpoint)
	}

	for {
	}
}

func scrapeMetrics(endpoint string) {
	for {
		time.Sleep(2 * time.Second)

		resp, err := http.Get(endpoint)
		if err != nil {
			log.Println(err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body) //body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("HAVE SCRAPED endpoint:: " + endpoint + "  body size: " + strconv.Itoa(len(body)))
	}
}

func publicToKafka(brokerList string, topicName *string) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokerList,
		"client.id":         "clientID",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	delivery_chan := make(chan kafka.Event, 10000)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: topicName, Partition: kafka.PartitionAny},
		Value:          []byte("test test test")},
		delivery_chan,
	)

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()
}

// func apiHandlerInternal(w http.ResponseWriter, r *http.Request) {
// 	start := time.Now()
// 	ti := 2
// 	if rand.Intn(3) == 1 {
// 		ti = 0
// 	}

// 	time.Sleep(time.Duration(ti) * time.Second)

// 	if r.Body != nil {
// 		body, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			log.Printf("Error reading body: %v", err)
// 			http.Error(w, "can't read body", http.StatusBadRequest)
// 			return
// 		}
// 		log.Print("Method:: ", r.Method+" - Body:: "+string(body))
// 	} else {
// 		log.Print("Method:: ", r.Method+" - No body")
// 	}

// 	latency := time.Since(start)

// 	text := "Sleep: " + strconv.FormatInt(int64(ti), 10) + "  latency: " + strconv.FormatInt(int64(latency), 10)

// 	log.Printf(text)
// 	w.Write([]byte(text))
// }
