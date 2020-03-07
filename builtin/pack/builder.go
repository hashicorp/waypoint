package pack

import (
	"context"
	"os"
	"os/exec"

	"github.com/mitchellh/devflow/sdk/component"
)

// Builder uses `pack` -- the frontend for CloudNative Buildpacks -- to build
// an artifact from source.
type Builder struct{}

// BuildFunc implements component.Builder
func (b *Builder) BuildFunc() interface{} {
	return b.Build
}

// Build
func (b *Builder) Build(ctx context.Context, src *component.Source) (*DockerImage, error) {
	// Build the image using `pack`. This doesn't give us any more information
	// unfortunately so we can only run the build with the image name
	// we want as a result.
	cmd := exec.CommandContext(ctx, "pack", "build", src.App)
	cmd.Dir = src.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// We don't even need to inspect Docker to verify we have the image.
	// If `pack` succeeded we can assume that it created an image for us.
	return &DockerImage{
		Image: src.App,
		Tag:   "latest", // It always tags latest
	}, nil
}
