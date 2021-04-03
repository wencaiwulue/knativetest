package client

import (
	"context"
	"fmt"
	tektoncd "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	serving "knative.dev/serving/pkg/client/clientset/versioned"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
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

func init() {
	initClient(true)
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
	clusterStatus, err := GetClusterStatus(k8sClient)
	if err != nil {
		fmt.Printf("GetClusterStatus error: %v\n", err)
	} else {
		fmt.Printf("clusterStatus: %v, endpoints: %v\n", clusterStatus, clusterStatus)
	}

	certificatesDir := filepath.Join("/run/config", "pki")

	tlsInfo := transport.TLSInfo{
		CertFile:      filepath.Join(certificatesDir, "etcd/ca.crt"),
		KeyFile:       filepath.Join(certificatesDir, "etcd/healthcheck-client.crt"),
		TrustedCAFile: filepath.Join(certificatesDir, "etcd/healthcheck-client.key"),
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		fmt.Printf("tcs error: %v, tlsConfig: %v\n", err, tlsConfig)
	}

	endpoints := []string{GetClientURLByIP(clusterStatus)}

	etcdClient, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if etcdClient == nil || err != nil {
		fmt.Printf("new client with DialNoWait should succeed, got %v, endpoints: %v", err, endpoints)
	} else {
		fmt.Printf("etcd cient: %v \n", etcdClient)
	}
	resp, err := etcdClient.Get(context.TODO(), "testKey", clientv3.WithFromKey())
	if err != nil {
		fmt.Printf("error info: %v", err)
	} else {
		fmt.Printf("response: %v\n", resp)
	}
	res, err := etcdClient.Put(context.TODO(), "testKey", "testValue")
	if err != nil {
		fmt.Printf("put data error: %v\n", err)
	} else {
		fmt.Printf("put data response: %v\n", res)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	resp, err = etcdClient.Get(ctx, "testKey", clientv3.WithRev(res.Header.Revision))
	defer cancel()
	if err != nil {
		fmt.Printf("Get data error: %v\n", err)
	} else {
		for _, ev := range resp.Kvs {
			fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		}
	}

	client = &Client{
		Config:        config,
		K8sClient:     k8sClient,
		ServingClient: servingClient,
		TektonClient:  tektonClient,
		EtcdClient:    etcdClient,
	}

}

func GetClusterStatus(client *kubernetes.Clientset) (string, error) {
	configMap, err := client.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get(context.Background(), "kubeadm-config", metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		fmt.Printf("not found")
		return "", nil
	}
	if err != nil {
		return "", err
	}

	s := configMap.Data["ClusterStatus"]
	fmt.Printf("configmap ClusterStatus: %v\n", s)

	clusterStatus := UnmarshalClusterStatus(configMap.Data)

	return clusterStatus, nil
}
func UnmarshalClusterStatus(data map[string]string) string {
	clusterStatusData, _ := data["ClusterStatus"]
	ip := clusterStatusData[strings.LastIndex(clusterStatusData, "advertiseAddress"):strings.LastIndex(clusterStatusData, "bindPort")]
	ip = strings.Split(ip, ":")[1]
	ip = strings.TrimRight(ip, " ")
	ip = strings.TrimLeft(ip, " ")
	return ip
}

func GetClientURLByIP(ip string) string {
	return "https://" + net.JoinHostPort(ip, strconv.Itoa(2379))
}
