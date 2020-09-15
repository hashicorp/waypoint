package cloudrun

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	run "google.golang.org/api/run/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/datadir"
	"github.com/hashicorp/waypoint/sdk/docs"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

const (
	// todo: make user configurable
	deployRegion = "us-central1"
)

// Platform is the Platform implementation for Google Cloud Run.
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

// ValidateAuthFunc implements component.Authenticator
func (p *Platform) ValidateAuthFunc() interface{} {
	return p.ValidateAuth
}

// AuthFunc implements component.Authenticator
func (p *Platform) AuthFunc() interface{} {
	return p.Auth
}

func (p *Platform) Auth() error {
	return nil
}

func (p *Platform) ValidateAuth(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	dir *datadir.Component,
	ui terminal.UI,
) error {
	apiService, err := run.NewService(ctx,
		option.WithEndpoint("https://"+deployRegion+"-run.googleapis.com"),
	)
	if err != nil {
		ui.Output("Error constructing api client: "+err.Error(), terminal.WithErrorStyle())
		return status.Errorf(codes.Aborted, err.Error())
	}

	// TODO auth doesn't work if the service isn't already created, which is a common case.
	// Until that is fixed, we disable the auth checks.
	return nil

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	client := run.NewProjectsLocationsServicesService(apiService)

	expectedPermissions := []string{
		"roles/run.admin",
	}

	// run.admin encompasses all the permissions we should need
	testReq := run.TestIamPermissionsRequest{
		Permissions: expectedPermissions,
	}

	// The resource we are checking permissions on
	apiResource := fmt.Sprintf("projects/%s/locations/%s/services/%s",
		p.config.Project,
		deployRegion,
		src.App,
	)

	st.Update("Testing IAM permissions...")
	result, err := client.TestIamPermissions(apiResource, &testReq).Do()
	if err != nil {
		st.Step(terminal.StatusError, "Error testing IAM permissions: "+err.Error())
		return err
	}

	// If our resulting permissions do not equal our expected permissions, auth does not validate
	if !reflect.DeepEqual(result.Permissions, expectedPermissions) {
		st.Step(terminal.StatusError, "Incorrect IAM permissions, received "+strings.Join(result.Permissions, ", "))
		return fmt.Errorf("incorrect IAM permissions, received %s", strings.Join(result.Permissions, ", "))
	}

	return nil
}

// Deploy deploys an image to GCR.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	dir *datadir.Component,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {
	// Start building our deployment since we use this information
	result := &Deployment{
		Resource: &Deployment_Resource{
			Location: deployRegion,
			Project:  p.config.Project,
			Name:     src.App,
		},
	}
	id, err := component.Id()
	if err != nil {
		return nil, err
	}
	result.Id = id

	apiService, err := run.NewService(ctx,
		option.WithEndpoint("https://"+result.Resource.Location+"-run.googleapis.com"),
	)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	// We need to determine if we're creating or updating a service. To
	// do this, we just query GCP directly. There is a bit of a race here
	// but we expect to own this service, so even if we get a delete/create
	// in the middle, we'll just error later.
	create := false
	client := run.NewNamespacesServicesService(apiService)
	log.Trace("checking if service already exists", "service", result.apiName())
	st.Update("Checking if service is already created")
	service, err := client.Get(result.apiName()).Context(ctx).Do()
	if err != nil {
		gerr, ok := err.(*googleapi.Error)
		if !ok {
			return nil, err
		}
		log.Trace("googleapi.Error value", "error", gerr)

		// If we have a 404 then we just haven't created it yet.
		if gerr.Code != 404 {
			return nil, err
		}

		create = true
		service = &run.Service{
			ApiVersion: "serving.knative.dev/v1",
			Kind:       "Service",
			Metadata: &run.ObjectMeta{
				Name: result.Resource.Name,
			},
			Spec: &run.ServiceSpec{},
		}
	}

	// If we're deploying to the "latest revision" then we want to enforce
	// we're only going to the last revision so that we don't release at the
	// same time. This happens because when we create the service we don't
	// specify a traffic target (because we don't have a revision ID yet) and
	// it defaults to latest revision.
	if len(service.Spec.Traffic) > 0 && service.Spec.Traffic[0].LatestRevision {
		service.Spec.Traffic = []*run.TrafficTarget{
			&run.TrafficTarget{
				RevisionName: service.Status.LatestCreatedRevisionName,
				Percent:      100,
			},
		}
	}

	// Create our env vars
	var env []*run.EnvVar
	for k, v := range deployConfig.Env() {
		env = append(env, &run.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	// Regardless of if we're creating or updating, we update our
	// spec to force a new revision.
	service.Spec.Template = &run.RevisionTemplate{
		Metadata: &run.ObjectMeta{
			Annotations: map[string]string{
				"waypoint.hashicorp.com/nonce": time.Now().UTC().Format(time.RFC3339Nano),
			},
		},
		Spec: &run.RevisionSpec{
			Containers: []*run.Container{
				&run.Container{
					Image: img.Name(),
					Env:   env,
				},
			},
		},
	}

	if create {
		// Create the service
		log.Info("creating the service")
		st.Update("Creating new Cloud Run service")
		service, err = client.Create("namespaces/"+result.Resource.Project, service).
			Context(ctx).Do()
		if err != nil {
			return nil, status.Errorf(codes.Aborted, err.Error())
		}
	} else {
		// Update
		log.Info("updating a pre-existing service", "service", result.apiName())
		st.Update("Deploying new Cloud Run revision")
		service, err = client.ReplaceService(result.apiName(), service).
			Context(ctx).Do()
		if err != nil {
			return nil, status.Errorf(codes.Aborted, err.Error())
		}
	}

	// Poll the service and wait for completion
	st.Update("Waiting for revision to be ready")
	service, err = result.pollServiceReady(ctx, log, apiService)
	if err != nil {
		return nil, err
	}

	// Now that the service is ready we can set the latest URL
	result.RevisionId = service.Status.LatestCreatedRevisionName
	result.Url = service.Status.Url

	// If we have tracing enabled we just dump the full service as we know it
	// in case we need to look up what the raw value is.
	if log.IsTrace() {
		bs, _ := json.Marshal(service)
		log.Trace("service JSON", "json", base64.StdEncoding.EncodeToString(bs))
	}

	return result, nil
}

// Destroy deletes the cloud run revision
func (p *Platform) Destroy(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) error {
	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	apiService, err := deployment.apiService(ctx)
	if err != nil {
		return err
	}
	client := run.NewNamespacesRevisionsService(apiService)

	st.Update("Deleting deployment...")

	_, err = client.Delete(deployment.apiRevisionName()).Context(ctx).Do()
	return err
}

// Config is the configuration structure for the Platform.
type Config struct {
	// Project is the project to deploy to.
	Project string `hcl:"project,attr"`

	// Unauthenticated, if set to true, will allow unauthenticated access
	// to your deployment. This defaults to true.
	Unauthenticated *bool `hcl:"unauthenticated,optional"`
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Deploy a container to Google Cloud Run")

	return doc, nil
}

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
)
