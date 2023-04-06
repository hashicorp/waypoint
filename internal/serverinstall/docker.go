// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"github.com/hashicorp/waypoint/internal/installutil"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/runnerinstall"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

type DockerInstaller struct {
	config dockerConfig
}

type dockerConfig struct {
	serverImage      string `hcl:"server_image,optional"`
	odrImage         string `hcl:"odr_image,optional"`
	runnerSocketPath string `hcl:"runner_socket_path,optional"`
}

var (
	grpcPort             = defaultGrpcPort
	httpPort             = defaultHttpPort
	containerLabel       = "waypoint-type=server"
	containerKey         = "waypoint-type"
	containerValue       = "server"
	containerValueRunner = "runner"
)

// Install is a method of DockerInstaller and implements the Installer interface to
// create a waypoint-server as a Docker container
func (i *DockerInstaller) Install(
	ctx context.Context,
	opts *InstallOpts,
) (*InstallResults, string, error) {
	if i.config.odrImage == "" {
		var err error
		i.config.odrImage, err = installutil.DeriveDefaultODRImage(i.config.serverImage)
		if err != nil {
			return nil, "", err
		}
	}

	ui := opts.UI
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing Docker client...")
	defer func() { s.Abort() }()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, "", err
	}
	cli.NegotiateAPIVersion(ctx)

	s.Update("Checking for existing installation...")

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: true, // include stopped containers
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: containerLabel,
		}),
	})
	if err != nil {
		return nil, "", err
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
		Platform:      "docker",
	}

	addr.Addr = serverName + ":" + grpcPort
	addr.Tls = true
	addr.TlsSkipVerify = true

	httpAddr = "localhost:" + httpPort

	// If we already have a server, bolt.
	if len(containers) > 0 {
		s.Update("Detected existing Waypoint server.")
		s.Status(terminal.StatusWarn)
		s.Done()

		// In the case where waypoint server container isn't running, the installer
		// will attempt to start the container. It does this for all containers
		// that match the 'containerLabel'. In the future case where we support
		// running multiple waypoint server containers, this loop will try to start
		// each container.
		for _, container := range containers {
			if container.State != "running" {
				s = sg.Add("Attempting to start container...")

				err = cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
				if err != nil {
					s.Update("Failed to start container %q", container.Names[0])
					s.Status(terminal.StatusError)
					s.Done()
					return nil, "", err
				}

				s.Update("Container %q started!", container.Names[0])
				s.Done()
			}
		}

		return &InstallResults{
			Context:       &clicfg,
			AdvertiseAddr: &addr,
			HTTPAddr:      httpAddr,
		}, "", nil
	}

	s.Update("Checking for Docker image: %s", i.config.serverImage)

	imageRef, err := reference.ParseNormalizedNamed(i.config.serverImage)
	if err != nil {
		return nil, "", fmt.Errorf("Error parsing Docker image: %s", err)
	}

	imageList, err := cli.ImageList(ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "reference",
			Value: reference.FamiliarString(imageRef),
		}),
	})
	if err != nil {
		return nil, "", err
	}

	if len(imageList) == 0 || i.config.serverImage == installutil.DefaultServerImage {
		s.Update("Pulling image: %s", i.config.serverImage)

		resp, err := cli.ImagePull(ctx, reference.FamiliarString(imageRef), types.ImagePullOptions{})
		if err != nil {
			return nil, "", err
		}
		defer resp.Close()

		stdout, _, err := ui.OutputWriters()
		if err != nil {
			return nil, "", err
		}

		var termFd uintptr
		if f, ok := stdout.(*os.File); ok {
			termFd = f.Fd()
		}

		err = jsonmessage.DisplayJSONMessagesStream(resp, s.TermOutput(), termFd, true, nil)
		if err != nil {
			return nil, "", fmt.Errorf("unable to stream pull logs to the terminal: %s", err)
		}

		s.Done()
		s = sg.Add("")
	}

	s.Update("Creating waypoint network...")

	nets, err := cli.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "use=waypoint")),
	})
	if err != nil {
		return nil, "", err
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
			return nil, "", err
		}

	}

	npGRPC, err := nat.NewPort("tcp", grpcPort)
	if err != nil {
		return nil, "", err
	}

	npHTTP, err := nat.NewPort("tcp", httpPort)
	if err != nil {
		return nil, "", err
	}

	s.Update("Installing Waypoint server to docker")

	cmd := []string{"server", "run", "-accept-tos", "-vv", "-db=/data/data.db", fmt.Sprintf("-listen-grpc=0.0.0.0:%s", grpcPort), fmt.Sprintf("-listen-http=0.0.0.0:%s", httpPort)}
	cmd = append(cmd, opts.ServerRunFlags...)
	cfg := container.Config{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		OpenStdin:    true,
		StdinOnce:    true,
		Image:        i.config.serverImage,
		ExposedPorts: nat.PortSet{npGRPC: struct{}{}, npHTTP: struct{}{}},
		Env:          []string{"PORT=" + grpcPort},
		Cmd:          cmd,
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
		Binds:        []string{fmt.Sprintf("%s:/data", serverName)},
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

	cr, err := cli.ContainerCreate(ctx, &cfg, &hostconfig, &netconfig, nil, serverName)
	if err != nil {
		return nil, "", err
	}

	err = cli.ContainerStart(ctx, cr.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, "", err
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
	}, "", nil
}

