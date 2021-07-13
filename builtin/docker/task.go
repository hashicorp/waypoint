package docker

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	goUnits "github.com/docker/go-units"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	wpdockerclient "github.com/hashicorp/waypoint/builtin/docker/client"
)

// TaskLauncher uses `docker build` to build a Docker iamge.
type TaskLauncher struct {
	config TaskLauncherConfig
}

// BuildFunc implements component.TaskLauncher
func (b *TaskLauncher) StartTaskFunc() interface{} {
	return b.StartTask
}

// BuildFunc implements component.TaskLauncher
func (b *TaskLauncher) StopTaskFunc() interface{} {
	return b.StopTask
}

type TaskResources struct {
	// How many CPU shares to allocate to each task
	CpuShares int64 `hcl:"cpu,optional"`

	// How much memory to allocate to each task
	MemoryLimit string `hcl:"memory,optional"`
}

// TaskLauncherConfig is the configuration structure for the task plugin.
type TaskLauncherConfig struct {
	// A list of folders to mount to the container.
	Binds []string `hcl:"binds,optional"`

	// ClientConfig allow the user to specify the connection to the Docker
	// engine. By default we try to load this from env vars:
	// DOCKER_HOST to set the url to the docker server.
	// DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	// DOCKER_CERT_PATH to load the TLS certificates from.
	// DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default.
	ClientConfig *ClientConfig `hcl:"client_config,block"`

	// Force pull the image from the remote repository
	ForcePull bool `hcl:"force_pull,optional"`

	// A map of key/value pairs, stored in docker as a string. Each key/value pair must
	// be unique. Validiation occurs at the docker layer, not in Waypoint. Label
	// keys are alphanumeric strings which may contain periods (.) and hyphens (-).
	// See the docker docs for more info: https://docs.docker.com/config/labels-custom-metadata/
	Labels map[string]string `hcl:"labels,optional"`

	// An array of strings with network names to connect the container to
	Networks []string `hcl:"networks,optional"`

	// Resources configures the resource constraints such as cpu and memory for the
	// created containers.
	Resources TaskResources `hcl:"resources,block"`

	// Environment variables that are meant to configure the application in a static
	// way. This might be start an image in a specific mode,
	// selected via environment variable. Most configuration should use the waypoint
	// config commands.
	StaticEnvVars map[string]string `hcl:"static_environment,optional"`
}

func (b *TaskLauncher) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(
		docs.FromConfig(&TaskLauncherConfig{}),
		docs.FromFunc(b.StartTaskFunc()),
	)
	if err != nil {
		return nil, err
	}

	doc.Description(`
Launch a Docker container as a task.

If a Docker server is available (either locally or via environment variables
such as "DOCKER_HOST"), then it will be used to start the container.
`)

	doc.Example(`
task {
  use "docker" {
		force_pull = true
  }
}
`)

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
		"A list of strings with network names to connect the container to.",
		docs.Default("waypoint"),
		docs.Summary(
			"A list of networks to connect the container to. By default the container",
			"will always connect to the `waypoint` network.",
		),
	)

	doc.SetField(
		"resources",
		"The resources that the tasks should use.",
		docs.SubFields(func(d *docs.SubFieldDoc) {
			d.SetField("cpu", "The cpu shares that the tasks should use")
			d.SetField("memory", "The amount of memory to use. Defined as '512MB', '44kB', etc.")
		}),
	)

	doc.SetField(
		"static_environment",
		"environment variables to expose to the application",
		docs.Summary(
			"These variables are used to control all of a container's modes,",
			"such as configuring it to start a web app vs a background worker.",
			"These environment variables should not be common",
			"configuration variables normally set in `waypoint config`.",
		),
	)

	return doc, nil
}

// TaskLauncher implements Configurable
func (b *TaskLauncher) Config() (interface{}, error) {
	return &b.config, nil
}

