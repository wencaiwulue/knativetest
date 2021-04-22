package network

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	sync2 "knativetest/pkg/dev/sync"
	"knativetest/pkg/dev/util"
	"knativetest/pkg/dev/watch"
	"log"
	"os/exec"
)

type forward struct {
	kubectl    *sync2.CLI
	client     *util.ClientSet
	namespace  string
	deployment string
	pod        string
	container  string
	localPort  int32
	remotePort int32
}

func (f *forward) GetKubeContext() string {
	return ""
}

func (f *forward) GetKubeConfig() string {
	return f.client.KubeConfig
}

func (f *forward) GetKubeNamespace() string {
	return f.namespace
}

func (f *forward) portForward(ctx context.Context) *exec.Cmd {
	args := []string{"--address", "0.0.0.0", "deployment/" + f.deployment, fmt.Sprintf("%d:%d", f.localPort, f.remotePort), "--namespace", f.GetKubeNamespace()}
	return f.kubectl.Command(ctx, "port-forward", args...)
}

func PortForward(clientSet *util.ClientSet, o watch.DevOptions) {
	f := &forward{
		client:     clientSet,
		namespace:  o.GetNamespace(),
		deployment: o.GetDeployment(),
		pod:        o.GetPod(),
		container:  o.GetContainer(),
	}
	f.kubectl = sync2.NewCLI(f, o.GetNamespace())
	f.localPort = f.getPort()
	f.remotePort = f.getPort()

	portForward := f.portForward(context.Background())
	fmt.Printf("Forwarding from 0.0.0.0:%d -> %d\n", f.localPort, f.remotePort)
	if b, err := sync2.RunCmdOut(portForward); err != nil {
		fmt.Printf("error port forward: %v, log: %s\n", err.Error(), string(b))
	}
}

func (f forward) getPort() int32 {
	if f.deployment != "" {
		deployment, err := f.client.ClientSet.AppsV1().Deployments(f.namespace).Get(context.TODO(), f.deployment, metav1.GetOptions{})
		if err != nil {
			log.Fatalf("get kubedev deployment error: %v", err)
		}
		return deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort
	} else if f.pod != "" {
		pods, err := f.client.ClientSet.CoreV1().Pods(f.namespace).Get(context.TODO(), f.pod, metav1.GetOptions{})
		if err != nil {
			log.Fatalf("get kubedev pod error: %v", err)
		}
		return pods.Spec.Containers[0].Ports[0].ContainerPort
	} else {
		pods, err := f.client.ClientSet.CoreV1().Pods(f.namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "kubedev=debug"})
		if err != nil {
			log.Fatalf("get kubedev pod error: %v", err)
		}
		if len(pods.Items) != 1 {
			log.Println("this should not happened")
		}
		return pods.Items[0].Spec.Containers[0].Ports[0].ContainerPort
	}
}
