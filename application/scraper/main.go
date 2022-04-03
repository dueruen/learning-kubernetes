package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var port = "80"
var scrapeEndpoints = "some,more"

func main() {
	portEnv := os.Getenv("PORT")
	scrapeEndpointsEnv := os.Getenv("SCRAPE_ENDPOINTS")
	if portEnv != "" {
		port = portEnv
	}
	if scrapeEndpointsEnv != "" {
		scrapeEndpoints = scrapeEndpointsEnv
	}
	log.Println("Service scraper is starting...")

	go startMetricScraping(scrapeEndpoints)

	log.Println("Service scraper api running listering at port: " + port)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/logs", handleLogs)
	r.Post("/traces", handleLogs)
	http.ListenAndServe(":"+port, r)

	//http.HandleFunc("/", apiHandlerInternal)
	// http.ListenAndServe(":"+port, nil)
}

func handleLogs(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body) //body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		log.Print("Method:: ", r.Method+" - READY TO SEND LOGS TO KAFKA   body size: "+strconv.Itoa(len(body)))
	} else {
		log.Print("Method:: ", r.Method+" - No body")
	}
}

func handleTraces(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body) //body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		log.Print("Method:: ", r.Method+" - READY TO SEND TRACES TO KAFKA   body size: "+strconv.Itoa(len(body)))
	} else {
		log.Print("Method:: ", r.Method+" - No body")
	}
}

func startMetricScraping(scrapeEndpoints string) {
	endpointList := strings.Split(scrapeEndpoints, ",")

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
