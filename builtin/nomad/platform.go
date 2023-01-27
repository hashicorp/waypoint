package nomad

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
)

const (
	metaId            = "waypoint.hashicorp.com/id"
	metaNonce         = "waypoint.hashicorp.com/nonce"
	rmResourceJobName = "job"
)

var (
	// default resources used for the deployed app. Can be overridden
	// through the resources stanza in a deploy. Note that these are the same defaults
	// used currently in Nomad if left unconfigured.
	defaultResourcesCPU      = 100
	defaultResourcesMemoryMB = 300
)

// Platform is the Platform implementation for Nomad.
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

func (p *Platform) ValidateAuth() error {
	return nil
}

// StatusFunc implements component.Status
func (p *Platform) StatusFunc() interface{} {
	return p.Status
}

func (p *Platform) resourceManager(log hclog.Logger) *resource.Manager {
	return resource.NewManager(
		resource.WithLogger(log.Named("resource_manager")),
		resource.WithValueProvider(getNomadClient),
		resource.WithResource(resource.NewResource(
			resource.WithName(rmResourceJobName),
			resource.WithState(&Resource_Job{}),
			resource.WithCreate(p.resourceJobCreate),
			resource.WithDestroy(p.resourceJobDestroy),
		)),
	)
}

// getNomadClient is a value provider for our resource manager and provides
// the client connection used by resources to interact with Nomad.
func getNomadClient() (*api.Client, error) {
	// Get our client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (p *Platform) resourceJobCreate(
	ctx context.Context,
	client *api.Client,
	result *Deployment,
	deployConfig *component.DeploymentConfig,
	img *docker.Image,
	st terminal.Status,
	state *Resource_Job,
) error {
	jobclient := client.Jobs()

	if p.config.ServicePort == 0 {
		p.config.ServicePort = 3000
	}

	if p.config.Datacenter == "" {
		p.config.Datacenter = "dc1"
	}

	// Determine if we have a job that we manage already
	job, _, err := jobclient.Info(result.Name, &api.QueryOptions{})
	if strings.Contains(err.Error(), "job not found") {
		job = api.NewServiceJob(result.Name, result.Name, p.config.Region, 10)
		job.Datacenters = []string{p.config.Datacenter}
		tg := api.NewTaskGroup(result.Name, 1)
		tg.Networks = []*api.NetworkResource{
			{
				Mode: "host",
				DynamicPorts: []api.Port{
					{
						Label: "waypoint",
						To:    int(p.config.ServicePort),
					},
				},
			},
		}

		// Register service with app deployment
		tg.Services = []*api.Service{
			{
				Name:      result.Name,
				PortLabel: "waypoint", // matches dynamic port label in NetworkResource
				Provider:  p.config.ServiceProvider,
			},
		}

		if p.config.Namespace == "" {
			p.config.Namespace = "default"
		}
		job.Namespace = &p.config.Namespace
		job.AddTaskGroup(tg)
		task := &api.Task{
			Name:   result.Name,
			Driver: "docker",
		}

		if p.config.Resources != nil {
			task.Resources = &api.Resources{
				CPU:      p.config.Resources.CPU,
				MemoryMB: p.config.Resources.MemoryMB,
			}
		}

		tg.AddTask(task)
		err = nil
	}
	if err != nil {
		return err
	}

	// Build our env vars
	env := map[string]string{
		"PORT": fmt.Sprint(p.config.ServicePort),
	}

	for k, v := range p.config.StaticEnvVars {
		env[k] = v
	}

	for k, v := range deployConfig.Env() {
		env[k] = v
	}

	// If no count is specified, presume that the user is managing the replica
	// count some other way (perhaps manual scaling, perhaps a pod autoscaler).
	// Either way if they don't specify a count, we should be sure we don't send one.
	if p.config.Count > 0 {
		job.TaskGroups[0].Count = &p.config.Count
	}

	// Set our ID on the meta.
	job.SetMeta(metaId, result.Id)
	job.SetMeta(metaNonce, time.Now().UTC().Format(time.RFC3339Nano))

	config := map[string]interface{}{
		"image": img.Name(),
		"ports": []string{"waypoint"},
	}

	if p.config.Auth != nil {
		config["auth"] = map[string]interface{}{
			"username": p.config.Auth.Username,
			"password": p.config.Auth.Password,
		}
	}

	job.TaskGroups[0].Tasks[0].Config = config
	job.TaskGroups[0].Tasks[0].Env = env

	// Get Consul ACL token from environment
	c, err := ConsulAuth()
	if err != nil {
		return err
	}
	job.ConsulToken = &c

	// Get Vault token from environment
	v, err := VaultAuth()
	if err != nil {
		return err
	}
	job.VaultToken = &v

	// Register job
	st.Update("Registering job...")
	regResult, _, err := jobclient.Register(job, nil)
	if err != nil {
		return err
	}

	// Store our state so we can destroy it properly
	state.Name = result.Name
	st.Step(terminal.StatusOK, "Job registration successful")

	// Wait on the allocation
	evalID := regResult.EvalID
	st.Update(fmt.Sprintf("Monitoring evaluation %q", evalID))

	if err := NewMonitor(st, client).Monitor(evalID); err != nil {
		return err
	}

	return nil
}

