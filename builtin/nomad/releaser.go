package nomad

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
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

// getNomadClient provides
// the client connection used by resources to interact with Nomad.
func (r *Release) getNomadClient() (*api.Client, error) {
	// Get our client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Release promotes the Nomad canary deployment
func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
	target *Deployment,
) (*Release, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	// Update user in real time
	st := ui.Status()
	defer st.Close()

	jobClient := client.Jobs()
	deploymentClient := client.Deployments()
    st.Update("Getting job...")
	jobs, _, err := jobClient.PrefixList(target.Name)
    if err != nil {
        return nil, status.Errorf(codes.Aborted, "Unable to fetch Nomad job: %s", err.Error())
    }

	q := &api.QueryOptions{Namespace: jobs[0].JobSummary.Namespace}
	st.Update("Getting latest deployments for job")
	deploy, _, err := jobClient.LatestDeployment(jobs[0].ID, q)
    if err != nil {
        return nil, status.Errorf(codes.Aborted, "Unable to fetch latest deployment for Nomad job: %s", err.Error())
    }

	if deploy == nil {
	    st.Update("No active deployment for Nomad job")
		return &Release{}, nil
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
    if !canaryDeployment {
        return &Release{}, nil
    }

	// Set write options
	wq := &api.WriteOptions{Namespace: jobs[0].JobSummary.Namespace}

	var u *api.DeploymentUpdateResponse
	//TODO: Create some mechanism to loop until the canary allocs are healthy
	//      Nomad prohibits promotion otherwise
	//TODO: Add logic to support promotion of specific group(s)
	u, _, err = deploymentClient.PromoteAll(deploy.ID, wq)
	st.Update(fmt.Sprintf("Monitoring evaluation %q", u.EvalID))

	if err := NewMonitor(st, client).Monitor(u.EvalID); err != nil {
		return nil, err
	}

	//TODO: If applicable, get Consul service from job. If multiple services, how to determine which service to use
	//      (maybe from ReleaserConfig)? Consul service URL structure may be ambiguous here as well:
	//      'service_name.service.consul' is common; however, `.consul` is default domain for Consul, but this is not
	//      mandatory. The service could also be an ingress gateway, where the name would be service_name.ingress.consul.
	//      The Consul data center may also be required, and/or tags, for FQDN of:
	//      tag_name.service_name.ingress/service.datacenter.consul
	//      https://www.consul.io/docs/discovery/dns#standard-lookup
	//      If no Consul service, select IP/Port of a random instance?
	return &Release{
		Url: "https://waypointproject.io",
	}, nil
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
