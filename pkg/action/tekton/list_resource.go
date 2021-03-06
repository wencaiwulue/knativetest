package tekton

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knativetest/pkg/action"
	"knativetest/pkg/action/knative"
	"knativetest/pkg/client"
)

type ListResourceAction struct {
	action.Action
	Namespace string
	Name      string
}

func (a *ListResourceAction) Process(ctx context.Context) interface{} {
	podList, _ := client.GetClient().K8sClient.CoreV1().Pods(a.Namespace).List(ctx, metav1.ListOptions{})
	tektonTasks, _ := client.GetClient().TektonClient.TektonV1beta1().Tasks(a.Namespace).List(ctx, metav1.ListOptions{})
	servingList, _ := client.GetClient().ServingClient.ServingV1().Services(a.Namespace).List(ctx, metav1.ListOptions{})
	result := map[string]string{}
	result["Namespace num"] = knative.ToJsonString(podList)
	result["tekton num"] = knative.ToJsonString(tektonTasks)
	result["serving num"] = knative.ToJsonString(servingList)
	return result
}
