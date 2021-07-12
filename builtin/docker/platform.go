package docker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	goUnits "github.com/docker/go-units"
	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	wpdockerclient "github.com/hashicorp/waypoint/builtin/docker/client"

	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
)

const (
	labelId    = "waypoint.hashicorp.com/id"
	labelNonce = "waypoint.hashicorp.com/nonce"
)

// Platform is the Platform implementation for Docker.
type Platform struct {
	config PlatformConfig
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
		resource.WithValueProvider(p.getDockerClient),
		resource.WithResource(resource.NewResource(
			resource.WithName("network"),
			resource.WithState(&Resource_Network{}),
			resource.WithCreate(p.resourceNetworkCreate),

			// networks have no destroy logic, we leave the network
			// lingering around for now. This was the logic prior to
			// refactoring into the resource manager so we kept it.
		)),

		resource.WithResource(resource.NewResource(
			resource.WithName("container"),
			resource.WithState(&Resource_Container{}),
			resource.WithCreate(p.resourceContainerCreate),
			resource.WithDestroy(p.resourceContainerDestroy),
		)),
	)
	return nil
}

func (p *Platform) Status(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) (*sdk.StatusReport, error) {
	cli, err := p.getDockerClient(ctx)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "unable to create Docker client: %s", err)
	}
	cli.NegotiateAPIVersion(ctx)

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Gathering health report for Docker platform...")
	defer s.Abort()

	// currently the docker platform only deploys 1 container
	containerInfo, err := cli.ContainerInspect(ctx, deployment.Container)
	if err != nil {
		return nil, err
	}

	// Create our status report
	var result sdk.StatusReport
	result.External = true

	// NOTE(briancain): The docker platform currently only deploys a single
	// container, so for now the status report makes the same assumption.

	log.Debug("querying docker for container health")

	resources := []*sdk.StatusReport_Resource{{
		Name: containerInfo.Name,
	}}

	if containerInfo.State.Health != nil {
		// Built-in Docker health reporting
		// NOTE: this only works if the container has configured health checks

		switch containerInfo.State.Health.Status {
		case "Healthy":
			resources[0].Health = sdk.StatusReport_READY
			resources[0].HealthMessage = "container is running"
		case "Unhealthy":
			resources[0].Health = sdk.StatusReport_DOWN
			resources[0].HealthMessage = "container is down"
		case "Starting":
			resources[0].Health = sdk.StatusReport_ALIVE
			resources[0].HealthMessage = "container is starting"
		default:
			resources[0].Health = sdk.StatusReport_UNKNOWN
			resources[0].HealthMessage = "unknown status reported by docker for container"
		}
	} else {
		// Waypoint container inspection

		if containerInfo.State.Running && containerInfo.State.ExitCode == 0 {
			resources[0].Health = sdk.StatusReport_READY
			resources[0].HealthMessage = "container is running"
		} else if containerInfo.State.Restarting || containerInfo.State.Status == "created" {
			resources[0].Health = sdk.StatusReport_ALIVE
			resources[0].HealthMessage = "container is still starting"
		} else if containerInfo.State.Dead || containerInfo.State.OOMKilled || containerInfo.State.ExitCode != 0 {
			resources[0].Health = sdk.StatusReport_DOWN
			resources[0].HealthMessage = "container is down"
		} else {
			resources[0].Health = sdk.StatusReport_UNKNOWN
			resources[0].HealthMessage = "unknown status for container"
		}
	}

	result.Resources = resources

	// Determine overall deployment health based on its resource health
	var ready, alive, down, unknown int
	for _, r := range result.Resources {
		switch r.Health {
		case sdk.StatusReport_DOWN:
			down++
		case sdk.StatusReport_UNKNOWN:
			unknown++
		case sdk.StatusReport_READY:
			ready++
		case sdk.StatusReport_ALIVE:
			alive++
		}
	}

	if ready == len(result.Resources) {
		result.Health = sdk.StatusReport_READY
		result.HealthMessage = fmt.Sprintf("Container %q is reporting ready!", containerInfo.Name)
	} else if down == len(result.Resources) {
		result.Health = sdk.StatusReport_DOWN
		result.HealthMessage = fmt.Sprintf("Container %q is reporting down!", containerInfo.Name)
	} else if unknown == len(result.Resources) {
		result.Health = sdk.StatusReport_UNKNOWN
		result.HealthMessage = fmt.Sprintf("Container %q is reporting unknown!", containerInfo.Name)
	} else if alive == len(result.Resources) {
		result.Health = sdk.StatusReport_ALIVE
		result.HealthMessage = fmt.Sprintf("Container %q is reporting alive!", containerInfo.Name)
	} else {
		result.Health = sdk.StatusReport_PARTIAL
		result.HealthMessage = fmt.Sprintf("Container %q is reporting partially available!", containerInfo.Name)
	}

	result.GeneratedTime = ptypes.TimestampNow()
	log.Debug("status report complete")

	// update output based on main health state
	s.Update("Finished building report for Docker platform")
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

