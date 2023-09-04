package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"worst-webhook/pkg/handler"
	"worst-webhook/pkg/patcher"
	"worst-webhook/pkg/server"

	admission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	name        = "reinvocation"
	description = "This scenario has two mutating webhooks with their reinvocationPolicy set to IfNeeded and both will mutate the object, triggering reinvocations of each other."
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

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

	annotations := pod.ObjectMeta.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations["webhook-updated-at"] = time.Now().String()
	annotations[randString(10)] = randString(10)

	patchBytes, err := patcher.New().
		Add(
			patcher.JSONPatchEntry{
				OP:    "add",
				Path:  "/metadata/annotations",
				Value: patcher.ToValueMap(annotations),
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
