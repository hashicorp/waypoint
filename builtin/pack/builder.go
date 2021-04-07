package pack

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/buildpacks/pack"
	"github.com/buildpacks/pack/logging"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	wpdockerclient "github.com/hashicorp/waypoint/builtin/docker/client"
	"github.com/hashicorp/waypoint/internal/assets"
	"github.com/hashicorp/waypoint/internal/pkg/epinject"
)

// Builder uses `pack` -- the frontend for CloudNative Buildpacks -- to build
// an artifact from source.
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

	// The Buildpack builder image to use, defaults to the standard heroku one.
	Builder string `hcl:"builder,optional"`

	// The exact buildpacks to use.
	Buildpacks []string `hcl:"buildpacks,optional"`

	// Environment variables that are meant to configure the application in a static
	// way. This might be control an image that has mulitple modes of operation,
	// selected via environment variable. Most configuration should use the waypoint
	// config commands.
	StaticEnvVars map[string]string `hcl:"static_environment,optional"`
}

const DefaultBuilder = "heroku/buildpacks:18"

// Config implements Configurable
func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

var skipBuildPacks = map[string]struct{}{
	"heroku/procfile": {},
}

// Build
func (b *Builder) Build(
	ctx context.Context,
	ui terminal.UI,
	jobInfo *component.JobInfo,
	src *component.Source,
	log hclog.Logger,
) (*DockerImage, error) {
	builder := b.config.Builder
	if builder == "" {
		builder = DefaultBuilder
	}

	dockerClient, err := wpdockerclient.NewClientWithOpts(
		client.FromEnv,
		// If we don't specify a version, the client will use too new an API, and users
		// will get and error of the form shown below. Note that when you don't pass a
		// client 'pack' does the same thing we're doing here:
		//
		// client version X.XX is too new. Maximum supported API
		client.WithVersion("1.38"),
	)
	if err != nil {
		return nil, err
	}

	// We now test if Docker is actually functional. Pack requires a Docker
	// daemon and we can't fallback to "img" or any other Dockerless solution.
	log.Debug("testing if Docker is available")
	if fallback, err := wpdockerclient.Fallback(ctx, log, dockerClient); err != nil {
		log.Warn("error during check if we should use Docker fallback", "err", err)
		return nil, status.Errorf(codes.Internal,
			"error validating Docker connection: %s", err)
	} else if fallback {
		ui.Output(
			`WARNING: `+
				`Docker daemon appears unavailable. The 'pack' builder requires access `+
				`to a Docker daemon. Pack does not support dockerless builds. We will `+
				`still attempt to run the build but it will likely fail. If you are `+
				`running this build locally, please install Docker. If you are running `+
				`this build remotely (in a Waypoint runner), the runner must be configured `+
				`to have access to the Docker daemon.`+"\n",
			terminal.WithWarningStyle(),
		)
	} else {
		log.Debug("Docker appears available")
	}

	ui.Output("Creating new buildpack-based image using builder: %s", builder)

	sg := ui.StepGroup()

	step := sg.Add("Creating pack client")
	defer step.Abort()

	build := sg.Add("Building image")
	defer build.Abort()

	client, err := pack.NewClient(
		pack.WithLogger(logging.New(build.TermOutput())),
		pack.WithDockerClient(dockerClient),
	)
	if err != nil {
		return nil, err
	}

	step.Done()

	err = client.Build(ctx, pack.BuildOptions{
		Image:      src.App,
		Builder:    builder,
		AppPath:    src.Path,
		Env:        b.config.StaticEnvVars,
		Buildpacks: b.config.Buildpacks,
		FileFilter: func(file string) bool {
			// Do not include the bolt.db or bolt.db.lock
			// These files hold the local state when Waypoint is running without a server
			// on Windows it will not be possible to copy these files due to a file lock.
			if jobInfo.Local {
				if strings.HasSuffix(file, "data.db") || strings.HasSuffix(file, "data.db.lock") {
					return false
				}
			}

			return true
		},
	})

	if err != nil {
		return nil, err
	}

	build.Done()

	info, err := client.InspectImage(src.App, true)
	if err != nil {
		return nil, err
	}

	labels := map[string]string{}

	var languages []string

	for _, bp := range info.Buildpacks {
		if _, ok := skipBuildPacks[bp.ID]; ok {
			continue
		}

		idx := strings.IndexByte(bp.ID, '/')
		if idx != -1 {
			languages = append(languages, bp.ID[idx+1:])
		} else {
			languages = append(languages, bp.ID)
		}
	}

	labels["common/languages"] = strings.Join(languages, ",")
	labels["common/buildpack-stack"] = info.StackID

	proc := info.Processes.DefaultProcess
	if proc != nil {
		cmd := proc.Command

		if len(proc.Args) > 0 {
			if len(cmd) > 0 {
				cmd = fmt.Sprintf("%s %s", cmd, strings.Join(proc.Args, " "))
			} else {
				cmd = strings.Join(proc.Args, " ")
			}
		}

		if cmd != "" {
			labels["common/command"] = cmd
			if proc.Type != "" {
				labels["common/command-type"] = proc.Type
			}
		}
	}

	if !b.config.DisableCEB {
		inject := sg.Add("Injecting entrypoint binary to image")
		defer inject.Abort()

		asset, err := assets.Asset("ceb/ceb")
		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to restore custom entry point binary: %s", err)
		}

		assetInfo, err := assets.AssetInfo("ceb/ceb")
		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to restore custom entry point binary: %s", err)
		}

		imageId, err := epinject.AlterEntrypoint(ctx, src.App+":latest", func(cur []string) (*epinject.NewEntrypoint, error) {
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
			return nil, err
		}

		labels["common/image-id"] = imageId

		inject.Done()
	}

	sg.Wait()

	ui.Output("")
	ui.Output("Generated new Docker image: %s:latest", src.App)

	// We don't even need to inspect Docker to verify we have the image.
	// If `pack` succeeded we can assume that it created an image for us.
	return &DockerImage{
		Image:       src.App,
		Tag:         "latest", // It always tags latest
		BuildLabels: labels,
	}, nil
}

