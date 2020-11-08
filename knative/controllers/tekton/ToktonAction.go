package tekton

import (
	"context"
	beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"test/knative/client"
	"test/knative/controllers"
)

type ToktonAction struct {
	controllers.Action
	Namespace string
}

func (a *ToktonAction) Process(ctx context.Context) interface{} {
	var _, _ = client.TektonClient.TektonV1beta1().Tasks(a.Namespace).List(metav1.ListOptions{})
	var _, _ = client.TektonClient.TektonV1beta1().ClusterTasks().List(metav1.ListOptions{})
	var _, _ = client.TektonClient.TektonV1beta1().PipelineRuns(a.Namespace).List(metav1.ListOptions{})
	var _, _ = client.TektonClient.TektonV1beta1().PipelineRuns(a.Namespace).List(metav1.ListOptions{})
	var _, _ = client.TektonClient.TektonV1beta1().Tasks(a.Namespace).Create(&beta1.Task{})
	return nil
}
