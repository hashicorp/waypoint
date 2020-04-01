package cloudrun

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/hashicorp/go-hclog"
	run "google.golang.org/api/run/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mitchellh/devflow/sdk/component"
	"github.com/mitchellh/devflow/sdk/terminal"
)

// Releaser is the ReleaseManager implementation for Google Cloud Run.
type Releaser struct {
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

// Release deploys an image to GCR.
func (r *Releaser) Release(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	targets []component.ReleaseTarget,
) (*Release, error) {
	// Get the deployments
	var deploys []*Deployment
	for _, t := range targets {
		var deploy Deployment
		if err := component.ProtoAnyUnmarshal(t.Deployment, &deploy); err != nil {
			return nil, err
		}

		deploys = append(deploys, &deploy)
	}

	// We use the most recent deploy for most things
	deploy := deploys[0]

	// Get the API service
	apiService, err := deploy.apiService(ctx)
	if err != nil {
		return nil, err
	}

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	// Get the service
	st.Update("Getting service information...")
	client := run.NewNamespacesServicesService(apiService)
	service, err := client.Get(deploy.apiName()).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	// Update the service with the traffic info
	service.Spec.Traffic = nil
	for i, t := range targets {
		log.Debug("setting traffic target",
			"revision", deploys[i].RevisionId,
			"percent", t.Percent)
		service.Spec.Traffic = append(service.Spec.Traffic, &run.TrafficTarget{
			RevisionName: deploys[i].RevisionId,
			Percent:      int64(t.Percent),
		})
	}

	// Replace the service
	st.Update("Deploying routing changes")
	service, err = client.ReplaceService(deploy.apiName(), service).
		Context(ctx).Do()
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	// Set the IAM policy so global traffic is allowed
	if err := r.setNoAuthPolicy(ctx, deploy, apiService); err != nil {
		return nil, err
	}

	// Poll the service and wait for completion
	st.Update("Waiting for revision to be ready")
	service, err = deploy.pollServiceReady(ctx, log, apiService)
	if err != nil {
		return nil, err
	}

	// If we have tracing enabled we just dump the full service as we know it
	// in case we need to look up what the raw value is.
	if log.IsTrace() {
		bs, _ := json.Marshal(service)
		log.Trace("service JSON", "json", base64.StdEncoding.EncodeToString(bs))
	}

	return &Release{
		Url: service.Status.Url,
	}, nil
}

// setNoAuthPolicy sets the IAM policy on the deployment so that anyone
// can access it (no auth required).
func (r *Releaser) setNoAuthPolicy(
	ctx context.Context,
	deployment *Deployment,
	apiService *run.APIService,
) error {
	client := run.NewProjectsLocationsServicesService(apiService)
	_, err := client.SetIamPolicy(deployment.apiResource(), &run.SetIamPolicyRequest{
		Policy: &run.Policy{
			Bindings: []*run.Binding{
				&run.Binding{
					Role:    "roles/run.invoker",
					Members: []string{"allUsers"},
				},
			},
		},
	}).Context(ctx).Do()
	if err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}

	return nil
}

// ReleaserConfig is the configuration structure for the Releaser.
type ReleaserConfig struct{}

func (r *Release) URL() string { return r.Url }

var (
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Configurable   = (*Releaser)(nil)
	_ component.Release        = (*Release)(nil)
)
