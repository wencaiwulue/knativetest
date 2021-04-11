package util

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var Clients *ClientSet

type ClientSet struct {
	Kubeconfig    string
	RestConfig    *restclient.Config
	ClientConfig  clientcmd.ClientConfig
	ClientSet     *kubernetes.Clientset
	DynamicClient dynamic.Interface
}

func InitClient(kubeconfig string) (*ClientSet, error) {
	client := &ClientSet{
		Kubeconfig: kubeconfig,
	}

	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
	client.ClientConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})

	var err error
	client.RestConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	if client.ClientSet, err = kubernetes.NewForConfig(client.RestConfig); err != nil {
		return nil, err
	}

	if client.DynamicClient, err = dynamic.NewForConfig(client.RestConfig); err != nil {
		return nil, err
	}

	Clients = client
	return client, nil
}
