package docker

import (
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
	"github.com/hashicorp/waypoint/sdk/terminal"
	"github.com/mattn/go-isatty"
)

// Registry represents access to a Docker registry.
type Registry struct {
	config Config
}

// Config implements Configurable
func (r *Registry) Config() (interface{}, error) {
	return &r.config, nil
}

// PushFunc implements component.Registry
func (r *Registry) PushFunc() interface{} {
	return r.Push
}

// Push pushes an image to the registry.
func (r *Registry) Push(
	ctx context.Context,
	img *Image,
	ui terminal.UI,
) (*Image, error) {
	ui.Output("Taging Docker image: %s => %s:%s", img.Name(), r.config.Image, r.config.Tag)

	stdout, stderr, err := ui.OutputWriters()
	if err != nil {
		return nil, err
	}

	target := &Image{Image: r.config.Image, Tag: r.config.Tag}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	cli.NegotiateAPIVersion(ctx)

	err = cli.ImageTag(ctx, img.Name(), target.Name())
	if err != nil {
		return nil, err
	}

	if r.config.Local {
		return target, nil
	}

	ref, err := reference.ParseNormalizedNamed(target.Name())
	if err != nil {
		return nil, err
	}

	encodedAuth := r.config.EncodedAuth

	if encodedAuth == "" {
		// Resolve the Repository name from fqn to RepositoryInfo
		repoInfo, err := registry.ParseRepositoryInfo(ref)
		if err != nil {
			return nil, err
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

		cf := config.LoadDefaultConfigFile(stderr)

		authConfig, _ := cf.GetAuthConfig(server)
		buf, err := json.Marshal(authConfig)
		if err != nil {
			return nil, err
		}
		encodedAuth = base64.URLEncoding.EncodeToString(buf)
	}

	options := types.ImagePushOptions{
		RegistryAuth: encodedAuth,
	}

	responseBody, err := cli.ImagePush(ctx, reference.FamiliarString(ref), options)
	if err != nil {
		return nil, err
	}

	defer responseBody.Close()

	var (
		termFd uintptr
		isTerm bool
	)

	if f, ok := stdout.(*os.File); ok {
		termFd = f.Fd()
		isTerm = isatty.IsTerminal(termFd)
	}

	err = jsonmessage.DisplayJSONMessagesStream(responseBody, stdout, termFd, isTerm, nil)
	if err != nil {
		return nil, err
	}

	ui.Output("Docker image pushed: %s:%s", r.config.Image, r.config.Tag)

	return target, nil
}

// Config is the configuration structure for the registry.
type Config struct {
	// Image is the name of the image plus tag that the image will be pushed as.
	Image string `hcl:"image,attr"`

	// Tag is the tag to apply to the image.
	Tag string `hcl:"tag,attr"`

	// Local if true will not push this image to a remote registry.
	Local bool `hcl:"local,optional"`

	// The docker specific encoded authentication string to use to talk to the registry.
	EncodedAuth string `hcl:"encoded_auth,optional"`
}
