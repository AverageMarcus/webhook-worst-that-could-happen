package main

import (
	"log"
	"net/http"

	"worst-webhook/pkg/server"
)

const (
	name        = "open-connection"
	description = "This scenario attempt to leave the incoming HTTP request open."
)

func main() {
	log.Printf("Scenario: %s", name)
	log.Println(description)

	server.StartServer(handler, handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Println("Failed to get http.Flusher")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	<-r.Context().Done()
	log.Println("Client hung up")
}
