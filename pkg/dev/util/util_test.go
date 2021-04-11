package util

import (
	v1 "k8s.io/api/apps/v1"
	"testing"
)

func TestWait(t *testing.T) {
	client, _ := InitClient("/Users/naison/codingtest")
	WaitToBeStatus(client.ClientSet, "test", "deployments", "kubedev=debug", func(i interface{}) bool {
		return i.(*v1.Deployment).Status.ReadyReplicas == 1
	})
}
