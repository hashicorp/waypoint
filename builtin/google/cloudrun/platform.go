// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudrun

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	run "google.golang.org/api/run/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"
)

const (
	rmResourceServiceName = "service"
	platformName          = "gcp"
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
		return status.Errorf(codes.InvalidArgument, "Invalid configuration, expected *cloudrun.Config, got %s", reflect.TypeOf(config))
	}

	return validateConfig(*c)
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

// DestroyWorkspaceFunc implements component.WorkspaceDestroyer
func (p *Platform) DestroyWorkspaceFunc() interface{} {
	return p.DestroyWorkspace
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

// DefaultReleaserFunc implements component.PlatformReleaser
func (p *Platform) DefaultReleaserFunc() interface{} {
	return func() *Releaser { return &Releaser{} }
}

func (p *Platform) ValidateAuth(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	ui terminal.UI,
) error {
	deployment := &Deployment{
		Resource: &Deployment_Resource{
			Location: p.config.Location,
			Project:  p.config.Project,
			Name:     src.App,
		},
	}

	apiService, err := deployment.apiService(ctx)
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
		p.config.Location,
		src.App,
	)

	st.Update("Testing Cloud Run IAM permissions...")
	result, err := client.TestIamPermissions(apiResource, &testReq).Do()
	if err != nil {
		st.Step(terminal.StatusError, "Error testing Cloud Run IAM permissions: "+err.Error())
		return err
	}

	// If our resulting permissions do not equal our expected permissions, auth does not validate
	if !reflect.DeepEqual(result.Permissions, expectedPermissions) {
		st.Step(terminal.StatusError, "Incorrect IAM permissions, received "+strings.Join(result.Permissions, ", "))
		return status.Errorf(codes.PermissionDenied, "incorrect IAM permissions, received %s", strings.Join(result.Permissions, ", "))
	}

	// Validate if user has access to the service account specified
	if p.config.ServiceAccountName != "" {

		iamAPIService, err := deployment.iamAPIService(ctx)
		if err != nil {
			ui.Output("Error constructing api client: "+err.Error(), terminal.WithErrorStyle())
			return status.Errorf(codes.Aborted, err.Error())
		}

		client := iam.NewProjectsServiceAccountsService(iamAPIService)

		expectedPermissions := []string{
			"iam.serviceAccounts.actAs",
		}

		// We need to ensure that the service creator has Service Account User role.
		testReq := iam.TestIamPermissionsRequest{
			Permissions: expectedPermissions,
		}

		apiResource := fmt.Sprintf("projects/%s/serviceAccounts/%s",
			p.config.Project,
			p.config.ServiceAccountName,
		)

		st.Update("Testing IAM permissions on the supplied service account...")
		result, err := client.TestIamPermissions(apiResource, &testReq).Do()
		if err != nil {
			st.Step(terminal.StatusError, "Error testing IAM permissions of the Service Account: "+err.Error())
			return err
		}

		// If our resulting permissions do not equal our expected permissions, auth does not validate
		if !reflect.DeepEqual(result.Permissions, expectedPermissions) {
			st.Step(terminal.StatusError, "Incorrect IAM permissions on the Service Account, received "+strings.Join(result.Permissions, ", "))
			return status.Errorf(codes.PermissionDenied, "Incorrect IAM permissions on the Service Account, received %s", strings.Join(result.Permissions, ", "))
		}
	}
	return nil
}

func (p *Platform) resourceManager(log hclog.Logger) *resource.Manager {
	return resource.NewManager(
		resource.WithLogger(log.Named("resource_manager")),
		resource.WithValueProvider(p.getApiService),
		resource.WithResource(resource.NewResource(
			resource.WithName(rmResourceServiceName),
			resource.WithState(&Resource_Service{}),
			resource.WithCreate(p.resourceServiceCreate),
			resource.WithDestroy(p.resourceServiceDestroy),

			resource.WithPlatform(platformName),
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER),
		)),
	)
}

func (p *Platform) getApiService(ctx context.Context, d *Deployment) (*run.APIService, error) {
	apiService, err := run.NewService(ctx,
		option.WithEndpoint("https://"+d.Resource.Location+"-run.googleapis.com"),
	)
	if err != nil {
		return nil, err
	}

	return apiService, nil
}

