package knative

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	"knativetest/pkg/action"
	"knativetest/pkg/client"
	"log"
	"net/http"
)

type UpdateServiceTrafficAction struct {
	action.Action
	Namespace     string
	Name          string
	TrafficTarget []servingv1.TrafficTarget
}

func (c *UpdateServiceTrafficAction) Process(ctx context.Context) interface{} {
	m := ctx.Value("http").(map[string]interface{})
	reqq := m["http.request"].(*http.Request)
	body := m["http.request.body"].(string)
	fmt.Printf("body: %v\n", body)
	fmt.Printf("header: %v\n", reqq.Header)

	var kservice = &servingv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "serving.knative.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name,
			Namespace: c.Namespace,
			Labels:    map[string]string{"name": c.Name},
		},
		Spec: servingv1.ServiceSpec{
			RouteSpec: servingv1.RouteSpec{
				Traffic: c.TrafficTarget,
			},
		},
	}

	option := metav1.UpdateOptions{}
	result, err := client.GetClient().ServingClient.ServingV1().Services(c.Namespace).Update(ctx, kservice, option)
	if err != nil {
		log.Printf("update service error info: %v\n", err.Error())
		log.Printf("service result: %v\n", ToJsonString(result))
	}

	return result
}
