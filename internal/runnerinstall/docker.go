package runnerinstall

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
)

type DockerConfig struct {
	RunnerImage string `hcl:"runner_image,optional"`
	Network     string `hcl:"network,optional"`
}

type DockerRunnerInstaller struct {
	Config DockerConfig
}

func (i *DockerRunnerInstaller) Install(ctx context.Context, opts *InstallOpts) error {
	sg := opts.UI.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing Docker client...")
	defer func() { s.Abort() }()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	cli.NegotiateAPIVersion(ctx)

	var runnerImage string
	if i.Config.RunnerImage == "" {
		runnerImage = "hashicorp/waypoint:latest"
	} else {
		runnerImage = i.Config.RunnerImage
	}

	var waypointNetwork network.NetworkingConfig
	if i.Config.Network != "" {
		waypointNetwork = network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				i.Config.Network: {},
			},
		}
	}

	// The key thing in the container creation below is that the environment
	// variables are set to the advertised address env vars which will
	// allow our runner to connect.
	cr, err := cli.ContainerCreate(ctx, &container.Config{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		OpenStdin:    true,
		StdinOnce:    true,
		User:         "root",
		Image:        runnerImage,
		Env:          opts.AdvertiseClient.Env(),
		Cmd:          []string{"runner", "agent", "-id=" + opts.Id, "-cookie=" + opts.Cookie, "-vv"},
		Labels: map[string]string{
			"waypoint-type": "runner",
		},
	}, &container.HostConfig{
		Privileged: true,
		CapAdd:     []string{"CAP_DAC_OVERRIDE"},
		Binds:      []string{"/var/run/docker.sock:/var/run/docker.sock"},
		// These security options are required for the runner so that
		// Docker daemonless image building works properly.
		SecurityOpt: []string{
			"seccomp=unconfined",
			"apparmor=unconfined",
		},
	}, &waypointNetwork, nil, "waypoint-runner-"+opts.Id)
	if err != nil {
		return err
	}

	err = cli.ContainerStart(ctx, cr.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	s.Update("Waypoint runner installed and started!")
	s.Done()

	return nil
}

func (i *DockerRunnerInstaller) InstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:    "docker-runner-image",
		Target:  &i.Config.RunnerImage,
		Usage:   "The Docker image for the Waypoint runner.",
		Default: "hashicorp/waypoint:latest",
	})

	set.StringVar(&flag.StringVar{
		Name:   "docker-runner-network",
		Target: &i.Config.Network,
		Usage:  "The Docker network in which to deploy the Waypoint runner.",
	})
}

func (d DockerRunnerInstaller) Uninstall(ctx context.Context, opts *InstallOpts) error {
	//TODO implement me
	panic("implement me")
}

func (d DockerRunnerInstaller) UninstallFlags(set *flag.Set) {
	//TODO implement me
	panic("implement me")
}
