package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
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

func (r *Registry) AccessInfoFunc() interface{} {
	return r.AccessInfo
}

// Push pushes an image to the registry.
func (r *Registry) AccessInfo() (*AccessInfo, error) {
	ai := &AccessInfo{
		Image:    r.config.Image,
		Tag:      r.config.Tag,
		Insecure: r.config.Insecure,
	}

	if r.config.EncodedAuth != "" {
		ai.Auth = &AccessInfo_Encoded{
			Encoded: r.config.EncodedAuth,
		}
	} else if r.config.Auth != nil {
		auth, err := json.Marshal(types.AuthConfig{
			Username:      r.config.Auth.Username,
			Password:      r.config.Auth.Password,
			Email:         r.config.Auth.Email,
			Auth:          r.config.Auth.Auth,
			ServerAddress: r.config.Auth.ServerAddress,
			IdentityToken: r.config.Auth.IdentityToken,
			RegistryToken: r.config.Auth.RegistryToken,
		})
		if err != nil {
			return ai, status.Errorf(codes.Internal, "failed to marshal auth info to json: %s", err)
		}
		ai.Auth = &AccessInfo_Encoded{
			Encoded: base64.URLEncoding.EncodeToString(auth),
		}
	} else if r.config.Password != "" {
		ai.Auth = &AccessInfo_UserPass_{
			UserPass: &AccessInfo_UserPass{
				Username: r.config.Username,
				Password: r.config.Password,
			},
		}
	}
	return ai, nil
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
		Image:        r.config.Image,
		Tag:          r.config.Tag,
		Architecture: img.Architecture,
	}
	auth := r.config.Auth
	if !r.config.Local {
		target.Location = &Image_Registry{Registry: &Image_RegistryLocation{}}
	}

	// Depending on whethere the image is, we diverge at this point.
	switch sv := img.Location.(type) {
	case *Image_Registry:
		if !sv.Registry.WaypointGenerated && (img.Image != r.config.Image || img.Tag != r.config.Tag) {
			return nil, status.Errorf(codes.FailedPrecondition,
				"Input image is not pulled locally and therefore can't be pushed. "+
					"Please pull the image or use a builder that pulls the image first.")
		}

		// This indicates that the builder used the AccessInfo and published the image
		// directly. Ergo we don't need to do anything and can just return the image as is.
		return img, nil

	case *Image_Docker, nil:
		// We support "nil" here for backwards compatibility. Images built
		// prior to supporting the Location field will set nil.
		if err := r.pushWithDocker(
			ctx,
			log,
			ui,
			img,
			target,
			auth,
		); err != nil {
			return nil, err
		}

	case *Image_UnusedImg:
		return nil, status.Errorf(codes.FailedPrecondition,
			"Input image is in `img` but Waypoint doesn't support img from "+
				"version 0.7 onwards. Please use a Waypoint 0.6 runner to complete "+
				"this job or rerun the build without img.")

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

	// Authenticates to private registry
	Auth *Auth `hcl:"auth,block"`

	// The docker specific encoded authentication string to use to talk to the registry.
	EncodedAuth string `hcl:"encoded_auth,optional"`

	// Insecure indicates if the registry should be accessed via http rather than https
	Insecure bool `hcl:"insecure,optional"`

	// Username is the username to use for authentication on the registry.
	Username string `hcl:"username,optional"`

	// Password is the authentication information assocated with username.
	Password string `hcl:"password,optional"`
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

	doc.SetField(
		"insecure",
		"access the registry via http rather than https",
		docs.Summary(
			"This indicates that the registry should be accessed via http rather than https.",
			"Not recommended for production usage.",
		),
	)

	doc.SetField(
		"username",
		"username to authenticate with the registry",
		docs.Summary(
			"This optional conflicts with encoded_auth and thusly only one can be used at a time.",
			"If both are used, encoded_auth takes precedence.",
		),
	)

	doc.SetField(
		"password",
		"password associated with username on the registry",
		docs.Summary(
			"This optional conflicts with encoded_auth and thusly only one can be used at a time.",
			"If both are used, encoded_auth takes precedence.",
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

var _ component.Registry = (*Registry)(nil)
var _ component.RegistryAccess = (*Registry)(nil)
