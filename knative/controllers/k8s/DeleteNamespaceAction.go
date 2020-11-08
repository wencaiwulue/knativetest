package k8s

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"test/knative/client"
	"test/knative/controllers"
)

type DeleteNamespaceAction struct {
	controllers.Action
	Namespace string
}

func (c *DeleteNamespaceAction) Process(ctx context.Context) interface{} {
	_ = client.K8sClient.CoreV1().Namespaces().Delete(c.Namespace, &metav1.DeleteOptions{})
	return true
}
