package files

import (
	"context"
	"path"

	"github.com/hashicorp/waypoint/internal/pkg/copy"
	"github.com/hashicorp/waypoint/sdk/terminal"
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
	dst := path.Join(files.Directory, r.config.Path)

	err := copy.CopyDir(files.Directory, dst)

	if err != nil {
		return nil, err
	}

	return &Files{Absolute: dst}, nil
}

// Config is the configuration structure for the registry.
type Config struct {
	// Path is the path that files are stored in
	Path string `hcl:"path,attr"`
}