func (p *Platform) resourceNetworkCreate(
	ctx context.Context,
	cli *client.Client,
	sg terminal.StepGroup,
	state *Resource_Network,
) error {
	s := sg.Add("Setting up network...")
	defer func() { s.Abort() }()

	nets, err := cli.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "use=waypoint")),
	})
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "unable to list Docker networks: %s", err)
	}

	// If we have a network already we're done. If we don't have a net, create it.
	if len(nets) == 0 {
		_, err = cli.NetworkCreate(ctx, "waypoint", types.NetworkCreate{
			Driver:         "bridge",
			CheckDuplicate: true,
			Internal:       false,
			Attachable:     true,
			Labels: map[string]string{
				"use": "waypoint",
			},
		})
		if err != nil {
			return status.Errorf(codes.FailedPrecondition, "unable to create Docker network: %s", err)
		}
	}
	s.Done()

	// Set our state
	state.Name = "waypoint"

	return nil
}

func (p *Platform) resourceContainerCreate(
	ctx context.Context,
	log hclog.Logger,
	cli *client.Client,
	src *component.Source,
	img *Image,
	job *component.JobInfo,
	deployConfig *component.DeploymentConfig,
	result *Deployment,
	sg terminal.StepGroup,
	ui terminal.UI,
	state *Resource_Container,
	netState *Resource_Network,
) error {
	// Pull the image
	err := p.pullImage(cli, log, ui, img, p.config.ForcePull)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition,
			"unable to pull image from Docker registry: %s", err)
	}

	s := sg.Add("Creating new container...")
	defer func() { s.Abort() }()

	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}
	for _, port := range append(p.config.ExtraPorts, p.config.ServicePort) {
		np, err := nat.NewPort("tcp", fmt.Sprint(port))
		if err != nil {
			return err
		}

		exposedPorts[np] = struct{}{}
		portBindings[np] = []nat.PortBinding{
			{
				HostPort: "", // this is intentionally left empty for a random host port assignment
			},
		}
	}

	cfg := container.Config{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		OpenStdin:    true,
		StdinOnce:    true,
		Image:        img.Image + ":" + img.Tag,
		ExposedPorts: exposedPorts,
		Env:          []string{"PORT=" + fmt.Sprint(p.config.ServicePort)},
	}
	if c := p.config.Command; len(c) > 0 {
		cfg.Cmd = c
	}

	// default container binds
	containerBinds := []string{src.App + "-scratch" + ":/input"}
	if p.config.Binds != nil {
		containerBinds = append(containerBinds, p.config.Binds...)
	}

	// Setup the resource requirements for the container if given
	var resources container.Resources
	if p.config.Resources != nil {
		memory, err := goUnits.FromHumanSize(p.config.Resources["memory"])
		if err != nil {
			return err
		}
		resources.Memory = memory

		cpu, err := strconv.ParseInt(p.config.Resources["cpu"], 10, 64)
		if err != nil {
			return err
		}
		resources.CPUShares = cpu
	}

	// Build our host configuration from the bindings, ports, and resources.
	hostconfig := container.HostConfig{
		Binds:        containerBinds,
		PortBindings: portBindings,
		Resources:    resources,
	}

	// Containers can only be connected to 1 network at creation time
	// Additional user defined networks will be connected after container is
	// created.
	netconfig := network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			netState.Name: {},
		},
	}

	for k, v := range p.config.StaticEnvVars {
		cfg.Env = append(cfg.Env, k+"="+v)
	}
	for k, v := range deployConfig.Env() {
		cfg.Env = append(cfg.Env, k+"="+v)
	}

	// Setup the labels. We setup a set of defaults and then override them
	// with any user configured labels.
	defaultLabels := map[string]string{
		labelId:     result.Id,
		"app":       src.App,
		"workspace": job.Workspace,
	}
	if p.config.Labels != nil {
		for k, v := range defaultLabels {
			p.config.Labels[k] = v
		}
	} else {
		p.config.Labels = defaultLabels
	}
	cfg.Labels = p.config.Labels

	// Create the container
	name := src.App + "-" + result.Id
	cr, err := cli.ContainerCreate(ctx, &cfg, &hostconfig, &netconfig, nil, name)
	if err != nil {
		return status.Errorf(codes.Internal, "unable to create Docker container: %s", err)
	}

	// Store our state so we can destroy it properly
	state.Id = cr.ID
	state.Name = name

	// Additional networks must be connected after container is created
	if p.config.Networks != nil {
		s.Update("Connecting additional networks to container...")
		for _, net := range p.config.Networks {
			err = cli.NetworkConnect(ctx, net, cr.ID, &network.EndpointSettings{})
			if err != nil {
				s.Update("Failed to connect additional network")
				s.Status(terminal.StatusError)
				s.Done()
				return status.Errorf(
					codes.Internal,
					"unable to connect container to additional networks: %s",
					err)
			}
		}
	}

	s.Update("Starting container")
	err = cli.ContainerStart(ctx, cr.ID, types.ContainerStartOptions{})
	if err != nil {
		return status.Errorf(codes.Internal, "unable to start Docker container: %s", err)
	}
	s.Done()

	return nil
}

