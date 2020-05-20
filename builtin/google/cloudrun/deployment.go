package cloudrun

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/api/option"
	run "google.golang.org/api/run/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/sdk/component"
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

// pollServiceReady waits for the service to become ready.
func (d *Deployment) pollServiceReady(
	ctx context.Context,
	log hclog.Logger,
	apiService *run.APIService,
) (*run.Service, error) {
	log = log.With("service", d.Resource.Name)
	log.Info("waiting for cloud run service to be ready")
	client := run.NewNamespacesServicesService(apiService)
	for {
		log.Trace("querying service")
		service, err := client.Get(d.apiName()).Context(ctx).Do()
		if err != nil {
			return nil, status.Errorf(codes.Aborted, err.Error())
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
				return nil, status.Errorf(codes.Aborted, "service failed to get ready")
			}
		}

		time.Sleep(1 * time.Second)
	}
}

var _ component.Deployment = (*Deployment)(nil)
