package tekton

import (
	"context"
	beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knativetest/pkg/action"
	"knativetest/pkg/client"
)

type ToktonAction struct {
	action.Action
	Namespace string
}

func (a *ToktonAction) Process(ctx context.Context) interface{} {
	var _, _ = client.GetClient().TektonClient.TektonV1beta1().Tasks(a.Namespace).List(ctx, metav1.ListOptions{})
	var _, _ = client.GetClient().TektonClient.TektonV1beta1().ClusterTasks().List(ctx, metav1.ListOptions{})
	var _, _ = client.GetClient().TektonClient.TektonV1beta1().PipelineRuns(a.Namespace).List(ctx, metav1.ListOptions{})
	var _, _ = client.GetClient().TektonClient.TektonV1beta1().PipelineRuns(a.Namespace).List(ctx, metav1.ListOptions{})
	var _, _ = client.GetClient().TektonClient.TektonV1beta1().Tasks(a.Namespace).Create(ctx, &beta1.Task{}, metav1.CreateOptions{})
	return nil
}
