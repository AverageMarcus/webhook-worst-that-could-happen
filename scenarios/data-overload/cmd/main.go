package main

import (
	"io"
	"log"
	"net/http"

	"worst-webhook/pkg/server"

	"crypto/rand"
)

const (
	name        = "data-overload"
	description = "This scenario attempt to overload the api-server by responding to the webhook request with a large amount of response data."
)

func main() {
	log.Printf("Scenario: %s", name)
	log.Println(description)

	server.StartServer(handler, handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, rand.Reader)
}
