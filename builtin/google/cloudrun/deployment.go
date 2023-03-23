// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudrun

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	run "google.golang.org/api/run/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
)

// apiResource returns the GCP API "resource" string format for API calls.
func (d *Deployment) apiResource() string {
	return fmt.Sprintf("projects/%s/locations/%s/services/%s",
		d.Resource.Project,
		d.Resource.Location,
		d.Resource.Name,
	)
}

// apiName returns the GCP API "name" string format for API calls.
func (d *Deployment) apiName() string {
	return fmt.Sprintf("namespaces/%s/services/%s",
		d.Resource.Project,
		d.Resource.Name,
	)
}

// apiRevisionName returns the GCP API "name" string format for API calls
// to the revisons api.
func (d *Deployment) apiRevisionName() string {
	return fmt.Sprintf("namespaces/%s/revisions/%s",
		d.Resource.Project,
		d.RevisionId,
	)
}

// apiService returns the API service for GCP client usage.
func (d *Deployment) apiService(ctx context.Context) (*run.APIService, error) {
	result, err := run.NewService(ctx,
		option.WithEndpoint("https://"+d.Resource.Location+"-run.googleapis.com"),
	)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	return result, nil
}

// iamAPIService returns the IAM API service for GCP client usage.
func (d *Deployment) iamAPIService(ctx context.Context) (*iam.Service, error) {
	result, err := iam.NewService(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	return result, nil
}

// getLocationsForProject returns the Cloud Run regions which are usable by this project
func (d *Deployment) getLocationsForProject(ctx context.Context) ([]*run.Location, error) {
	apiService, err := run.NewService(ctx)
	if err != nil {
		return nil, err
	}

	// Validate that the given region is available
	pClient := run.NewProjectsLocationsService(apiService)
	pls, err := pClient.List(fmt.Sprintf("projects/%s", d.Resource.Project)).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to list regions for project %s: %s", d.Resource.Project, err)
	}

	return pls.Locations, nil
}

// findServices finds the Cloud Run services in all regions
func (d *Deployment) findServicesForLocations(ctx context.Context, locations []*run.Location) (map[string]*run.Service, error) {
	services := map[string]*run.Service{}

	for _, l := range locations {
		apiService, err := run.NewService(ctx,
			option.WithEndpoint("https://"+l.LocationId+"-run.googleapis.com"),
		)
		if err != nil {
			return nil, err
		}

		client := run.NewNamespacesServicesService(apiService)

		service, err := client.Get(d.apiName()).Context(ctx).Do()
		if err != nil {
			gerr, ok := err.(*googleapi.Error)
			if !ok {
				return nil, err
			}

			// If we have a 404 then we just haven't created it yet, everything else is an API error
			if gerr.Code != 404 {
				return nil, err
			}
		}

		if service != nil {
			services[l.LocationId] = service
		}
	}

	return services, nil
}

func (d *Deployment) replaceService(ctx context.Context, s *run.Service) (*run.Service, error) {
	apiService, err := d.apiService(ctx)
	if err != nil {
		return nil, err
	}

	client := run.NewNamespacesServicesService(apiService)

	return client.ReplaceService(d.apiName(), s).Context(ctx).Do()
}

func (d *Deployment) createService(ctx context.Context, s *run.Service) (*run.Service, error) {
	apiService, err := d.apiService(ctx)
	if err != nil {
		return nil, err
	}

	client := run.NewNamespacesServicesService(apiService)

	return client.Create("namespaces/"+d.Resource.Project, s).Context(ctx).Do()
}

// pollServiceReady waits for the service to become ready.
func (d *Deployment) pollServiceReady(
	ctx context.Context,
	log hclog.Logger,
) (*run.Service, error) {
	log = log.With("service", d.Resource.Name)
	log.Info("waiting for cloud run service to be ready")

	apiClient, err := d.apiService(ctx)
	if err != nil {
		return nil, err
	}

	client := run.NewNamespacesServicesService(apiClient)
	for {
		log.Trace("querying service")
		service, err := client.Get(d.apiName()).Context(ctx).Do()
		if err != nil {
			return nil, err
		}

		for _, cond := range service.Status.Conditions {
			if cond.Type != "Ready" {
				continue
			}

			log.Debug("ready status", "status", cond.Status)
			switch cond.Status {
			case "True":
				log.Info("service is ready")
				return service, nil

			case "False":
				return nil, fmt.Errorf("Cloud Run service failed to get ready")
			}
		}

		time.Sleep(1 * time.Second)
	}
}

var _ component.Deployment = (*Deployment)(nil)
