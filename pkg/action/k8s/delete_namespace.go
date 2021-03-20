package k8s

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knativetest/pkg/action"
	"knativetest/pkg/client"
)

type DeleteNamespaceAction struct {
	action.Action
	Name string
}

func (c *DeleteNamespaceAction) Process(ctx context.Context) interface{} {
	result := client.GetClient().K8sClient.CoreV1().Namespaces().Delete(ctx, c.Name, metav1.DeleteOptions{})
	return result
}
