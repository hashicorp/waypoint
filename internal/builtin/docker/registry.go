package docker

import (
	"context"
)

// Registry represents access to a Docker registry.
type Registry struct {
	config *Config
}

// Config implements Configurable
func (r *Registry) Config() *Config {
	return &r.config
}

// PushFunc implements component.Registry
func (r *Registry) PushFunc() interface{} {
	return r.Push
}

// Push pushes an image to the registry.
func (r *Registry) Push(ctx context.Context, img *Image) (*Image, error) {
	// Re-tag the image to our target value
	cmd := exec.CommandContext(ctx, "pack", "build", b.source.App)
	cmd.Dir = b.source.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return nil, nil
}

// Config is the configuration structure for the registry.
type Config struct {
	// Image is the name of the image plus tag that the image will be pushed as.
	Image string `hcl:"image,attr"`
}
