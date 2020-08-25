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
	"github.com/kr/pretty"
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

// Platform is the Platform implementation for Google Cloud Run.
type Platform struct {
	config Config
}

// ConfigSet is called after a configuration has been decoded
// we can use this to validate the config
func (p *Platform) ConfigSet(config interface{}) error {
	c, ok := config.(*Config)
	if !ok {
		// this should never happen
		return fmt.Errorf("Invalid configuration, expected *cloudrun.Config, got %s", reflect.TypeOf(config))
	}

	return ValidateConfig(*c)
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

	apiService, err := getAPIService(ctx, p.config.Region)
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
		p.config.Region,
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
	// Validate that the Docker image is stored in a GCP registry
	// It is not possible to deploy to Cloud Run using external container registries
	err := ValidateImageName(img.Image, p.config.Project)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	apiService, err := getAPIService(ctx, "")
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	// Validate that the given region is available
	pClient := run.NewProjectsLocationsService(apiService)
	pls, err := pClient.List(fmt.Sprintf("projects/%s", p.config.Project)).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to list regions for project %s: %s", p.config.Project, err)
	}

	// validate the region defined in the config is available for deployment
	err = ValidateRegionAvailable(p.config.Region, pls.Locations)
	if err != nil {
		return nil, err
	}

	// Start building our deployment since we use this information
	result := &Deployment{
		Resource: &Deployment_Resource{
			Location: p.config.Region,
			Project:  p.config.Project,
			Name:     src.App,
		},
	}
	id, err := component.Id()
	if err != nil {
		return nil, err
	}
	result.Id = id

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	// We need to determine if we're creating or updating a service. To
	// do this, we just query GCP directly. There is a bit of a race here
	// but we expect to own this service, so even if we get a delete/create
	// in the middle, we'll just error later.
	create := false
	var service *run.Service

	log.Trace("checking if service has already exists", "service", result.apiName())
	st.Update("Checking if service already exists")

	// Is there a deployment for this service
	services, err := findServices(ctx, result.apiName(), pls.Locations)
	if err != nil {
		return nil, fmt.Errorf("Unable to query existing Cloud Run services: %s", err)
	}

	// no services found create a new one
	if len(services) == 0 {
		create = true
		service = &run.Service{
			ApiVersion: "serving.knative.dev/v1",
			Kind:       "Service",
			Metadata: &run.ObjectMeta{
				Name: result.Resource.Name,
			},
			Spec: &run.ServiceSpec{},
		}
	} else if len(services) != 1 {
		// the service should only exist in a single region
		return nil, fmt.Errorf("Cloud Run services named '%s' exist in multiple regions. Please remove any manually created services.", src.App)
	} else {
		// Loop through the regions which contain services and ensure that
		// the current deployment is not in a different region to an existing service
		for k := range services {
			if k != p.config.Region {
				// Waypoint can not change the region of a service so return an error.
				return nil, fmt.Errorf("The Cloud Run service '%s' already exists in the region '%s', Waypoint is unable to change the region of a deployed service", src.App, k)
			}
		}

		service = services[p.config.Region]
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

	// define resources
	// values must adhere to: https: //github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/api/resource/quantity.go
	resources := &run.ResourceRequirements{
		Limits: map[string]string{},
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
					Image:     img.Name(),
					Env:       env,
					Resources: resources,
				},
			},
		},
	}

	// override the defaults if provided in config
	if p.config.Port > 0 {
		service.Spec.Template.Spec.Containers[0].Ports = []*run.ContainerPort{{ContainerPort: int64(p.config.Port)}}
	}

	if p.config.Capacity.MaxRequestsPerContainer > 0 {
		service.Spec.Template.Spec.ContainerConcurrency = int64(p.config.Capacity.MaxRequestsPerContainer)
	}

	if p.config.Capacity.Memory != "" {
		resources.Limits["memory"] = p.config.Capacity.Memory
	}

	if p.config.Capacity.CPUCount > 0 {
		resources.Limits["cpu"] = fmt.Sprintf("%d", p.config.Capacity.CPUCount)
	}

	if p.config.Capacity.RequestTimeout > 0 {
		service.Spec.Template.Spec.TimeoutSeconds = int64(p.config.Capacity.RequestTimeout)
	}

	if p.config.AutoScaling.Max > 0 {
		service.Spec.Template.Metadata.Annotations["autoscaling.knative.dev/maxScale"] = fmt.Sprintf("%d", p.config.AutoScaling.Max)
	}

	if p.config.Env != nil {
		env = service.Spec.Template.Spec.Containers[0].Env
		for k, v := range p.config.Env {
			env = append(env, &run.EnvVar{Name: k, Value: v})
		}
		service.Spec.Template.Spec.Containers[0].Env = env
	}

	/*
		// Not yet implemented by Cloud Run
		if p.config.AutoScaling.Min > 0 {
			service.Spec.Template.Metadata.Annotations["autoscaling.knative.dev/minScale"] = fmt.Sprintf("%d", p.config.AutoScaling.Min)
		}
	*/

	apiService, err = getAPIService(ctx, p.config.Region)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}
	client := run.NewNamespacesServicesService(apiService)

	if create {
		// Create the service
		log.Info("creating the service")
		st.Update("Creating new Cloud Run service")

		service, err = client.Create("namespaces/"+result.Resource.Project, service).
			Context(ctx).Do()
		if err != nil {
			return nil, status.Errorf(codes.Aborted, fmt.Sprintf("Unable to create Cloud Run service: %s", err.Error()))
		}
	} else {
		// Update
		log.Info("updating a pre-existing service", "service", result.apiName())
		st.Update("Deploying new Cloud Run revision")
		service, err = client.ReplaceService(result.apiName(), service).
			Context(ctx).Do()
		if err != nil {
			if gerr, ok := err.(*googleapi.Error); ok {
				log.Debug("Google error", "error", pretty.Sprint(gerr))
			}

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

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}))
	if err != nil {
		return nil, err
	}

	doc.Description("Deploy a container to Google Cloud Run")

	return doc, nil
}

// Config is the configuration structure for the Platform.
// Validation tags are provided by Go Pkg Validator
// https://pkg.go.dev/gopkg.in/go-playground/validator.v10?tab=doc
type Config struct {
	// Project is the project to deploy to.
	Project string `hcl:"project,attr"`
	// Region	is the GCP region to deploy to
	Region string `hcl:"region,attr"`

	// Unauthenticated, if set to true, will allow unauthenticated access
	// to your deployment. This defaults to true.
	Unauthenticated *bool `hcl:"unauthenticated,optional"`

	Port int `hcl:"port,optional"`

	Env map[string]string `hcl:"env,optional"`

	// Capacity details for cloud run container
	Capacity *Capacity `hcl:"capacity,block"`

	// AutoScaling details
	AutoScaling *AutoScaling `hcl:"auto_scaling,block"`
}

type Capacity struct {
	Memory                  string `hcl:"memory,attr" validate:"kubernetes-memory"`
	CPUCount                int    `hcl:"cpu_count,attr" validate:"gte=0,lte=2"`
	RequestTimeout          int    `hcl:"request_timeout,attr" validate:"gte=0,lte=900"`
	MaxRequestsPerContainer int    `hcl:"max_requests_per_container,attr" validate:"gte=0"`
}

type AutoScaling struct {
	//Min int `hcl:"min,attr"` // not yet supported by cloud run
	Max int `hcl:"max,attr" validate:"gte=0"`
}