func (b *Builder) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&BuilderConfig{}), docs.FromFunc(b.BuildFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description(`
Create a Docker image using CloudNative Buildpacks.

**Pack requires access to a Docker daemon.** For remote builds, such as those
triggered by [Git polling](/docs/projects/git), the
[runner](/docs/runner) needs to have access to a Docker daemon such
as exposing the Docker socket, enabling Docker-in-Docker, etc. Unfortunately,
pack doesn't support dockerless builds. Configuring Docker access within
a Docker container is outside the scope of these docs, please search the
internet for "Docker in Docker" or other terms for more information.
`)

	doc.Example(`
build {
  use "pack" {
	builder     = "heroku/buildpacks:18"
	disable_entrypoint = false
  }
}
`)

	doc.Input("component.Source")
	doc.Output("pack.Image")
	doc.AddMapper(
		"pack.Image",
		"docker.Image",
		"Allow pack images to be used as normal docker images",
	)

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
		"builder",
		"The buildpack builder image to use",
		docs.Default(DefaultBuilder),
	)

	doc.SetField(
		"buildpacks",
		"The exact buildpacks to use",
		docs.Summary(
			"If set, the builder will run these buildpacks in the specified order.\n\n",
			"They can be listed using several [URI formats](https://buildpacks.io/docs/app-developer-guide/specific-buildpacks).",
		),
	)

	doc.SetField(
		"static_environment",
		"environment variables to expose to the buildpack",
		docs.Summary(
			"these environment variables should not be run of the mill",
			"configuration variables, use waypoint config for that.",
			"These variables are used to control over all container modes,",
			"such as configuring it to start a web app vs a background worker",
		),
	)

	return doc, nil
}
