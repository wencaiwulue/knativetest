package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"knativetest/pkg/action/admissionwebhook"
	"log"
	"net/http"
)

const (
	keyFile  = "/etc/webhook/certs/server.key"
	certFile = "/etc/webhook/certs/server.crt"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/add-label", serveAddLabel)
	mux.HandleFunc("/mutating-pods", serveMutatePods)
	log.Fatal(http.ListenAndServeTLS(":80", certFile, keyFile, mux))
}

// admitFunc is the type we use for all of our validators and mutators
type admitFunc func(v1.AdmissionReview) *v1.AdmissionResponse

// serve handles the http portion of a request prior to handing to an admit
// function
func serve(w http.ResponseWriter, r *http.Request, admit admitFunc) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		klog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	klog.V(2).Info(fmt.Sprintf("handling request: %s", body))

	// The AdmissionReview that was sent to the webhook
	requestedAdmissionReview := v1.AdmissionReview{}

	// The AdmissionReview that will be returned
	responseAdmissionReview := v1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
	}

	deserializer := admissionwebhook.Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &requestedAdmissionReview); err != nil {
		klog.Error(err)
		responseAdmissionReview.Response = admissionwebhook.ToAdmissionResponse(err)
	} else {
		// pass to admitFunc
		responseAdmissionReview.Response = admit(requestedAdmissionReview)
	}

	// Return the same UID
	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID

	klog.V(2).Info(fmt.Sprintf("sending response: %v", responseAdmissionReview.Response))

	respBytes, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		klog.Error(err)
	}
	if _, err := w.Write(respBytes); err != nil {
		klog.Error(err)
	}
}

func serveAddLabel(w http.ResponseWriter, r *http.Request) {
	serve(w, r, admissionwebhook.AddLabel)
}

func serveMutatePods(w http.ResponseWriter, r *http.Request) {
	serve(w, r, admissionwebhook.MutatePods)
}