func (p *Platform) resourceJobDestroy(
	client *api.Client,
	state *Resource_Job,
	st terminal.Status,
) error {
	st.Step("Deleting job: %s", state.Name)
	_, _, err := client.Jobs().Deregister(state.Name, true, nil)
	return err
}

// Deploy deploys an image to Nomad.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {
	// Create our deployment and set an initial ID
	var result Deployment
	id, err := component.Id()
	if err != nil {
		return nil, err
	}

	// TODO(briancain): Update to use sequence number instead, and also append
	// project name to prevent any app name collisions like having two apps with two versions (go-2)
	result.Id = id
	result.Name = strings.ToLower(fmt.Sprintf("%s-%s", src.App, id))

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	rm := p.resourceManager(log)
	if err := rm.CreateAll(
		ctx, deployConfig, &result, img, st,
	); err != nil {
		return nil, err
	}

	// Store our resource state
	result.ResourceState = rm.State()

	// Get our service state
	servState := rm.Resource(rmResourceJobName).State().(*Resource_Job)
	if servState == nil {
		return nil, status.Errorf(codes.Internal,
			"service state is nil, this should never happen")
	}

	st.Step(terminal.StatusOK, "Deployment successfully rolled out!")

	return &result, nil
}

// Destroy deletes the Nomad job.
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
		rm.Resource(rmResourceJobName).SetState(&Resource_Job{
			Name: deployment.Name,
		})
	} else {
		// Load state
		if err := rm.LoadState(deployment.ResourceState); err != nil {
			return err
		}
	}

	return rm.DestroyAll(st)
}

func (p *Platform) Status(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) (*sdk.StatusReport, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	jobclient := client.Jobs()

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Gathering health report for Nomad platform...")
	defer func() { s.Abort() }()

	// Create our status report
	var result sdk.StatusReport
	result.External = true

	log.Debug("querying nomad for job health")

	job, _, err := jobclient.Info(deployment.Name, &api.QueryOptions{})
	if err != nil {
		return nil, err
	}
	if *job.Status == "running" {
		result.Health = sdk.StatusReport_READY
		result.HealthMessage = fmt.Sprintf("Job %q is reporting ready!", deployment.Name)
	} else if *job.Status == "queued" || *job.Status == "started" {
		result.Health = sdk.StatusReport_ALIVE
		result.HealthMessage = fmt.Sprintf("Job %q is reporting alive!", deployment.Name)
	} else if *job.Status == "completed" {
		result.Health = sdk.StatusReport_PARTIAL
		result.HealthMessage = fmt.Sprintf("Job %q is reporting partially available!", deployment.Name)
	} else if *job.Status == "failed" || *job.Status == "lost" {
		result.Health = sdk.StatusReport_DOWN
		result.HealthMessage = fmt.Sprintf("Job %q is reporting down!", deployment.Name)
	} else {
		result.Health = sdk.StatusReport_UNKNOWN
		result.HealthMessage = fmt.Sprintf("Job %q is reporting unknown!", deployment.Name)
	}

	if *job.StatusDescription != "" {
		result.HealthMessage = *job.StatusDescription
	}

	result.GeneratedTime = timestamppb.Now()

	s.Update("Finished building report for Nomad platform")
	s.Done()

	// NOTE(briancain): Replace ui.Status with StepGroups once this bug
	// has been fixed: https://github.com/hashicorp/waypoint/issues/1536
	st := ui.Status()
	defer st.Close()

	st.Update("Determining overall container health...")
	if result.Health == sdk.StatusReport_READY {
		st.Step(terminal.StatusOK, result.HealthMessage)
	} else {
		if result.Health == sdk.StatusReport_PARTIAL {
			st.Step(terminal.StatusWarn, result.HealthMessage)
		} else {
			st.Step(terminal.StatusError, result.HealthMessage)
		}

		// Extra advisory wording to let user know that the deployment could be still starting up
		// if the report was generated immediately after it was deployed or released.
		st.Step(terminal.StatusWarn, mixedHealthWarn)
	}

	return &result, nil
}