func (p *Platform) resourceServiceCreate(
	ctx context.Context,
	deployment *Deployment,
	log hclog.Logger,
	st terminal.Status,
	pls []*run.Location,
	deployConfig *component.DeploymentConfig,
	img *docker.Image,
	src *component.Source,
	state *Resource_Service,
) error {

	// We need to determine if we're creating or updating a service. To
	// do this, we just query GCP directly. There is a bit of a race here
	// but we expect to own this service, so even if we get a delete/create
	// in the middle, we'll just error later.
	create := false
	var service *run.Service

	log.Trace("checking if service has already exists", "service", deployment.apiName())
	st.Update("Checking if service already exists")

	// Is there a deployment for this service
	services, err := deployment.findServicesForLocations(ctx, pls)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "Unable to query existing Cloud Run services: %s", err)
	}

	// no services found create a new one
	if len(services) == 0 {
		create = true
		service = &run.Service{
			ApiVersion: "serving.knative.dev/v1",
			Kind:       "Service",
			Metadata: &run.ObjectMeta{
				Name: deployment.Resource.Name,
			},
			Spec: &run.ServiceSpec{},
		}
	} else if len(services) != 1 {
		// the service should only exist in a single region
		return status.Errorf(codes.FailedPrecondition, "Cloud Run services named '%s' exist in multiple regions. Please remove any manually created services.", src.App)
	} else {
		// Loop through the regions which contain services and ensure that
		// the current deployment is not in a different region to an existing service
		for k := range services {
			if k != p.config.Location {
				// Waypoint can not change the region of a service so return an error.
				return status.Errorf(codes.AlreadyExists, "The Cloud Run service '%s' already exists in the region '%s', Waypoint is unable to change the region of a deployed service", src.App, k)
			}
		}

		service = services[p.config.Location]
	}

	// If we're deploying to the "latest revision" then we want to enforce
	// we're only going to the last revision so that we don't release at the
	// same time. This happens because when we create the service we don't
	// specify a traffic target (because we don't have a revision ID yet) and
	// it defaults to latest revision.
	if len(service.Spec.Traffic) > 0 && service.Spec.Traffic[0].LatestRevision {
		service.Spec.Traffic = []*run.TrafficTarget{
			{
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
				{
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

	if p.config.Capacity != nil {
		if p.config.Capacity.MaxRequestsPerContainer > 0 {
			service.Spec.Template.Spec.ContainerConcurrency = int64(p.config.Capacity.MaxRequestsPerContainer)
		}

		if p.config.Capacity.Memory > 0 {
			// Requires value expressed as Kubernetes Quantity
			// (https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/api/resource/quantity.go)
			resources.Limits["memory"] = fmt.Sprintf("%dMi", p.config.Capacity.Memory)
		}

		if p.config.Capacity.CPUCount > 0 {
			// Can only be 1 or 2
			resources.Limits["cpu"] = fmt.Sprintf("%d", p.config.Capacity.CPUCount)
		}

		if p.config.Capacity.RequestTimeout > 0 {
			// Max value of 900
			service.Spec.Template.Spec.TimeoutSeconds = int64(p.config.Capacity.RequestTimeout)
		}
	}

	if p.config.StaticEnvVars != nil {
		env = service.Spec.Template.Spec.Containers[0].Env
		for k, v := range p.config.StaticEnvVars {
			env = append(env, &run.EnvVar{Name: k, Value: v})
		}
		service.Spec.Template.Spec.Containers[0].Env = env
	}

	if p.config.AutoScaling != nil {
		if p.config.AutoScaling.Max > 0 {
			service.Spec.Template.Metadata.Annotations["autoscaling.knative.dev/maxScale"] = fmt.Sprintf("%d", p.config.AutoScaling.Max)
		}

		/*
			// Not yet implemented by Cloud Run
			if p.config.AutoScaling.Min > 0 {
				service.Spec.Template.Metadata.Annotations["autoscaling.knative.dev/minScale"] = fmt.Sprintf("%d", p.config.AutoScaling.Min)
			}
		*/
	}

	if p.config.ServiceAccountName != "" {
		service.Spec.Template.Spec.ServiceAccountName = p.config.ServiceAccountName
	}

	if len(p.config.CloudSQLInstances) > 0 {
		service.Spec.Template.Metadata.Annotations["run.googleapis.com/cloudsql-instances"] = strings.Join(p.config.CloudSQLInstances, ",")
	}

	if p.config.VPCAccess != nil {
		if p.config.VPCAccess.Connector != "" {
			service.Spec.Template.Metadata.Annotations["run.googleapis.com/vpc-access-connector"] = p.config.VPCAccess.Connector
		}
		if p.config.VPCAccess.Egress != "" {
			service.Spec.Template.Metadata.Annotations["run.googleapis.com/vpc-access-egress"] = p.config.VPCAccess.Egress
		}
	}

	if create {
		// Create the service
		log.Info("creating the service")
		st.Update("Creating new Cloud Run service")

		service, err = deployment.createService(ctx, service)
		if err != nil {
			return status.Errorf(codes.Aborted, "Unable to create Cloud Run service: %s", err.Error())
		}
	} else {
		// Update
		log.Info("updating a pre-existing service", "service", deployment.apiName())
		st.Update("Deploying new Cloud Run revision")

		service, err = deployment.replaceService(ctx, service)
		if err != nil {
			return status.Errorf(codes.Aborted, "Unable to deploy new Cloud Run revision: %s", err.Error())
		}
	}

	// Poll the service and wait for completion
	st.Update("Waiting for revision to be ready")
	service, err = deployment.pollServiceReady(ctx, log)
	if err != nil {
		return err
	}

	// Now that the service is ready we can set the latest URL
	deployment.RevisionId = service.Status.LatestCreatedRevisionName
	deployment.Url = service.Status.Url

	// Set state
	state.ApiName = deployment.apiName()
	state.RevisionName = deployment.apiRevisionName()

	// If we have tracing enabled we just dump the full service as we know it
	// in case we need to look up what the raw value is.
	if log.IsTrace() {
		bs, _ := json.Marshal(service)
		log.Trace("service JSON", "json", base64.StdEncoding.EncodeToString(bs))
	}
	return nil
}

func (p *Platform) resourceServiceDestroy(
	ctx context.Context,
	apiService *run.APIService,
	st terminal.Status,
	isWorkspaceDestroy bool,
	d *Deployment,
	state *Resource_Service,
) error {
	if isWorkspaceDestroy {
		client := run.NewNamespacesServicesService(apiService)
		st.Step("Deleting service: %s", state.ApiName)

		_, err := client.Delete(d.apiName()).Context(ctx).Do()
		return err
	}
	rn := state.RevisionName
	client := run.NewNamespacesRevisionsService(apiService)
	st.Step("Deleting deployment with revision: %s", rn)

	// First check if it's currently serving traffic
	r, err := client.Get(rn).Context(ctx).Do()
	if err != nil {
		return err
	}
	// the slice of status conditions returned by GCP stores the reason
	// for the most recent status transition at index [1]. We check that
	// the revision is serving by verifying the operating status is active
	// https://cloud.google.com/run/docs/reference/rest/v1/namespaces.revisions
	cs := r.Status.Conditions[1]
	if cs.Status == "True" {
		st.Step(
			terminal.StatusWarn,
			fmt.Sprintf("Cannot destroy deployment with revision %q, as revision is actively receiving traffic. Disregard if destroying workspace.", rn),
		)
		return nil
	}

	_, err = client.Delete(state.RevisionName).Context(ctx).Do()
	return err
}

// Deploy deploys an image to Cloud Run.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {
	// Start building our deployment since we use this information
	deployment := &Deployment{
		Resource: &Deployment_Resource{
			Location: p.config.Location,
			Project:  p.config.Project,
			Name:     src.App,
		},
	}
	id, err := component.Id()
	if err != nil {
		return nil, err
	}
	deployment.Id = id

	// Validate that the Docker image is stored in a GCP registry
	// It is not possible to deploy to Cloud Run using external container registries
	err = validateImageName(img.Image)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	// Get the Cloud Run Locations available for this project
	// Cloud Run is only available in a limited number of Locations, this may be further restricted
	// by the users access
	pls, err := deployment.getLocationsForProject(ctx)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	// Validate that the Location specified for the deployment is available for the project
	err = validateLocationAvailable(p.config.Location, pls)
	if err != nil {
		return nil, err
	}

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	rm := p.resourceManager(log)
	if err := rm.CreateAll(
		ctx, deployment, log, st, pls,
		deployConfig, img, src,
	); err != nil {
		return nil, err
	}

	// Store our resource state
	deployment.ResourceState = rm.State()

	// Get our service state
	servState := rm.Resource(rmResourceServiceName).State().(*Resource_Service)
	if servState == nil {
		return nil, status.Errorf(codes.Internal,
			"service state is nil, this should never happen")
	}

	return deployment, nil
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

	rm := p.resourceManager(log)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if deployment.ResourceState == nil {
		rm.Resource(rmResourceServiceName).SetState(&Resource_Service{
			ApiName:      deployment.apiName(),
			RevisionName: deployment.apiRevisionName(),
		})
	} else {
		// Load state
		if err := rm.LoadState(deployment.ResourceState); err != nil {
			return err
		}
	}

	return rm.DestroyAll(ctx, st, false, deployment)
}

func (p *Platform) DestroyWorkspace(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) error {
	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	rm := p.resourceManager(log)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if deployment.ResourceState == nil {
		rm.Resource(rmResourceServiceName).SetState(&Resource_Service{
			ApiName:      deployment.apiName(),
			RevisionName: deployment.apiRevisionName(),
		})
	} else {
		// Load state
		if err := rm.LoadState(deployment.ResourceState); err != nil {
			return err
		}
	}

	return rm.DestroyAll(ctx, st, true, deployment)
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}), docs.FromFunc(p.DeployFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description("Deploy a container to Google Cloud Run")
	doc.Input("docker.Image")
	doc.Output("google.cloudrun.Deployment")
	doc.Example(
		`
project = "wpmini"

app "wpmini" {
  labels = {
    "service" = "wpmini",
    "env"     = "dev"
  }

  build {
    use "pack" {}

    registry {
      use "docker" {
        image = "gcr.io/waypoint-project-id/wpmini"
        tag   = "latest"
      }
    }
  }

  deploy {
    use "google-cloud-run" {
      project  = "waypoint-project-id"
      location = "europe-north1"

      port = 5000

      static_environment = {
        "NAME" : "World"
      }

      capacity {
        memory                     = 128
        cpu_count                  = 2
        max_requests_per_container = 10
        request_timeout            = 300
      }

	  service_account_name = "cloudrun@waypoint-project-id.iam.gserviceaccount.com"

      auto_scaling {
        max = 10
      }

      cloudsql_instances = ["waypoint-project-id:europe-north1:sql-instance"]

      vpc_access {
        connector = "custom-vpc-connector"
        egress = "all"
      }
    }
  }

  release {
    use "google-cloud-run" {}
  }
}
`)

	doc.SetField(
		"project",
		"GCP project ID where the Cloud Run instance will be deployed.",
	)

	doc.SetField(
		"location",
		"GCP location, e.g. europe-north-1.",
	)

	doc.SetField(
		"unauthenticated",
		"Is public unauthenticated access allowed for the Cloud Run instance?",
	)

	doc.SetField(
		"port",
		"The port your application listens on.",
	)

	doc.SetField(
		"static_environment",
		"Additional environment variables to be added to the Cloud Run instance.",
	)

	doc.SetField(
		"capacity",
		"CPU, Memory, and resource limits for each Cloud Run instance.",

		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"memory",
				"Memory to allocate the Cloud Run instance specified in MB, min 128, max 4096.",
				docs.Default("128"),
			)

			doc.SetField(
				"cpu_count",
				"Number of CPUs to allocate the Cloud Run instance, min 1, max 2.",
				docs.Default("1"),
			)

			doc.SetField(
				"request_timeout",
				"Maximum time a request can take before timing out, max 900.",
				docs.Default("300"),
			)

			doc.SetField(
				"max_requests_per_container",
				"Maximum number of concurrent requests each instance can handle. When the maximum requests are exceeded, Cloud Run will create an additional instance.",
				docs.Default("80"),
			)
		}),
	)

	doc.SetField(
		"auto_scaling",
		"Configuration to control the auto scaling parameters for Cloud Run.",

		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"max",
				`Maximum number of Cloud Run instances. When the maximum requests per container is exceeded, Cloud Run will create an additional container instance to handle load.
				This parameter controls the maximum number of instances that can be created.`,
				docs.Default("1000"),
			)
		}),
	)

	doc.SetField(
		"service_account_name",
		"Specify a service account email that Cloud Run will use to run the service. You must have the `iam.serviceAccounts.actAs` permission on the service account.",
	)

	doc.SetField(
		"cloudsql_instances",
		"Specify list of CloudSQL instances that the Cloud Run instance will have access to.",
	)

	doc.SetField(
		"vpc_access",
		"VPCAccess details",
		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"connector",
				"Set VPC Access Connector for the Cloud Run instance.",
			)
			doc.SetField(
				"egress",
				"Set VPC egress. Supported values are 'all' and 'private-ranges-only'.",
			)
		}),
	)

	return doc, nil
}

