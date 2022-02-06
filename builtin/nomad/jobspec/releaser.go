package jobspec

import (
	"context"
	"encoding/json"
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
	"github.com/hashicorp/waypoint/builtin/nomad"
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
		jobResource.HealthMessage = fmt.Sprintf("Job %q exists and is ready", *job.Name)
		stateJson, err := json.Marshal(map[string]interface{}{
			"deployment": job,
		})
		if err != nil {
			jobResource.StateJson = string(stateJson)
		}
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
	// Set up clients
	jobClient := client.NomadClient.Jobs()
	deploymentClient := client.NomadClient.Deployments()

	// TODO: Use step group
	st.Update("Getting job...")
	// TODO: Error if target.Name doesn't exactly match the IDs of any of the jobs returned
	jobs, _, err := jobClient.PrefixList(target.Name)
	if err != nil {
		return status.Errorf(codes.Aborted, "Unable to fetch Nomad jobs: %s", err.Error())
	}

	// TODO: Add ACL token here
	q := &api.QueryOptions{
		Namespace: jobs[0].JobSummary.Namespace,
	}

	job, _, err := jobClient.Info(jobs[0].ID, q)

	if err != nil {
		return status.Errorf(codes.Aborted, "Unable to fetch Nomad job: %s", err.Error())
	}

	st.Update("Getting latest deployments for job")
	deploy, _, err := jobClient.LatestDeployment(*job.ID, q)
	if err != nil {
		return status.Errorf(codes.Aborted, "Unable to fetch latest deployment for Nomad job: %s", err.Error())
	} else if deploy == nil {
		st.Update("No active deployment for Nomad job")
		return err
	}

	canaryDeployment := false
	groupsToPromote := make([]string, len(deploy.TaskGroups))
	for taskGroupName, taskGroup := range deploy.TaskGroups {
		if r.config.Groups != nil {
			if isElementExist(r.config.Groups, taskGroupName) && taskGroup.DesiredCanaries != 0 {
				canaryDeployment = true
				groupsToPromote = append(groupsToPromote, taskGroupName)
				continue
			}
		} else if taskGroup.DesiredCanaries != 0 {
			canaryDeployment = true
			// if no groups to promote are specified in the config, promote all groups
			// that have canaries
			groupsToPromote = append(groupsToPromote, taskGroupName)
		}
	}
	if !canaryDeployment {
		log.Info("Canaries not detected")
		return nil
	}

	// check each task group to be promoted; if the canary allocs aren't healthy,
	//   check again in 5 seconds
	// TODO: Force timeout if exceeds healthy deadline or progress deadline of job
	var currentTaskGroupState *api.DeploymentState
	var groupHealthy bool
	for _, group := range groupsToPromote {
		if group != "" {
			currentTaskGroupState = deploy.TaskGroups[group]
			groupHealthy = false
			for !groupHealthy {
				if currentTaskGroupState.HealthyAllocs < len(currentTaskGroupState.PlacedCanaries) {
					time.Sleep(5 * time.Second)
					deploy, _, err = jobClient.LatestDeployment(*job.ID, q)
					currentTaskGroupState = deploy.TaskGroups[group]
				} else {
					groupHealthy = true
				}
			}
		}
	}

	// TODO: Add ACL token here
	wq := &api.WriteOptions{Namespace: *job.Namespace}

	var u *api.DeploymentUpdateResponse
	if r.config.FailDeployment {
		u, _, err = deploymentClient.Fail(deploy.ID, wq)
	} else {
		u, _, err = deploymentClient.PromoteGroups(deploy.ID, groupsToPromote, wq)
	}
	if err != nil {
		return err
	}

	st.Update(fmt.Sprintf("Monitoring evaluation %s", string(u.EvalID)))
	if err := nomad.NewMonitor(st, client.NomadClient).Monitor(u.EvalID); err != nil {
		return err
	}

	// TODO: Automatically search for Consul service, determine FQDN for service
	// TODO: Automatically search for ingress gateway, determine FQDN for service
	// TODO: Automatically search for IP and port of random Nomad alloc in job
	// If meta not set, URL is empty
	result.Url = job.Meta["waypoint.hashicorp.com/release_url"]
	return nil
}

func (r *Releaser) resourceJobDestroy(
	client *nomadClient,
	state *Resource_Job,
	sg terminal.StepGroup,
) error {
	// Do nothing because the platform will destroy the job for us
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
		ctx, log, st, &result, target,
	); err != nil {
		return nil, err
	}

	result.ResourceState = rm.State()

	st.Step(terminal.StatusOK, "Release successfully rolled out!")
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
			Name: rmResourceJobName,
		})
	} else {
		// Load our set state
		if err := rm.LoadState(release.ResourceState); err != nil {
			return err
		}
	}

	return rm.DestroyAll(sg, ui)
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
			Name: rmResourceJobName,
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
	Groups         []string `hcl:"groups,optional"`
	FailDeployment bool     `hcl:"fail_deployment,optional"`
	//TODO: Support option to revert to a previous version?
	//      Should something like this (rollbacks) be accommodated by a releaser?
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

func isElementExist(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func (r *Release) URL() string { return r.Url }

var (
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Configurable   = (*Releaser)(nil)
)
