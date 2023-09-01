package server

import (
	"log"
	"net/http"
	"os"
)

var (
	port    string
	crtPath string
	keyPath string
)

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	crtPath = os.Getenv("CERT_PATH")
	if crtPath == "" {
		crtPath = "/tls/tls.crt"
	}

	keyPath = os.Getenv("KEY_PATH")
	if keyPath == "" {
		keyPath = "/tls/tls.key"
	}
}

func logRequest(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		f(w, r)
	}
}
func StartServer(mutate, validate func(w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc("/mutate", logRequest(mutate))
	http.HandleFunc("/validate", logRequest(validate))
	log.Printf("Webhook server listening on port %s", port)
	log.Fatal(http.ListenAndServeTLS(":"+port, crtPath, keyPath, nil))
}