// Config is the configuration structure for the Platform.
// Validation tags are provided by Go Pkg Validator
// https://pkg.go.dev/gopkg.in/go-playground/validator.v10?tab=doc
type Config struct {
	// Project is the project to deploy to
	Project string `hcl:"project,attr"`

	// Location	represents the Google Cloud location where the application will be deployed
	// e.g. us-west1
	Location string `hcl:"location,attr"`

	// Unauthenticated, if set to true, will allow unauthenticated access
	// to your deployment. This defaults to true.
	Unauthenticated *bool `hcl:"unauthenticated,optional"`

	// Port the applications is listening on.
	Port int `hcl:"port,optional"`

	// Environment variables that are meant to configure the application in a static
	// way. This might be control an image that has multiple modes of operation,
	// selected via environment variable. Most configuration should use the waypoint
	// config commands.
	StaticEnvVars map[string]string `hcl:"static_environment,optional"`

	// Capacity details for cloud run container.
	Capacity *Capacity `hcl:"capacity,block"`

	// AutoScaling details.
	AutoScaling *AutoScaling `hcl:"auto_scaling,block"`

	// Service Account details
	ServiceAccountName string `hcl:"service_account_name,optional"`

	// CloudSQLInstances that the application can connect to.
	CloudSQLInstances []string `hcl:"cloudsql_instances,optional"`

	// VPCAccess details.
	VPCAccess *VPCAccess `hcl:"vpc_access,block"`
}

