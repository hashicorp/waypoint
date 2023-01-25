package apply

import (
	"context"
	"os/exec"

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

	deployment := &Deployment{
		PruneLabel:     p.config.PruneLabel,
		PruneWhitelist: p.config.PruneWhitelist,
	}

	s.Update("Executing kubectl apply ...")
	args := []string{
		"-R",
		"-f", p.config.Path,
		"--prune",
		"-l", deployment.PruneLabel,
	}

	for _, v := range deployment.PruneWhitelist {
		args = append(args, []string{
			"--prune-whitelist",
			v,
		}...)
	}

	cmd, err := p.cmd(ctx, s, "apply", args...)
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

func (p *Platform) cmd(
	ctx context.Context,
	step terminal.Step,
	subcmd string, args ...string,
) (*exec.Cmd, error) {
	path, err := exec.LookPath("kubectl")
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"kubectl could not be found on the PATH")
	}

	// Build the args we'll send to kubectl, which are:
	// <subcmd> <flags> <flags + args provided to this func>
	realArgs := []string{subcmd}
	if p.config.KubeconfigPath != "" {
		realArgs = append(realArgs, "--kubeconfig="+p.config.KubeconfigPath)
	}
	if p.config.Context != "" {
		realArgs = append(realArgs, "--context="+p.config.Context)
	}
	realArgs = append(realArgs, args...)

	// Build our command
	cmd := exec.CommandContext(ctx, path, realArgs...)

	// Ensure output goes to our step
	cmd.Stdout = step.TermOutput()
	cmd.Stderr = cmd.Stdout
	return cmd, nil
}

// Config is the configuration structure for the Platform.
type Config struct {
	// The path to the job specification to load.
	Path string `hcl:"path,attr"`

	// Prune label is the label to use to destroy resources that don't match.
	PruneLabel string `hcl:"prune_label,attr"`

	// PruneWhitelist is a list of Kubernetes Objects that are allowed to be pruned
	// An empty list means the defaults. Specify them as group/version/kind (e.g: apps/v1/Deployment)
	// (see https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands --prune-whitelist)
	PruneWhitelist []string `hcl:"prune_whitelist,attr"`

	// KubeconfigPath is the path to the kubeconfig file.
	KubeconfigPath string `hcl:"kubeconfig,optional"`

	// Kubernetes context to use in the kubeconfig
	Context string `hcl:"context,optional"`
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}), docs.FromFunc(p.DeployFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description(`
Deploy Kubernetes resources directly from a single file or a directory of YAML
or JSON files.

This plugin lets you use any pre-existing set of Kubernetes resource files
to deploy to Kubernetes. This plugin supports all the features of Waypoint.
You may use Waypoint's [templating features](/waypoint/docs/waypoint-hcl/functions/template)
to template the resources with information such as the artifact from
a previous build step, entrypoint environment variables, etc.

### Requirements

This plugin requires "kubectl" to be installed since this plugin works by
subprocessing to "kubectl apply". Other Waypoint Kubernetes plugins sometimes
use the API directly but this plugin requires "kubectl".

"kubectl" must also be configured to access your Kubernetes cluster. You
may specify an alternate kubeconfig file using the "kubeconfig" configuration
parameter. If this isn't specified, the default kubectl lookup paths will be
used.

### Artifact Access

You may use Waypoint's [templating features](/waypoint/docs/waypoint-hcl/functions/template)
to access information such as the artifact from the build or push stages.
An example below shows this by using ` + "`templatedir`" + ` mixed with
variables such as ` + "`artifact.image`" + ` to dynamically configure the
Docker image within a Kubernetes Deployment.

### Entrypoint Functionality

Waypoint [entrypoint functionality](/waypoint/docs/entrypoint#functionality) such
as logs, exec, app configuration, and more require two properties to be true:

1. The running image must already have the Waypoint entrypoint installed
  and configured as the entrypoint. This should happen in the build stage.

2. Proper environment variables must be set so the entrypoint knows how
  to communicate to the Waypoint server. **This step happens in this
  deployment stage.**

**Step 2 does not happen automatically.** You must manually set the entrypoint
environment variables using the [templating feature](/waypoint/docs/waypoint-hcl/functions/template).
One of the examples below shows the entrypoint environment variables being
injected.

### URL Service

If you want your workload to be accessible by the
[Waypoint URL service](/waypoint/docs/url), you must set the PORT environment variable
within the pod with your web service and also be using the Waypoint
entrypoint (documented in the previous section).

The PORT environment variable should be the port that your web service
is listening on that the URL service will connect to. See one of the examples
below for more details.

`)

	doc.Input("FIXME")
	doc.Output("k8sapply.Deployment")

	doc.Example(`
deploy {
  use "kubernetes-apply" {
    path = "${path.app}/k8s"
  }
}
`)

	doc.Example(`
// The waypoint.hcl file
deploy {
  use "kubernetes-apply" {
    // Templated to perhaps bring in the artifact from a previous
    // build/registry, entrypoint env vars, etc.
    path        = templatedir("${path.app}/k8s")
    prune_label = "app=myapp"
	prune_whitelist = [
		"apps/v1/Deployment",
		"apps/v1/ReplicaSet"
  	]
  }
}

// ./k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
  labels:
    app: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: ${artifact.image}:${artifact.tag}
        env:
          %{ for k,v in entrypoint.env ~}
          - name: ${k}
            value: "${v}"
          %{ endfor ~}

          # Ensure we set PORT for the URL service. This is only necessary
          # if we want the URL service to function.
          - name: PORT
            value: "3000"
`)

	doc.SetField(
		"path",
		"Path to a file or directory of YAML or  JSON files.",
		docs.Summary(
			"This will be used for `kubectl apply` to create a set of",
			"Kubernetes resources. Pair this with `templatefile` or `templatedir`",
			"[templating functions](/waypoint/docs/waypoint-hcl/functions/template)",
			"to inject dynamic elements into your Kubernetes resources.",
			"Subdirectories are included recursively.",
		),
	)

	doc.SetField(
		"prune_label",
		"Label selector to prune resources that aren't present in the `path`.",
		docs.Summary(
			"This is a label selector that is used to search for any resources",
			"that are NOT present in the configured `path` and delete them.",
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
