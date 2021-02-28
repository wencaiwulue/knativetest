package knative

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	"log"
	"net/http"
	"strconv"
	"test/pkg/action"
	"test/pkg/client"
)

const IMAGE = "test:latest"

type CreateServiceAction struct {
	action.Action
	Namespace string
	Name      string
}

func (c *CreateServiceAction) Process(ctx context.Context) interface{} {
	m := ctx.Value("http").(map[string]interface{})
	reqq := m["http.request"].(*http.Request)
	body := m["http.request.body"].(string)
	fmt.Printf("body: %v\n", body)
	fmt.Printf("header: %v\n", reqq.Header)
	var containerConcurrency, _ = strconv.ParseInt("4", 10, 64)
	var trafficMin, _ = strconv.ParseInt("0", 10, 64)
	var trafficMax, _ = strconv.ParseInt("100", 10, 64)
	resourceMin := map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:              resource.MustParse("100m"),
		corev1.ResourceMemory:           resource.MustParse("128Mi"),
		corev1.ResourceEphemeralStorage: resource.MustParse("128Mi"),
	}

	resourceMax := map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:              resource.MustParse("200m"),
		corev1.ResourceMemory:           resource.MustParse("256Mi"),
		corev1.ResourceEphemeralStorage: resource.MustParse("256Mi"),
	}
	var kservice = &servingv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "serving.knative.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        c.Name,
			Namespace:   c.Namespace,
			Labels:      map[string]string{"service": c.Name},
			Annotations: map[string]string{"k1": "v1", "k2": "v2"},
		},
		Spec: servingv1.ServiceSpec{
			ConfigurationSpec: servingv1.ConfigurationSpec{
				Template: servingv1.RevisionTemplateSpec{
					Spec: servingv1.RevisionSpec{
						PodSpec: corev1.PodSpec{
							Containers: []corev1.Container{{
								Image:           IMAGE,
								ImagePullPolicy: corev1.PullIfNotPresent,
								Ports:           []corev1.ContainerPort{{ContainerPort: 80}},
								Resources: corev1.ResourceRequirements{
									Limits:   resourceMax,
									Requests: resourceMin,
								}},
							},
						},
						// the system decides the target concurrency for the autoscaler.
						//  For per-revision concurrency, you must configure both `autoscaling.knative.dev/metric`
						//  and `autoscaling.knative.dev/target` for maxPolicySelect [soft limit](#soft-limit),
						//  or `containerConcurrency` for maxPolicySelect [hard limit](#hard-limit).
						ContainerConcurrency: &containerConcurrency,
					},
				},
			},
			RouteSpec: servingv1.RouteSpec{
				Traffic: []servingv1.TrafficTarget{{
					// need enable tagHeaderBasedRouting, reference https://knative.dev/docs/serving/samples/tag-header-based-routing/
					Tag:          "rev1",
					RevisionName: c.Name + "-1", //attention keep the same with revision name
					Percent:      &trafficMin,
				}, {
					RevisionName: c.Name + "-2", //attention keep the same with revision name
					Percent:      &trafficMax,
				}},
			},
		},
	}

	option := metav1.CreateOptions{}
	result, err := client.GetClient().ServingClient.ServingV1().Services(c.Namespace).Create(ctx, kservice, option)
	if err != nil {
		log.Printf("create service error info: %v\n", err.Error())
		println("service result: " + ToJsonString(result))
	}

	return result
}
