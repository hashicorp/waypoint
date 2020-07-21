package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/go-hclog"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

const (
	labelId    = "waypoint.hashicorp.com/id"
	labelNonce = "waypoint.hashicorp.com/nonce"
)

// Platform is the Platform implementation for Kubernetes.
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

// Deploy deploys an image to Kubernetes.
func (p *Platform) Deploy(
	ctx context.Context,
	log hclog.Logger,
	src *component.Source,
	img *Image,
	deployConfig *component.DeploymentConfig,
	ui terminal.UI,
) (*Deployment, error) {
	// We'll update the user in real time
	sg := ui.StepGroup()

	s1 := sg.Add("Checking for existing containers")
	defer s1.Abort()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	cli.NegotiateAPIVersion(ctx)

	// Create our deployment and set an initial ID
	var result Deployment
	id, err := component.Id()
	if err != nil {
		return nil, err
	}
	result.Id = id
	result.Name = src.App

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "app=" + src.App,
		}),
	})

	if err != nil {
		return nil, err
	}

	if len(containers) > 0 {
		s1.Update("Found an existing containers")
		s1.Done()

		s2 := sg.Add("Deleting existing container: " + containers[0].ID)
		defer s2.Abort()

		err = cli.ContainerRemove(ctx, containers[0].ID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			return nil, err
		}

		s2.Done()
	} else {
		s1.Update("No existing containers detected")
	}

	s3 := sg.Add("Creating new container")
	defer s3.Abort()

	port := "3000"

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
		Image:        img.Image,
		ExposedPorts: nat.PortSet{np: struct{}{}},
		Env:          []string{"PORT=" + port},
	}

	if p.config.Command != "" {
		cfg.Cmd = append(cfg.Cmd, "/bin/sh", "-c", p.config.Command)
	}

	bindings := nat.PortMap{}
	bindings[np] = []nat.PortBinding{
		{
			HostPort: port,
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
		labelId: result.Id,
		"app":   src.App,
	}

	name := src.App + "-" + id

	cr, err := cli.ContainerCreate(ctx, &cfg, &hostconfig, &netconfig, name)
	if err != nil {
		return nil, err
	}

	s3.Update("Starting container")
	err = cli.ContainerStart(ctx, cr.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	s3.Done()

	s4 := sg.Add("App deployed as container container: " + name)
	s4.Done()

	result.Container = cr.ID

	sg.Wait()

	return &result, nil
}

// Destroy deletes the K8S deployment.
func (p *Platform) Destroy(
	ctx context.Context,
	log hclog.Logger,
	deployment *Deployment,
	ui terminal.UI,
) error {
	// We'll update the user in real time
	st := ui.Status()
	defer st.Close()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	cli.NegotiateAPIVersion(ctx)

	st.Update("Deleting container...")

	return cli.ContainerRemove(ctx, deployment.Container, types.ContainerRemoveOptions{
		Force: true,
	})
}

// Config is the configuration structure for the Platform.
type PlatformConfig struct {
	// The command to run in the container
	Command string `hcl:"command,optional"`

	// A path to a directory that will be created for the service to store
	// temporary data.
	ScratchSpace string `hcl:"scratch_path,optional"`

	// Environment variables that are meant to configure the application in a static
	// way. This might be control an image that has mulitple modes of operation,
	// selected via environment variable. Most configuration should use the waypoint
	// config commands.
	StaticEnvVars map[string]string `hcl:"static_environment,optional"`
}

var (
	_ component.Platform     = (*Platform)(nil)
	_ component.Configurable = (*Platform)(nil)
	_ component.Destroyer    = (*Platform)(nil)
)
