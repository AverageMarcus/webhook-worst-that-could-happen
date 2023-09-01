package main

import (
	"context"
	"log"
	"os"

	"worst-webhook/pkg/client"
	"worst-webhook/pkg/handler"
	"worst-webhook/pkg/server"

	admission "k8s.io/api/admission/v1"
	eventsv1 "k8s.io/api/events/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	name        = "fork-bomb"
	description = "This scenario creates Kubernetes resources within the webhook handler logic which in turn trigger the webhook again causing a 'fork bomb' like situation where the about of Events in the cluster keep increasing"
)

var (
	kubeClient *client.Client
	scheme     *runtime.Scheme
)

func main() {
	log.Printf("Scenario: %s", name)
	log.Println(description)

	kubeClient = client.New()
	scheme = runtime.NewScheme()
	eventsv1.AddToScheme(scheme)

	server.StartServer(
		handler.Handle(mutate),
		handler.Handle(validate),
	)
}

func mutate(review *admission.AdmissionReview, response *admission.AdmissionResponse) error {
	go generateEvent(review)

	response.Allowed = true
	response.Result = &metav1.Status{
		Message: "Allowed",
	}
	response.Patch = []byte("")
	return nil
}

func validate(review *admission.AdmissionReview, response *admission.AdmissionResponse) error {
	go generateEvent(review)

	response.Allowed = true
	response.Result = &metav1.Status{
		Message: "Allowed",
	}
	return nil
}

func generateEvent(review *admission.AdmissionReview) error {
	log.Println("Creating event...")
	_, err := kubeClient.EventsV1().Events("default").Create(context.Background(), &eventsv1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "webhook-event-",
			Namespace:    "default",
		},
		EventTime:           metav1.NowMicro(),
		ReportingController: "worst-webhook/fork-bomb",
		ReportingInstance:   os.Getenv("HOSTNAME"),
		Action:              "webhook",
		Reason:              "webhook triggered",
		Note:                "This Event was created in response to another Event being created that triggered the fork-bomb webhook",
		Type:                "Normal",
	}, metav1.CreateOptions{})

	if err != nil {
		log.Printf("Failed to create event: %s", err)
		return err
	}

	log.Println("Created event")
	return nil
}
