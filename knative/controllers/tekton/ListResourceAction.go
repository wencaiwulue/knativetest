package tekton

import (
	"context"
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"test/knative/client"
	"test/knative/controllers"
)

type ListResourceAction struct {
	controllers.Action
	Namespace string
	Name      string
}

func (a *ListResourceAction) Process(ctx context.Context) interface{} {
	var podList, _ = client.K8sClient.CoreV1().Pods(a.Namespace).List(metav1.ListOptions{})
	var tektonTasks, _ = client.TektonClient.TektonV1beta1().Tasks(a.Namespace).List(metav1.ListOptions{})
	var servingList, _ = client.ServingClient.ServingV1().Services(a.Namespace).List(metav1.ListOptions{})
	var result = map[string]string{}
	var bytes, _ = json.Marshal(podList)
	result["Namespace num"] = string(bytes)
	bytes, _ = json.Marshal(tektonTasks)
	result["tekton num"] = string(bytes)
	bytes, _ = json.Marshal(servingList)
	result["serving num"] = string(bytes)
	return result
}
