package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	admission "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)

	patchType = admission.PatchTypeJSONPatch
)

type AdmissionFunc func(*admission.AdmissionReview, *admission.AdmissionResponse) error

func Handle(admissionFn AdmissionFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			UID: review.Request.UID,
		}

		if err := admissionFn(review, response); err != nil {
			http.Error(w, fmt.Sprintf("could not process request: %v", err), http.StatusBadRequest)
			return
		}

		review.Response = response
		review.Response.PatchType = &patchType
		responseBytes, err := json.Marshal(review)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not write response body: %v", err), http.StatusBadRequest)
			return
		}
		w.Write(responseBytes)
	}
}