// Upgrade is a method of DockerInstaller and implements the Installer interface to
// upgrade a waypoint-server as a Docker container
func (i *DockerInstaller) Upgrade(
	ctx context.Context, opts *InstallOpts, serverCfg serverconfig.Client) (
	*InstallResults, error,
) {
	if i.config.odrImage == "" {
		var err error
		i.config.odrImage, err = installutil.DeriveDefaultODRImage(i.config.serverImage)
		if err != nil {
			return nil, err
		}
	}

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

	s.Update("Checking for an existing Waypoint server installation...")
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: true, // include stopped containers
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "waypoint-type=server",
		}),
	})
	if err != nil {
		return nil, err
	}

	grpcPort := defaultGrpcPort
	httpPort := defaultHttpPort

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

	addr.Addr = serverName + ":" + grpcPort
	addr.Tls = true
	addr.TlsSkipVerify = true

	httpAddr = "localhost:" + httpPort

	if len(containers) == 0 {
		s.Update("No waypoint server detected. Nothing to upgrade.")
		s.Status(terminal.StatusWarn)
		s.Done()
		return nil, fmt.Errorf("No waypoint server container detected")
	}

	// Assume waypoint-server is the first container with the waypoint-type label
	waypointServerContainer := containers[0]

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

	if len(imageList) == 0 || i.config.serverImage == installutil.DefaultServerImage {
		s.Done()
		s = sg.Add("Pulling image: %s", i.config.serverImage)

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

	s.Update(
		"Upgrading Waypoint server image from %q to %q",
		waypointServerContainer.Image,
		i.config.serverImage,
	)
	s.Done()

	s = sg.Add("Removing and restarting current server container")
	// stop and remove container
	err = cli.ContainerStop(ctx, waypointServerContainer.ID, nil)
	if err != nil {
		return nil, err
	}
	err = cli.ContainerRemove(ctx, waypointServerContainer.ID, types.ContainerRemoveOptions{
		Force:         true,
		RemoveVolumes: false,
	})
	if err != nil {
		return nil, err
	}

	npGRPC, err := nat.NewPort("tcp", grpcPort)
	if err != nil {
		return nil, err
	}

	npHTTP, err := nat.NewPort("tcp", httpPort)
	if err != nil {
		return nil, err
	}

	cmd := []string{"server", "run", "-accept-tos", "-vv", "-db=/data/data.db", fmt.Sprintf("-listen-grpc=0.0.0.0:%s", grpcPort), fmt.Sprintf("-listen-http=0.0.0.0:%s", httpPort)}
	cmd = append(cmd, opts.ServerRunFlags...)
	cfg := container.Config{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		OpenStdin:    true,
		StdinOnce:    true,
		Image:        i.config.serverImage,
		ExposedPorts: nat.PortSet{npGRPC: struct{}{}, npHTTP: struct{}{}},
		Env:          []string{"PORT=" + grpcPort},
		Cmd:          cmd,
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
		Binds:        []string{fmt.Sprintf("%s:/data", serverName)},
		PortBindings: bindings,
	}

	netconfig := network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"waypoint": {},
		},
	}

	cfg.Labels = map[string]string{
		"waypoint-type": "server",
	}
	s.Update("Creating and starting container")
	//
	cr, err := cli.ContainerCreate(ctx, &cfg, &hostconfig, &netconfig, nil, serverName)
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

	s.Update("Server container started!")
	s.Done()

	return &InstallResults{
		Context:       &clicfg,
		AdvertiseAddr: &addr,
		HTTPAddr:      httpAddr,
	}, nil
}

