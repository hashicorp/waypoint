package pack

import (
	"context"
	"os/exec"

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/terminal"
)

// Builder uses `pack` -- the frontend for CloudNative Buildpacks -- to build
// an artifact from source.
type Builder struct{}

// BuildFunc implements component.Builder
func (b *Builder) BuildFunc() interface{} {
	return b.Build
}

// Build
func (b *Builder) Build(
	ctx context.Context,
	ui terminal.UI,
	src *component.Source,
) (*DockerImage, error) {
	stdout, stderr, err := ui.OutputWriters()
	if err != nil {
		return nil, err
	}

	// Build the image using `pack`. This doesn't give us any more information
	// unfortunately so we can only run the build with the image name
	// we want as a result.
	cmd := exec.CommandContext(ctx, "pack", "build", src.App)
	cmd.Dir = src.Path
	cmd.Stdout = stdout
	cmd.Stderr = stderr
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
