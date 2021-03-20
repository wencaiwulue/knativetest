package knative

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/networking/pkg/apis/networking"
	"knative.dev/serving/pkg/apis/autoscaling/v1alpha1"
	"knativetest/pkg/action"
	"knativetest/pkg/client"
	"log"
)

type CreateKPAAction struct {
	action.Action
	Namespace string
	Name      string
}

func (c *CreateKPAAction) Process(ctx context.Context) interface{} {
	// this can create by annotation automatically, of course you can create it by handle
	var kpa = v1alpha1.PodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PodAutoscaler",
			APIVersion: "autoscaling.internal.knative.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name + "-" + "hpa",
			Namespace: c.Namespace,
		},
		Spec: v1alpha1.PodAutoscalerSpec{
			ContainerConcurrency: 3,
			ScaleTargetRef: corev1.ObjectReference{
				Kind:       "Revision",
				Name:       c.Name + "-" + "b",
				APIVersion: "apps/v1",
			},
			ProtocolType: networking.ProtocolHTTP1,
		},
	}
	option := metav1.CreateOptions{}
	result, err := client.GetClient().ServingClient.AutoscalingV1alpha1().PodAutoscalers(c.Namespace).Create(ctx, &kpa, option)
	if err != nil {
		log.Printf("create service hpa scale info: %v\n", err.Error())
		println("hpa scale result: " + ToJsonString(result))
	}
	return result
}
