package extra

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/util/term"
	"knativetest/pkg/dev/util"
	"knativetest/pkg/dev/watch"
	"log"
	"os"
)

func Shell(clientSet *util.ClientSet, options watch.DevOptions) error {
	podName := options.GetPod()
	if podName == "" {
		pods, err := clientSet.ClientSet.CoreV1().Pods(options.GetNamespace()).List(context.TODO(), metav1.ListOptions{LabelSelector: "kubedev=debug"})
		if err != nil {
			log.Fatalf("get kubedev pod error: %v", err)
		}
		if len(pods.Items) != 1 {
			log.Println("this should not happened")
		}
		podName = pods.Items[0].Name
	}

	tty := term.TTY{
		Out: os.Stdout,
		In:  os.Stdin,
		Raw: true,
	}

	if !tty.IsTerminalIn() {
		log.Println("Unable to use a TTY - input is not a terminal or the right kind of file")
	}

	var terminalSizeQueue remotecommand.TerminalSizeQueue
	if tty.Raw {
		terminalSizeQueue = tty.MonitorSize(tty.GetSize())
	}
	f := func() error {
		rc := clientSet.RestConfig
		restClient, err := restclient.RESTClientFor(rc)
		if err != nil {
			return err
		}

		req := restClient.Post().
			Resource("pods").
			Name(podName).
			Namespace(options.GetDeployment()).
			SubResource("exec").
			VersionedParams(&v1.PodExecOptions{
				Container: options.GetContainer(),
				Command:   []string{"bash"},
				Stdin:     true,
				Stdout:    true,
				Stderr:    false,
				TTY:       true,
			}, scheme.ParameterCodec)

		executor, err := remotecommand.NewSPDYExecutor(rc, "post", req.URL())
		if err != nil {
			return err
		}
		return executor.Stream(remotecommand.StreamOptions{
			Stdin:             tty.In,
			Stdout:            tty.Out,
			Stderr:            os.Stderr,
			Tty:               true,
			TerminalSizeQueue: terminalSizeQueue,
		})
	}

	if err := tty.Safe(f); err != nil {
		return err
	}
	return nil
}
