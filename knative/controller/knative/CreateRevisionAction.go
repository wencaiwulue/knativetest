package knative

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	"log"
	"strconv"
	"test/knative/client"
	"test/knative/controller"
)

type CreateRevisionAction struct {
	controller.Action
	Namespace string
	Name      string
}

// use service to control traffic, if create maxPolicySelect new revision, should update service traffic manager
// todo traffic manager RevisionName should equals to revision name
func (c *CreateRevisionAction) Process(ctx context.Context) interface{} {
	var containerConcurrency, _ = strconv.ParseInt("4", 10, 64)
	var containerTimeout, _ = strconv.ParseInt("4", 10, 64)
	resourceMax := map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse("200m"),
		corev1.ResourceMemory: resource.MustParse("256Mi"),
		//corev1.ResourceEphemeralStorage: resource.MustParse("256Mi"),
	}
	resourceMin := map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse("100m"),
		corev1.ResourceMemory: resource.MustParse("128Mi"),
		//corev1.ResourceEphemeralStorage: resource.MustParse("256Mi"),
	}

	var serviceRevision1 = servingv1.Revision{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Revision",
			APIVersion: "serving.knative.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name + "-" + "1", //todo same, use for loop to create revision
			Namespace: c.Namespace,
			Labels:    map[string]string{"service": c.Name},
			Annotations: map[string]string{
				// reference https://knative.dev/docs/serving/autoscaling/autoscaling-metrics/
				//  For per-revision concurrency, you must configure both `autoscaling.knative.dev/metric`
				//  and `autoscaling.knative.dev/target` for maxPolicySelect [soft limit](#soft-limit),
				//  or `containerConcurrency` for maxPolicySelect [hard limit](#hard-limit).
				"autoscaling.knative.dev/class":  "kpa.autoscaling.knative.dev",
				"autoscaling.knative.dev/metric": "concurrency",
				// concurrency > 100 for 60s, it will scale number of revision
				"autoscaling.knative.dev/target": "100",
			},
		},
		Spec: servingv1.RevisionSpec{
			PodSpec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Image:           IMAGE,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Ports:           []corev1.ContainerPort{{ContainerPort: 80}},
					Resources: corev1.ResourceRequirements{
						Limits:   resourceMax,
						Requests: resourceMin,
					},
				}},
				Overhead: resourceMax,
			},
			ContainerConcurrency: &containerConcurrency,
			TimeoutSeconds:       &containerTimeout,
		},
	}

	var serviceRevision2 = servingv1.Revision{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Revision",
			APIVersion: "serving.knative.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name + "-" + "2", //todo same, use for loop to create revision
			Namespace: c.Namespace,
			Labels:    map[string]string{"service": c.Name},
			Annotations: map[string]string{
				// reference https://knative.dev/docs/serving/autoscaling/autoscaling-metrics/
				// documents said hpa support cpu and concurrency, but my test result is hpa not support concurrency
				// kpa not support cpu, on knative 0.18 and kubernetes 1.19.3, so needs to change metrics class
				"autoscaling.knative.dev/metric": "cpu",
				"autoscaling.knative.dev/class":  "hpa.autoscaling.knative.dev",
			},
		},
		Spec: servingv1.RevisionSpec{
			PodSpec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Image:           IMAGE,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Ports:           []corev1.ContainerPort{{ContainerPort: 80}},
					Resources: corev1.ResourceRequirements{
						Limits:   resourceMax,
						Requests: resourceMin,
					},
				}},
				Overhead: resourceMax,
			},
			ContainerConcurrency: &containerConcurrency,
			TimeoutSeconds:       &containerTimeout,
		},
	}
	option := metav1.CreateOptions{}
	result, err := client.GetClient().ServingClient.ServingV1().Revisions(c.Namespace).Create(ctx, &serviceRevision1, option)
	if err != nil {
		log.Printf("create revision1 error info: %v\n", err.Error())
		log.Printf("create revision1 result info: %v\n", ToJsonString(result))
	}
	result, err = client.GetClient().ServingClient.ServingV1().Revisions(c.Namespace).Create(ctx, &serviceRevision2, option)
	if err != nil {
		log.Printf("create revision2 error info: %v\n", err.Error())
		log.Printf("create revision2 result info: %v\n", ToJsonString(result))
	}
	return result
}
