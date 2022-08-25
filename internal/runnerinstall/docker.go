package runnerinstall

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
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
	ui := opts.UI
	sg := opts.UI.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing Docker client...")
	defer func() { s.Abort() }()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	cli.NegotiateAPIVersion(ctx)

	runnerImage := i.Config.RunnerImage
	imageRef, err := reference.ParseNormalizedNamed(runnerImage)
	if err != nil {
		return err
	}

	imageList, err := cli.ImageList(ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "reference",
			Value: reference.FamiliarString(imageRef),
		}),
	})
	if err != nil {
		return err
	}

	if len(imageList) == 0 {
		s.Update("Pulling image %s", runnerImage)

		resp, err := cli.ImagePull(ctx, reference.FamiliarName(imageRef), types.ImagePullOptions{})
		if err != nil {
			s.Update("Unable to pull waypoint image")
			return err
		}
		defer resp.Close()

		stdout, _, err := ui.OutputWriters()
		if err != nil {
			return err
		}

		var termFd uintptr
		if f, ok := stdout.(*os.File); ok {
			termFd = f.Fd()
		}

		err = jsonmessage.DisplayJSONMessagesStream(resp, s.TermOutput(), termFd, true, nil)
		if err != nil {
			return fmt.Errorf("unable to stream pull logs to the terminal: %s", err)
		}

	}

	// We have no opinion on the EndpointSettings for the Docker network which
	// the user can specify, so we assign an empty struct to the single element
	// in the map of networks (which is the name of the network, if provided)
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
		Cmd:          append([]string{"runner", "agent", "-id=" + opts.Id, "-cookie=" + opts.Cookie, "-vv"}, opts.RunnerRunFlags...),
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
	}, &waypointNetwork, nil, "waypoint-"+opts.Id+"-runner")
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
		Default: defaultRunnerImage,
	})

	set.StringVar(&flag.StringVar{
		Name:   "docker-runner-network",
		Target: &i.Config.Network,
		Usage:  "The Docker network in which to deploy the Waypoint runner.",
	})
}

func (d DockerRunnerInstaller) Uninstall(ctx context.Context, opts *InstallOpts) error {
	sg := opts.UI.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing Docker client...")
	defer func() { s.Abort() }()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	cli.NegotiateAPIVersion(ctx)

	s.Update("Finding runner container")
	containerNames := []string{
		defaultRunnerName(opts.Id),
		DefaultRunnerTagName,
	}
	var foundContainer types.Container
	for _, containerName := range containerNames {
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
			Filters: filters.NewArgs(filters.KeyValuePair{
				Key:   "name",
				Value: containerName,
			}),
		})
		if err != nil {
			s.Update("Could not get container list")
			return err
		}
		if len(containers) > 0 {
			foundContainer = containers[0]
			break
		}
	}
	if foundContainer.ID == "" {
		s.Update("Could not find runner in docker.")
		return fmt.Errorf("Runner not found.")
	}

	stopTimeout := time.Second * 30
	s.Update("Stopping runner...")
	err = cli.ContainerStop(ctx, foundContainer.ID, &stopTimeout)
	if err != nil {
		return err
	}

	s.Update("Removing runner container")
	err = cli.ContainerRemove(ctx, foundContainer.ID, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	s.Update("Waypoint Runner uninstalled")
	s.Done()
	return nil
}

func (d DockerRunnerInstaller) UninstallFlags(set *flag.Set) {
}