func (p *Platform) resourceContainerDestroy(
	ctx context.Context,
	cli *client.Client,
	state *Resource_Container,
	sg terminal.StepGroup,
) error {
	// Check if the container exists
	_, err := cli.ContainerInspect(ctx, state.Id)
	if client.IsErrNotFound(err) {
		return nil
	}

	s := sg.Add("Deleting container: %s", state.Id)
	defer func() { s.Abort() }()

	// Remove it
	err = cli.ContainerRemove(ctx, state.Id, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		return err
	}

	s.Done()
	return nil
}

// Deploy deploys an image to Docker.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	job *component.JobInfo,
	img *Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {
	// We'll update the user in real time
	sg := ui.StepGroup()
	defer sg.Wait()

	if p.config.ServicePort == 0 {
		p.config.ServicePort = 3000
	}

	// Create our deployment and set an initial ID. This just creates
	// the initial structure this doesn't persist any state yet.
	var result Deployment
	id, err := component.Id()
	if err != nil {
		return nil, err
	}
	result.Id = id
	result.Name = src.App

	// Create our resource manager and create
	rm := p.resourceManager(log)
	if err := rm.CreateAll(
		ctx, log, sg, ui,
		src, job, img, deployConfig, &result,
	); err != nil {
		return nil, err
	}

	// Store our resource state
	result.ResourceState = rm.State()

	// Get our container state
	crState := rm.Resource("container").State().(*Resource_Container)
	if crState == nil {
		return nil, status.Errorf(codes.Internal,
			"container state is nil, this should never happen")
	}

	s := sg.Add("App deployed as container: " + crState.Name)
	s.Done()

	result.Container = crState.Id
	return &result, nil
}

// Destroy deletes a Docker deployment.
func (p *Platform) Destroy(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	rm := p.resourceManager(log)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if deployment.ResourceState == nil {
		rm.Resource("container").SetState(&Resource_Container{
			Id: deployment.Container,
		})
	} else {
		// Load our set state
		if err := rm.LoadState(deployment.ResourceState); err != nil {
			return err
		}
	}

	// Destroy
	return rm.DestroyAll(ctx, log, sg, ui)
}

func (p *Platform) getDockerClient(ctx context.Context) (*client.Client, error) {
	if p.config.ClientConfig == nil {
		return wpdockerclient.NewClientWithOpts(client.FromEnv)
	}

	opts := []client.Opt{}

	if host := p.config.ClientConfig.Host; host != "" {
		opts = append(opts, client.WithHost(host))
	}

	if path := p.config.ClientConfig.CertPath; path != "" {
		opts = append(opts, client.WithTLSClientConfig(
			filepath.Join(path, "ca.pem"),
			filepath.Join(path, "cert.pem"),
			filepath.Join(path, "key.pem"),
		))
	}

	if version := p.config.ClientConfig.APIVersion; version != "" {
		opts = append(opts, client.WithVersion(version))
	}

	cli, err := wpdockerclient.NewClientWithOpts(opts...)
	if err != nil {
		return nil, err
	}

	cli.NegotiateAPIVersion(ctx)
	return cli, nil
}