// Capacity defines configuration for deployed Cloud Run resources
type Capacity struct {
	// Memory to allocate to the container specified in MB, min 128, max 4096.
	// Default value of 0 sets memory to 128MB which is default Cloud Run behaviour
	Memory int `hcl:"memory,attr" validate:"eq=0|gte=128,lte=4096"`
	// CPUCount is the number CPUs to allocate to a Cloud Run instance.
	CPUCount int `hcl:"cpu_count,attr" validate:"gte=0,lte=2"`
	// Maximum request time in seconds, max 900.
	RequestTimeout int `hcl:"request_timeout,attr" validate:"gte=0,lte=900"`
	// Maximum number of concurrent requests per container instance.
	// When max requests is exceeded Cloud Run will scale the number of containers.
	MaxRequestsPerContainer int `hcl:"max_requests_per_container,attr" validate:"gte=0"`
}

// AutoScaling defines the parameters which the Cloud Run instance can AutoScale.
// Currently only the maximum bound is supported
type AutoScaling struct {
	//Min int `hcl:"min,attr"` // not yet supported by cloud run
	Max int `hcl:"max,attr" validate:"gte=0"`
}

type VPCAccess struct {
	// Connector sets Serverless VPC Access connector for the application.
	Connector string `hcl:"connector,optional"`

	// Egress sets VPC egress type. Supported values are `all` and `private-ranges-only`.
	Egress string `hcl:"egress,optional" validate:"oneof=all private-ranges-only"`
}

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
)
