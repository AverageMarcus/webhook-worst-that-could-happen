package main

import (
	"log"
	"net/http"

	"worst-webhook/pkg/handler"
	"worst-webhook/pkg/patcher"
	"worst-webhook/pkg/server"

	admission "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	name        = "cut-off-response"
	description = "In this scenario we're going to set the `Content-Length` response header to be longer than the actual response body"
)

func main() {
	log.Printf("Scenario: %s", name)
	log.Println(description)

	server.StartServer(
		handler.Handle(mutate, contentLengthHook),
		handler.Handle(validate, contentLengthHook),
	)
}

func contentLengthHook(w http.ResponseWriter, r *http.Request) {
	log.Printf("Setting content-length header")
	w.Header().Set("Content-Length", "18000")
}

func mutate(review *admission.AdmissionReview, response *admission.AdmissionResponse) error {
	patchBytes, err := patcher.New().
		Add(
			patcher.JSONPatchEntry{
				OP:    "add",
				Path:  "/metadata/annotations",
				Value: patcher.ToValueMap(map[string]string{}),
			},
		).
		ToBytes()
	if err != nil {
		log.Println(err)
		return err
	}

	response.Allowed = true
	response.Result = &metav1.Status{
		Message: "Allowed",
	}
	response.Patch = patchBytes
	return nil
}

func validate(review *admission.AdmissionReview, response *admission.AdmissionResponse) error {
	response.Allowed = true
	response.Result = &metav1.Status{
		Message: "Allowed",
	}
	return nil
}
