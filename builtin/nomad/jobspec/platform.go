// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package jobspec

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/hashicorp/waypoint/builtin/nomad"
)

const (
	metaId            = "waypoint.hashicorp.com/id"
	rmResourceJobName = "job"
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

// StatusFunc implements component.Status
func (p *Platform) StatusFunc() interface{} {
	return p.Status
}

func (p *Platform) resourceManager(log hclog.Logger, dcr *component.DeclaredResourcesResp) *resource.Manager {
	return resource.NewManager(
		resource.WithLogger(log.Named("resource_manager")),
		resource.WithValueProvider(p.getNomadClient),
		resource.WithDeclaredResourcesResp(dcr),
		resource.WithResource(resource.NewResource(
			resource.WithName(rmResourceJobName),
			resource.WithState(&Resource_Job{}),
			resource.WithCreate(p.resourceJobCreate),
			resource.WithDestroy(p.resourceJobDestroy),
			resource.WithStatus(p.resourceJobStatus),
			resource.WithPlatform("nomad-jobspec"),
		)),
	)
}

// getNomadJobspecClient is a value provider for our resource manager and provides
// the client connection used by resources to interact with Nomad.
func (p *Platform) getNomadClient() (*api.Client, error) {
	// Get our client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (p *Platform) resourceJobCreate(
	ctx context.Context,
	client *api.Client,
	result *Deployment,
	deployConfig *component.DeploymentConfig,
	img *docker.Image,
	ui terminal.UI,
	state *Resource_Job,
) error {
	st := ui.Status()
	defer st.Close()
	jobclient := client.Jobs()

	// Parse the HCL
	st.Update("Parsing the job specification...")
	job, err := p.jobspec(client, p.config.Jobspec, p.config.Hcl1)
	if err != nil {
		return err
	}

	result.Id = deployConfig.Id
	result.Name = *job.ID

	// Set our deployment ID on the meta.
	job.SetMeta(metaId, result.Id)

	// Update our client to use the Namespace set in the jobspec
	client.SetNamespace(*job.Namespace)

	// Get Consul ACL token from environment
	*job.ConsulToken, err = nomad.ConsulAuth()
	if err != nil {
		return err
	}

	// Get Vault token from environment
	*job.VaultToken, err = nomad.VaultAuth()
	if err != nil {
		return err
	}

	// Register job
	st.Update("Registering job " + *job.Name + "...")
	regResult, _, err := jobclient.Register(job, nil)
	if err != nil {
		return err
	}

	// Store our state so we can destroy it properly
	state.Name = result.Name
	st.Step(terminal.StatusOK, "Job registration successful")

	// Wait on the allocation. Periodic Nomad jobs will not get an evaluation,
	// so we don't monitor an evaluation if we don't have one.
	evalID := regResult.EvalID
	if evalID != "" {
		st.Update("Monitoring evaluation " + evalID)
		if err := nomad.NewMonitor(st, client).Monitor(evalID); err != nil {
			return err
		}
	}
	st.Step(terminal.StatusOK, "Deployment successfully rolled out!")

	return nil
}

func (p *Platform) resourceJobDestroy(
	ctx context.Context,
	state *Resource_Job,
	sg terminal.StepGroup,
	client *api.Client,
) error {
	step := sg.Add("")
	defer func() { step.Abort() }()
	step.Update("Deleting job: %s", state.Name)
	step.Done()
	_, _, err := client.Jobs().Deregister(state.Name, true, nil)
	return err
}

func (p *Platform) resourceJobStatus(
	ctx context.Context,
	log hclog.Logger,
	sg terminal.StepGroup,
	state *Resource_Job,
	client *api.Client,
	sr *resource.StatusResponse,
	ui terminal.UI,
) error {
	s := sg.Add("Gathering health report for Nomad job...")
	defer s.Abort()

	jobClient := client.Jobs()

	s.Update("Parsing the job specification...")
	jobspec, err := p.jobspec(client, p.config.Jobspec, p.config.Hcl1)
	if err != nil {
		return err
	}

	jobResource := sdk.StatusReport_Resource{
		Type:                "job",
		CategoryDisplayHint: sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER,
	}
	sr.Resources = append(sr.Resources, &jobResource)

	s.Update("Getting job info...")
	q := &api.QueryOptions{Namespace: *jobspec.Namespace}
	job, _, err := jobClient.Info(state.Name, q)
	if err != nil {
		return err
	}

	jobResource.Id = *job.ID
	jobResource.Name = *job.Name
	jobResource.CreatedTime = timestamppb.New(time.Unix(0, *job.SubmitTime))
	stateJson, err := json.Marshal(map[string]interface{}{
		"deployment": job,
	})
	if err != nil {
		return err
	}

	jobResource.StateJson = string(stateJson)

	// If job is running, start checking evals, then allocs
	if *job.Status == "running" {
		// Get list of evaluations for job
		evals, _, err := jobClient.Evaluations(*job.ID, q)
		if err != nil {
			return err
		}
		hasSquashedEvals := false
		for _, eval := range evals {
			switch eval.Status {
			case "blocked":
				hasSquashedEvals = true
				break
			case "pending":
				break
			case "complete":
				break
			case "failed":
				break
			case "canceled":
				break
			}
		}

		allocs, _, err := jobClient.Allocations(*job.ID, false, q)
		if err != nil {
			return err
		}
		pending, running, complete, failed, lost, unknown, currentJobVersionAllocs := 0, 0, 0, 0, 0, 0, 0
		for _, alloc := range allocs {
			// Check alloc only if it is for the current job version
			if *job.Version == alloc.JobVersion {
				switch alloc.ClientStatus {
				case api.AllocClientStatusPending:
					pending += 1
				case api.AllocClientStatusComplete:
					complete += 1
				case api.AllocClientStatusRunning:
					running += 1
				case api.AllocClientStatusFailed:
					failed += 1
				case api.AllocClientStatusLost:
					lost += 1
				default:
					unknown += 1
				}
				currentJobVersionAllocs += 1
			}
		}

		// Need to subtract # of canaries in the update stanza from
		// "completed". Canary allocs will end up in the "completed"
		// state after the deployment, and thusly throw off the count
		// of otherwise "completed" allocs, resulting in a partial
		// state, when it's actually healthy
		if complete > 0 {
			complete = complete - *job.Update.Canary
		}

		if running == currentJobVersionAllocs && hasSquashedEvals == false {
			jobResource.Health = sdk.StatusReport_READY
			jobResource.HealthMessage = fmt.Sprintf("Job %q is reporting ready!", state.Name)
		} else if running == currentJobVersionAllocs && hasSquashedEvals == true {
			jobResource.Health = sdk.StatusReport_PARTIAL
			jobResource.HealthMessage = fmt.Sprintf("Allocs for job %q are running, but there is at least one blocked evaluation!", state.Name)
		} else if running > 0 && (complete > 0 || failed > 0 || pending > 0 || lost > 0) {
			jobResource.Health = sdk.StatusReport_PARTIAL
			jobResource.HealthMessage = fmt.Sprintf("Some allocations are running for job %q!", state.Name)
		} else if (complete + pending + failed + lost) == currentJobVersionAllocs {
			jobResource.Health = sdk.StatusReport_DOWN
			jobResource.HealthMessage = fmt.Sprintf("No allocations for job %q are running!", state.Name)
		} else if unknown > 0 {
			jobResource.Health = sdk.StatusReport_UNKNOWN
			jobResource.HealthMessage = fmt.Sprintf("Unknown allocation status for job %q!", state.Name)
		}
	} else if *job.Status == "pending" {
		jobResource.Health = sdk.StatusReport_PARTIAL
		jobResource.HealthMessage = fmt.Sprintf("Job %q is not scheduled!", state.Name)
	} else if *job.Status == "dead" {
		jobResource.Health = sdk.StatusReport_DOWN
		jobResource.HealthMessage = fmt.Sprintf("Job %q is down!", state.Name)
	}

	if *job.StatusDescription != "" {
		jobResource.HealthMessage = *job.StatusDescription
	}

	s.Done()

	// NOTE(briancain): Replace ui.Status with StepGroups once this bug
	// has been fixed: https://github.com/hashicorp/waypoint/issues/1536
	st := ui.Status()
	defer st.Close()

	st.Update("Determining overall container health...")
	if jobResource.Health == sdk.StatusReport_READY {
		st.Step(terminal.StatusOK, jobResource.HealthMessage)
	} else {
		if jobResource.Health == sdk.StatusReport_PARTIAL {
			st.Step(terminal.StatusWarn, jobResource.HealthMessage)
		} else {
			st.Step(terminal.StatusError, jobResource.HealthMessage)
		}

		// Extra advisory wording to let user know that the deployment could be still starting up
		// if the report was generated immediately after it was deployed or released.
		st.Step(terminal.StatusWarn, mixedHealthWarn)
	}

	return nil
}

// Deploy deploys an image to Nomad.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	dcr *component.DeclaredResourcesResp,
	ui terminal.UI,
) (*Deployment, error) {
	var result Deployment
	// Get our client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	// Parse the HCL
	job, err := p.jobspec(client, p.config.Jobspec, p.config.Hcl1)
	if err != nil {
		return nil, err
	}

	result.Name = *job.ID

	// We'll update the user in real time
	sg := ui.StepGroup()
	defer sg.Wait()
	rm := p.resourceManager(log, dcr)
	if err := rm.CreateAll(
		ctx, log, sg, ui,
		src, img, deployConfig, &result,
	); err != nil {
		return nil, err
	}

	// Store our resource state
	result.ResourceState = rm.State()

	return &result, nil
}

