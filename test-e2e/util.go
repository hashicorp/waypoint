// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
)

// Test config settings used by the tests
var (
	wpBinary             = Getenv("WP_BINARY", "waypoint")
	wpServerImage        = Getenv("WP_SERVERIMAGE", "hashicorp/waypoint:latest")
	wpOdrImage           = Getenv("WP_ODRIMAGE", "hashicorp/waypoint-odr:latest")
	wpServerImageUpgrade = Getenv("WP_SERVERIMAGE_UPGRADE", "ghcr.io/hashicorp/waypoint/alpha:latest")
	wpOdrImageUpgrade    = Getenv("WP_ODRIMAGE_UPGRADE", "ghcr.io/hashicorp/waypoint/alpha-odr:latest")

	examplesRootDir = Getenv("WP_EXAMPLES_PATH", "waypoint-examples")
)

// A struct representation of the waypoint binary
type binary struct {
	t          *testing.T
	binaryPath string
	workingDir string
}

func NewBinary(t *testing.T, binaryPath string, workingDir string) *binary {
	return &binary{
		t:          t,
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

// Super minimal table representation
type TableOutput struct {
	header string
	rows   []string
}

// Sets up a new project test directory in a temp folder, and renames the
// project to have a random suffix. Provides project isolation between runs.
func SetupTestProject(templateDir string) (projectName string, projectDir string, err error) {
	projectRandomId := randSuffix()

	tempDir, err := os.MkdirTemp("", randSuffix())
	if err != nil {
		return "", "", errors.Wrap(err, "failed to create temp dir")
	}
	templateDir = templateDir + "/"
	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("rsync", "-av", templateDir, tempDir)
	} else {
		cmd = exec.Command("cp", "-r", templateDir, tempDir)
	}
	err = cmd.Run()
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to copy %s to %s", templateDir, tempDir)
	}

	waypointHclFilePath := tempDir + "/waypoint.hcl"

	hclContents, err := ioutil.ReadFile(waypointHclFilePath)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to read %s", waypointHclFilePath)
	}

	r := regexp.MustCompile(`project = "(.*)"`)

	newHclContents := r.ReplaceAll(hclContents, []byte("project = \"$1-"+projectRandomId+"\""))

	newProjectName := r.FindString(string(newHclContents))

	err = ioutil.WriteFile(waypointHclFilePath, newHclContents, 0644)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to write %s", waypointHclFilePath)
	}

	// Just the project name
	newProjectName = newProjectName[strings.Index(newProjectName, "\"")+1 : len(newProjectName)-1]

	return newProjectName, tempDir, nil
}

// Runs the test, ensures output contains outputContains
func (b *binary) RunWithOutput(outputContains string, args string) {
	stdout := b.Run(args)
	if !strings.Contains(stdout, outputContains) {
		b.t.Fatalf("output of command %q is %q, does not contain %q", args, stdout, outputContains)
	}
}

// Appends -o json to the args, runs the test, returns the unmarshalled json
func (b *binary) RunJson(args string) map[string]interface{} {
	args = args + " -o json"
	stdout := b.Run(args)
	ret := map[string]interface{}{}
	err := json.Unmarshal([]byte(stdout), &ret)
	if err != nil {
		b.t.Fatalf("failed to marshal json output for command %s: %s", args, err)
	}
	return ret
}

// Runs the test, formats as a table, and requires expectedLength rows
func (b *binary) RunTableExpectLength(expectedLength int, args string) {
	table := b.RunTable(args)
	if len(table.rows) != expectedLength {
		b.t.Fatalf("command %s returned %d rows, not expected %d.\nRows: %s", args, len(table.rows), expectedLength, strings.Join(table.rows, "\n"))
	}
}

// Runs the command, presents output as a table, fails test on error
func (b *binary) RunTable(args string) TableOutput {
	stdout := b.Run(args)
	lines := strings.Split(stdout, "\n")
	if len(lines) == 0 {
		b.t.Fatalf("command %s returned empty", args)
	}

	// Remove empty lines
	var validLines []string
	for _, line := range lines {
		if line != "" {
			validLines = append(validLines, line)
		}
	}

	headerSeparatorIndex := -1

	for i, line := range validLines {
		if strings.HasPrefix(line, "---") {
			headerSeparatorIndex = i
		}
	}

	if headerSeparatorIndex == -1 {
		b.t.Fatalf("command %s returned a non-table:\n%s", args, stdout)
	}

	return TableOutput{
		header: validLines[headerSeparatorIndex-1],
		rows:   validLines[headerSeparatorIndex+1:],
	}
}

func splitArgs(args string) []string {
	return strings.Split(args, " ")
}

// Runs the command, fails the test on errors
func (b *binary) Run(args string) (stdout string) {
	fmt.Printf("running %s ...\n", args)
	stdout, stderr, err := b.RunRaw(splitArgs(args)...)
	if err != nil {
		b.t.Fatalf("unexpected error running %q inside %q\nERROR:\n%s\n\nSTDERR:\n%s\n\nSTDOUT:\n%s", args, b.workingDir, err, stderr, stdout)
	}
	if stderr != "" {
		b.t.Fatalf("unexpected stderr output running %s:\n%s", args, stderr)
	}
	return stdout
}

// Runs a command with the arguments specified
func (b *binary) RunRaw(args ...string) (stdout, stderr string, err error) {
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

func randSuffix() string {
	charset := "abcdefghijklmnopqrstuvwxyz"
	length := 5

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
