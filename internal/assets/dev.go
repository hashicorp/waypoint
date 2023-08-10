// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:generate go-bindata -dev -pkg assets -o dev_assets.go -tags !assetsembedded ceb

//go:build !assetsembedded
// +build !assetsembedded

package assets

import (
	"os"
	"path/filepath"
)

var rootDir string

func init() {
	// Set a reasonable default in the event we somehow fail to find the root
	// directory
	rootDir = "./internal/assets"
	dir, err := os.Getwd()
	if err != nil {
		// There is some strange circumstance that would cause this to panic,
		// but would only happen in a dev environment anyway.
		panic(err)
	}

	for dir != "/" {
		path := filepath.Join(dir, "internal/assets")
		if _, err := os.Stat(path); err == nil {
			rootDir = path
			return
		}

		nextDir := filepath.Dir(dir)
		if nextDir == dir {
			break
		}

		dir = nextDir
	}
}
