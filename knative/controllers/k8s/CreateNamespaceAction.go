package k8s

import (
	"context"
	k8sV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"test/knative/client"
	"test/knative/controllers"
)

type CreateNamespaceAction struct {
	controllers.Action
	Namespace string
}

func (c *CreateNamespaceAction) Process(ctx context.Context) interface{} {
	var v1Namespace, _ = client.K8sClient.CoreV1().Namespaces().Create(&k8sV1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: c.Namespace},
	})
	return v1Namespace
}
