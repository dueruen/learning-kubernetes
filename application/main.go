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
	ServiceName     string `json:"name"`
	HttpPort        string `json:"httpport"`
	ExternalGetURI  string `json:"externalGetURI"`
	ExternalPostURI string `json:"externalPostURI"`
	PostSleep       int    `json:"postSleep"`
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
		if serviceData.ExternalGetURI == "" {
			log.Println("No externalURI found")
			apiHandlerInternal(w, r)
			return
		}

		uri := serviceData.ExternalGetURI + "/api/external"

		resp, err := http.Get(uri)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(body)

	case "POST":
		log.Println(serviceData.ServiceName, " External: ", r.Method)
		if serviceData.ExternalPostURI == "" {
			log.Println("No externalURI found")
			apiHandlerInternal(w, r)
			return
		}

		uri := serviceData.ExternalPostURI + "/api/external"

		postBody, _ := json.Marshal(map[string]string{
			"name":  "test",
			"email": "test@example.com",
		})
		responseBody := bytes.NewBuffer(postBody)

		resp, err := http.Post(uri, "application/json", responseBody)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
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
	var externalGetURI = os.Getenv("GET_URI")
	var externalPostURI = os.Getenv("POST_URI")
	var postSleep = os.Getenv("POST_SLEEP")

	if serviceName != "" {
		serviceData.ServiceName = serviceName
	}
	if httpPort != "" {
		serviceData.HttpPort = httpPort
	}

	log.Println("Service ", serviceData.ServiceName, " is starting...")

	serviceData.ExternalGetURI = externalGetURI
	log.Println("ExternalGetURI: ", serviceData.ExternalGetURI)

	serviceData.ExternalPostURI = externalPostURI
	log.Println("ExternalPostURI: ", serviceData.ExternalPostURI)

	if postSleep != "" {
		sleepNumber, err := strconv.Atoi(postSleep)
		if err != nil {
			log.Println(err)
			return
		}
		serviceData.PostSleep = sleepNumber
	}

	log.Println("Alpha api running listering at localhost:", serviceData.HttpPort)

	http.HandleFunc("/api", apiHandlerInternal)
	http.HandleFunc("/api/external", apiHandlerExternal)
	http.ListenAndServe(":"+serviceData.HttpPort, nil)
}
