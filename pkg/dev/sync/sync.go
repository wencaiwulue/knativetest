package sync

import (
	"context"
	"io"
	v1 "k8s.io/api/core/v1"
	"os/exec"
)

type podSyncer struct {
	kubectl CLI
}

func NewPodSyncer(namespace, kubeconfig, kubecontext string) *podSyncer {
	return &podSyncer{kubectl: CLI{
		KubeContext: kubecontext,
		KubeConfig:  kubeconfig,
		Namespace:   namespace,
	}}
}

func (s *podSyncer) DeleteFileFn(ctx context.Context, pod v1.Pod, container v1.Container, files map[string][]string) *exec.Cmd {
	args := make([]string, 0, 9+len(files))
	args = append(args, pod.Name, "--namespace", pod.Namespace, "-c", container.Name, "--", "rm", "-rf", "--")
	for _, dsts := range files {
		args = append(args, dsts...)
	}
	return s.kubectl.Command(ctx, "exec", args...)
}

func (s *podSyncer) CopyFileFn(ctx context.Context, pod v1.Pod, container v1.Container, files map[string][]string) *exec.Cmd {
	// Use "m" flag to touch the files as they are copied.
	reader, writer := io.Pipe()
	go func() {
		if err := CreateMappedTar(writer, "/", files); err != nil {
			writer.CloseWithError(err)
		} else {
			writer.Close()
		}
	}()

	copyCmd := s.kubectl.Command(ctx, "exec", pod.Name, "--namespace", pod.Namespace, "-c", container.Name, "-i", "--", "tar", "xmf", "-", "-C", "/", "--no-same-owner")
	// attention here:
	// the most valuable code, using `kubectl exec` to extract the file from pipeline to pod directory
	copyCmd.Stdin = reader
	return copyCmd
}

func (s *podSyncer) CopyFolderFn(ctx context.Context, pod v1.Pod, container v1.Container, files map[string][]string) *exec.Cmd {
	// Use "m" flag to touch the files as they are copied.
	reader, writer := io.Pipe()
	go func() {
		if err := CreateMappedTar(writer, "/", files); err != nil {
			writer.CloseWithError(err)
		} else {
			writer.Close()
		}
	}()

	copyCmd := s.kubectl.Command(ctx, "exec", pod.Name, "--namespace", pod.Namespace, "-c", container.Name, "-i", "--", "tar", "xmf", "-", "-C", "/", "--no-same-owner")
	// attention here:
	// the most valuable code, using `kubectl exec` to extract the file from pipeline to pod directory
	copyCmd.Stdin = reader
	return copyCmd
}
