package canary

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
	"github.com/hashicorp/waypoint/builtin/nomad"
	"github.com/hashicorp/waypoint/builtin/nomad/jobspec"
)

const (
	rmResourcePromotedJobName = "promoted-job"
)

// Releaser is the ReleaseManager implementation for Nomad.
type Releaser struct {
	p      *jobspec.Platform
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
			resource.WithName(rmResourcePromotedJobName),
			resource.WithState(&jobspec.Resource_Job{}),
			resource.WithCreate(r.resourceJobCreate),
			resource.WithDestroy(r.resourceJobDestroy),
			resource.WithStatus(r.resourceJobStatus),
			resource.WithPlatform("nomad-jobspec"),
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER),
		)),
	)
}

// getNomadClient provides
// the client connection used by resources to interact with Nomad.
func (r *Releaser) getNomadClient() (*api.Client, error) {
	// Get our client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (r *Releaser) resourceJobCreate(
	ctx context.Context,
	log hclog.Logger,
	target *jobspec.Deployment,
	result *Release,
	state *jobspec.Resource_Job,
	client *api.Client,
	st terminal.Status,
	sg terminal.StepGroup,
) error {
	// Set up clients
	jobClient := client.Jobs()
	deploymentClient := client.Deployments()

	st.Update("Getting job...")
	jobs, _, err := jobClient.PrefixList(target.Name)
	if err != nil {
		return status.Errorf(codes.Aborted, "Unable to fetch Nomad jobs: %s", err.Error())
	}

	if len(jobs) > 0 {
		if target.Name != jobs[0].ID {
			st.Step(terminal.StatusError, fmt.Sprintf("Job could not be found, did you mean to promote %q?", jobs[0].ID))
			return nil
		}
	} else {
		st.Step(terminal.StatusError, "Job not found.")
		return nil
	}

	q := &api.QueryOptions{
		Namespace: jobs[0].JobSummary.Namespace,
	}

	job, _, err := jobClient.Info(jobs[0].ID, q)

	if err != nil {
		return status.Errorf(codes.Aborted, "Unable to fetch Nomad job: %s", err.Error())
	}

	// if first deployment of the job, no chance of it being a canary deployment
	if *job.Version == 0 {
		st.Step(terminal.StatusOK, "This is the first deployment of the job - no canaries to promote.")
		return nil
	}

	st.Update("Getting latest deployments for job")
	deploy, _, err := jobClient.LatestDeployment(*job.ID, q)
	if err != nil {
		return status.Errorf(codes.Aborted, "Unable to fetch latest deployment for Nomad job: %s", err.Error())
	} else if deploy == nil {
		st.Update("No active deployment for Nomad job.")
		return nil
	} else if deploy.JobVersion != *job.Version {
		st.Update("Job version does not match deployment's job version.")
		return nil
	} else if deploy.Status != "running" {
		st.Update("Deployment for job is no longer running.")
		return nil
	}

	canaryDeployment := false
	groupsToPromote := make([]string, len(deploy.TaskGroups))
	for taskGroupName, taskGroup := range deploy.TaskGroups {
		if r.config.Groups != nil {
			if isElementExist(r.config.Groups, taskGroupName) && taskGroup.DesiredCanaries > 0 {
				canaryDeployment = true
				groupsToPromote = append(groupsToPromote, taskGroupName)
				continue
			}
		} else if taskGroup.DesiredCanaries > 0 {
			canaryDeployment = true
			// If no groups to promote are specified in the config, promote all groups
			//   that have canaries
			groupsToPromote = append(groupsToPromote, taskGroupName)
		}
	}
	if !canaryDeployment {
		st.Step(terminal.StatusWarn, "Nomad canary allocations not detected in job task groups.")
		log.Info("Canaries not detected")
		return nil
	}

	var currentTaskGroupState *api.DeploymentState
	for _, group := range groupsToPromote {
		if group != "" {
			st.Update("Checking task group: " + group)
			// TODO: Update to pair a task group's name with its healthy deadline, so
			//       we can set the deadline accordingly
			//       d := time.Now().Add(time.Nanosecond * time.Duration(*job.TaskGroups[indexOfTaskGroupInSliceOfTaskGroupsOfJob].Update.HealthyDeadline))
			d := time.Now().Add(time.Minute * time.Duration(5))
			log.Debug(fmt.Sprintf("Healthy deadline: %s", d.String()))
			ctx, cancel := context.WithDeadline(ctx, d)
			defer cancel()
			ticker := time.NewTicker(5 * time.Second)
			groupHealthy := false
			for !groupHealthy {
				currentTaskGroupState = deploy.TaskGroups[group]
				if currentTaskGroupState.HealthyAllocs < len(currentTaskGroupState.PlacedCanaries) {
					st.Update("Waiting on allocations to become healthy: healthy allocs=" + strconv.Itoa(currentTaskGroupState.HealthyAllocs) + " placed canaries=" + strconv.Itoa(len(currentTaskGroupState.PlacedCanaries)))
					select {
					case <-ticker.C:
					case <-ctx.Done(): // cancelled
						return status.Errorf(codes.Aborted, "Context cancelled from timeout checking health of task group %q: %s", group, ctx.Err())
					}
					deploy, _, err = jobClient.LatestDeployment(*job.ID, q)
					if err != nil {
						return status.Errorf(codes.Aborted, "Unable to fetch latest deployment: %s", err.Error())
					}
					currentTaskGroupState = deploy.TaskGroups[group]
					log.Info(fmt.Sprintf("Task group not healthy: %s", group))
				} else {
					groupHealthy = true
				}
			}
		}
	}

	wq := &api.WriteOptions{
		Namespace: *job.Namespace,
	}

	var u *api.DeploymentUpdateResponse
	if r.config.FailDeployment {
		u, _, err = deploymentClient.Fail(deploy.ID, wq)
	} else {
		u, _, err = deploymentClient.PromoteGroups(deploy.ID, groupsToPromote, wq)
	}
	if err != nil {
		return err
	}

	st.Update("Monitoring evaluation " + u.EvalID)
	if err := nomad.NewMonitor(st, client).Monitor(u.EvalID); err != nil {
		return err
	}

	state.Name = *job.Name
	state.Id = *job.ID
	state.Namespace = *job.Namespace

	// TODO: Automatically search for Consul service, determine FQDN for service
	// TODO: Automatically search for ingress gateway, determine FQDN for service
	// TODO: Automatically search for IP and port of random Nomad alloc in job
	// If meta not set, URL is empty
	result.Url = job.Meta["waypoint.hashicorp.com/release_url"]
	return nil
}

