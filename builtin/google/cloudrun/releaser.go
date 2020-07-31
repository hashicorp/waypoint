package cloudrun

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/hashicorp/go-hclog"
	run "google.golang.org/api/run/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/terminal"
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
	target *Deployment,
) (*Release, error) {
	// Get the API service
	apiService, err := target.apiService(ctx)
	if err != nil {
		return nil, err
	}

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	// Get the service
	st.Update("Getting service information...")
	client := run.NewNamespacesServicesService(apiService)
	service, err := client.Get(target.apiName()).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	// Update the service with the traffic info. This code is laid out this
	// way to make it more trivial in the future to add traffic splitting.
	service.Spec.Traffic = nil
	log.Debug("setting traffic target",
		"revision", target.RevisionId,
		"percent", 100)
	service.Spec.Traffic = append(service.Spec.Traffic, &run.TrafficTarget{
		RevisionName: target.RevisionId,
		Percent:      100,
	})

	// Replace the service
	st.Update("Deploying routing changes")
	service, err = client.ReplaceService(target.apiName(), service).
		Context(ctx).Do()
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	// Set the IAM policy so global traffic is allowed
	if err := r.setNoAuthPolicy(ctx, target, apiService); err != nil {
		return nil, err
	}

	// Poll the service and wait for completion
	st.Update("Waiting for revision to be ready")
	service, err = target.pollServiceReady(ctx, log, apiService)
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
