package util

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClientSet struct {
	KubeConfig    string
	KubeContext   string
	RestConfig    *restclient.Config
	ClientConfig  clientcmd.ClientConfig
	ClientSet     *kubernetes.Clientset
	DynamicClient dynamic.Interface
}

func InitClient(kubeConfig string) (*ClientSet, error) {
	client := &ClientSet{
		KubeConfig: kubeConfig,
	}

	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig}
	client.ClientConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})

	var err error
	client.RestConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	client.RestConfig.GroupVersion = &schema.GroupVersion{
		Group:   "apps",
		Version: "v1",
	}
	//DirectClientConfig.getContext
	client.RestConfig.NegotiatedSerializer = scheme.Codecs
	if err != nil {
		return nil, err
	}

	if client.ClientSet, err = kubernetes.NewForConfig(client.RestConfig); err != nil {
		return nil, err
	}

	if client.DynamicClient, err = dynamic.NewForConfig(client.RestConfig); err != nil {
		return nil, err
	}

	return client, nil
}
