package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/javdevapp/docker-chaos-simulator/service"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/container", service.ListContainers)
	r.Handle("/chaos-simulator/container/{containerid}", new(service.ChaosPauseHandler))

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
