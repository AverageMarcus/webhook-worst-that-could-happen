package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"worst-webhook/pkg/handler"
	"worst-webhook/pkg/patcher"
	"worst-webhook/pkg/server"

	admission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

const (
	name        = "different-content-type"
	description = "In this scenario we're going to return a valid `Content-Type` header but one that doesn't match the content being returned"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func main() {
	log.Printf("Scenario: %s", name)
	log.Println(description)

	server.StartServer(
		handler.Handle(mutate, xmlContentTypeHook),
		validate,
	)
}

func xmlContentTypeHook(w http.ResponseWriter, r *http.Request) {
	log.Printf("Setting content-type to `application/xml`")
	w.Header().Set("Content-Type", "application/xml")
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

func validate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Setting content-type to `application/json` but returned content will be xml")
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

	w.Write([]byte(fmt.Sprintf(`<apiVersion>admission.k8s.io/v1</apiVersion>
<kind>AdmissionReview</kind>
<response>
	<allowed>true</allowed>
	<status>
		<message>Allowed</message>
		<metadata/>
	</status>
	<uid>%s</uid>
</response>`, review.Request.UID)))
}
