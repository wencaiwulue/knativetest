package admissionwebhook

import (
	"fmt"
	v1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"

	"k8s.io/klog"
)

const (
	podsAddContainerPatch string = `[
		 {"op":"add","path":"/spec/containers","value":[{"image":"webhook-added-image","name":"webhook-added-init-container","resources":{}}]}
	]`

	podsRemoveContainerPatch string = `[
		 {"op":"replace","path":"/spec/template/spec/containers/0/image","value":"empty-container:latest"}
	]`
)

// only allow pods to pull images from specific registry.
func AdmitPods(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.V(2).Info("admitting pods")
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if ar.Request.Resource != podResource {
		err := fmt.Errorf("expect resource to be %s", podResource)
		klog.Error(err)
		return ToAdmissionResponse(err)
	}

	raw := ar.Request.Object.Raw
	pod := corev1.Pod{}
	deserializer := Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		klog.Error(err)
		return ToAdmissionResponse(err)
	}
	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true

	var msg string
	if v, ok := pod.Labels["webhook-e2e-test"]; ok {
		if v == "webhook-disallow" {
			reviewResponse.Allowed = false
			msg = msg + "the pod contains unwanted label; "
		}
		if v == "wait-forever" {
			reviewResponse.Allowed = false
			msg = msg + "the pod response should not be sent; "
			<-make(chan int) // Sleep forever - no one sends to this channel
		}
	}
	for _, container := range pod.Spec.Containers {
		if strings.Contains(container.Name, "webhook-disallow") {
			reviewResponse.Allowed = false
			msg = msg + "the pod contains unwanted container name; "
		}
	}
	if !reviewResponse.Allowed {
		reviewResponse.Result = &metav1.Status{Message: strings.TrimSpace(msg)}
	}
	return &reviewResponse
}

func MutatePods(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.V(2).Info("mutating deployments")
	podResource := metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	if ar.Request.Resource != podResource {
		klog.Errorf("expect resource to be %s", podResource)
		return nil
	}

	raw := ar.Request.Object.Raw
	pod := appsv1.Deployment{}
	deserializer := Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		klog.Error(err)
		return ToAdmissionResponse(err)
	}
	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true
	//if pod.Name == "webhook-to-be-mutated" {
	reviewResponse.Patch = []byte(podsRemoveContainerPatch)
	pt := v1.PatchTypeJSONPatch
	reviewResponse.PatchType = &pt
	//}
	return &reviewResponse
}
