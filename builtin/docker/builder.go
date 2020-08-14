package docker

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/hashicorp/waypoint/internal/assets"
	"github.com/hashicorp/waypoint/internal/pkg/epinject"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/terminal"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
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
	DisableCEB bool `hcl:"disable_ceb,optional"`

	// Controls whether or not the image should be build with buildkit or docker v1
	UseBuildKit bool `hcl:"buildkit,optional"`
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

	result := &Image{
		Image: fmt.Sprintf("waypoint.local/%s", src.App),
		Tag:   "latest",
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	cli.NegotiateAPIVersion(ctx)

	contextDir, relDockerfile, err := build.GetContextFromLocalDir(src.Path, "")
	if err != nil {
		return nil, err
	}

	excludes, err := build.ReadDockerignore(contextDir)
	if err != nil {
		return nil, err
	}

	if err := build.ValidateContextDirectory(contextDir, excludes); err != nil {
		return nil, errors.Errorf("error checking context: '%s'.", err)
	}

	// And canonicalize dockerfile name to a platform-independent one
	relDockerfile = archive.CanonicalTarNameForPath(relDockerfile)

	excludes = build.TrimBuildFilesFromExcludes(excludes, relDockerfile, false)
	buildCtx, err := archive.TarWithOptions(contextDir, &archive.TarOptions{
		ExcludePatterns: excludes,
		ChownOpts:       &idtools.Identity{UID: 0, GID: 0},
	})

	if err != nil {
		return nil, err
	}

	ver := types.BuilderV1
	if b.config.UseBuildKit {
		ver = types.BuilderBuildKit
	}

	resp, err := cli.ImageBuild(ctx, buildCtx, types.ImageBuildOptions{
		Version:    ver,
		Dockerfile: relDockerfile,
		Tags:       []string{result.Name()},
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var (
		termFd uintptr
		isTerm bool
	)

	if f, ok := stdout.(*os.File); ok {
		termFd = f.Fd()
		isTerm = isatty.IsTerminal(termFd)
	}

	err = jsonmessage.DisplayJSONMessagesStream(resp.Body, stdout, termFd, isTerm, nil)
	if err != nil {
		return nil, err
	}

	if !b.config.DisableCEB {
		tmpdir, err := ioutil.TempDir("", "waypoint")
		if err != nil {
			return nil, err
		}

		defer os.RemoveAll(tmpdir)

		err = assets.RestoreAsset(tmpdir, "ceb/ceb")
		if err != nil {
			return nil, err
		}

		_, err = epinject.AlterEntrypoint(ctx, result.Name(), func(cur []string) (*epinject.NewEntrypoint, error) {
			ep := &epinject.NewEntrypoint{
				Entrypoint: append([]string{"/bin/wpceb"}, cur...),
				InjectFiles: map[string]string{
					filepath.Join(tmpdir, "ceb/ceb"): "/bin/wpceb",
				},
			}

			return ep, nil
		})

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
