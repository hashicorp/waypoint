// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cloudrun

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/hashicorp/go-hclog"
	run "google.golang.org/api/run/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
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
		return nil, status.Errorf(codes.Aborted, "Unable to fetch service information from Google Cloud: %s", err.Error())
	}

	// Update the service with the traffic info. This code is laid out this
	// way to make it more trivial in the future to add traffic splitting.
	service.Spec.Traffic = nil
	log.Debug("Setting traffic target", "revision", target.RevisionId, "percent", 100)
	service.Spec.Traffic = append(service.Spec.Traffic, &run.TrafficTarget{
		RevisionName: target.RevisionId,
		Percent:      100,
	})

	// Replace the service
	st.Update("Deploying routing changes")
	service, err = client.ReplaceService(target.apiName(), service).Context(ctx).Do()
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Unable to deploy routing changes: %s", err.Error())
	}

	// Set the IAM policy so global traffic is allowed
	if err := r.setNoAuthPolicy(ctx, target, apiService); err != nil {
		return nil, status.Errorf(codes.Aborted, "Unable to set no auth policy: %s", err)
	}

	// Poll the service and wait for completion
	st.Update("Waiting for revision to be ready")
	service, err = target.pollServiceReady(ctx, log)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Timeout waiting for revision to be ready: %s", err)
	}

	// If we have tracing enabled we just dump the full service as we know it
	// in case we need to look up what the raw value is.
	if log.IsTrace() {
		bs, _ := json.Marshal(service)
		log.Trace("Service JSON", "json", base64.StdEncoding.EncodeToString(bs))
	}

	// Note: it is not possible to add custom domain mappings as there are a number of steps
	// such as adding the DNS record which is out of control of Waypoint

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
				{
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

func (r *Releaser) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&ReleaserConfig{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Manipulates the Cloud Run APIs to make deployments active")
	doc.Input("google.cloudrun.Deployment")
	doc.Output("google.cloudrun.Release")

	return doc, nil
}

func (r *Release) URL() string { return r.Url }

var (
	_ component.ReleaseManager = (*Releaser)(nil)
	_ component.Configurable   = (*Releaser)(nil)
	_ component.Release        = (*Release)(nil)
)
