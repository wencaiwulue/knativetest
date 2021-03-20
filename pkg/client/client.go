package client

import (
	"fmt"
	tektoncd "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"go.etcd.io/etcd/clientv3"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	serving "knative.dev/serving/pkg/client/clientset/versioned"
	"sync"
)

var client *Client
var lock sync.Mutex

type Client struct {
	Config        *rest.Config
	K8sClient     *kubernetes.Clientset
	ServingClient *serving.Clientset
	TektonClient  *tektoncd.Clientset
	EtcdClient    *clientv3.Client
}

func GetClient() *Client {
	if client == nil {
		lock.Lock()
		if client == nil {
			initClient(true)
		}
		defer lock.Unlock()
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
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("error create k8s client: %v", err)
	}
	servingClient, err := serving.NewForConfig(config)
	if err != nil {
		fmt.Printf("error create knative serving client: %v", err)
	}
	tektonClient, err := tektoncd.NewForConfig(config)
	if err != nil {
		fmt.Printf("error create tekton client: %v", err)
	}
	//internalcfg, err := configutil.DefaultedInitConfiguration(&kubeadmapiv1beta1.InitConfiguration{})
	//if err != nil {
	//	fmt.Printf("unexpected error getting default config: %v", err)
	//}
	//fromCluster, err := etcd.NewFromCluster(k8sClient, internalcfg.CertificatesDir)
	//if fromCluster == nil || err != nil {
	//	fmt.Printf("new client with DialNoWait should succeed, got %v", err)
	//}
	//etcdClient, err := clientv3.New(clientv3.Config{Endpoints: fromCluster.Endpoints, TLS: fromCluster.TLS})
	//if etcdClient == nil || err != nil {
	//	fmt.Printf("new client with DialNoWait should succeed, got %v", err)
	//}

	client = &Client{
		Config:        config,
		K8sClient:     k8sClient,
		ServingClient: servingClient,
		TektonClient:  tektonClient,
		//EtcdClient:    etcdClient,
	}
}
