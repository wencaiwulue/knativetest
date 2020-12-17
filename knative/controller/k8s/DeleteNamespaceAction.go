package k8s

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"test/knative/client"
	"test/knative/controller"
)

type DeleteNamespaceAction struct {
	controller.Action
	Namespace string
}

func (c *DeleteNamespaceAction) Process(ctx context.Context) interface{} {
	_ = client.GetClient().K8sClient.CoreV1().Namespaces().Delete(ctx, c.Namespace, metav1.DeleteOptions{})
	return true
}
