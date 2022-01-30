package nomad

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

// Releaser is the ReleaseManager implementation for Nomad.
type Releaser struct {
	p      *Platform
	config ReleaserConfig
}

// Config implements Configurable
func (r *Releaser) Config() (interface{}, error) {
	return &r.config, nil
}

// ReleaseFunc implements component.ReleaseManager
func (r *Releaser) ReleaseFunc() interface{} {
	return r.Release
}

// DestroyFunc implements component.Destroyer
func (r *Releaser) DestroyFunc() interface{} {
	return r.Destroy
}

// StatusFunc implements component.Status
func (r *Releaser) StatusFunc() interface{} {
	return r.Status
}

func (r *Releaser) resourceManager(log hclog.Logger, dcr *component.DeclaredResourcesResp) *resource.Manager {
	return resource.NewManager(
		resource.WithLogger(log.Named("resource_manager")),
		resource.WithValueProvider(r.getNomadClient),
		resource.WithDeclaredResourcesResp(dcr),
		resource.WithResource(resource.NewResource(
			resource.WithName(rmResourceJobName),
			resource.WithState(&Resource_Job{}),
			resource.WithCreate(r.resourceJobCreate),
			resource.WithDestroy(r.resourceJobDestroy),
			resource.WithStatus(r.resourceJobStatus),
			resource.WithPlatform("nomad"),
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER),
		)),
	)
}

// getNomadClient provides
// the client connection used by resources to interact with Nomad.
func (r *Releaser) getNomadClient() (*nomadClient, error) {
	// Get our client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return &nomadClient{
		NomadClient: client,
	}, nil
}

func (r *Releaser) resourceJobStatus(
	ctx context.Context,
	log hclog.Logger,
	sg terminal.StepGroup,
	state *Resource_Job,
	client *nomadClient,
	sr *resource.StatusResponse,
) error {
	s := sg.Add("Checking status of Nomad job resource %q...", state.Name)
	defer s.Abort()

	jobClient := client.NomadClient.Jobs()
	s.Update("Getting job...")
	jobs, _, err := jobClient.PrefixList(state.Name)
	q := &api.QueryOptions{Namespace: jobs[0].JobSummary.Namespace}

	jobResource := sdk.StatusReport_Resource{
		CategoryDisplayHint: sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER,
	}
	sr.Resources = append(sr.Resources, &jobResource)

	job, _, err := jobClient.Info(jobs[0].ID, q)

	if jobs == nil {
		return status.Errorf(codes.FailedPrecondition, "Nomad job response cannot be empty")
	} else if err != nil {
		s.Update("No job was found")
		s.Status(terminal.StatusError)
		s.Done()
		s = sg.Add("")

		jobResource.Name = state.Name
		jobResource.Health = sdk.StatusReport_MISSING
		jobResource.HealthMessage = sdk.StatusReport_MISSING.String()
	} else {
		jobResource.Id = *job.ID
		jobResource.Name = *job.Name
		jobResource.CreatedTime = timestamppb.New(time.Unix(0, *job.SubmitTime))
		jobResource.Health = sdk.StatusReport_READY
		jobResource.HealthMessage = fmt.Sprintf("Job %q exists and is ready", job.Name)
		//jobResource.StateJson =
	}

	s.Update("Finished building report for Nomad job resource")
	s.Done()
	return nil
}

// might be a better name than JobCreate for promoting a canary
func (r *Releaser) resourceJobCreate(
	ctx context.Context,
	log hclog.Logger,
	target *Deployment,
	result *Release,
	state *Resource_Job,
	client *nomadClient,
	st terminal.Status,
) error {
    //TODO: Use step group
	//step := sg.Add("Initializing Nomad client...")
	//defer func() { step.Abort() }()

	jobClient := client.NomadClient.Jobs()
	deploymentClient := client.NomadClient.Deployments()

	st.Update("Getting job...")
	jobs, _, err := jobClient.PrefixList(target.Name)
	if err != nil {
		return status.Errorf(codes.Aborted, "Unable to fetch Nomad job: %s", err.Error())
	}

	q := &api.QueryOptions{Namespace: jobs[0].JobSummary.Namespace}
	st.Update("Getting latest deployments for job")
	deploy, _, err := jobClient.LatestDeployment(jobs[0].ID, q)
	if err != nil {
		return status.Errorf(codes.Aborted, "Unable to fetch latest deployment for Nomad job: %s", err.Error())
	}

	if deploy == nil {
		st.Update("No active deployment for Nomad job")
		return err
	}

	//Check if any of the task groups are canary deployments
	//TODO: Match up specified 'groups' in ReleaserConfig to group names found in the Deployment
	//      Verify that they 1) exist and 2) have canaries
	canaryDeployment := false
	for _, taskGroup := range deploy.TaskGroups {
		if taskGroup.DesiredCanaries != 0 {
			canaryDeployment = true
		}
	}
	//return errorf here
	if !canaryDeployment {
		return nil
	}

	// Set write options
	wq := &api.WriteOptions{Namespace: jobs[0].JobSummary.Namespace}

	var u *api.DeploymentUpdateResponse
	//TODO: Add logic to support promotion of specific group(s)
	u, _, err = deploymentClient.PromoteAll(deploy.ID, wq)
	st.Update(fmt.Sprintf("Monitoring evaluation %q", u.EvalID))

	if err := NewMonitor(st, client.NomadClient).Monitor(u.EvalID); err != nil {
		return err
	}

	//TODO: If applicable, get Consul service from job. If multiple services, how to determine which service to use
	//      (maybe from ReleaserConfig)? Consul service URL structure may be ambiguous here as well:
	//      'service_name.service.consul' is common; however, `.consul` is default domain for Consul, but this is not
	//      mandatory. The service could also be an ingress gateway, where the name would be service_name.ingress.consul.
	//      The Consul data center may also be required, and/or tags, for FQDN of:
	//      tag_name.service_name.ingress/service.datacenter.consul
	//      https://www.consul.io/docs/discovery/dns#standard-lookup
	//      If no Consul service, select IP/Port of a random instance?
	result.Url = "https://waypointproject.io"
	result.Id = jobs[0].ID
	result.Name = jobs[0].Name
	return nil
}

