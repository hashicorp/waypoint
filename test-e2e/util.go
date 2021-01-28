package test

import (
	"bytes"
	"os"
	"os/exec"
)

// Test config settings used by the tests
var (
	wpBinary             = Getenv("WP_BINARY", "waypoint")
	wpServerImage        = Getenv("WP_SERVERIMAGE", "hashicorp/waypoint:0.2.0")
	wpServerImageUpgrade = Getenv("WP_SERVERIMAGE_UPGRADE", "hashicorp/waypoint:latest")

	examplesRootDir = Getenv("WP_EXAMPLES_PATH", "waypoint-examples")
)

// A struct representation of the waypoint binary
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

// Obtains the env var key. If unset, it returns the default
func Getenv(key, def string) string {
	result := os.Getenv(key)
	if result == "" {
		result = def
	}
	return result
}
