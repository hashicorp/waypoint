package dockerpull

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/docker/cli/cli/config"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/registry"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	wpdocker "github.com/hashicorp/waypoint/builtin/docker"
	wpdockerclient "github.com/hashicorp/waypoint/builtin/docker/client"
	"github.com/hashicorp/waypoint/internal/assets"
	"github.com/hashicorp/waypoint/internal/pkg/epinject"
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
	// Image to pull
	Image string `hcl:"image,attr"`
	Tag   string `hcl:"tag,attr"`

	// Control whether or not to inject the entrypoint binary into the resulting image
	DisableCEB bool `hcl:"disable_entrypoint,optional"`

	// The docker specific encoded authentication string to use to talk to the registry.
	EncodedAuth string `hcl:"encoded_auth,optional"`
}

func (b *Builder) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&BuilderConfig{}), docs.FromFunc(b.BuildFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description(`
Use an existing, pre-built Docker image.

This builder will automatically inject the Waypoint entrypoint. You
can disable this with the "disable_entrypoint" configuration.

If you wish to rename or retag an image, use this along with the
"docker" registry option which will rename/retag the image and then
push it to the specified registry.

If Docker isn't available (the Docker daemon isn't running or a DOCKER_HOST
isn't set), a daemonless solution will be used instead.

If "disable_entrypoint" is set to true and the Waypoint configuration
has no registry, this builder will not physically pull the image. This enables
Waypoint to work in environments where the image is built outside of Waypoint
(such as in a CI pipeline).
`)

	doc.Example(`
build {
  use "docker-pull" {
    image = "gcr.io/my-project/my-image"
    tag   = "abcd1234"
  }
}
`)

	doc.Input("component.Source")
	doc.Output("docker.Image")

	doc.SetField(
		"image",
		"The image to pull.",
		docs.Summary(
			"This should NOT include the tag (the value following the ':' in a Docker image).",
			"Use `tag` to define the image tag.",
		),
	)

	doc.SetField(
		"tag",
		"The tag of the image to pull.",
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
		"encoded_auth",
		"the authentication information to log into the docker repository",
		docs.Summary(
			"WARNING: be very careful to not leak the authentication information",
			"by hardcoding it here. Use a helper function like `file()` to read",
			"the information from a file not stored in VCS",
		),
	)

	return doc, nil
}

// Config implements Configurable
func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

// We use the struct form of arguments so that we can access named
// values (such as "HasRegistry").
type buildArgs struct {
	argmapper.Struct

	Ctx         context.Context
	UI          terminal.UI
	Log         hclog.Logger
	HasRegistry bool
}

// Build
func (b *Builder) Build(args buildArgs) (*wpdocker.Image, error) {
	// Pull all the args out to top-level values. This is mostly done
	// cause the struct was added later, but also because these are very common.
	ctx := args.Ctx
	ui := args.UI
	log := args.Log

	sg := ui.StepGroup()
	defer sg.Wait()
	step := sg.Add("")
	defer func() {
		if step != nil {
			step.Abort()
		}
	}()

	result := &wpdocker.Image{
		Image:    b.config.Image,
		Tag:      b.config.Tag,
		Location: &wpdocker.Image_Docker{Docker: &empty.Empty{}},
	}

	// If we aren't injected the entrypoint AND we don't have a registry
	// defined, then we don't pull the image at all. We do this so that
	// Waypoint can work in an environment where Docker doesn't exist, img
	// doesn't work, and we're just using an image reference that was built
	// outside of Waypoint.
	if b.config.DisableCEB && !args.HasRegistry {
		step.Update("Using Docker image in remote registry: %s", result.Name())
		step.Done()

		result.Location = &wpdocker.Image_Registry{Registry: &empty.Empty{}}
		return result, nil
	}

	step.Update("Initializing Docker client...")
	cli, err := wpdockerclient.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "unable to create Docker client: %s", err)
	}
	cli.NegotiateAPIVersion(ctx)

	// We now test if Docker is actually functional. We do this here because we
	// need all of the above to complete the actual build.
	log.Debug("testing if we should use a Docker fallback")
	useImg := false
	if fallback, err := wpdockerclient.Fallback(ctx, log, cli); err != nil {
		log.Warn("error during check if we should use Docker fallback", "err", err)
		return nil, status.Errorf(codes.Internal,
			"error validating Docker connection: %s", err)
	} else if fallback && wpdocker.HasImg() {
		// If we're falling back and have "img" available, use that. If we
		// don't have "img" available, we continue to try to use Docker. We'll
		// fail but that error message should help the user.
		step.Update("Docker isn't available. Falling back to daemonless image pull...")
		step.Done()
		step = nil
		if err := b.buildWithImg(
			ctx, log, sg, result,
		); err != nil {
			return nil, err
		}

		// Our image is in the img registry now. We set this so that
		// future users of this result type know where to look.
		result.Location = &wpdocker.Image_Img{Img: &empty.Empty{}}

		// We set this to true so we use the img-based injector later
		useImg = true
	} else {
		// No fallback, build with Docker
		step.Done()
		step = nil
		if err := b.buildWithDocker(
			ctx, log, ui, sg, cli, result,
		); err != nil {
			return nil, err
		}
	}

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

		callback := func(cur []string) (*epinject.NewEntrypoint, error) {
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
		}

		if !useImg {
			_, err = epinject.AlterEntrypoint(ctx, result.Name(), callback)
		} else {
			_, err = epinject.AlterEntrypointImg(ctx, result.Name(), callback)
		}
		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to set modify Docker entrypoint: %s", err)
		}

		step.Done()
	}

	return result, nil
}

