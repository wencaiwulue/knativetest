package sync

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type cmdError struct {
	args   []string
	stdout []byte
	stderr []byte
	cause  error
}

func (e *cmdError) Error() string {
	return fmt.Sprintf("running %s\n - stdout: %q\n - stderr: %q\n - cause: %s", e.args, e.stdout, e.stderr, e.cause)
}

func (e *cmdError) Unwrap() error {
	return e.cause
}

func (e *cmdError) ExitCode() int {
	if exitError, ok := e.cause.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}
	return 0
}

// DefaultExecCommand runs commands using exec.Cmd
var DefaultExecCommand Command = &Commander{}

// Command is an interface used to run commands. All packages should use this
// interface instead of calling exec.Cmd directly.
type Command interface {
	RunCmdOut(cmd *exec.Cmd) ([]byte, error)
	RunCmd(cmd *exec.Cmd) error
}

func RunCmdOut(cmd *exec.Cmd) ([]byte, error) {
	return DefaultExecCommand.RunCmdOut(cmd)
}

func RunCmd(cmd *exec.Cmd) error {
	return DefaultExecCommand.RunCmd(cmd)
}

// Commander is the exec.Cmd implementation of the Command interface
type Commander struct{}

// RunCmdOut runs an exec.Command and returns the stdout and error.
func (*Commander) RunCmdOut(cmd *exec.Cmd) ([]byte, error) {
	logrus.Debugf("Running command: %s", cmd.Args)

	stdout := bytes.Buffer{}
	cmd.Stdout = &stdout
	stderr := bytes.Buffer{}
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting command %v: %w", cmd, err)
	}

	if err := cmd.Wait(); err != nil {
		return stdout.Bytes(), &cmdError{
			args:   cmd.Args,
			stdout: stdout.Bytes(),
			stderr: stderr.Bytes(),
			cause:  err,
		}
	}

	if stderr.Len() > 0 {
		logrus.Debugf("Command output: [%s], stderr: %s", stdout.String(), stderr.String())
	} else {
		logrus.Debugf("Command output: [%s]", stdout.String())
	}

	return stdout.Bytes(), nil
}

// RunCmd runs an exec.Command.
func (*Commander) RunCmd(cmd *exec.Cmd) error {
	logrus.Debugf("Running command: %s", cmd.Args)
	return cmd.Run()
}