// Install is a method of DockerInstaller and implements the Installer interface to
// remove the waypoint-server Docker container and associated image and volume
func (i *DockerInstaller) Uninstall(
	ctx context.Context,
	opts *InstallOpts,
) error {
	sg := opts.UI.StepGroup()
	defer sg.Wait()

	// used base functionality from PR#660
	s := sg.Add("Initializing Docker client...")
	defer func() { s.Abort() }()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	defer cli.Close()

	cli.NegotiateAPIVersion(ctx)

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: true, // include stopped containers
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: containerLabel,
		}),
	})

	if err != nil {
		return err
	}

	if len(containers) < 1 {
		return fmt.Errorf(
			"cannot find a Waypoint Docker container; Waypoint may already be uninstalled.",
		)
	}

	// Pick the first container, as there should be only one.
	containerId := containers[0].ID
	image := containers[0].Image

	imageRef, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return fmt.Errorf("Error parsing Docker image: %s", err)
	}

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
	s.Update("Docker container %q removed", serverName)
	s.Done()
	s = sg.Add("")

	s.Update("Removing Waypoint Docker volume...")
	// Find volume of the server
	vl, err := cli.VolumeList(ctx, filters.NewArgs(filters.KeyValuePair{
		Key:   "name",
		Value: serverName,
	}))
	if err != nil {
		return err
	}
	volumeExists := len(vl.Volumes) > 0

	// If the Waypoint Docker volume does not exist, we keep going and just warn
	if !volumeExists {
		s.Update("Couldn't find Waypoint Docker volume %q; not removing", serverName)
		s.Status(terminal.StatusWarn)
		s.Done()
	} else {
		if err := cli.VolumeRemove(ctx, serverName, true); err != nil {
			return err
		}
		s.Update("Docker volume %q removed", serverName)
		s.Done()
	}

	s = sg.Add("")

	imageList, err := cli.ImageList(ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "reference",
			Value: reference.FamiliarString(imageRef),
		}),
	})
	if err != nil {
		return err
	}
	if len(imageList) < 1 {
		s.Update("Could not find image %q, not removing", imageRef.Name())
		s.Status(terminal.StatusWarn)
		s.Done()
		return nil
	}

	// Pick the first image, as there should be only one.
	imageId := imageList[0].ID
	_, err = cli.ImageRemove(ctx, imageId, types.ImageRemoveOptions{})
	// If we can't remove the image, we keep going and just warn
	if err != nil {
		s.Update("Could not remove image %q: %s", imageRef.Name(), err)
		s.Status(terminal.StatusWarn)
		s.Done()
		return nil
	}

	s.Update("Docker image %q removed", imageRef.Name())
	s.Done()

	return nil
}

// InstallRunner implements Installer by starting a single runner container.
func (i *DockerInstaller) InstallRunner(
	ctx context.Context,
	opts *runnerinstall.InstallOpts,
) error {
	runnerInstaller := runnerinstall.DockerRunnerInstaller{Config: runnerinstall.DockerConfig{
		RunnerImage: i.config.serverImage,
		Network:     "waypoint",
		SocketPath:  i.config.runnerSocketPath,
	}}
	err := runnerInstaller.Install(ctx, opts)
	if err != nil {
		return err
	}
	return nil
}

func (i *DockerInstaller) OnDemandRunnerConfig() *pb.OnDemandRunnerConfig {
	return &pb.OnDemandRunnerConfig{
		Name:       "docker",
		OciUrl:     i.config.odrImage,
		PluginType: "docker",
		Default:    true,
	}
}

// UninstallRunner implements Installer.
func (i *DockerInstaller) UninstallRunner(
	ctx context.Context,
	opts *runnerinstall.InstallOpts,
) error {
	runnerInstaller := runnerinstall.DockerRunnerInstaller{
		Config: runnerinstall.DockerConfig{},
	}

	err := runnerInstaller.Uninstall(ctx, opts)
	if err != nil {
		return err
	}
	return nil
}

// HasRunner implements Installer.
func (i *DockerInstaller) HasRunner(
	ctx context.Context,
	opts *InstallOpts,
) (bool, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false, err
	}
	defer cli.Close()
	cli.NegotiateAPIVersion(ctx)

	// Find and delete any runners. There could be zero, 1, or more.
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: true, // include stopped containers
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: containerKey + "=" + containerValueRunner,
		}),
	})
	if err != nil {
		return false, err
	}

	return len(containers) > 0, nil
}

func (i *DockerInstaller) InstallFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:    "docker-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: installutil.DefaultServerImage,
	})

	set.StringVar(&flag.StringVar{
		Name:   "docker-odr-image",
		Target: &i.config.odrImage,
		Usage: "Docker image for the Waypoint On-Demand Runners. This will " +
			"default to the server image with the name (not label) suffixed with '-odr'.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "docker-runner-socket-path",
		Target:  &i.config.runnerSocketPath,
		Usage:   "The path of the Docker socket that will be bound in runner",
		Default: "/var/run/docker.sock",
	})
}

func (i *DockerInstaller) UpgradeFlags(set *flag.Set) {
	set.StringVar(&flag.StringVar{
		Name:    "docker-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: installutil.DefaultServerImage,
	})

	set.StringVar(&flag.StringVar{
		Name:   "docker-odr-image",
		Target: &i.config.odrImage,
		Usage: "Docker image for the Waypoint On-Demand Runners. This will " +
			"default to the server image with the name (not label) suffixed with '-odr'.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "docker-runner-socket-path",
		Target:  &i.config.runnerSocketPath,
		Usage:   "The path of the Docker socket that will be bound in runner",
		Default: "/var/run/docker.sock",
	})
}

func (i *DockerInstaller) UninstallFlags(set *flag.Set) {
	// Purposely empty, no flags
}
