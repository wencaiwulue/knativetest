package knative

import (
	"context"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/serving/pkg/apis/autoscaling/v1alpha1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	"log"
	"strconv"
	"test/knative/client"
	"test/knative/controller"
)

type CreateServiceAction struct {
	controller.Action
	Namespace string
	Name      string
}

func (c *CreateServiceAction) Process(ctx context.Context) interface{} {
	var containerConcurrency, _ = strconv.ParseInt("4", 10, 64)
	var trafficMin, _ = strconv.ParseInt("0", 10, 64)
	var trafficMax, _ = strconv.ParseInt("100", 10, 64)
	resourceMin := map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse("100m"),
		corev1.ResourceMemory: resource.MustParse("128Mi"),
		//corev1.ResourceStorage: resource.MustParse("100Mi"),
	}

	resourceMax := map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse("200m"),
		corev1.ResourceMemory: resource.MustParse("256Mi"),
		//corev1.ResourceStorage: resource.MustParse("200Mi"),
	}
	const image = "test:latest"
	var svc = &servingv1.Service{
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
							Containers: []corev1.Container{
								{Image: image, ImagePullPolicy: corev1.PullIfNotPresent, Ports: []corev1.ContainerPort{{
									ContainerPort: 80,
								}}, Resources: corev1.ResourceRequirements{
									Limits:   resourceMax,
									Requests: resourceMin,
								}},
							},
						},
						ContainerConcurrency: &containerConcurrency,
					},
				},
			},
			RouteSpec: servingv1.RouteSpec{
				Traffic: []servingv1.TrafficTarget{{
					RevisionName: c.Name + "-" + "a",
					Percent:      &trafficMin,
				}, {
					RevisionName: c.Name + "-" + "b",
					Percent:      &trafficMax,
				}},
			},
		},
	}

	// todo use service to control traffic, if create a new revision, should update service traffic manager
	// todo traffic manager RevisionName should equals to revisoin name

	var serviceRevisionA = servingv1.Revision{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Revision",
			APIVersion: "serving.knative.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        c.Name + "-" + "a", //todo same, use for loop to create revision
			Namespace:   c.Namespace,
			Labels:      map[string]string{"service": c.Name},
			Annotations: map[string]string{"k1": "v1", "k2": "v2"},
		},
		Spec: servingv1.RevisionSpec{
			PodSpec: corev1.PodSpec{
				Containers: []corev1.Container{{Image: image, ImagePullPolicy: corev1.PullIfNotPresent, Ports: []corev1.ContainerPort{
					{ContainerPort: 80},
				}},
				},
				Overhead: resourceMax,
			},
			ContainerConcurrency: &containerConcurrency,
		},
	}

	var serviceRevisionB = servingv1.Revision{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Revision",
			APIVersion: "serving.knative.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        c.Name + "-" + "b", //todo same, use for loop to create revision
			Namespace:   c.Namespace,
			Labels:      map[string]string{"service": c.Name},
			Annotations: map[string]string{"k1": "v1", "k2": "v2"},
		},
		Spec: servingv1.RevisionSpec{
			PodSpec: corev1.PodSpec{
				Containers: []corev1.Container{{Image: image, ImagePullPolicy: corev1.PullIfNotPresent, Ports: []corev1.ContainerPort{
					{ContainerPort: 80},
				}},
				},
				Overhead: resourceMax,
			},
			ContainerConcurrency: &containerConcurrency,
		},
	}

	var hpa = v1alpha1.PodAutoscaler{
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
		},
	}
	option := metav1.CreateOptions{}
	service, err := client.GetClient().ServingClient.ServingV1().Services(c.Namespace).Create(ctx, svc, option)
	if err != nil {
		log.Printf("create service error info: %v", err.Error())
	}
	revisionA, err1 := client.GetClient().ServingClient.ServingV1().Revisions(c.Namespace).Create(ctx, &serviceRevisionA, option)
	if err1 != nil {
		log.Printf("create revision error info: %v", err1.Error())
	}
	revisionB, err1 := client.GetClient().ServingClient.ServingV1().Revisions(c.Namespace).Create(ctx, &serviceRevisionB, option)
	if err1 != nil {
		log.Printf("create revision error info: %v", err1.Error())
	}
	result, err2 := client.GetClient().ServingClient.AutoscalingV1alpha1().PodAutoscalers(c.Namespace).Create(ctx, &hpa, option)
	if err2 != nil {
		log.Printf("create service hpa scale info: %v", err2.Error())
	}
	var b, _ = json.Marshal(service)
	println("service result: " + string(b))
	b, _ = json.Marshal(revisionA)
	println("revisionA result: " + string(b))
	b, _ = json.Marshal(revisionB)
	println("revisionB result: " + string(b))
	b, _ = json.Marshal(result)
	println("hpa scale result: " + string(b))

	return map[string]string{"service result": string(b), "revision result": string(b), "hpa scale result": string(b)}
}
