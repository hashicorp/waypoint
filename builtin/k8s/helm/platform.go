package helm

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Platform is the Platform implementation
type Platform struct {
	config Config
}

// Config implements Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

// DeployFunc implements component.Platform
func (p *Platform) DeployFunc() interface{} {
	return p.Deploy
}

// DestroyFunc implements component.Destroyer
func (p *Platform) DestroyFunc() interface{} {
	return p.Destroy
}

// GenerationFunc implements component.Generation
func (p *Platform) GenerationFunc() interface{} {
	return p.Generation
}

// Deploy deploys to Kubernetes
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {
	sg := ui.StepGroup()
	defer sg.Wait()
	s := sg.Add("")
	defer func() { s.Abort() }()

	s.Update("Initializing Helm...")
	actionConfig, err := p.actionInit(log)
	if err != nil {
		return nil, err
	}
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	s.Done()

	return deployment, nil
}

// Destroy
func (p *Platform) Destroy(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) error {
	sg := ui.StepGroup()
	defer sg.Wait()
	s := sg.Add("Executing kubectl to destroy...")
	defer func() { s.Abort() }()

	if deployment.PruneLabel == "" {
		s.Update("No prune label on deployment. Not destroying.")
		s.Done()
		return nil
	}

	// Apply it
	cmd, err := p.cmd(ctx, s, "delete", "all", "-l", deployment.PruneLabel)
	if err != nil {
		return err
	}
	if err := cmd.Run(); err != nil {
		return err
	}

	s.Done()
	return nil
}

// Generation returns the generation ID.
func (p *Platform) Generation(
	ctx context.Context,
) ([]byte, error) {
	// Static generation since we will always use the `prune_label` to
	// automatically delete unused resources.
	return []byte("kubernetes-apply"), nil
}

// Config is the configuration structure for the Platform.
//
// For docs on the fields, please see the Documentation function.
type Config struct {
	Name       string   `hcl:"name,attr"`
	Repository string   `hcl:"repository,optional"`
	Chart      string   `hcl:"chart,attr"`
	Version    string   `hcl:"version,optional"`
	Values     []string `hcl:"values,optional"`
	Set        []*struct {
		Name  string `hcl:"name,attr"`
		Value string `hcl:"value,attr"`
		Type  string `hcl:"type,optional"`
	} `hcl:"set,block`
	Driver string `hcl:"driver,optional"`

	KubeconfigPath string `hcl:"kubeconfig,optional"`
	Context        string `hcl:"context,optional"`
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}), docs.FromFunc(p.DeployFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description(`
Deploy to Kubernetes from a Helm chart. The Helm chart can be a local path
or a chart in a repository.

### Entrypoint Functionality

Waypoint [entrypoint functionality](/docs/entrypoint#functionality) such
as logs, exec, app configuration, and more require two properties to be true:

1. The running image must already have the Waypoint entrypoint installed
  and configured as the entrypoint. This should happen in the build stage.

2. Proper environment variables must be set so the entrypoint knows how
  to communicate to the Waypoint server. **This step happens in this
  deployment stage.**

**Step 2 does not happen automatically.** You must manually set the entrypoint
environment variables using the [templating feature](/docs/waypoint-hcl/functions/template).
These must be passed in using Helm values (i.e. the chart must make
environment variables configurable).

#### URL Service

If you want your workload to be accessible by the
[Waypoint URL service](/docs/url), you must set the PORT environment variable
within the pod with your web service and also be using the Waypoint
entrypoint (documented in the previous section).

The PORT environment variable should be the port that your web service
is listening on that the URL service will connect to. See one of the examples
below for more details.

`)

	doc.Example(`
// A local helm chart relative to the app.
deploy {
  use "helm" {
    chart = "${path.app}/chart"
  }
}
`)

	doc.SetField(
		"name",
		"Name of the Helm release.",
		docs.Summary(
			"This must be globally unique within the context of your Helm installation.",
		),
	)

	doc.SetField(
		"repository",
		"URL of the Helm repository that contains the chart.",
		docs.Summary(
			"This only needs to be set if you're NOT using a local chart.",
		),
	)

	doc.SetField(
		"chart",
		"The name or path of the chart to install.",
		docs.Summary(
			"If you're installing a local chart, this is the path to the chart.",
			"If you're installing a chart from a repository (have the `repository`",
			"configuration set), then this is the name of the chart in the repository.",
		),
	)

	doc.SetField(
		"version",
		"The version of the chart to install.",
	)

	doc.SetField(
		"values",
		"Values in raw YAML to configure the Helm chart.",
		docs.Summary(
			"These values are usually loaded from files using HCL functions such as",
			"`file` or templating with `templatefile`. Multiple values will be merged",
			"using the same logic as the `-f` flag with Helm.",
		),
	)

	doc.SetField(
		"set",
		"A single value to set. This can be repeated multiple times.",
		docs.Summary(
			"This sets a single value. Separate nested values with a `.`. This is",
			"the same as the `--set` flag on `helm install`.",
		),
	)

	doc.SetField(
		"driver",
		"The Helm storage driver to use.",
		docs.Summary(
			"This can be one of `configmap`, `secret`, `memory`, or `sql`.",
			"The SQL connection string can not be set currently so this must",
			"be set on the runners.",
		),
		docs.Default("secret"),
	)

	doc.SetField(
		"kubeconfig",
		"Path to the kubeconfig file to use.",
		docs.Summary(
			"If this isn't set, the default lookup used by `kubectl` will be used.",
		),
		docs.EnvVar("KUBECONFIG"),
	)

	doc.SetField(
		"context",
		"The kubectl context to use, as defined in the kubeconfig file.",
	)

	return doc, nil
}

var (
	_ component.Generation   = (*Platform)(nil)
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
	_ component.Destroyer    = (*Platform)(nil)
)
