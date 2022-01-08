package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	Trace *log.Logger
	Info  *log.Logger
)

type alpha struct {
	Name string `json:"name"`
	Id   int    `json:"number"`
}

var alphaOne *alpha = &alpha{
	Name: "Alpha",
	Id:   42,
}

func apiHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		log.Println("From api: Get alpha")
		j, _ := json.Marshal(alphaOne)
		w.Write(j)

	case "POST":
		log.Println("From api: Post alpha")
		d := json.NewDecoder(r.Body)
		p := &alpha{}
		err := d.Decode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		alphaOne = p

	default:
		log.Println("From api: endpoint not found")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func main() {
	log.Println("Alpha api running")
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/api/test", apiHandler)
	http.ListenAndServe(":8080", nil)
}
