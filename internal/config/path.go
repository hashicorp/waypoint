// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"os"
	"path/filepath"
)

// Filename is the default filename for the Waypoint configuration.
const Filename = "waypoint.hcl"

// FindPath looks for our configuration file starting at "start" and
// traversing parent directories until it is found. If it is found, the
// path is returned. If it is not found, an empty string is returned.
// Error will be non-nil only if an error occurred.
//
// If start is empty, start will be the current working directory. If
// filename is empty, it will default to the Filename constant.
//
// If searchParent is false, then we will not search parent directories
// and require the Waypoint configuration file be directly in the "start"
// directory.
func FindPath(start, filename string, searchParent bool) (string, error) {
	var err error
	if start == "" {
		start, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	if filename == "" {
		filename = Filename
	}

	for {
		// Look for HCL syntax
		path := filepath.Join(start, filename)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}

		// Look for JSON
		path += ".json"
		if _, err := os.Stat(path); err == nil {
			return path, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}

		if !searchParent {
			return "", nil
		}

		next := filepath.Dir(start)
		if next == start {
			return "", nil
		}

		start = next
	}
}