func (b *TaskLauncher) setupImage(
	ctx context.Context,
	log hclog.Logger,
	cli *client.Client,
	img string,
) error {
	args := filters.NewArgs()
	args.Add("reference", img)

	// only pull if image is not in current registry so check to see if the image is present
	// if force then skip this check
	if !b.config.ForcePull {
		sum, err := cli.ImageList(context.Background(), types.ImageListOptions{Filters: args})
		if err != nil {
			return status.Errorf(codes.FailedPrecondition, "unable to list images in local Docker cache: %s", err)
		}

		log.Debug("image list", "images", len(sum))

		// if we have images do not pull
		if len(sum) > 0 {
			log.Info("reusing existing image for task", "image", img, "id", sum[0].ID)
			return nil
		}
	}

	img = makeImageCanonical(img)

	out, err := cli.ImagePull(context.Background(), img, types.ImagePullOptions{})
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "unable to pull image: %s", err)
	}

	var stdout bytes.Buffer

	err = jsonmessage.DisplayJSONMessagesStream(out, &stdout, 0, false, nil)
	if err != nil {
		log.Error("error pulling image for task", "image", img, "output", stdout.String())
		return status.Errorf(codes.Internal, "unable to stream build logs to the terminal: %s", err)
	} else {
		log.Debug("finished pulling image for task", "image", img, "output", stdout.String())
	}

	log.Info("pulled image for task", "image", img)

	return nil
}

func (b *TaskLauncher) setupNetworking(
	ctx context.Context,
	cli *client.Client,
) (string, error) {
	nets, err := cli.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "use=waypoint")),
	})
	if err != nil {
		return "", status.Errorf(codes.FailedPrecondition, "unable to list Docker networks: %s", err)
	}

	if len(nets) > 1 {
		// We use whichever network has the use=waypoint label, allowing the user to configure
		// a network themselves with whatever name they wish.
		return nets[0].Name, nil
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
			return "", status.Errorf(codes.FailedPrecondition, "unable to create Docker network: %s", err)
		}
	}

	return "waypoint", nil
}

// StopTask signals to docker to stop the container created previously
func (b *TaskLauncher) StopTask(
	ctx context.Context,
	log hclog.Logger,
	ti *TaskInfo,
) error {
	cli, err := wpdockerclient.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "unable to create Docker client: %s", err)
	}
	cli.NegotiateAPIVersion(ctx)

	return cli.ContainerStop(ctx, ti.Id, nil)
}

// StartTask creates a docker container for the task.
func (b *TaskLauncher) StartTask(
	ctx context.Context,
	log hclog.Logger,
	tli *component.TaskLaunchInfo,
) (*TaskInfo, error) {
	cli, err := wpdockerclient.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "unable to create Docker client: %s", err)
	}
	cli.NegotiateAPIVersion(ctx)

	err = b.setupImage(ctx, log, cli, tli.OciUrl)
	if err != nil {
		return nil, err
	}

	netName, err := b.setupNetworking(ctx, cli)
	if err != nil {
		return nil, err
	}

	randId, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("waypoint-task-%s", randId)

	var env []string

	for k, v := range tli.EnvironmentVariables {
		env = append(env, k+"="+v)
	}

	log.Debug(
		"spawn docker container for task",
		"oci-url", tli.OciUrl,
		"arguments", tli.Arguments,
		"environment", env,
	)

	var memory int64

	if b.config.Resources.MemoryLimit != "" {
		memory, err = goUnits.FromHumanSize(b.config.Resources.MemoryLimit)
		if err != nil {
			return nil, err
		}
	}

	cc, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Env:   env,
			Cmd:   tli.Arguments,
			Image: tli.OciUrl,
		},
		&container.HostConfig{
			Binds:      b.config.Binds,
			AutoRemove: true,

			Resources: container.Resources{
				CPUShares: b.config.Resources.CpuShares,
				Memory:    memory,
			},
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				netName: {},
			},
		},
		nil,
		name,
	)
	if err != nil {
		return nil, err
	}

	if b.config.Networks != nil {
		for _, net := range b.config.Networks {
			err = cli.NetworkConnect(ctx, net, cc.ID, &network.EndpointSettings{})
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					"unable to connect container to additional networks: %s",
					err)
			}
		}
	}

	err = cli.ContainerStart(ctx, cc.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	log.Info("launched task container", "id", cc.ID, "name", name)

	ti := &TaskInfo{
		Id: cc.ID,
	}

	return ti, nil
}
