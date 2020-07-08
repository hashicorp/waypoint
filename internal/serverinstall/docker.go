package serverinstall

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/waypoint/internal/clicontext"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

func InstallDocker(
	ctx context.Context, ui terminal.UI, st terminal.Status, scfg *Config) (
	*clicontext.Config, *pb.ServerConfig_AdvertiseAddr, error,
) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, nil, err
	}

	cli.NegotiateAPIVersion(ctx)

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "waypoint-type=server",
		}),
	})

	if err != nil {
		return nil, nil, err
	}

	port := "9701"

	var (
		clicfg clicontext.Config
		addr   pb.ServerConfig_AdvertiseAddr
	)

	clicfg.Server = configpkg.Server{
		Address:  "localhost:" + port,
		Insecure: true,
	}

	addr.Addr = "waypoint-server:" + port
	addr.Insecure = true

	// If we already have a server, bolt.
	if len(containers) > 0 {
		st.Step(terminal.StatusWarn, "Detected existing waypoint server")
		return &clicfg, &addr, nil
	}

	st.Update("Creating waypoint network...")

	nets, err := cli.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "use=waypoint")),
	})

	if err != nil {
		return nil, nil, err
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
			return nil, nil, err
		}
	}

	np, err := nat.NewPort("tcp", port)
	if err != nil {
		return nil, nil, err
	}

	st.Update("Installing waypoint server to docker")

	cfg := container.Config{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		OpenStdin:    true,
		StdinOnce:    true,
		Image:        scfg.ServerImage,
		ExposedPorts: nat.PortSet{np: struct{}{}},
		Env:          []string{"PORT=" + port},
		Cmd:          []string{"server", "-vvv", "-db=/data/data.db", "-listen-grpc=0.0.0.0:9701"},
	}

	bindings := nat.PortMap{}
	bindings[np] = []nat.PortBinding{
		{
			HostPort: port,
		},
	}

	hostconfig := container.HostConfig{
		Binds:        []string{"waypoint-server:/data"},
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

	cr, err := cli.ContainerCreate(ctx, &cfg, &hostconfig, &netconfig, "waypoint-server")
	if err != nil {
		return nil, nil, err
	}

	err = cli.ContainerStart(ctx, cr.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, nil, err
	}

	// KLUDGE: There isn't a way to find out if the container is up or not, so we just give it 5 seconds
	// to normalize before trying to use it.
	time.Sleep(5 * time.Second)

	st.Step(terminal.StatusOK, "Server container started")

	return &clicfg, &addr, nil
}