func (p *Platform) pullImage(cli *client.Client, log hclog.Logger, ui terminal.UI, img *Image, force bool) error {
	in := fmt.Sprintf("%s:%s", img.Image, img.Tag)
	args := filters.NewArgs()
	args.Add("reference", in)

	sg := ui.StepGroup()
	s := sg.Add("")
	defer func() { s.Abort() }()

	// only pull if image is not in current registry so check to see if the image is present
	// if force then skip this check
	if force == false {
		s.Update("Checking Docker image cache for Image " + in)

		sum, err := cli.ImageList(context.Background(), types.ImageListOptions{Filters: args})
		if err != nil {
			return fmt.Errorf("unable to list images in local Docker cache: %w", err)
		}

		// if we have images do not pull
		if len(sum) > 0 {
			s.Update("Docker image %q up to date!", in)
			s.Done()
			return nil
		}
	}

	s.Update("Pulling Docker Image " + in)

	ipo := types.ImagePullOptions{}

	// if the username and password is not null make an authenticated
	// image pull
	/*
		if image.Username != "" && image.Password != "" {
			ipo.RegistryAuth = createRegistryAuth(image.Username, image.Password)
		}
	*/

	in = makeImageCanonical(in)
	log.Debug("pulling image", "image", in)

	out, err := cli.ImagePull(context.Background(), in, ipo)
	if err != nil {
		return fmt.Errorf("unable to pull image: %w", err)
	}

	stdout, _, err := ui.OutputWriters()
	if err != nil {
		return fmt.Errorf("unable to get output writers: %s", err)
	}

	var termFd uintptr
	if f, ok := stdout.(*os.File); ok {
		termFd = f.Fd()
	}

	err = jsonmessage.DisplayJSONMessagesStream(out, s.TermOutput(), termFd, true, nil)
	if err != nil {
		return status.Errorf(codes.Internal, "unable to stream build logs to the terminal: %s", err)
	}

	s.Done()

	return nil
}

// makeImageCanonical makes sure the image reference uses full canonical name i.e.
// consul:1.6.1 -> docker.io/library/consul:1.6.1
func makeImageCanonical(image string) string {
	imageParts := strings.Split(image, "/")
	switch len(imageParts) {
	case 1:
		return fmt.Sprintf("docker.io/library/%s", imageParts[0])
	case 2:
		return fmt.Sprintf("docker.io/%s/%s", imageParts[0], imageParts[1])
	}

	return image
}

// Config is the configuration structure for the Platform.
type PlatformConfig struct {
	// A list of folders to mount to the container.
	Binds []string `hcl:"binds,optional"`

	// ClientConfig allow the user to specify the connection to the Docker
	// engine. By default we try to load this from env vars:
	// DOCKER_HOST to set the url to the docker server.
	// DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	// DOCKER_CERT_PATH to load the TLS certificates from.
	// DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default.
	ClientConfig *ClientConfig `hcl:"client_config,block"`

	// The command to run in the container. This is an array of arguments
	// that are executed directly. These are not executed in the context of
	// a shell. If you want to use a shell, add that to this command manually.
	Command []string `hcl:"command,optional"`

	// Force pull the image from the remote repository
	ForcePull bool `hcl:"force_pull,optional"`

	// A map of key/value pairs, stored in docker as a string. Each key/value pair must
	// be unique. Validiation occurs at the docker layer, not in Waypoint. Label
	// keys are alphanumeric strings which may contain periods (.) and hyphens (-).
	// See the docker docs for more info: https://docs.docker.com/config/labels-custom-metadata/
	Labels map[string]string `hcl:"labels,optional"`

	// An array of strings with network names to connect the container to
	Networks []string `hcl:"networks,optional"`

	// A map of resources to configure the container with such as memory and cpu
	// limits.
	Resources map[string]string `hcl:"resources,optional"`

	// A path to a directory that will be created for the service to store
	// temporary data.
	ScratchSpace string `hcl:"scratch_path,optional"`

	// Environment variables that are meant to configure the application in a static
	// way. This might be control an image that has mulitple modes of operation,
	// selected via environment variable. Most configuration should use the waypoint
	// config commands.
	StaticEnvVars map[string]string `hcl:"static_environment,optional"`

	// Additional ports the application is listening on to expose on the container
	ExtraPorts []uint `hcl:"extra_ports,optional"`

	// Port that your service is running on within the actual container.
	// Defaults to port 3000.
	// TODO Evaluate if this should remain as a default 3000, should be a required field,
	// or default to another port.
	ServicePort uint `hcl:"service_port,optional"`
}