func (b *Builder) buildWithDocker(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	sg terminal.StepGroup,
	cli *client.Client,
	result *wpdocker.Image,
) error {
	ref, err := reference.ParseNormalizedNamed(result.Name())
	if err != nil {
		return status.Errorf(codes.Internal, "unable to parse image name: %s", err)
	}

	encodedAuth := b.config.EncodedAuth
	if encodedAuth == "" {
		// Resolve the Repository name from fqn to RepositoryInfo
		repoInfo, err := registry.ParseRepositoryInfo(ref)
		if err != nil {
			return status.Errorf(codes.Internal, "unable to parse repository info from image name: %s", err)
		}

		var server string

		if repoInfo.Index.Official {
			info, err := cli.Info(ctx)
			if err != nil || info.IndexServerAddress == "" {
				server = registry.IndexServer
			} else {
				server = info.IndexServerAddress
			}
		} else {
			server = repoInfo.Index.Name
		}

		var errBuf bytes.Buffer
		cf := config.LoadDefaultConfigFile(&errBuf)
		if errBuf.Len() > 0 {
			// NOTE(mitchellh): I don't know why we ignore this, but we always have.
			log.Warn("error loading Docker config file", "err", err)
		}

		authConfig, _ := cf.GetAuthConfig(server)
		buf, err := json.Marshal(authConfig)
		if err != nil {
			return status.Errorf(codes.Internal, "unable to generate authentication info for registry: %s", err)
		}
		encodedAuth = base64.URLEncoding.EncodeToString(buf)
	}

	step := sg.Add("Pulling image...")
	defer step.Abort()

	resp, err := cli.ImagePull(ctx, reference.FamiliarString(ref), types.ImagePullOptions{
		RegistryAuth: encodedAuth,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "error pulling image: %s", err)
	}
	defer resp.Close()

	stdout, _, err := ui.OutputWriters()
	if err != nil {
		return err
	}

	var termFd uintptr
	if f, ok := stdout.(*os.File); ok {
		termFd = f.Fd()
	}

	err = jsonmessage.DisplayJSONMessagesStream(resp, step.TermOutput(), termFd, true, nil)
	if err != nil {
		return status.Errorf(codes.Internal, "unable to stream pull logs to the terminal: %s", err)
	}

	step.Done()
	return nil
}

func (b *Builder) buildWithImg(
	ctx context.Context,
	log hclog.Logger,
	sg terminal.StepGroup,
	target *wpdocker.Image,
) error {
	var step terminal.Step
	defer func() {
		if step != nil {
			step.Abort()
		}
	}()

	step = sg.Add("Preparing Docker configuration...")
	env := os.Environ()
	if path, err := wpdocker.TempDockerConfig(log, target, b.config.EncodedAuth); err != nil {
		return err
	} else if path != "" {
		defer os.RemoveAll(path)
		env = append(env, "DOCKER_CONFIG="+path)
	}

	step.Done()
	step = sg.Add("Pulling Docker image with img...")

	// NOTE(mitchellh): we can probably use the img Go pkg directly one day.
	cmd := exec.CommandContext(ctx,
		"img",
		"pull",
		target.Name(),
	)
	cmd.Env = env
	cmd.Stdout = step.TermOutput()
	cmd.Stderr = cmd.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}

	step.Done()
	return nil
}
