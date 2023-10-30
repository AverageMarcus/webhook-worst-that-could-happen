package main

import (
	"encoding/json"
	"log"

	"worst-webhook/pkg/handler"
	"worst-webhook/pkg/patcher"
	"worst-webhook/pkg/server"

	admission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	name        = "label-overwrite"
	description = "This scenaio configures a mutating webhook that updates a pods labels with a new entry, but accidentally replaces the existing."
)

func main() {
	log.Printf("Scenario: %s", name)
	log.Println(description)

	server.StartServer(
		handler.Handle(mutate),
		handler.Handle(validate),
	)
}

func mutate(review *admission.AdmissionReview, response *admission.AdmissionResponse) error {
	pod := &corev1.Pod{}
	if err := json.Unmarshal(review.Request.Object.Raw, pod); err != nil {
		log.Println(err)
		return err
	}

	labels := map[string]string{}
	labels["webhook"] = "label-overwrite"

	patchBytes, err := patcher.New().
		Add(
			patcher.JSONPatchEntry{
				OP:    "add",
				Path:  "/metadata/labels",
				Value: patcher.ToValueMap(labels),
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
