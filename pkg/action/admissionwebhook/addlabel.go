package admissionwebhook

import (
	"encoding/json"
	v1 "k8s.io/api/admission/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

const (
	addFirstLabelPatch string = `[
         { "op": "add", "path": "/metadata/labels", "value": {"added-label": "yes"}}
     ]`
	addAdditionalLabelPatch string = `[
         { "op": "add", "path": "/metadata/labels/added-label", "value": "yes" }
     ]`
)

// Add a label {"added-label": "yes"} to the object
func AddLabel(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.V(2).Info("calling add-label")
	obj := struct {
		metav1.ObjectMeta
		Data map[string]string
	}{}
	raw := ar.Request.Object.Raw
	err := json.Unmarshal(raw, &obj)
	if err != nil {
		klog.Error(err)
		return ToAdmissionResponse(err)
	}

	reviewResponse := v1.AdmissionResponse{}
	reviewResponse.Allowed = true
	if len(obj.ObjectMeta.Labels) == 0 {
		reviewResponse.Patch = []byte(addFirstLabelPatch)
	} else {
		reviewResponse.Patch = []byte(addAdditionalLabelPatch)
	}
	pt := v1.PatchTypeJSONPatch
	reviewResponse.PatchType = &pt
	return &reviewResponse
}

// ToAdmissionResponse is a helper function to create an AdmissionResponse
// with an embedded error
func ToAdmissionResponse(err error) *v1.AdmissionResponse {
	return &v1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}
