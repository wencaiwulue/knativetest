package sync

import (
	"context"
	"os/exec"
)

// Cmd is a wrapper on exec.Cmd
type Cmd struct {
	*exec.Cmd
}

// CommandContext creates a new Cmd
func CommandContext(ctx context.Context, name string, arg ...string) *Cmd {
	return &Cmd{Cmd: exec.CommandContext(ctx, name, arg...)}
}

// Terminate kills the underlying process
func (c *Cmd) Terminate() error {
	return c.Process.Kill()
}
