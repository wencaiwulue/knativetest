package k8s

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"test/knative/client"
	"test/knative/controller"
)

type DeleteNamespaceAction struct {
	controller.Action
	Name string
}

func (c *DeleteNamespaceAction) Process(ctx context.Context) interface{} {
	result := client.GetClient().K8sClient.CoreV1().Namespaces().Delete(ctx, c.Name, metav1.DeleteOptions{})
	return result
}
