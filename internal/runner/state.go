// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package runner

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
)

// verifyStateDir verifies that the directory exists and can be read
// and written to.
func verifyStateDir(L hclog.Logger, dir string) error {
	L = L.With("state_dir", dir)
	L.Debug("verifying the state directory")

	// create if it doesn't exist
	if _, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		L.Info("state directory does not exist, creating it")
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}

	// read and write a file to verify we have write privs
	tf, err := ioutil.TempFile(dir, "rwtest")
	if err != nil {
		return err
	}
	defer os.Remove(tf.Name())
	defer tf.Close()
	if _, err := tf.Write([]byte("test")); err != nil {
		return err
	}

	return nil
}

func (r *Runner) statePutId(v string) error {
	if r.stateDir == "" {
		return nil
	}

	path := filepath.Join(r.stateDir, "id")
	return ioutil.WriteFile(path, []byte(v), 0600)
}

func (r *Runner) stateGetId() (string, error) {
	path := filepath.Join(r.stateDir, "id")
	data, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (r *Runner) statePutToken(v string) error {
	if r.stateDir == "" {
		return nil
	}

	path := filepath.Join(r.stateDir, "token")
	return ioutil.WriteFile(path, []byte(v), 0600)
}

func (r *Runner) stateGetToken() (string, error) {
	path := filepath.Join(r.stateDir, "token")
	data, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return string(data), nil
}
