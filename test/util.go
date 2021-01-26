package test

import (
	"bytes"
	"os"
	"os/exec"
)

type binary struct {
	binaryPath string
	workingDir string
}

func NewBinary(binaryPath string, workingDir string) *binary {
	return &binary{
		binaryPath: binaryPath,
		workingDir: workingDir,
	}
}

// Builds a generic execer for running waypoint commands
func (b *binary) NewCmd(args ...string) *exec.Cmd {
	cmd := exec.Command(b.binaryPath, args...)
	cmd.Dir = b.workingDir
	cmd.Env = os.Environ()

	cmd.Env = append(cmd.Env, "CHECKPOINT_DISABLE=1")
	return cmd
}

// Runs a command with the arguments specified
func (b *binary) Run(args ...string) (stdout, stderr string, err error) {
	cmd := b.NewCmd(args...)
	cmd.Stdin = nil
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}
	err = cmd.Run()
	stdout = cmd.Stdout.(*bytes.Buffer).String()
	stderr = cmd.Stderr.(*bytes.Buffer).String()
	return
}
