package dockerpull

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/docker/cli/cli/config"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/registry"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
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
Use an existing, pre-built Docker image

This builder will automatically inject the Waypoint entrypoint. You
can disable this with the "disable_entrypoint" configuration.

If you wish to rename or retag an image, use this along with the
"docker" registry option which will rename/retag the image and then
push it to the specified registry.
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
			"Use `tag` to define the  image tag.",
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

// Build
func (b *Builder) Build(
	ctx context.Context,
	ui terminal.UI,
	src *component.Source,
	log hclog.Logger,
) (*wpdocker.Image, error) {
	stdout, _, err := ui.OutputWriters()
	if err != nil {
		return nil, err
	}

	sg := ui.StepGroup()
	step := sg.Add("Initializing Docker client...")
	defer step.Abort()

	result := &wpdocker.Image{
		Image: b.config.Image,
		Tag:   b.config.Tag,
	}

	cli, err := wpdockerclient.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "unable to create Docker client: %s", err)
	}
	cli.NegotiateAPIVersion(ctx)

	step.Done()

	ref, err := reference.ParseNormalizedNamed(result.Name())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to parse image name: %s", err)
	}

	encodedAuth := b.config.EncodedAuth
	if encodedAuth == "" {
		// Resolve the Repository name from fqn to RepositoryInfo
		repoInfo, err := registry.ParseRepositoryInfo(ref)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to parse repository info from image name: %s", err)
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
			return nil, status.Errorf(codes.Internal, "unable to generate authentication info for registry: %s", err)
		}
		encodedAuth = base64.URLEncoding.EncodeToString(buf)
	}

	step = sg.Add("Pulling image...")

	resp, err := cli.ImagePull(ctx, reference.FamiliarString(ref), types.ImagePullOptions{
		RegistryAuth: encodedAuth,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error pulling image: %s", err)
	}
	defer resp.Close()

	var termFd uintptr
	if f, ok := stdout.(*os.File); ok {
		termFd = f.Fd()
	}

	err = jsonmessage.DisplayJSONMessagesStream(resp, step.TermOutput(), termFd, true, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to stream pull logs to the terminal: %s", err)
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
			return nil, status.Errorf(codes.Internal, "unable to modify Docker entrypoint: %s", err)
		}

		step.Done()
	}

	return result, nil
}
