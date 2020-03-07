package docker

import (
	"context"
	"os"
	"os/exec"
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
func (r *Registry) Push(ctx context.Context, img *Image) (*Image, error) {
	target := &Image{Image: r.config.Image, Tag: r.config.Tag}

	{
		// Re-tag the image to our target value
		cmd := exec.CommandContext(ctx, "docker", "tag", img.Name(), target.Name())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return nil, err
		}
	}

	{
		// Push it
		cmd := exec.CommandContext(ctx, "docker", "push", target.Name())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
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
}