// Config is the configuration structure for the Platform.
type Config struct {
	// The credential of docker registry.
	Auth *AuthConfig `hcl:"auth,block"`

	// The number of replicas of the service to maintain. If this number is maintained
	// outside waypoint, do not set this variable.
	Count int `hcl:"replicas,optional"`

	// The datacenters to deploy to, defaults to ["dc1"]
	Datacenter string `hcl:"datacenter,optional"`

	// The namespace of the job
	Namespace string `hcl:"namespace,optional"`

	// The Nomad region to deploy to, defaults to "global"
	Region string `hcl:"region,optional"`

	// The amount of resources to allocate to the Nomad task for the deployed
	// application
	Resources *Resources `hcl:"resources,block"`

	// Port that your service is running on within the actual container.
	// Defaults to port 3000.
	// TODO Evaluate if this should remain as a default 3000, should be a required field,
	// or default to another port.
	ServicePort uint `hcl:"service_port,optional"`

	// Specifies the service registration provider to use for service registrations
	ServiceProvider string `hcl:"service_provider,optional"`

	// Environment variables that are meant to configure the application in a static
	// way. This might be control an image that has multiple modes of operation,
	// selected via environment variable. Most configuration should use the waypoint
	// config commands.
	StaticEnvVars map[string]string `hcl:"static_environment,optional"`
}

type Resources struct {
	CPU      *int `hcl:"cpu,optional"`
	MemoryMB *int `hcl:"memorymb,optional"`
}

// AuthConfig maps the the Nomad Docker driver 'auth' config block
// and is used to set credentials for pulling images from the registry
type AuthConfig struct {
	Username string `hcl:"username"`
	Password string `hcl:"password"`
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&Config{}), docs.FromFunc(p.DeployFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description("Deploy to a nomad cluster as a service using docker")
	doc.Input("docker.Image")
	doc.Output("nomad.Deployment")

	doc.Example(
		`
deploy {
        use "nomad" {
          region = "global"
          datacenter = "dc1"
          auth {
            username = "username"
            password = "password"
          }
          static_environment = {
            "environment": "production",
            "LOG_LEVEL": "debug"
          }
          service_port = 3000
          replicas = 1
        }
}
`)

	doc.SetField(
		"region",
		"The Nomad region to deploy the job to.",
		docs.Default("global"),
	)

	doc.SetField(
		"datacenter",
		"The Nomad datacenter to deploy the job to.",
		docs.Default("dc1"),
	)

	doc.SetField(
		"namespace",
		"The Nomad namespace to deploy the job to.",
	)

	doc.SetField(
		"replicas",
		"The replica count for the job.",
		docs.Default("1"),
	)

	doc.SetField(
		"resources",
		"The amount of resources to allocate to the deployed allocation.",
		docs.SubFields(func(doc *docs.SubFieldDoc) {
			doc.SetField(
				"cpu",
				"Amount of CPU in MHz to allocate to this task",
				docs.Default(strconv.Itoa(defaultResourcesCPU)),
			)

			doc.SetField(
				"memorymb",
				"Amount of memory in MB to allocate to this task.",
				docs.Default(strconv.Itoa(defaultResourcesMemoryMB)),
			)
		}),
	)

	doc.SetField(
		"auth",
		"The credentials for docker registry.",
	)

	doc.SetField(
		"consul_token",
		"The Consul ACL token used to register services with the Nomad job.",
		docs.Summary("Uses the runner config environment variable CONSUL_HTTP_TOKEN."),
		docs.EnvVar("CONSUL_HTTP_TOKEN"),
	)

	doc.SetField(
		"vault_token",
		"The Vault token used to deploy the Nomad job with a token having specific Vault policies attached.",
		docs.Summary("Uses the runner config environment variable VAULT_TOKEN."),
		docs.EnvVar("VAULT_TOKEN"),
	)

	doc.SetField(
		"static_environment",
		"Environment variables to add to the job.",
	)

	doc.SetField(
		"service_provider",
		"Specifies the service registration provider to use for registering a service for the job",
		docs.Default("consul"),
	)

	doc.SetField(
		"service_port",
		"TCP port the job is listening on.",
	)

	return doc, nil
}

var (
	mixedHealthWarn = strings.TrimSpace(`
Waypoint detected that the current deployment is not ready, however your application
might be available or still starting up.
`)
)

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
	_ component.Destroyer    = (*Platform)(nil)
)
