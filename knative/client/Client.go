package client

import (
	"fmt"
	tektoncd "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	serving "knative.dev/serving/pkg/client/clientset/versioned"
	"sync"
)

var client *Client
var lock sync.Mutex

type Client struct {
	K8sClient     *kubernetes.Clientset
	ServingClient *serving.Clientset
	TektonClient  *tektoncd.Clientset
}

func GetClient() *Client {
	if client == nil {
		lock.Lock()
		if client == nil {
			initClient(true)
		}
		lock.Unlock()
	}
	return client
}

func initClient(inCluster bool) {
	var config *rest.Config
	if inCluster {
		config, _ = rest.InClusterConfig()
	} else {
		config, _ = clientcmd.BuildConfigFromFlags("", "~/.kube/config")
	}
	//var err error
	K8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("error create k8s client: %v", err)
	}
	ServingClient, err := serving.NewForConfig(config)
	if err != nil {
		fmt.Printf("error create knative serving client: %v", err)
	}
	TektonClient, err := tektoncd.NewForConfig(config)
	if err != nil {
		fmt.Printf("error create tekton client: %v", err)
	}

	client = &Client{K8sClient: K8sClient, ServingClient: ServingClient, TektonClient: TektonClient}
}
