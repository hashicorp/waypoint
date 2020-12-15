package serverinstall

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
	"github.com/docker/go-connections/nat"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverconfig"
)

type DockerInstaller struct {
	config dockerConfig
}

type dockerConfig struct {
	serverImage string `hcl:"server_image,optional"`
}

var (
	grpcPort       = "9701"
	httpPort       = "9702"
	containerLabel = "waypoint-type=server"
	containerKey   = "waypoint-type"
	containerValue = "server"
	containerName  = "waypoint-server"
)

// Install is a method of DockerInstaller and implements the Installer interface to
// create a waypoint-server as a Docker container
func (i *DockerInstaller) Install(
	ctx context.Context,
	opts *InstallOpts,
) (*InstallResults, error) {
	ui := opts.UI

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing Docker client...")
	defer func() { s.Abort() }()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(ctx)

	s.Update("Checking for existing installation...")

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: containerLabel,
		}),
	})
	if err != nil {
		return nil, err
	}

	var (
		clicfg   clicontext.Config
		addr     pb.ServerConfig_AdvertiseAddr
		httpAddr string
	)

	clicfg.Server = serverconfig.Client{
		Address:       "localhost:" + grpcPort,
		Tls:           true,
		TlsSkipVerify: true,
	}

	addr.Addr = containerName + ":" + grpcPort
	addr.Tls = true
	addr.TlsSkipVerify = true

	httpAddr = "localhost:" + httpPort

	// If we already have a server, bolt.
	if len(containers) > 0 {
		s.Update("Detected existing Waypoint server.")
		s.Status(terminal.StatusWarn)
		s.Done()
		return &InstallResults{
			Context:       &clicfg,
			AdvertiseAddr: &addr,
		}, nil
	}

	s.Update("Checking for Docker image: %s", i.config.serverImage)

	imageRef, err := reference.ParseNormalizedNamed(i.config.serverImage)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Docker image: %s", err)
	}

	imageList, err := cli.ImageList(ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "reference",
			Value: reference.FamiliarString(imageRef),
		}),
	})
	if err != nil {
		return nil, err
	}

	if len(imageList) == 0 {
		s.Update("Pulling image: %s", i.config.serverImage)

		resp, err := cli.ImagePull(ctx, reference.FamiliarString(imageRef), types.ImagePullOptions{})
		if err != nil {
			return nil, err
		}
		defer resp.Close()

		stdout, _, err := ui.OutputWriters()
		if err != nil {
			return nil, err
		}

		var termFd uintptr
		if f, ok := stdout.(*os.File); ok {
			termFd = f.Fd()
		}

		err = jsonmessage.DisplayJSONMessagesStream(resp, s.TermOutput(), termFd, true, nil)
		if err != nil {
			return nil, fmt.Errorf("unable to stream pull logs to the terminal: %s", err)
		}

		s.Done()
		s = sg.Add("")
	}

	s.Update("Creating waypoint network...")

	nets, err := cli.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "use=waypoint")),
	})
	if err != nil {
		return nil, err
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
			return nil, err
		}

	}

	npGRPC, err := nat.NewPort("tcp", grpcPort)
	if err != nil {
		return nil, err
	}

	npHTTP, err := nat.NewPort("tcp", httpPort)
	if err != nil {
		return nil, err
	}

	s.Update("Installing Waypoint server to docker")

	cfg := container.Config{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		OpenStdin:    true,
		StdinOnce:    true,
		Image:        i.config.serverImage,
		ExposedPorts: nat.PortSet{npGRPC: struct{}{}, npHTTP: struct{}{}},
		Env:          []string{"PORT=" + grpcPort},
		Cmd:          []string{"server", "run", "-accept-tos", "-vvv", "-db=/data/data.db", "-listen-grpc=0.0.0.0:9701", "-listen-http=0.0.0.0:9702"},
	}

	bindings := nat.PortMap{}
	bindings[npGRPC] = []nat.PortBinding{
		{
			HostPort: grpcPort,
		},
	}
	bindings[npHTTP] = []nat.PortBinding{
		{
			HostPort: httpPort,
		},
	}
	hostconfig := container.HostConfig{
		Binds:        []string{containerName + ":/data"},
		PortBindings: bindings,
	}

	netconfig := network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"waypoint": {},
		},
	}

	cfg.Labels = map[string]string{
		containerKey: containerValue,
	}

	cr, err := cli.ContainerCreate(ctx, &cfg, &hostconfig, &netconfig, containerName)
	if err != nil {
		return nil, err
	}

	err = cli.ContainerStart(ctx, cr.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	// KLUDGE: There isn't a way to find out if the container is up or not,
	// so we just give it 5 seconds to normalize before trying to use it.
	time.Sleep(5 * time.Second)

	s.Done()
	s = sg.Add("Server container started!")
	s.Done()

	return &InstallResults{
		Context:       &clicfg,
		AdvertiseAddr: &addr,
		HTTPAddr:      httpAddr,
	}, nil
}

