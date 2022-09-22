package helm

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"helm.sh/helm/v3/pkg/action"
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
	settings, err := p.settingsInit()
	if err != nil {
		return nil, err
	}

	actionConfig, err := p.actionInit(log)
	if err != nil {
		return nil, err
	}
	s.Done()

	// We need to look up the previous release if it exists because
	// if it does then we are upgrading.
	s = sg.Add("Checking for previous release...")
	prevRel, err := getRelease(actionConfig, p.config.Name)
	if err != nil {
		return nil, err
	}

	s.Update("Loading Helm chart...")
	cpo, chartName, err := p.chartPathOptions()
	if err != nil {
		return nil, err
	}

	c, _, err := getChart(chartName, cpo, settings)
	if err != nil {
		return nil, err
	}
	s.Update("Loaded Chart: %s (version: %s)", c.Metadata.Name, c.Metadata.Version)
	s.Done()
	s = sg.Add("")

	// Parse our values
	values, err := p.chartValues()
	if err != nil {
		return nil, err
	}

	if p.config.Namespace == "" {
		// default the namespace to "default"
		p.config.Namespace = "default"
	}

	// From here on out, we will always return a partial deployment if we error.
	result := &Deployment{Release: p.config.Name}

	// If we have no previous release, install.
	if prevRel == nil {
		// Initialize our installation settings. These defaults are safe defaults
		// and are mostly taken from the Terraform provider.
		client := action.NewInstall(actionConfig)
		client.ChartPathOptions = *cpo
		client.ClientOnly = false
		client.DryRun = false
		client.DisableHooks = false
		client.Wait = true
		client.WaitForJobs = false
		client.Devel = p.config.Devel
		client.DependencyUpdate = false
		client.Timeout = 300 * time.Second
		client.Namespace = p.config.Namespace
		client.ReleaseName = p.config.Name
		client.GenerateName = false
		client.NameTemplate = ""
		client.OutputDir = ""
		client.Atomic = false
		client.SkipCRDs = false
		client.SubNotes = true
		client.DisableOpenAPIValidation = false
		client.Replace = false
		client.Description = ""
		client.CreateNamespace = true

		s.Update("Installing Chart...")
		rel, err := client.Run(c, values)
		if err != nil {
			return result, err
		}
		s.Done()

		// Ensure our release name matches
		result.Release = rel.Name

		return result, nil
	}

	// We have a previous release, upgrade.
	client := action.NewUpgrade(actionConfig)
	client.ChartPathOptions = *cpo
	client.DryRun = false
	client.DisableHooks = false
	client.Wait = true
	client.WaitForJobs = false
	client.Devel = p.config.Devel
	client.DependencyUpdate = false
	client.Timeout = 300 * time.Second
	client.Namespace = p.config.Namespace
	client.Atomic = false
	client.SkipCRDs = false
	client.SubNotes = true
	client.DisableOpenAPIValidation = false
	client.Description = ""
	client.ResetValues = false
	client.ReuseValues = false
	client.Recreate = false
	client.MaxHistory = 0
	client.CleanupOnFail = false
	client.Force = false

	s.Update("Upgrading release...")
	rel, err := client.Run(prevRel.Name, c, values)
	if err != nil {
		return result, err
	}
	s.Done()

	// Ensure our release name matches
	result.Release = rel.Name

	return result, nil
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
	s := sg.Add("Uninstalling Helm release...")
	defer func() { s.Abort() }()

	actionConfig, err := p.actionInit(log)
	if err != nil {
		return err
	}

	_, err = action.NewUninstall(actionConfig).Run(deployment.Release)
	if err != nil {
		if strings.Contains(err.Error(), "release: not found") {
			err = nil
		}

		if err != nil {
			return err
		}
	}

	s.Done()
	return nil
}

// Generation returns the generation ID.
func (p *Platform) Generation(
	ctx context.Context,
) ([]byte, error) {
	// The generation is the release name.
	return []byte(p.config.Name), nil
}

// Config is the configuration structure for the Platform.
//
// For docs on the fields, please see the Documentation function.
type Config struct {
	Name       string   `hcl:"name,attr"`
	Repository string   `hcl:"repository,optional"`
	Chart      string   `hcl:"chart,attr"`
	Version    string   `hcl:"version,optional"`
	Devel      bool     `hcl:"devel,optional"`
	Values     []string `hcl:"values,optional"`
	Set        []*struct {
		Name  string `hcl:"name,attr"`
		Value string `hcl:"value,attr"`
		Type  string `hcl:"type,optional"`
	} `hcl:"set,block"`
	Driver    string `hcl:"driver,optional"`
	Namespace string `hcl:"namespace,optional"`

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

This is documented in more detail with a full example in the
[Kubernetes Helm Deployment documentation](/docs/platforms/kubernetes/helm-deploy).

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

	doc.Example(`
// Configuring an image to match the build. This assumes the chart
// has a value named "deployment.image".
deploy {
  use "helm" {
    chart = "${path.app}/chart"

    set {
      name  = "deployment.image"
      value = artifact.name
    }
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
		"devel",
		"True to considered non-released chart versions for installation.",
		docs.Summary(
			"This is equivalent to the `--devel` flag to `helm install`.",
		),
		docs.Default("false"),
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
		"namespace",
		"Namespace to deploy the Helm chart.",
		docs.Summary(
			"This will be created if it does not exist. This defaults to the ",
			"current namespace of the auth settings.",
		),
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