func (r *Releaser) resourceJobDestroy(
	log hclog.Logger,
	client *api.Client,
	state *jobspec.Resource_Job,
	sg terminal.StepGroup,
) error {
	log.Trace("No resource destroyed")
	return nil
}

func (r *Releaser) resourceJobStatus(
	ctx context.Context,
	log hclog.Logger,
	sg terminal.StepGroup,
	state *jobspec.Resource_Job,
	client *api.Client,
	sr *resource.StatusResponse,
) error {
	s := sg.Add("Checking status of Nomad job resource %q...", state.Name)
	defer s.Abort()

	jobClient := client.Jobs()
	s.Update("Getting job...")
	// TODO: Because we don't have the namespace from the jobspec, we rely on the
	//   NOMAD_NAMESPACE env var/searching for job via prefix- consider passing namespace
	//   from deploy phase
	jobResource := sdk.StatusReport_Resource{
		CategoryDisplayHint: sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER,
	}
	sr.Resources = append(sr.Resources, &jobResource)

	s.Update("Getting job info...")
	q := &api.QueryOptions{Namespace: state.Namespace}
	job, _, err := jobClient.Info(state.Id, q)

	if err != nil && job == nil {
		jobResource.Name = state.Name
		jobResource.Health = sdk.StatusReport_UNKNOWN
		jobResource.HealthMessage = sdk.StatusReport_UNKNOWN.String()
		return err
	} else if job == nil {
		s.Update("No job was found")
		s.Status(terminal.StatusError)
		s.Done()
		s = sg.Add("")

		jobResource.Name = state.Name
		jobResource.Health = sdk.StatusReport_UNKNOWN
		jobResource.HealthMessage = sdk.StatusReport_UNKNOWN.String()
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

// Release promotes the Nomad canary deployment
func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	job *component.JobInfo,
	ui terminal.UI,
	target *jobspec.Deployment,
	dcr *component.DeclaredResourcesResp,
) (*Release, error) {
	var result Release

	// We'll update the user in real time
	// TODO: Replace ui.Status with StepGroups once this bug
	// has been fixed: https://github.com/hashicorp/waypoint/issues/1536
	st := ui.Status()
	sg := ui.StepGroup()
	defer st.Close()
	defer sg.Wait()

	rm := r.resourceManager(log, dcr)
	if err := rm.CreateAll(
		ctx, log, st, sg, &result, target,
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
		rm.Resource(rmResourcePromotedJobName).SetState(&jobspec.Resource_Job{
			Name: rmResourcePromotedJobName,
		})
	} else {
		// Load our set state
		if err := rm.LoadState(release.ResourceState); err != nil {
			return err
		}
	}

	return rm.DestroyAll(log, sg, ui)
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
		rm.Resource(rmResourcePromotedJobName).SetState(&jobspec.Resource_Job{
			Name: rmResourcePromotedJobName,
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
		if r.Type == rmResourcePromotedJobName {
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
	// List of task group names which are to be promoted
	Groups []string `hcl:"groups,optional"`

	// If true, marks the deployment as failed
	FailDeployment bool `hcl:"fail_deployment,optional"`
}

func (r *Releaser) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&ReleaserConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description(`
Promotes a Nomad canary deployment initiated by a Nomad jobspec deployment.

If your Nomad deployment is configured to use canaries, this releaser plugin lets
you promote (or fail) the canary deployment. You may also target specific task
groups within your job for promotion, if you have multiple task groups in your canary
deployment.

-> **Note:** Using the ` + "`-prune=false`" + ` flag is recommended for this releaser. By default,
Waypoint prunes and destroys all unreleased deployments and keeps only one previous
deployment. Therefore, if ` + "`-prune=false`" + ` is not set, Waypoint may delete
your job via "pruning" a previous version. See [deployment pruning](/waypoint/docs/lifecycle/release#deployment-pruning)
for more information.

### Release URL

If you want the URL of the release of your deployment to be published in Waypoint,
you must set the meta 'waypoint.hashicorp.com/release_url' in your jobspec. The
value specified in this meta field will be published as the release URL for your
application. In the future, this may source from Consul.

`)

	doc.Example(`
// The waypoint.hcl file
release {
  use "nomad-jobspec-canary" {
    groups = [
      "app"
    ]
  }
}

// The app.nomad.tpl file
job "web" {
  datacenters = ["dc1"]

  group "app" {
    network {
      mode = "bridge"
      port "http" {
        to = 80
      }
    }

    // Setting a canary in the update stanza indicates a canary deployment
    update {
      max_parallel = 1
      canary       = 1
      auto_revert  = true
      auto_promote = false
      health_check = "task_states"
    }

    service {
      name = "app"
      port = 80
      connect {
        sidecar_service {}
      }
    }

    task "app" {
      driver = "docker"
      config {
        image = "${artifact.image}:${artifact.tag}"
        ports  = ["http"]
      }

      env {
        %{ for k,v in entrypoint.env ~}
        ${k} = "${v}"
        %{ endfor ~}

        // Ensure we set PORT for the URL service. This is only necessary
        // if we want the URL service to function.
        PORT = 80
      }
    }
  }

  group "app-gateway" {
    network {
      mode = "bridge"
      port "inbound" {
        static = 8080
        to     = 8080
      }
    }

    service {
      name = "gateway"
      port = "8080"

      connect {
        gateway {
          proxy {}

          ingress {
            listener {
              port = 8080
              protocol = "http"
              service {
                name  = "app"
                hosts = [ "*" ]
              }
            }
          }
        }
      }
    }
  }
  meta = {
    // Ensure we set meta for Waypoint to detect the release URL
    "waypoint.hashicorp.com/release_url" = "http://app.ingress.dc1.consul:8080"
  }
}
`)

	doc.SetField(
		"groups",
		"List of task group names which are to be promoted.",
	)

	doc.SetField(
		"fail_deployment",
		"If true, marks the deployment as failed.",
	)

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
