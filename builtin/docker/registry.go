package docker

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
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
	target := &Image{
		Image: r.config.Image,
		Tag:   r.config.Tag,
	}
	if !r.config.Local {
		target.Location = &Image_Registry{Registry: &empty.Empty{}}
	}

	// Depending on whethere the image is, we diverge at this point.
	switch img.Location.(type) {
	case *Image_Registry:
		// We can't push an image that isn't pulled locally in some form.
		return nil, status.Errorf(codes.FailedPrecondition,
			"Input image is not pulled locally and therefore can't be pushed. "+
				"Please pull the image or use a builder that pulls the image first.")

	case *Image_Img:
		// If the image is already in img, we have to use `img push`.
		if err := r.pushWithImg(
			ctx,
			log,
			ui,
			img,
			target,
		); err != nil {
			return nil, err
		}

	case *Image_Docker, nil:
		// We support "nil" here for backwards compatibility. Images built
		// prior to supporting the Location field will set nil.
		if err := r.pushWithDocker(
			ctx,
			log,
			ui,
			img,
			target,
		); err != nil {
			return nil, err
		}
	}

	sg := ui.StepGroup()
	step := sg.Add("Docker image pushed: %s:%s", r.config.Image, r.config.Tag)
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
			"The format of this is base64-encoded JSON. The structure is the ",
			"[`AuthConfig`](https://pkg.go.dev/github.com/docker/cli/cli/config/types#AuthConfig)",
			"structure used by Docker.",
			"",
			"WARNING: be very careful to not leak the authentication information",
			"by hardcoding it here. Use a helper function like `file()` to read",
			"the information from a file not stored in VCS",
		),
	)

	return doc, nil
}
