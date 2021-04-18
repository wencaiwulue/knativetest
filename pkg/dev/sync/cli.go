package sync

import (
	"context"
	"io"
	"os/exec"
)

type CLI struct {
	KubeContext string
	KubeConfig  string
	Namespace   string
}

type Config interface {
	GetKubeContext() string
	GetKubeConfig() string
	GetKubeNamespace() string
}

func NewCLI(cfg Config, defaultNamespace string) *CLI {
	ns := defaultNamespace
	if nsFromOpts := cfg.GetKubeNamespace(); nsFromOpts != "" {
		ns = nsFromOpts
	}
	return &CLI{
		KubeContext: cfg.GetKubeContext(),
		KubeConfig:  cfg.GetKubeConfig(),
		Namespace:   ns,
	}
}

// Command creates the underlying exec.CommandContext. This allows low-level control of the executed command.
func (c *CLI) Command(ctx context.Context, command string, arg ...string) *exec.Cmd {
	args := c.args(command, "", arg...)
	return exec.CommandContext(ctx, "kubectl", args...)
}

// Command creates the underlying exec.CommandContext with namespace. This allows low-level control of the executed command.
func (c *CLI) CommandWithNamespaceArg(ctx context.Context, command string, namespace string, arg ...string) *exec.Cmd {
	args := c.args(command, namespace, arg...)
	return exec.CommandContext(ctx, "kubectl", args...)
}

// Run shells out kubectl CLI.
func (c *CLI) Run(ctx context.Context, in io.Reader, out io.Writer, command string, arg ...string) error {
	cmd := c.Command(ctx, command, arg...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = out
	return RunCmd(cmd)
}

// RunInNamespace shells out kubectl CLI with given namespace
func (c *CLI) RunInNamespace(ctx context.Context, in io.Reader, out io.Writer, command string, namespace string, arg ...string) error {
	cmd := c.CommandWithNamespaceArg(ctx, command, namespace, arg...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = out
	return RunCmd(cmd)
}

// RunOut shells out kubectl CLI.
func (c *CLI) RunOut(ctx context.Context, command string, arg ...string) ([]byte, error) {
	cmd := c.Command(ctx, command, arg...)
	return RunCmdOut(cmd)
}

// RunOutInput shells out kubectl CLI with a given input stream.
func (c *CLI) RunOutInput(ctx context.Context, in io.Reader, command string, arg ...string) ([]byte, error) {
	cmd := c.Command(ctx, command, arg...)
	cmd.Stdin = in
	return RunCmdOut(cmd)
}

// CommandWithStrictCancellation ensures for windows OS that all child process get terminated on cancellation
func (c *CLI) CommandWithStrictCancellation(ctx context.Context, command string, arg ...string) *Cmd {
	args := c.args(command, "", arg...)
	return CommandContext(ctx, "kubectl", args...)
}

// args builds an argument list for calling kubectl and consistently
// adds the `--context` and `--namespace` flags.
func (c *CLI) args(command string, namespace string, arg ...string) []string {
	args := []string{"--context", c.KubeContext}
	namespace = c.resolveNamespace(namespace)
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	if c.KubeConfig != "" {
		args = append(args, "--kubeconfig", c.KubeConfig)
	}
	args = append(args, command)
	args = append(args, arg...)
	return args
}

func (c *CLI) resolveNamespace(ns string) string {
	if ns != "" {
		return ns
	}
	return c.Namespace
}
