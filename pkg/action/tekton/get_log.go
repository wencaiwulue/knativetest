package tekton

import (
	"bytes"
	"context"
	"io"
	corev1 "k8s.io/api/core/v1"
	"knativetest/pkg/action"
	"knativetest/pkg/client"
)

type GetLogAction struct {
	action.Action
	Namespace string
	Name      string
	Container string
}

func (a *GetLogAction) Process(ctx context.Context) interface{} {
	coreV1 := client.GetClient().K8sClient.CoreV1()
	request := coreV1.Pods(a.Namespace).GetLogs(a.Name, &corev1.PodLogOptions{
		Container: a.Container,
	})

	logs, err := request.Stream(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = logs.Close()
	}()

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, logs); err != nil {
		return err
	}

	return buf.String()
}
