package docker

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	wpdockerclient "github.com/hashicorp/waypoint/builtin/docker/client"
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

	cli, err := p.getDockerClient()
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "unable to create Docker client: %s", err)
	}

	cli.NegotiateAPIVersion(ctx)

	// Pull the image
	err = p.pullImage(cli, log, ui, img, p.config.ForcePull)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "unable to pull image from Docker registry: %s", err)
	}

	// Create our deployment and set an initial ID
	var result Deployment
	id, err := component.Id()
	if err != nil {
		return nil, err
	}
	result.Id = id
	result.Name = src.App

	s := sg.Add("Setting up waypoint network")
	defer func() { s.Abort() }()

	nets, err := cli.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "use=waypoint")),
	})

	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "unable to list Docker networks: %s", err)
	}

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
			return nil, status.Errorf(codes.FailedPrecondition, "unable to create Docker network: %s", err)
		}
	}

	s.Done()

	s = sg.Add("Creating new container")

	port := fmt.Sprint(p.config.ServicePort)
	np, err := nat.NewPort("tcp", port)
	if err != nil {
		return nil, err
	}

	cfg := container.Config{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		OpenStdin:    true,
		StdinOnce:    true,
		Image:        img.Image + ":" + img.Tag,
		ExposedPorts: nat.PortSet{np: struct{}{}},
		Env:          []string{"PORT=" + port},
	}

	if c := p.config.Command; len(c) > 0 {
		cfg.Cmd = c
	}

	bindings := nat.PortMap{}
	bindings[np] = []nat.PortBinding{
		{
			HostPort: "",
		},
	}

	hostconfig := container.HostConfig{
		Binds:        []string{src.App + "-scratch" + ":/input"},
		PortBindings: bindings,
	}

	netconfig := network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"waypoint": {},
		},
	}

	for k, v := range p.config.StaticEnvVars {
		cfg.Env = append(cfg.Env, k+"="+v)
	}

	for k, v := range deployConfig.Env() {
		cfg.Env = append(cfg.Env, k+"="+v)
	}

	cfg.Labels = map[string]string{
		labelId:     result.Id,
		"app":       src.App,
		"workspace": job.Workspace,
	}

	name := src.App + "-" + id

	cr, err := cli.ContainerCreate(ctx, &cfg, &hostconfig, &netconfig, name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to create Docker container: %s", err)
	}

	s.Update("Starting container")
	err = cli.ContainerStart(ctx, cr.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to start Docker container: %s", err)
	}
	s.Done()

	s = sg.Add("App deployed as container: " + name)
	s.Done()

	result.Container = cr.ID

	return &result, nil
}

// Destroy deletes the K8S deployment.
func (p *Platform) Destroy(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) error {
	cli, err := p.getDockerClient()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "unable to create Docker client: %s", err)
	}

	cli.NegotiateAPIVersion(ctx)

	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()
	st.Update("Deleting container...")

	// Check if the container exists
	_, err = cli.ContainerInspect(ctx, deployment.Container)
	if client.IsErrNotFound(err) {
		return nil
	}

	// Remove it
	return cli.ContainerRemove(ctx, deployment.Container, types.ContainerRemoveOptions{
		Force: true,
	})
}

func (p *Platform) getDockerClient() (*client.Client, error) {
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

	return wpdockerclient.NewClientWithOpts(opts...)
}

func (p *Platform) pullImage(cli *client.Client, log hclog.Logger, ui terminal.UI, img *Image, force bool) error {
	in := fmt.Sprintf("%s:%s", img.Image, img.Tag)
	args := filters.NewArgs()
	args.Add("reference", in)

	sg := ui.StepGroup()

	// only pull if image is not in current registry so check to see if the image is present
	// if force then skip this check
	if force == false {
		sg.Add("Checking Docker image cache for Image " + in)

		sum, err := cli.ImageList(context.Background(), types.ImageListOptions{Filters: args})
		if err != nil {
			return fmt.Errorf("unable to list images in local Docker cache: %w", err)
		}

		// if we have images do not pull
		if len(sum) > 0 {
			return nil
		}
	}

	step := sg.Add("Pulling Docker Image " + in)
	defer step.Abort()

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

	err = jsonmessage.DisplayJSONMessagesStream(out, step.TermOutput(), termFd, true, nil)
	if err != nil {
		return status.Errorf(codes.Internal, "unable to stream build logs to the terminal: %s", err)
	}

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
	// The command to run in the container. This is an array of arguments
	// that are executed directly. These are not executed in the context of
	// a shell. If you want to use a shell, add that to this command manually.
	Command []string `hcl:"command,optional"`

	// A path to a directory that will be created for the service to store
	// temporary data.
	ScratchSpace string `hcl:"scratch_path,optional"`

	// Environment variables that are meant to configure the application in a static
	// way. This might be control an image that has mulitple modes of operation,
	// selected via environment variable. Most configuration should use the waypoint
	// config commands.
	StaticEnvVars map[string]string `hcl:"static_environment,optional"`

	// Port that your service is running on within the actual container.
	// Defaults to port 3000.
	// TODO Evaluate if this should remain as a default 3000, should be a required field,
	// or default to another port.
	ServicePort uint `hcl:"service_port,optional"`

	// Force pull the image from the remote repository
	ForcePull bool `hcl:"force_pull,optional"`

	// ClientConfig allow the user to specify the connection to the Docker
	// engine. By default we try to load this from env vars:
	// DOCKER_HOST to set the url to the docker server.
	// DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	// DOCKER_CERT_PATH to load the TLS certificates from.
	// DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default.
	ClientConfig *ClientConfig `hcl:"client_config,block"`
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
		"command",
		"the command to run to start the application in the container",
		docs.Default("the image entrypoint"),
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
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
	_ component.Destroyer    = (*Platform)(nil)
)
