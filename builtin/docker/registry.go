package docker

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
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	wpdockerclient "github.com/hashicorp/waypoint/builtin/docker/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	log hclog.Logger,
) (*Image, error) {
	stdout, _, err := ui.OutputWriters()
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "unable to create output for logs:%s", err)
	}

	sg := ui.StepGroup()
	step := sg.Add("Initializing Docker client...")
	defer func() { step.Abort() }()

	cli, err := wpdockerclient.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "unable to create Docker client:%s", err)
	}
	cli.NegotiateAPIVersion(ctx)

	step.Update("Tagging Docker image: %s => %s:%s", img.Name(), r.config.Image, r.config.Tag)

	target := &Image{Image: r.config.Image, Tag: r.config.Tag}
	err = cli.ImageTag(ctx, img.Name(), target.Name())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to tag image:%s", err)
	}

	step.Done()

	if r.config.Local {
		return target, nil
	}

	ref, err := reference.ParseNormalizedNamed(target.Name())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to parse image name: %s", err)
	}

	encodedAuth := r.config.EncodedAuth
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

	step = sg.Add("Pushing Docker image...")

	options := types.ImagePushOptions{
		RegistryAuth: encodedAuth,
	}

	responseBody, err := cli.ImagePush(ctx, reference.FamiliarString(ref), options)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to push image to registry: %s", err)
	}

	defer responseBody.Close()

	var termFd uintptr
	if f, ok := stdout.(*os.File); ok {
		termFd = f.Fd()
	}

	err = jsonmessage.DisplayJSONMessagesStream(responseBody, step.TermOutput(), termFd, true, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to stream Docker logs to terminal: %s", err)
	}

	step.Done()

	step = sg.Add("Docker image pushed: %s:%s", r.config.Image, r.config.Tag)
	step.Done()

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

func (r *Registry) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}), docs.FromFunc(r.PushFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description("Push a Docker image to a Docker compatible registry")

	doc.Example(`
build {
  use "docker" {}
  registry {
    use "docker" {
      image = "hashicorp/http-echo"
      tag   = "latest"
    }
  }
}
`)

	doc.Input("docker.Image")
	doc.Output("docker.Image")

	doc.SetField(
		"image",
		"the image to push the local image to, fully qualified",
		docs.Summary(
			"this value must be the fully qualified name to the image.",
			"for example: gcr.io/waypoint-demo/demo",
		),
	)

	doc.SetField(
		"tag",
		"the tag for the new image",
		docs.Summary(
			"this is added to image to provide the full image reference",
		),
	)

	doc.SetField(
		"local",
		"if set, the image will only be tagged locally and not pushed to a remote repository",
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