// Destroy deletes the Nomad job.
func (p *Platform) Destroy(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) error {
	sg := ui.StepGroup()
	defer sg.Wait()
	rm := p.resourceManager(log, nil)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if deployment.ResourceState == nil {
		rm.Resource(rmResourceJobName).SetState(&Resource_Job{
			Name: deployment.Name,
		})
	} else {
		// Load state
		if err := rm.LoadState(deployment.ResourceState); err != nil {
			return err
		}
	}

	return rm.DestroyAll(ctx, log, sg, ui)
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
	job, err := p.jobspec(client, p.config.Jobspec, p.config.Hcl1)
	if err != nil {
		return nil, err
	}

	canaryDeployment := false
	// If we have canaries, generate random ID, otherwise keep gen ID as job ID.
	// Periodic jobs and system jobs currently don't support canaries, so we don't
	// do this check if our job fits either case.
	if !job.IsPeriodic() && *job.Type != "system" {
		for _, taskGroup := range job.TaskGroups {
			if *taskGroup.Update.Canary > 0 {
				canaryDeployment = true
			}
		}
	}

	if !canaryDeployment {
		return []byte(*job.ID), nil
	} else {
		return nil, nil
	}

}

func (p *Platform) jobspec(client *api.Client, path string, hcl1 bool) (*api.Job, error) {
	jobspec, err := ioutil.ReadFile(p.config.Jobspec)
	if err != nil {
		return nil, err
	}
	job, err := client.Jobs().ParseHCLOpts(&api.JobsParseRequest{
		JobHCL:       string(jobspec),
		HCLv1:        hcl1,
		Canonicalize: true,
	})
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

func (p *Platform) Status(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) (*sdk.StatusReport, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	rm := p.resourceManager(log, nil)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if deployment.ResourceState == nil {
		rm.Resource("job").SetState(&Resource_Job{
			Name: deployment.Name,
		})
	} else {
		// Load our set state
		if err := rm.LoadState(deployment.ResourceState); err != nil {
			return nil, err
		}
	}

	step := sg.Add("Gathering health report for Nomad platform...")
	defer step.Abort()

	resources, err := rm.StatusAll(ctx, log, sg, ui)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "resource manager failed to generate resource statuses: %s", err)
	}

	if len(resources) == 0 {
		// This shouldn't happen - the status func for the releaser should always return a resource or an error.
		return nil, status.Errorf(codes.Internal, "no resources generated for release - cannot determine status.")
	}

	var jobResource *sdk.StatusReport_Resource
	for _, r := range resources {
		if r.Type == "job" {
			jobResource = r
			break
		}
	}
	if jobResource == nil {
		return nil, status.Errorf(codes.Internal, "no job resource found - cannot determine overall health")
	}

	// Create our status report
	result := sdk.StatusReport{
		External:      true,
		GeneratedTime: timestamppb.Now(),
		Resources:     resources,
		Health:        jobResource.Health,
		HealthMessage: jobResource.HealthMessage,
	}

	log.Debug("status report complete")

	// update output based on main health state
	step.Update("Finished building report for Nomad platform")
	step.Done()

	// NOTE(briancain): Replace ui.Status with StepGroups once this bug
	// has been fixed: https://github.com/hashicorp/waypoint/issues/1536
	st := ui.Status()
	defer st.Close()

	// More UI detail for non-ready resources
	for _, resource := range result.Resources {
		if resource.Health != sdk.StatusReport_READY {
			st.Step(terminal.StatusWarn, fmt.Sprintf("Resource %q is reporting %q", resource.Name, resource.Health.String()))
		}
	}

	return &result, nil
}

