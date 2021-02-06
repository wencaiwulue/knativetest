package knative

import (
	"context"
	"encoding/json"
	"k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"test/knative/client"
	"test/knative/controller"
)

type CreateHPAAction struct {
	controller.Action
	Namespace string
	Name      string
}

func (c *CreateHPAAction) Process(ctx context.Context) interface{} {
	minPod := int32(1)
	maxPolicySelect := v2beta2.MaxPolicySelect
	second15 := int32(15)
	second30 := int32(30)
	hpa := &v2beta2.HorizontalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling/v2beta2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name,
			Namespace: c.Namespace,
		},
		Spec: v2beta2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: v2beta2.CrossVersionObjectReference{
				Kind:       "Revision",
				APIVersion: "serving.knative.dev/v1",
				Name:       "rev1",
			},
			MinReplicas: &minPod,
			MaxReplicas: 100,
			Metrics: []v2beta2.MetricSpec{{
				Type: v2beta2.PodsMetricSourceType,
				Pods: &v2beta2.PodsMetricSource{
					Metric: v2beta2.MetricIdentifier{
						Name: "cpu",
					},
					Target: v2beta2.MetricTarget{
						Type:         v2beta2.AverageValueMetricType,
						AverageValue: resource.NewMilliQuantity(100, resource.DecimalSI),
					},
				},
			}, {
				Type: v2beta2.PodsMetricSourceType,
				Pods: &v2beta2.PodsMetricSource{
					Metric: v2beta2.MetricIdentifier{
						Name: "memory",
					},
					Target: v2beta2.MetricTarget{
						Type:         v2beta2.AverageValueMetricType,
						AverageValue: resource.NewMilliQuantity(100, resource.DecimalSI),
					},
				},
			}},
			Behavior: &v2beta2.HorizontalPodAutoscalerBehavior{
				ScaleUp: &v2beta2.HPAScalingRules{
					StabilizationWindowSeconds: &second15,
					SelectPolicy:               &maxPolicySelect,
					Policies: []v2beta2.HPAScalingPolicy{{
						Type:          v2beta2.PodsScalingPolicy,
						Value:         12,
						PeriodSeconds: 60,
					}},
				},
				ScaleDown: &v2beta2.HPAScalingRules{
					StabilizationWindowSeconds: &second30,
					SelectPolicy:               &maxPolicySelect,
					Policies: []v2beta2.HPAScalingPolicy{{
						Type:          v2beta2.PodsScalingPolicy,
						Value:         12,
						PeriodSeconds: 60,
					}},
				},
			},
		},
	}

	option := metav1.CreateOptions{}
	result, err := client.GetClient().K8sClient.AutoscalingV2beta2().HorizontalPodAutoscalers(c.Namespace).Create(ctx, hpa, option)
	if err != nil {
		println("hpa scale result: " + ToJsonString(result))
	}
	return result
}

func ToJsonString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}
