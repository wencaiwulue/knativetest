package knative

import (
	"context"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/serving/pkg/apis/autoscaling/v1alpha1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	"strconv"
	"test/knative/client"
	"test/knative/controllers"
)

type CreateServiceAction struct {
	controllers.Action
	Namespace string
	Name      string
}

func (c *CreateServiceAction) Process(ctx context.Context) interface{} {
	var containerConcurrency, _ = strconv.ParseInt("4", 10, 64)
	var trafficMin, _ = strconv.ParseInt("0", 10, 64)
	var trafficMax, _ = strconv.ParseInt("100", 10, 64)
	var svc = &servingv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "service",
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
								{Image: "test:latest", ImagePullPolicy: corev1.PullIfNotPresent, Ports: []corev1.ContainerPort{{
									ContainerPort: 80,
								}}},
								{Image: "test:latest", ImagePullPolicy: corev1.PullIfNotPresent, Ports: []corev1.ContainerPort{{
									ContainerPort: 81,
								}}},
							},
						},
						ContainerConcurrency: &containerConcurrency,
					},
				},
			},
			RouteSpec: servingv1.RouteSpec{
				Traffic: []servingv1.TrafficTarget{{Percent: &trafficMin}, {Percent: &trafficMax}},
			},
		},
	}

	var serviceRevision = servingv1.Revision{
		TypeMeta: metav1.TypeMeta{
			Kind:       "revisions.serving.knative.dev",
			APIVersion: "serving.knative.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        c.Name + "_" + "revision",
			Namespace:   c.Namespace,
			Labels:      map[string]string{"service": c.Name},
			Annotations: map[string]string{"k1": "v1", "k2": "v2"},
		},
		Spec: servingv1.RevisionSpec{
			PodSpec: corev1.PodSpec{
				Containers: []corev1.Container{
					{Image: "test:latest", ImagePullPolicy: corev1.PullIfNotPresent, Ports: []corev1.ContainerPort{
						{ContainerPort: 80},
					}},
				},
			},
			ContainerConcurrency: &containerConcurrency,
		},
	}

	var auto = v1alpha1.PodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling.internal.knative.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name + "_autoscaler",
			Namespace: c.Namespace,
		},
		Spec: v1alpha1.PodAutoscalerSpec{
			ContainerConcurrency: 3,
			ScaleTargetRef: corev1.ObjectReference{
				Kind:      "Pod",
				Namespace: c.Namespace,
				Name:      c.Name,
			},
		},
	}
	var service, _ = client.ServingClient.ServingV1().Services(c.Namespace).Create(svc)
	var revision, _ = client.ServingClient.ServingV1().Revisions(c.Namespace).Create(&serviceRevision)
	var result, _ = client.ServingClient.AutoscalingV1alpha1().PodAutoscalers(c.Namespace).Create(&auto)
	var b, _ = json.Marshal(service)
	println("service result: " + string(b))
	b, _ = json.Marshal(revision)
	println("revision result: " + string(b))
	b, _ = json.Marshal(result)
	println("auto scale result: " + string(b))

	return map[string]string{"service result": string(b), "revision result": string(b), "auto scale result": string(b)}
}