// Config is the configuration structure for the Platform.
type Config struct {
	// The path to the job specification to load.
	Jobspec string `hcl:"jobspec,attr"`

	// Signifies whether the jobspec should be parsed as HCL1 or not
	Hcl1 bool `hcl:"hcl1,optional"`
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
You may use Waypoint's [templating features](/waypoint/docs/waypoint-hcl/functions/template)
to template the Nomad jobspec with information such as the artifact from
a previous build step, entrypoint environment variables, etc.

### Artifact Access

You may use Waypoint's [templating features](/waypoint/docs/waypoint-hcl/functions/template)
to access information such as the artifact from the build or push stages.
An example below shows this by using ` + "`templatefile`" + ` mixed with
variables such as ` + "`artifact.image`" + ` to dynamically configure the
Docker image within the Nomad job specification.

-> **Note:** If using [Nomad interpolation](/nomad/docs/runtime/interpolation) in your jobspec file,
and the ` + "`templatefile`" + ` function in your waypoint.hcl file, any interpolated values must be escaped with a second 
` + "`$`" + `. For example: ` + "`$${meta.metadata}`" + ` instead of ` + "`${meta.metadata}`" + `.

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

-> **Note:** The Waypoint entrypoint and the [Nomad entrypoint functionality](/nomad/docs/drivers/docker#entrypoint) 
cannot be used simultaneously. In order to use the features of the Waypoint entrypoint, the Nomad entrypoint must not be used in your jobspec.

### URL Service

If you want your workload to be accessible by the
[Waypoint URL service](/waypoint/docs/url), you must set the PORT environment variable
within your job and be using the Waypoint entrypoint (documented in the
previous section).

The PORT environment variable should be the port that your web service
is listening on that the URL service will connect to. See one of the examples
below for more details.

`)

	doc.Input("docker.Image")
	doc.Output("jobspec.Deployment")

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

	doc.SetField(
		"hcl1",
		"Parses jobspec as HCL1 instead of HCL2.",
		docs.Default("false"),
	)

	doc.SetField(
		"consul_token",
		"The Consul ACL token used to register services with the Nomad job.",
		docs.Summary("Uses the runner config environment variable CONSUL_HTTP_TOKEN."),
		docs.EnvVar("CONSUL_HTTP_TOKEN"),
	)

	doc.SetField(
		"vault_token",
		"The Vault token used to deploy the Nomad job with a token having specific Vault policies attached.",
		docs.Summary("Uses the runner config environment variable VAULT_TOKEN."),
		docs.EnvVar("VAULT_TOKEN"),
	)

	return doc, nil
}

var (
	mixedHealthWarn = strings.TrimSpace(`
Waypoint detected that the current deployment is not ready, however your application
might be available or still starting up.
`)
)

var (
	_ component.Generation   = (*Platform)(nil)
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
	_ component.Destroyer    = (*Platform)(nil)
)
