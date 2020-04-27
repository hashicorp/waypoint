package docker

import (
	"context"
	"os/exec"

	"github.com/hashicorp/waypoint/sdk/terminal"
)

// Registry represents access to a Docker registry.
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
	img *Image,
	ui terminal.UI,
) (*Image, error) {
	stdout, stderr, err := ui.OutputWriters()
	if err != nil {
		return nil, err
	}

	target := &Image{Image: r.config.Image, Tag: r.config.Tag}

	{
		// Re-tag the image to our target value
		cmd := exec.CommandContext(ctx, "docker", "tag", img.Name(), target.Name())
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			return nil, err
		}
	}

	if !r.config.Local {
		// Push it
		cmd := exec.CommandContext(ctx, "docker", "push", target.Name())
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			return nil, err
		}
	}

	return target, nil
}

// Config is the configuration structure for the registry.
type Config struct {
	// Image is the name of the image plus tag that the image will be pushed as.
	Image string `hcl:"image,attr"`

	// Tag is the tag to apply to the image.
	Tag string `hcl:"tag,attr"`

	// Local if true will not push this image to a remote registry.
	Local bool `hcl:"local,optional"`
}
