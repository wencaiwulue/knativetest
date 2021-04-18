package sync

import (
	"context"
	"io"
	v1 "k8s.io/api/core/v1"
	"os/exec"
)

type podSyncer struct {
	kubectl *CLI
}

type Syncer interface {
	DeleteFileCmd(ctx context.Context, pod v1.Pod, container v1.Container, files map[string][]string) *exec.Cmd
	CopyFileCmd(ctx context.Context, pod v1.Pod, container v1.Container, files map[string][]string) *exec.Cmd
	CopyFolderCmd(ctx context.Context, pod v1.Pod, container v1.Container, local string, remote string) *exec.Cmd
}

func NewPodSyncer(cfg Config, namespace string) *podSyncer {
	return &podSyncer{kubectl: NewCLI(cfg, namespace)}
}

func (s *podSyncer) DeleteFileCmd(ctx context.Context, pod v1.Pod, container v1.Container, files map[string][]string) *exec.Cmd {
	args := make([]string, 0, 9+len(files))
	args = append(args, pod.Name, "--namespace", pod.Namespace, "-c", container.Name, "--", "rm", "-rf", "--")
	for _, dsts := range files {
		args = append(args, dsts...)
	}
	return s.kubectl.Command(ctx, "exec", args...)
}

func (s *podSyncer) CopyFileCmd(ctx context.Context, pod v1.Pod, container v1.Container, files map[string][]string) *exec.Cmd {
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

func (s *podSyncer) CopyFolderCmd(ctx context.Context, pod v1.Pod, container v1.Container, local string, remote string) *exec.Cmd {
	copyCmd := s.kubectl.Command(ctx, "cp", "--namespace", pod.Namespace, local, pod.Name+":"+remote, "-c", container.Name)
	return copyCmd
}
