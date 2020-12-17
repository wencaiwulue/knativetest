package k8s

import (
	"context"
	k8sV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"test/knative/client"
	"test/knative/controller"
)

type CreateNamespaceAction struct {
	controller.Action
	Namespace string
}

func (c *CreateNamespaceAction) Process(ctx context.Context) interface{} {
	var v1Namespace, _ = client.GetClient().K8sClient.CoreV1().Namespaces().Create(ctx, &k8sV1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: c.Namespace},
	}, metav1.CreateOptions{})
	return v1Namespace
}
