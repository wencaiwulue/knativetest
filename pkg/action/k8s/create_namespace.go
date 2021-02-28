package k8s

import (
	"context"
	k8sV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"test/pkg/action"
	"test/pkg/client"
)

type CreateNamespaceAction struct {
	action.Action
	Name string
}

func (c *CreateNamespaceAction) Process(ctx context.Context) interface{} {
	var result, _ = client.GetClient().K8sClient.CoreV1().Namespaces().Create(ctx, &k8sV1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: c.Name},
	}, metav1.CreateOptions{})
	return result
}
