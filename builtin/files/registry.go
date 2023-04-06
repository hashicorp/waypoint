// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package files

import (
	"context"
	"crypto/rand"
	"path"

	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/pkg/copy"
	"github.com/oklog/ulid"
)

// Registry represents access to a Files registry.
type Registry struct {
	config Config
}

// Config implements Configurable
func (r *Registry) Config() (interface{}, error) {
	return &r.config, nil
}

// PushFunc implements component.Registry
func (r *Registry) PushFunc() interface{} {
	return r.Push
}

// Push pushes an image to the registry.
func (r *Registry) Push(
	ctx context.Context,
	files *Files,
	ui terminal.UI,
) (*Files, error) {
	// Generate a unique path for the destination file
	dstID, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}

	dst := path.Join(r.config.Path, files.Path, dstID.String())

	err = copy.CopyDir(files.Path, dst)

	if err != nil {
		return nil, err
	}

	return &Files{Path: dst}, nil
}

// Config is the configuration structure for the registry.
type Config struct {
	// Path is the path that files are stored in
	Path string `hcl:"path,attr"`
}

func (r *Registry) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}), docs.FromFunc(r.PushFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description("Copies files to a specific directory")
	doc.Input("files.Files")
	doc.Output("files.Files")

	doc.Example(`
build {
  use "files" {}
  registry {
	use "files" {
	  path = "/path/to/file"
	}
  }
}
`)

	doc.SetField(
		"path",
		"the filesystem path to store the files into",
	)

	return doc, nil
}
