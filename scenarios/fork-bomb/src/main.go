package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	admission "k8s.io/api/admission/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	port    string
	crtPath string
	keyPath string

	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
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

func main() {
	http.HandleFunc("/mutate", logRequest(mutate))
	http.HandleFunc("/validate", logRequest(validate))
	log.Printf("Webhook server listening on port %s", port)
	log.Fatal(http.ListenAndServeTLS(":"+port, crtPath, keyPath, nil))
}

func mutate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not read request body: %v", err), http.StatusBadRequest)
		return
	}

	review := &admission.AdmissionReview{}
	if _, _, err := codecs.UniversalDeserializer().Decode(body, nil, review); err != nil {
		http.Error(w, fmt.Sprintf("could not deserialize request: %v", err), http.StatusBadRequest)
		return
	}

	response := &admission.AdmissionResponse{
		UID:     review.Request.UID,
		Allowed: true,
		Result: &v1.Status{
			Message: "Allowed",
		},
		Patch: []byte(""),
	}

	review.Response = response
	responseBytes, err := json.Marshal(review)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not write response body: %v", err), http.StatusBadRequest)
		return
	}
	w.Write(responseBytes)
}

func validate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not read request body: %v", err), http.StatusBadRequest)
		return
	}

	review := &admission.AdmissionReview{}
	if _, _, err := codecs.UniversalDeserializer().Decode(body, nil, review); err != nil {
		http.Error(w, fmt.Sprintf("could not deserialize request: %v", err), http.StatusBadRequest)
		return
	}

	response := &admission.AdmissionResponse{
		UID:     review.Request.UID,
		Allowed: true,
		Result: &v1.Status{
			Message: "Allowed",
		},
	}

	review.Response = response
	responseBytes, err := json.Marshal(review)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not write response body: %v", err), http.StatusBadRequest)
		return
	}
	w.Write(responseBytes)
}
