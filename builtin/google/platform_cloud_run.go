package google

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/mitchellh/devflow/builtin/docker"
	"github.com/mitchellh/devflow/internal/terraform"
	"github.com/mitchellh/devflow/sdk"
	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/datadir"
)

// CloudRunOptions are the SDK options to use for instantiation for
// the Google Cloud Run plugin.
var CloudRunOptions = []sdk.Option{
	sdk.WithComponents(&CloudRunPlatform{}),
}

// CloudRunPlatform is the Platform implementation for Google Cloud Run.
type CloudRunPlatform struct {
	config Config
}

// Config implements Configurable
func (p *CloudRunPlatform) Config() (interface{}, error) {
	return &p.config, nil
}

// DeployFunc implements component.Platform
func (p *CloudRunPlatform) DeployFunc() interface{} {
	return p.Deploy
}

// Deploy deploys an image to GCR.
func (p *CloudRunPlatform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	dir *datadir.Component,
) (*CloudRunDeployment, error) {
	// We need to create a scoped directory so that our Terraform run
	// happens in isolation of our data.
	log.Debug("preparing scoped directory for Terraform")
	tfDir, err := datadir.NewScopedDir(dir, "terraform")
	if err != nil {
		return nil, err
	}

	// Build our Terraform run
	tf := &terraform.Terraform{
		Context:    ctx,
		Logger:     log,
		Dir:        tfDir,
		ConfigFS:   AssetFile(),
		ConfigPath: "terraform-cloud-run-0",
		Vars: map[string]interface{}{
			"name":    src.App,
			"project": p.config.Project,
			"image":   img.Name(),
		},
	}

	// Apply!
	outputs, err := tf.Apply()
	if err != nil {
		return nil, err
	}

	return &CloudRunDeployment{
		Url: outputs["url"].(string),
	}, nil
}

// Config is the configuration structure for the CloudRunPlatform.
type Config struct {
	// Project is the project to deploy to.
	Project string `hcl:"project,attr"`

	// Unauthenticated, if set to true, will allow unauthenticated access
	// to your deployment. This defaults to true.
	Unauthenticated *bool `hcl:"unauthenticated,optional"`
}

// MarshalText implements encoding.TextMarshaler so that protobuf generates
// the correct string version.
func (d *CloudRunDeployment) MarshalText() ([]byte, error) {
	return []byte("URL: " + d.Url), nil
}

var (
	_ component.Platform     = (*CloudRunPlatform)(nil)
	_ component.Configurable = (*CloudRunPlatform)(nil)
	_ component.Deployment   = (*CloudRunDeployment)(nil)
)