func (r *Releaser) resourceJobDestroy(
	ctx context.Context,
	state *Resource_Job,
	sg terminal.StepGroup,
	client *nomadClient,
) error {
	step := sg.Add("Initializing Nomad client...")
	defer func() { step.Abort() }()

	nomadClient := client.NomadClient
	jobClient := nomadClient.Jobs()

	step.Update("Getting job...")
	jobs, _, err := jobClient.PrefixList(state.Name)
	if err != nil {
		return status.Errorf(codes.Aborted, "Unable to fetch Nomad job: %s", err.Error())
	}

	// Set write options
	wq := &api.WriteOptions{Namespace: jobs[0].JobSummary.Namespace}

	step.Update("Deleting job: %s", state.Name)
	_, _, err = jobClient.Deregister(state.Name, true, wq)

	if err != nil {
	  return err
	}

	step.Update("Job deleted")
	step.Done()

	return nil
}

// Release promotes the Nomad canary deployment
func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	job *component.JobInfo,
	ui terminal.UI,
	target *Deployment,
	dcr *component.DeclaredResourcesResp,
) (*Release, error) {
	var result Release

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	rm := r.resourceManager(log, dcr)
	if err := rm.CreateAll(
		ctx, log, st, ui,
		target, &result,
	); err != nil {
		return nil, err
	}

	result.ResourceState = rm.State()

	return &result, nil
}

func (r *Releaser) Destroy(
	ctx context.Context,
	log hclog.Logger,
	release *Release,
	ui terminal.UI,
) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	rm := r.resourceManager(log, nil)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if release.ResourceState == nil {
		rm.Resource(rmResourceJobName).SetState(&Resource_Job{
			Name: release.Name,
		})
	} else {
		// Load our set state
		if err := rm.LoadState(release.ResourceState); err != nil {
			return err
		}
	}

	return rm.DestroyAll(ctx, log, sg, ui)
}

func (r *Releaser) Status(
	ctx context.Context,
	log hclog.Logger,
	release *Release,
	ui terminal.UI,
) (*sdk.StatusReport, error) {
	sg := ui.StepGroup()
	defer sg.Wait()

	rm := r.resourceManager(log, nil)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if release.ResourceState == nil {
		rm.Resource(rmResourceJobName).SetState(&Resource_Job{
			Name: release.Name,
		})
	} else {
		// Load our set state
		if err := rm.LoadState(release.ResourceState); err != nil {
			return nil, err
		}
	}

	step := sg.Add("Getting status of Nomad release...")
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

// ReleaserConfig is the configuration structure for the Releaser.
type ReleaserConfig struct {
	//Groups only applies to the nomad-jobspec platform since the nomad platform (currently) uses only one task group
	Groups []string `hcl:"groups,optional"`
	//TODO: Support option to fail canary deployment?
	//TODO: Support option to revert to a previous version?
	//      Should something like this (rollbacks) be accommodated by a releaser?
	//TODO: Support option to scale count?
	//      This may warrant a different releaser plugin, or a more generic name for this releaser plugin
	//      Note: Scaling a deployment doesn't require canaries (hence the generic name idea)
}

type nomadClient struct {
	NomadClient *api.Client
}

func (r *Releaser) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&ReleaserConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Promotes a Nomad canary deployment")

	doc.Input("nomad.Deployment")
	doc.Output("nomad.Release")

	return doc, nil
}

var (
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Configurable   = (*Releaser)(nil)
)