// InstallRunner implements Installer by starting a single runner container.
func (i *DockerInstaller) InstallRunner(
	ctx context.Context,
	opts *InstallRunnerOpts,
) error {
	ui := opts.UI

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing Docker client...")
	defer func() { s.Abort() }()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	cli.NegotiateAPIVersion(ctx)

	s.Update("Checking for an existing runner...")
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "waypoint-type=runner",
		}),
	})
	if err != nil {
		return err
	}
	if len(containers) > 0 {
		s.Update("Detected existing Waypoint runner.")
		s.Status(terminal.StatusWarn)
		s.Done()
		return nil
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
		Image:        i.config.serverImage,
		Env:          opts.AdvertiseClient.Env(),
		Cmd:          []string{"runner", "agent", "-vvv"},
		Labels: map[string]string{
			"waypoint-type": "runner",
		},
	}, nil, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"waypoint": {},
		},
	}, "waypoint-runner")
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

func (i *DockerInstaller) InstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:    "docker-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: "hashicorp/waypoint:latest",
	})
}

func (i *DockerInstaller) Uninstall(
	ctx context.Context,
	opts *InstallOpts,
) error {
	sg := opts.UI.StepGroup()
	defer sg.Wait()

	// bulk of this copied from PR#660
	s := sg.Add("Initializing Docker client...")
	defer s.Abort()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	// TODO do we need this?
	// defer func() {
	// 	_ = cli.Close()
	// }()

	cli.NegotiateAPIVersion(ctx)

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: containerLabel,
		}),
	})

	if err != nil {
		return err
	}

	if len(containers) < 1 {
		return fmt.Errorf("cannot find a Waypoint Docker container")
	}

	// Pick the first container, as there should be only one.
	containerId := containers[0].ID

	s.Update("Stopping Waypoint Docker container...")

	// Stop the container gracefully, respecting the Engine's default timeout.
	if err := cli.ContainerStop(ctx, containerId, nil); err != nil {
		return err
	}

	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	if err := cli.ContainerRemove(ctx, containerId, removeOptions); err != nil {
		return err
	}
	s.Update("Docker container %q removed", containerName)
	s.Done()
	s = sg.Add("")

	volumeExists, err := volumeExists(ctx, cli)
	if err != nil {
		return err
	}

	s.Update("Removing Waypoint Docker volume...")
	// If the Waypoint Docker volume does not exist, return
	if !volumeExists {
		s.Update("Couldn't find Waypoint Docker volume %q; not removing", containerName)
		s.Status(terminal.StatusWarn)
		s.Done()
		return nil
	}

	if err := cli.VolumeRemove(ctx, containerName, true); err != nil {
		return err
	}
	s.Update("Docker volume %q removed", containerName)
	s.Done()

	return nil
}

func (i *DockerInstaller) UninstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:    "docker-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: "hashicorp/waypoint:latest",
	})
}

// volumeExists determines whether the Waypoint Docker volume exists.
func volumeExists(ctx context.Context, cli *client.Client) (bool, error) {
	listBody, err := cli.VolumeList(ctx, filters.NewArgs(filters.KeyValuePair{
		Key:   "name",
		Value: containerName,
	}))

	if err != nil {
		return false, err
	}

	exists := len(listBody.Volumes) > 0

	return exists, nil
}
