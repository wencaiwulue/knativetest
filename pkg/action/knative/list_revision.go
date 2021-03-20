package knative

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knativetest/pkg/action"
	"knativetest/pkg/client"
	"log"
)

type ListRevisionAction struct {
	action.Action
	Namespace string
	Name      string
}

func (c *ListRevisionAction) Process(ctx context.Context) interface{} {
	option := metav1.ListOptions{}
	if c.Name != "" {
		option.LabelSelector = fmt.Sprintf("name=%s", c.Name)
	}
	if c.Name != "" {
		option.LabelSelector = fmt.Sprintf("name=%s", c.Name)
	}
	var service, err = client.GetClient().ServingClient.ServingV1().Revisions(c.Namespace).List(ctx, option)
	if err != nil || service == nil {
		log.Printf("create service error info: %v", err)
	}
	result := Result{List: make([]Item, len(service.Items))}
	for i, e := range service.Items {
		result.List[i] = Item{Name: e.Name}
	}
	return result
}
