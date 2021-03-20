package tekton

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knativetest/pkg/action"
	"knativetest/pkg/client"
)

type CreateKanikoTaskRun struct {
	action.Action
	Namespace string
	Name      string
}

func (a *CreateKanikoTaskRun) Process(ctx context.Context) interface{} {
	run := v1beta1.TaskRun{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TaskRun",
			APIVersion: "tekton.dev/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "TaskRun-" + uuid.New().String(),
			Namespace:   a.Namespace,
			Labels:      nil,
			Annotations: nil,
		},
		Spec: v1beta1.TaskRunSpec{
			ServiceAccountName: "",
			TaskRef: &v1beta1.TaskRef{
				Name:       "kaniko",
				Kind:       "ClusterTask",
				APIVersion: "tekton.dev/v1beta1",
			},
		},
	}

	r, err := client.GetClient().TektonClient.TektonV1beta1().TaskRuns(a.Namespace).Create(ctx, &run, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("error info: %v\n", err)
	}
	fmt.Printf("result: %v\n", r)
	return r
}
