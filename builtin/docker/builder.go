package docker

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/assets"
	"github.com/hashicorp/waypoint/internal/pkg/epinject"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Builder uses `docker build` to build a Docker iamge.
type Builder struct {
	config BuilderConfig
}

// BuildFunc implements component.Builder
func (b *Builder) BuildFunc() interface{} {
	return b.Build
}

// Config is the configuration structure for the registry.
type BuilderConfig struct {
	// Control whether or not to inject the entrypoint binary into the resulting image
	DisableCEB bool `hcl:"disable_entrypoint,optional"`

	// Controls whether or not the image should be build with buildkit or docker v1
	UseBuildKit bool `hcl:"buildkit,optional"`

	// The name/path to the Dockerfile if it is not the root of the project
	Dockerfile string `hcl:"dockerfile,optional"`
}

func (b *Builder) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&BuilderConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Build a Docker image using the `docker build` protocol")

	doc.Example(`
build {
  use "docker" {
	buildkit    = false
	disable_entrypoint = false
  }
}
`)

	doc.Input("component.Source")
	doc.Output("docker.Image")

	doc.SetField(
		"disable_entrypoint",
		"if set, the entrypoint binary won't be injected into the image",
		docs.Summary(
			"The entrypoint binary is what provides extended functionality",
			"such as logs and exec. If it is not injected at build time",
			"the expectation is that the image already contains it",
		),
	)

	doc.SetField(
		"buildkit",
		"if set, use the buildkit builder from Docker",
	)

	return doc, nil
}

// Config implements Configurable
func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

// Build
func (b *Builder) Build(
	ctx context.Context,
	ui terminal.UI,
	src *component.Source,
) (*Image, error) {
	stdout, _, err := ui.OutputWriters()
	if err != nil {
		return nil, err
	}

	sg := ui.StepGroup()
	step := sg.Add("Initializing Docker client...")
	defer step.Abort()

	result := &Image{
		Image: fmt.Sprintf("waypoint.local/%s", src.App),
		Tag:   "latest",
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "unable to create Docker client: %s", err)
	}

	cli.NegotiateAPIVersion(ctx)

	dockerfile := b.config.Dockerfile
	if dockerfile != "" {
		dockerfile = path.Join(src.Path, dockerfile)
	}

	contextDir, relDockerfile, err := build.GetContextFromLocalDir(src.Path, dockerfile)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "unable to create Docker context: %s", err)
	}

	excludes, err := build.ReadDockerignore(contextDir)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to read .dockerignore: %s", err)
	}

	if err := build.ValidateContextDirectory(contextDir, excludes); err != nil {
		return nil, status.Errorf(codes.Internal, "error checking context: %s", err)
	}

	// And canonicalize dockerfile name to a platform-independent one
	relDockerfile = archive.CanonicalTarNameForPath(relDockerfile)

	excludes = build.TrimBuildFilesFromExcludes(excludes, relDockerfile, false)
	buildCtx, err := archive.TarWithOptions(contextDir, &archive.TarOptions{
		ExcludePatterns: excludes,
		ChownOpts:       &idtools.Identity{UID: 0, GID: 0},
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to compress context: %s", err)
	}

	ver := types.BuilderV1
	if b.config.UseBuildKit {
		ver = types.BuilderBuildKit
	}

	step.Done()
	step = sg.Add("Building image...")

	resp, err := cli.ImageBuild(ctx, buildCtx, types.ImageBuildOptions{
		Version:    ver,
		Dockerfile: relDockerfile,
		Tags:       []string{result.Name()},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error building image: %s", err)
	}
	defer resp.Body.Close()

	var termFd uintptr
	if f, ok := stdout.(*os.File); ok {
		termFd = f.Fd()
	}

	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, step.TermOutput(), termFd, true, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to stream build logs to the terminal: %s", err)
	}

	step.Done()

	if !b.config.DisableCEB {
		step = sg.Add("Injecting Waypoint Entrypoint...")

		asset, err := assets.Asset("ceb/ceb")
		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to restore custom entry point binary: %s", err)
		}

		assetInfo, err := assets.AssetInfo("ceb/ceb")
		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to restore custom entry point binary: %s", err)
		}

		_, err = epinject.AlterEntrypoint(ctx, result.Name(), func(cur []string) (*epinject.NewEntrypoint, error) {
			ep := &epinject.NewEntrypoint{
				Entrypoint: append([]string{"/waypoint-entrypoint"}, cur...),
				InjectFiles: map[string]epinject.InjectFile{
					"/waypoint-entrypoint": {
						Reader: bytes.NewReader(asset),
						Info:   assetInfo,
					},
				},
			}

			return ep, nil
		})

		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to set modify Docker entrypoint: %s", err)
		}

		step.Done()
	}

	return result, nil
}
