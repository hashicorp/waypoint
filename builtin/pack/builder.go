package pack

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/buildpacks/pack"
	"github.com/buildpacks/pack/logging"
	"github.com/buildpacks/pack/project"
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

	// Files patterns to prevent from being pulled into the build.
	Ignore []string `hcl:"ignore,optional"`

	// Process type that will be used when setting container start command.
	ProcessType string `hcl:"process_type,optional" default:"web"`
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

	dockerClient, err := wpdockerclient.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	dockerClient.NegotiateAPIVersion(ctx)

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

	bo := pack.BuildOptions{
		Image:      src.App,
		Builder:    builder,
		AppPath:    src.Path,
		Env:        b.config.StaticEnvVars,
		Buildpacks: b.config.Buildpacks,
		ProjectDescriptor: project.Descriptor{
			Build: project.Build{
				Exclude: b.config.Ignore,
			},
		},
		DefaultProcessType: b.config.ProcessType,
	}

	err = client.Build(ctx, bo)
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

	doc.SetField(
		"ignore",
		"file patterns to match files which will not be included in the build",
		docs.Summary(
			`Each pattern follows the semantics of .gitignore. This is a summarized version:

1. A blank line matches no files, so it can serve as a separator
	 for readability.

2. A line starting with # serves as a comment. Put a backslash ("\")
	 in front of the first hash for patterns that begin with a hash.

3. Trailing spaces are ignored unless they are quoted with backslash ("\").

4. An optional prefix "!" which negates the pattern; any matching file
	 excluded by a previous pattern will become included again. It is not
	 possible to re-include a file if a parent directory of that file is
	 excluded. Git doesn’t list excluded directories for performance reasons,
	 so any patterns on contained files have no effect, no matter where they
	 are defined. Put a backslash ("\") in front of the first "!" for
	 patterns that begin with a literal "!", for example, "\!important!.txt".

5. If the pattern ends with a slash, it is removed for the purpose of the
	 following description, but it would only find a match with a directory.
	 In other words, foo/ will match a directory foo and paths underneath it,
	 but will not match a regular file or a symbolic link foo (this is
	 consistent with the way how pathspec works in general in Git).

6. If the pattern does not contain a slash /, Git treats it as a shell glob
	 pattern and checks for a match against the pathname relative to the
	 location of the .gitignore file (relative to the top level of the work
	 tree if not from a .gitignore file).

7. Otherwise, Git treats the pattern as a shell glob suitable for
	 consumption by fnmatch(3) with the FNM_PATHNAME flag: wildcards in the
	 pattern will not match a / in the pathname. For example,
	 "Documentation/*.html" matches "Documentation/git.html" but not
	 "Documentation/ppc/ppc.html" or "tools/perf/Documentation/perf.html".

8. A leading slash matches the beginning of the pathname. For example,
	 "/*.c" matches "cat-file.c" but not "mozilla-sha1/sha1.c".

9. Two consecutive asterisks ("**") in patterns matched against full
	 pathname may have special meaning:

		i.   A leading "**" followed by a slash means match in all directories.
				 For example, "** /foo" matches file or directory "foo" anywhere,
				 the same as pattern "foo". "** /foo/bar" matches file or directory
				 "bar" anywhere that is directly under directory "foo".

		ii.  A trailing "/**" matches everything inside. For example, "abc/**"
				 matches all files inside directory "abc", relative to the location
				 of the .gitignore file, with infinite depth.

		iii. A slash followed by two consecutive asterisks then a slash matches
				 zero or more directories. For example, "a/** /b" matches "a/b",
				 "a/x/b", "a/x/y/b" and so on.

		iv.  Other consecutive asterisks are considered invalid.`),
	)

	doc.SetField(
		"process_type",
		"The process type to use from your Procfile. if not set, defaults to `web`",
		docs.Summary(
			"The process type is used to control over all container modes,",
			"such as configuring it to start a web app vs a background worker",
		),
	)

	return doc, nil
}
