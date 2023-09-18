package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"worst-webhook/pkg/server"
)

const (
	name        = "redirect"
	description = "This scenario responds to all webhook requests with a redirect to a service that infinitely redirects the client."
)

func main() {
	log.Printf("Scenario: %s", name)
	log.Println(description)

	server.StartServer(handler, handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", fmt.Sprintf("https://redirect.default.svc/mutate?%d", time.Now().Unix()))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
