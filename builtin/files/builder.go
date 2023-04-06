// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package files

import (
	"context"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

// Builder validates the files in a directory or in the filesystem
type Builder struct {
	config BuilderConfig
}

// BuildFunc implements component.Builder
func (b *Builder) BuildFunc() interface{} {
	return b.Build
}

// BuilderConfig is the configuration structure for the builder.
type BuilderConfig struct{}

// Config implements Configurable
func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

// Build
func (b *Builder) Build(
	ctx context.Context,
	ui terminal.UI,
	src *component.Source,
) (*Files, error) {
	return &Files{
		Path: src.Path,
	}, nil
}

func (b *Builder) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&BuilderConfig{}), docs.FromFunc(b.BuildFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description("Generates a value representing a path on disk")
	doc.Input("component.Source")
	doc.Output("files.Files")

	doc.Example(`
build {
  use "files" {}
}
`)

	return doc, nil
}
