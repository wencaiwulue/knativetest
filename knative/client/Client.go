package client

import (
	tektoncd "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	serving "knative.dev/serving/pkg/client/clientset/versioned"
	"sync"
)

var onceInit sync.Once
var K8sClient *kubernetes.Clientset
var ServingClient *serving.Clientset
var TektonClient *tektoncd.Clientset

func init() {
	onceInit.Do(func() {
		initClient(true)
	})
}

func initClient(inCluster bool) {
	var config *rest.Config
	if inCluster {
		config, _ = rest.InClusterConfig()
	} else {
		config, _ = clientcmd.BuildConfigFromFlags("", "~/.kube/config")
	}
	var err error
	K8sClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		println("error1: " + err.Error())
	}
	ServingClient, err = serving.NewForConfig(config)
	if err != nil {
		println("error2: " + err.Error())
	}
	TektonClient, err = tektoncd.NewForConfig(config)
	if err != nil {
		println("error3: " + err.Error())
	}
}
