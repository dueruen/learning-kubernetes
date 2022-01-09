package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Service struct {
	ServiceName      string `json:"name"`
	HttpPort         string `json:"httpport"`
	ExternalGetPort  string `json:"getport"`
	ExternalPostPort string `json:"postport"`
	PostSleep        int    `json:"postSleep"`
}

var serviceData *Service = &Service{
	ServiceName: "DefaultName",
	HttpPort:    "8080",
	PostSleep:   1,
}

func apiHandlerInternal(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		log.Println(serviceData.ServiceName, " Internal: ", r.Method)
		j, _ := json.Marshal(serviceData)
		w.Write(j)

	case "POST":
		log.Println(serviceData.ServiceName, " Internal: ", r.Method)
		time.Sleep(time.Duration(serviceData.PostSleep) * time.Second)

	default:
		log.Println(serviceData.ServiceName, " Internal: ", r.Method, "  was not found")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func apiHandlerExternal(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		log.Println(serviceData.ServiceName, " External: ", r.Method)
		var uri = ""
		if serviceData.ExternalGetPort == "" {
			uri = "http://localhost:" + serviceData.HttpPort + "/api"
		} else {
			uri = "http://localhost:" + serviceData.ExternalGetPort + "/api/external"
		}

		resp, err := http.Get(uri)
		if err != nil {
			log.Fatalln(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(body)

	case "POST":
		log.Println(serviceData.ServiceName, " External: ", r.Method)
		var uri = ""
		if serviceData.ExternalPostPort == "" {
			uri = "http://localhost:" + serviceData.HttpPort + "/api"
		} else {
			uri = "http://localhost:" + serviceData.ExternalPostPort + "/api/external"
		}

		postBody, _ := json.Marshal(map[string]string{
			"name":  "test",
			"email": "test@example.com",
		})
		responseBody := bytes.NewBuffer(postBody)

		resp, err := http.Post(uri, "application/json", responseBody)
		if err != nil {
			log.Fatalln(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(body)

	default:
		log.Println(serviceData.ServiceName, " External: ", r.Method, "  was not found")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func main() {
	var serviceName = os.Getenv("NAME")
	var httpPort = os.Getenv("HTTP_PORT")
	var externalGetPort = os.Getenv("GET_PORT")
	var externalPostPort = os.Getenv("POST_PORT")
	var postSleep = os.Getenv("POST_SLEEP")

	if serviceName != "" {
		serviceData.ServiceName = serviceName
	}
	if httpPort != "" {
		serviceData.HttpPort = httpPort
	}

	log.Println("Service ", serviceData.ServiceName, " is starting...")

	serviceData.ExternalGetPort = externalGetPort
	serviceData.ExternalPostPort = externalPostPort
	if postSleep != "" {
		sleepNumber, err := strconv.Atoi(postSleep)
		if err != nil {
			log.Fatalln(err)
			return
		}
		serviceData.PostSleep = sleepNumber
	}

	log.Println("Alpha api running listering at localhost:", serviceData.HttpPort)

	http.HandleFunc("/api", apiHandlerInternal)
	http.HandleFunc("/api/external", apiHandlerExternal)
	http.ListenAndServe(":"+serviceData.HttpPort, nil)
}
