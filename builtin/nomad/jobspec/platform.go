package jobspec

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/hashicorp/waypoint/builtin/nomad"
)

const (
	metaId = "waypoint.hashicorp.com/id"
)

// Platform is the Platform implementation for Nomad.
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

// Deploy deploys an image to Nomad.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*nomad.Deployment, error) {
	// We'll update the user in real time
	sg := ui.StepGroup()
	defer sg.Wait()
	s := sg.Add("Initializing the Nomad client...")
	defer func() { s.Abort() }()

	// Get our client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	jobclient := client.Jobs()

	// Parse the HCL
	s.Update("Parsing the job specification...")
	job, err := p.jobspec(client, p.config.Jobspec)
	if err != nil {
		return nil, err
	}

	// Create our deployment and set an initial ID
	var result nomad.Deployment
	result.Id = deployConfig.Id
	result.Name = *job.ID

	// Set our deployment ID on the meta.
	job.SetMeta(metaId, result.Id)

	// Register job
	s.Update("Registering job %q...", *job.Name)
	regResult, _, err := jobclient.Register(job, nil)
	if err != nil {
		return nil, err
	}
	s.Done()

	// Wait on the allocation
	st := ui.Status()
	defer st.Close()
	evalID := regResult.EvalID
	st.Update(fmt.Sprintf("Monitoring evaluation %q", evalID))
	if err := nomad.NewMonitor(st, client).Monitor(evalID); err != nil {
		return nil, err
	}
	st.Step(terminal.StatusOK, "Deployment successfully rolled out!")

	return &result, nil
}

// Destroy deletes the Nomad job.
func (p *Platform) Destroy(
	ctx context.Context,
	log hclog.Logger,
	deployment *nomad.Deployment,
	ui terminal.UI,
) error {
	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	st.Update("Deleting job...")
	_, _, err = client.Jobs().Deregister(deployment.Name, true, nil)
	return err
}

// Generation returns the generation ID. The ID we use is the name of the
// job since this is the unique ID that determines insert vs. update behavior
// for Nomad.
func (p *Platform) Generation(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) ([]byte, error) {
	// Get our client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	// Parse the HCL
	job, err := p.jobspec(client, p.config.Jobspec)
	if err != nil {
		return nil, err
	}

	return []byte(*job.ID), nil
}

func (p *Platform) jobspec(client *api.Client, path string) (*api.Job, error) {
	jobspec, err := ioutil.ReadFile(p.config.Jobspec)
	if err != nil {
		return nil, err
	}
	job, err := client.Jobs().ParseHCL(string(jobspec), true)
	if err != nil {
		return nil, err
	}
	if job.ID == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "job ID must not be empty")
	}
	if job.Name == nil {
		job.Name = job.ID
	}

	return job, nil
}

// Config is the configuration structure for the Platform.
type Config struct {
	// The path to the job specification to load.
	Jobspec string `hcl:"jobspec,attr"`
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}), docs.FromFunc(p.DeployFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description(`
Deploy to a Nomad cluster from a pre-existing Nomad job specification file.

This plugin lets you use any pre-existing Nomad job specification file to
deploy to Nomad. This deployment is able to support all the features of Waypoint.
You may use Waypoint's [templating features](/docs/waypoint-hcl/functions/template)
to template the Nomad jobspec with information such as the artifact from
a previous build step, entrypoint environment variables, etc.

### Artifact Access

You may use Waypoint's [templating features](/docs/waypoint-hcl/functions/template)
to access information such as the artifact from the build or push stages.
An example below shows this by using ` + "`templatefile`" + ` mixed with
variables such as ` + "`artifact.image`" + ` to dynamically configure the
Docker image within the Nomad job specification.

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
One of the examples below shows the entrypoint environment variables being
injected.

### URL Service

If you want your workload to be accessible by the
[Waypoint URL service](/docs/url), you must set the PORT environment variable
within your job and be using the Waypoint entrypoint (documented in the
previous section).

The PORT environment variable should be the port that your web service
is listening on that the URL service will connect to. See one of the examples
below for more details.

`)

	doc.Example(`
deploy {
  use "nomad-jobspec" {
    jobspec = "${path.app}/app.nomad"
  }
}
`)

	doc.Example(`
// The waypoint.hcl file
deploy {
  use "nomad-jobspec" {
    // Templated to perhaps bring in the artifact from a previous
    // build/registry, entrypoint env vars, etc.
    jobspec = templatefile("${path.app}/app.nomad.tpl")
  }
}

// The app.nomad.tpl file
job "web" {
  datacenters = ["dc1"]

  group "app" {
    task "app" {
      driver = "docker"

      config {
        image = "${artifact.image}:${artifact.tag}"
      }

      env {
        %{ for k,v in entrypoint.env ~}
        ${k} = "${v}"
        %{ endfor ~}

        // Ensure we set PORT for the URL service. This is only necessary
        // if we want the URL service to function.
        PORT = 3000
      }
    }
  }
}
`)

	doc.SetField(
		"jobspec",
		"Path to a Nomad job specification file.",
	)

	return doc, nil
}

var (
	_ component.Generation   = (*Platform)(nil)
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
	_ component.Destroyer    = (*Platform)(nil)
)
