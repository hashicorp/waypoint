// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package dockerpull

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/docker/cli/cli/config"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/registry"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"os"
	"os/exec"

	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"
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

// BuildFunc implements component.BuilderODR
func (b *Builder) BuildODRFunc() interface{} {
	return b.BuildODR
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

	// Authenticates to private registry
	Auth *docker.Auth `hcl:"auth,block"`
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

	doc.SetField(
		"auth",
		"the authentication information to log into the docker repository",
		docs.SubFields(func(d *docs.SubFieldDoc) {
			d.SetField("hostname", "Hostname of Docker registry")
			d.SetField("username", "Username of Docker registry account")
			d.SetField("password", "Password of Docker registry account")
			d.SetField("serverAddress", "Address of Docker registry")
			d.SetField("identityToken", "Token used to authenticate user")
			d.SetField("registryToken", "Bearer tokens to be sent to Docker registry")
		}),
	)

	return doc, nil
}

// Config implements Configurable
func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

// We use the struct form of arguments so that we can access named
// values (such as "HasRegistry").
type BuildArgs struct {
	argmapper.Struct

	Ctx         context.Context
	UI          terminal.UI
	Log         hclog.Logger
	HasRegistry bool
}

// Build
func (b *Builder) BuildODR(
	ctx context.Context,
	ui terminal.UI,
	src *component.Source,
	log hclog.Logger,
	ai *wpdocker.AccessInfo,
) (*wpdocker.Image, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	result, err := b.pullWithKaniko(ctx, ui, sg, log, ai)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Build
func (b *Builder) Build(args BuildArgs) (*wpdocker.Image, error) {
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

	auth := b.config.Auth

	// If we aren't injected the entrypoint AND we don't have a registry
	// defined, then we don't pull the image at all. We do this so that
	// Waypoint can work in an environment where Docker doesn't exist, img
	// doesn't work, and we're just using an image reference that was built
	// outside of Waypoint.
	if b.config.DisableCEB && !args.HasRegistry {
		step.Update("Using Docker image in remote registry: %s", result.Name())
		step.Done()

		result.Location = &wpdocker.Image_Registry{Registry: &wpdocker.Image_RegistryLocation{}}
		return result, nil
	}

	step.Update("Initializing Docker client...")
	cli, err := wpdockerclient.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "unable to create Docker client: %s", err)
	}
	cli.NegotiateAPIVersion(ctx)

	// Build
	step.Done()
	step = nil
	if err := b.buildWithDocker(
		ctx, log, ui, sg, cli, result, auth,
	); err != nil {
		return nil, err
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

		_, err = epinject.AlterEntrypoint(ctx, result.Name(), callback)
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
	authConfig *docker.Auth,
) error {
	ref, err := reference.ParseNormalizedNamed(result.Name())
	if err != nil {
		return status.Errorf(codes.Internal, "unable to parse image name: %s", err)
	}

	var encodedAuth = ""

	if b.config.EncodedAuth != "" {
		//If EncodedAuth is set, use that
		encodedAuth = b.config.EncodedAuth
	} else if b.config.Auth == nil && b.config.EncodedAuth == "" {
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
	} else if *b.config.Auth != (docker.Auth{}) {
		//If EncodedAuth is not set, and Auth is, use Auth
		authBytes, err := json.Marshal(types.AuthConfig{
			Username:      authConfig.Username,
			Password:      authConfig.Password,
			Email:         authConfig.Email,
			Auth:          authConfig.Auth,
			ServerAddress: authConfig.ServerAddress,
			IdentityToken: authConfig.IdentityToken,
			RegistryToken: authConfig.RegistryToken,
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to marshal auth info to json: %s", err)
		}
		encodedAuth = base64.URLEncoding.EncodeToString(authBytes)
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
	authConfig *docker.Auth,
) error {
	var step terminal.Step
	defer func() {
		if step != nil {
			step.Abort()
		}
	}()

	step = sg.Add("Preparing Docker configuration...")
	env := os.Environ()

	//Check if auth configuration is not null
	if (*b.config.Auth != docker.Auth{}) {
		auth, err := json.Marshal(types.AuthConfig{
			Username:      authConfig.Username,
			Password:      authConfig.Password,
			Email:         authConfig.Email,
			Auth:          authConfig.Auth,
			ServerAddress: authConfig.ServerAddress,
			IdentityToken: authConfig.IdentityToken,
			RegistryToken: authConfig.RegistryToken,
		})
		if err != nil {
			return status.Errorf(codes.Internal, "failed to marshal auth info to json: %s", err)
		}
		if path, err := wpdocker.TempDockerConfig(log, target, base64.URLEncoding.EncodeToString(auth)); err != nil {
			return err
		} else if path != "" {
			defer os.RemoveAll(path)
			env = append(env, "DOCKER_CONFIG="+path)
		}
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