type ClientConfig struct {
	// Host to use when connecting to Docker
	// This can be used to connect to remote Docker instances
	Host string `hcl:"host,optional"`

	// Path to load the certificates for the Docker Engine
	CertPath string `hcl:"cert_path,optional"`

	// Docker API version to use for connection
	APIVersion string `hcl:"api_version,optional"`
}

func (p *Platform) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(docs.FromConfig(&PlatformConfig{}), docs.FromFunc(p.DeployFunc()))
	if err != nil {
		return nil, err
	}

	doc.Description("Deploy a container to Docker, local or remote")

	doc.Example(`
deploy {
  use "docker" {
	command      = ["ps"]
	service_port = 3000
	static_environment = {
	  "environment": "production",
	  "LOG_LEVEL": "debug"
	}
  }
}
`)

	doc.Input("docker.Image")
	doc.Output("docker.Deployment")

	doc.SetField(
		"binds",
		"A 'source:destination' list of folders to mount onto the container from the host.",
		docs.Summary(
			"A list of folders to mount onto the container from the host. The expected",
			"format for each string entry in the list is `source:destination`. So",
			"for example: `binds: [\"host_folder/scripts:/scripts\"]",
		),
	)

	doc.SetField(
		"command",
		"the command to run to start the application in the container",
		docs.Default("the image entrypoint"),
	)

	doc.SetField(
		"labels",
		"A map of key/value pairs to label the docker container with.",
		docs.Summary(
			"A map of key/value pair(s), stored in docker as a string. Each key/value pair must",
			"be unique. Validiation occurs at the docker layer, not in Waypoint. Label",
			"keys are alphanumeric strings which may contain periods (.) and hyphens (-).",
		),
	)

	doc.SetField(
		"networks",
		"An list of strings with network names to connect the container to.",
		docs.Default("waypoint"),
		docs.Summary(
			"A list of networks to connect the container to. By default the container",
			"will always connect to the `waypoint` network.",
		),
	)

	doc.SetField(
		"resources",
		"A map of resources to configure the container with, such as memory or cpu limits.",
		docs.Summary(
			"these options are used to configure the container used when deploying",
			"with docker. Currently, the supported resources are 'memory' and 'cpu' limits.",
			"The field 'memory' is expected to be defined as \"512MB\", \"44kB\", etc.",
		),
	)

	doc.SetField(
		"scratch_path",
		"a path within the container to store temporary data",
		docs.Summary(
			"docker will mount a tmpfs at this path",
		),
	)

	doc.SetField(
		"static_environment",
		"environment variables to expose to the application",
		docs.Summary(
			"these environment variables should not be run of the mill",
			"configuration variables, use waypoint config for that.",
			"These variables are used to control over all container modes,",
			"such as configuring it to start a web app vs a background worker",
		),
	)

	doc.SetField(
		"extra_ports",
		"additional TCP ports the application is listening on to expose on the container",
		docs.Summary(
			"Used to define and expose multiple ports that the application is listening on for the container in use.",
			"These ports will get merged with service_port when creating the container if defined.",
		),
	)

	doc.SetField(
		"service_port",
		"port that your service is running on in the container",
		docs.Default("3000"),
	)

	doc.SetField(
		"force_pull",
		"always pull the docker container from the registry",
		docs.Default("false"),
	)

	doc.SetField(
		"client_config",
		"client config for remote Docker engine",
		docs.Summary(
			"this config block can be used to configure",
			"a remote Docker engine.",
			"By default Waypoint will attempt to discover this configuration",
			"using the environment variables:",
			"`DOCKER_HOST` to set the url to the docker server.",
			"`DOCKER_API_VERSION` to set the version of the API to reach, leave empty for latest.",
			"`DOCKER_CERT_PATH` to load the TLS certificates from.",
			"`DOCKER_TLS_VERIFY` to enable or disable TLS verification, off by default.",
		),
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
	_ component.Status       = (*Platform)(nil)
)
