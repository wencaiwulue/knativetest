package tekton

import (
	"context"
	"fmt"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knativetest/pkg/action"
	"knativetest/pkg/client"
)

type CreateKanikoClusterTask struct {
	action.Action
	Namespace string
	Name      string
}

func (a *CreateKanikoClusterTask) Process(ctx context.Context) interface{} {
	kaniko := v1beta1.ClusterTask{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterTask",
			APIVersion: "tekton.dev/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kaniko",
			Namespace: a.Namespace,
		},
		Spec: v1beta1.TaskSpec{
			Resources: &v1beta1.TaskResources{
				Inputs: []v1beta1.TaskResource{{ResourceDeclaration: v1beta1.ResourceDeclaration{
					Name: "docker-source",
					Type: "git",
				}}},
				Outputs: []v1beta1.TaskResource{{ResourceDeclaration: v1beta1.ResourceDeclaration{
					Name: "builtImage",
					Type: "image",
				}}},
			},
			Params: []v1beta1.ParamSpec{{
				Name:        "",
				Type:        "",
				Description: "",
				Default:     nil,
			}, v1beta1.ParamSpec{
				Name:        "",
				Type:        "",
				Description: "",
				Default:     nil,
			}, v1beta1.ParamSpec{
				Name:        "",
				Type:        "",
				Description: "",
				Default:     nil,
			}},
			Description: "this is a test",
			Steps: []v1beta1.Step{{Container: corev1.Container{
				Name:    "kaniko",
				Image:   "gcr.io/kaniko-project/executor:latest",
				Command: []string{"/kaniko/executor"},
				//Args:                     []string{"dockerfile"},
				Ports:                    nil,
				StartupProbe:             nil,
				TerminationMessagePolicy: "OnFailure",
				ImagePullPolicy:          "IfNotPresent",
			}}},
		},
	}
	r, err := client.GetClient().TektonClient.TektonV1beta1().ClusterTasks().Create(ctx, &kaniko, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("error info: %v\n", err)
	}
	fmt.Printf("result: %v\n", r)
	return r
}
