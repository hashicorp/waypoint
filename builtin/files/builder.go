package files

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/terminal"
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
	err := filepath.Walk(src.Path, func(path string, info os.FileInfo, err error) error {
		// Check each file
		if !info.IsDir() && info.Mode().IsRegular() {
			o, err := os.Open(path)
			defer o.Close()
			if err != nil {
				return err
			}

			ui.Output("Have file %s", path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Files{
		Path: src.Path,
	}, nil
}
