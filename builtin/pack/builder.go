package pack

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/buildpacks/pack"
	"github.com/buildpacks/pack/logging"
	"github.com/buildpacks/pack/project"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/builtin/docker"
	wpdockerclient "github.com/hashicorp/waypoint/builtin/docker/client"
	"github.com/hashicorp/waypoint/internal/assets"
	"github.com/hashicorp/waypoint/internal/pkg/epinject"
	"github.com/hashicorp/waypoint/internal/pkg/epinject/ociregistry"
)

const (
	// This is legacy and comes from heroku, which cnb continued with.
	DefaultProcessType = "web"
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

// BuildFunc implements component.BuilderODR
func (b *Builder) BuildODRFunc() interface{} {
	return b.BuildODR
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

const DefaultBuilder = "heroku/buildpacks:20"

// Config implements Configurable
func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

var skipBuildPacks = map[string]struct{}{
	"heroku/procfile": {},
}

// Build
func (b *Builder) BuildODR(
	ctx context.Context,
	ui terminal.UI,
	jobInfo *component.JobInfo,
	src *component.Source,
	log hclog.Logger,
	ai *docker.AccessInfo,
) (*DockerImage, error) {
	sg := ui.StepGroup()

	builder := b.config.Builder
	if builder == "" {
		builder = DefaultBuilder
	}

	log.Info("executing the ODR version of pack")

	// We don't even need to inspect Docker to verify we have the image.
	// If `pack` succeeded we can assume that it created an image for us.
	// return &DockerImage{
	// Image: ai.Image,
	// Tag:   ai.Tag,
	// }, nil

	step := sg.Add("Building Buildpack with kaniko...")
	defer func() {
		if step != nil {
			step.Abort()
		}
	}()

	target := &docker.Image{
		Image: ai.Image,
		Tag:   ai.Tag,
	}

	var ocis ociregistry.Server

	if ai.Auth != nil {
		switch sv := ai.Auth.(type) {
		case *docker.AccessInfo_Encoded:
			user, pass, err := docker.CredentialsFromConfig(sv.Encoded)
			if err != nil {
				return nil, err
			}
			ocis.AuthConfig.Username = user
			ocis.AuthConfig.Password = pass
		case *docker.AccessInfo_Header:
			ocis.AuthConfig.Auth = sv.Header
		case *docker.AccessInfo_UserPass_:
			ocis.AuthConfig.Username = sv.UserPass.Username
			ocis.AuthConfig.Password = sv.UserPass.Password
		}
	}

	// Determine the host that we're setting auth for. We have to parse the
	// image for this cause it may not contain a host. Luckily Docker has
	// libs to normalize this all for us.
	log.Trace("determining host for auth configuration", "image", target.Name())
	ref, err := reference.ParseNormalizedNamed(target.Image)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to parse image name: %s", err)
	}
	host := reference.Domain(ref)
	if host == "docker.io" {
		// The normalized name parse above will turn short names like "foo/bar"
		// into "docker.io/foo/bar" but the actual registry host for these
		// is "index.docker.io".
		host = "index.docker.io"
	}
	log.Trace("auth host", "host", host)

	ocis.DisableEntrypoint = b.config.DisableCEB
	ocis.Logger = log

	if ai.Insecure {
		ocis.Upstream = "http://" + host
	} else {
		ocis.Upstream = "https://" + host
	}

	err = ocis.Negotiate(ref.Name())
	if err != nil {
		return nil, errors.Wrapf(err, "unable to negotiate with upstream")
	}

	refPath := reference.Path(ref)

	if !b.config.DisableCEB {
		data, err := assets.Asset("ceb/ceb")
		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to restore custom entry point binary: %s", err)
		}

		step.Done()
		step = sg.Add("Testing registry and uploading entrypoint layer")

		err = ocis.SetupEntrypointLayer(refPath, data)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error setting up entrypoint layer to host: %q, err: %s", ocis.Upstream, err)
		}
	}

	li, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}

	defer li.Close()
	go http.Serve(li, &ocis)

	port := li.Addr().(*net.TCPAddr).Port

	localRef := fmt.Sprintf("localhost:%d/%s:%s", port, refPath, ai.Tag)

	// The patterns are the same os docker's patterns, so we just populate .dockerignore
	// which we know kaniko will honor.
	if len(b.config.Ignore) > 0 {

		// We open this in append mode to preserve what's there in before we add entries to it.
		f, err := os.OpenFile(filepath.Join(src.Path, ".dockerignore"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, status.Errorf(codes.Unavailable, "unable to write dockerignore: %s", err)
		}

		fmt.Fprintln(f, "")

		for _, pattern := range b.config.Ignore {
			fmt.Fprintln(f, pattern)
		}

		f.Close()
	}

	f, err := os.Create("Dockerfile.kaniko")
	if err != nil {
		return nil, err
	}

	processType := b.config.ProcessType
	if processType == "" {
		processType = DefaultProcessType
	}

	// TODO: buildpacks. They require us downloading data remotely and writing out /cnb/order.toml
	// to reference which buildpacks should be execute in which order.
	if len(b.config.Buildpacks) != 0 {
		return nil, status.Errorf(codes.Unavailable, "explicit buildpacks are not yet implemented")
	}

	fmt.Fprintf(f, `FROM %s

ADD --chown=1000 . /app

WORKDIR /app

USER 1000

RUN mkdir /tmp/cache /tmp/layers && \
	mkdir -p /app/bin && \
		/cnb/lifecycle/creator \
			"-app=/app" \
      "-cache-dir=/tmp/cache" \
      "-gid=1000" \
      "-layers=/tmp/layers" \
      "-platform=/platform" \
      "-previous-image=%s" \
      "-uid=1000" \
			"-process-type=%s" \
      "%s"

`, builder, localRef, processType, localRef)

	err = f.Close()
	if err != nil {
		return nil, err
	}

	dir := src.Path

	step.Update("Repository is available and ready: %s:%s", ai.Image, ai.Tag)
	step.Done()
	step = sg.Add("Executing kaniko...")

	// Start constructing our arg string for img
	args := []string{
		"/kaniko/executor",
		"-f", "Dockerfile.kaniko",
		"--no-push",
		"--context=dir:///" + dir,
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	// Command output should go to the step
	cmd.Stdout = step.TermOutput()
	cmd.Stderr = cmd.Stdout

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	step.Done()

	step = sg.Add("Image pushed to '%s:%s'", ai.Image, ai.Tag)
	step.Done()

	return &DockerImage{
		Image:  ai.Image,
		Tag:    ai.Tag,
		Remote: true,
	}, nil
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

	// We need to test if we're running in arm64 for the Docker server.
	// Buildpacks have issues with arm64: https://github.com/buildpacks/pack/issues/907
	// We just do a warning in case buildpacks support arm64 and magically
	// work later.
	serverInfo, err := dockerClient.Info(ctx)
	if err != nil {
		return nil, err
	}
	if serverInfo.Architecture != "amd64" && serverInfo.Architecture != "x86_64" {
		ui.Output(
			"Warning! Buildpacks are known to have issues on architectures "+
				"other than amd64. The architecure being reported by the Docker "+
				"server is %q. We will still attempt to build the image, but "+
				"may run into issues.",
			serverInfo.Architecture,
			terminal.WithWarningStyle(),
		)
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

		// Use the server architecture to determine our entrypoint architecture.
		assetName, ok := assets.CEBArch[strings.ToLower(serverInfo.Architecture)]
		if !ok {
			return nil, status.Errorf(codes.FailedPrecondition,
				"Automatic Waypoint entrypoint injection only supports amd64 and arm64 "+
					"image architectures. Got: %s", serverInfo.Architecture)
		}

		asset, err := assets.Asset(assetName)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to restore custom entry point binary: %s", err)
		}

		assetInfo, err := assets.AssetInfo(assetName)
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

**This plugin must either be run via Docker or inside an ondemand runner**.
`)

	doc.Example(`
build {
  use "pack" {
	builder     = "heroku/buildpacks:20"
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
