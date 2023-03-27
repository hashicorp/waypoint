package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	defaultEditor = "vi"
	defaultShell  = "/bin/bash"
)

var (
	defaultEnvEditor = []string{"EDITOR"}
)

type Editor struct {
	// Various arguments required to launch $EDITOR
	Args []string
}

func NewEditor(args []string) (*Editor, error) {
	return &Editor{
		Args: args,
	}, nil
}

// Run will launch a editor to use a system defined editor such as vim to edit
// configs in place. It saves that content to a temp file for use as well as
// returning the raw bytes from the edit. It can optionally take an original
// bytes of content which can be used to compare if any edits were made.
func Run(o []byte) ([]byte, string, error) {
	var (
		original = []byte{}
		edited   = []byte{}
		filePath string
		err      error
	)

	// set an original if it exists
	if o != nil {
		original = o
	}

	edit, err := NewEditor(defaultEnvEditor)
	if err != nil {
		return nil, "", err
	}

	// generate the file to edit
	buf := &bytes.Buffer{}

	prefix := fmt.Sprintf("%s-edit-", filepath.Base(os.Args[0]))
	suffix := ".hcl"
	edited, filePath, err = edit.LaunchWithTmp(prefix, suffix, original, buf)
	if err != nil {
		return nil, "", err
	}

	if o != nil && bytes.Equal(original, edited) {
		return nil, "", fmt.Errorf("edited file matches original content")
	}

	// TODO(briancain): Get this working properly with json and hcl
	/*
		// Validate that the edited content can be parsed
		_, diag := hclsyntax.ParseConfig(edited, "<waypoint-hcl>", hcl.Pos{})
		if diag.HasErrors() {
			return nil, "", fmt.Errorf("the edited file has errors: %s", diag)
		}
	*/

	return edited, filePath, nil
}

func (e *Editor) LaunchEditor(filePath string) error {
	if len(e.Args) == 0 {
		return fmt.Errorf("No arguments given for launching editor tool")
	}
	abs, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	args := make([]string, len(e.Args))
	copy(args, e.Args)
	args = append(args, abs)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// launch the configured editor
	if err := cmd.Run(); err != nil {
		if err, ok := err.(*exec.Error); ok {
			if err.Err == exec.ErrNotFound {
				return fmt.Errorf("unable to launch editor %q with error %s",
					strings.Join(args, " "), err)
			}
		}
		return fmt.Errorf("an error was encountered while launching the editor %q with error %s",
			strings.Join(args, " "), err)
	}
	return nil
}

func (e *Editor) LaunchWithTmp(prefix, suffix string, original []byte, r io.Reader) ([]byte, string, error) {
	f, err := os.CreateTemp("", prefix+"*"+suffix)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	path := f.Name()
	if _, err := io.Copy(f, r); err != nil {
		os.Remove(path)
		return nil, path, err
	}
	if original != nil {
		_, err = f.Write(original)
		if err != nil {
			return nil, "", fmt.Errorf("failed to write original content to tmp file: %s", err)
		}
	}
	// This file descriptor needs to close so the next process (Launch) can claim it.
	f.Close()
	if err := e.LaunchEditor(path); err != nil {
		return nil, path, err
	}
	bytes, err := os.ReadFile(path)
	return bytes, path, err
}
