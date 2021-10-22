package nomad

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
)

// TaskLauncher implements the TaskLauncher plugin interface to support
// launching on-demand tasks for the Waypoint server.
type TaskLauncher struct {
	config TaskLauncherConfig
}

// StartTaskFunc implements component.TaskLauncher
func (p *TaskLauncher) StartTaskFunc() interface{} {
	return p.StartTask
}

// StopTaskFunc implements component.TaskLauncher
func (p *TaskLauncher) StopTaskFunc() interface{} {
	return p.StopTask
}

// TaskLauncherConfig is the configuration structure for the task plugin.
type TaskLauncherConfig struct {
}

func (p *TaskLauncher) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(
		docs.FromConfig(&TaskLauncherConfig{}),
		docs.FromFunc(p.StartTaskFunc()),
	)
	if err != nil {
		return nil, err
	}

	doc.Description(`this`)
	return &docs.Documentation{}, nil
}

// Config implements Configurable
func (p *TaskLauncher) Config() (interface{}, error) {
	return &p.config, nil
}

// StopTask signals to docker to stop the container created previously
func (p *TaskLauncher) StopTask(
	ctx context.Context,
	log hclog.Logger,
	ti *TaskInfo,
) error {
	return nil
}

// StartTask creates a docker container for the task.
func (p *TaskLauncher) StartTask(
	ctx context.Context,
	log hclog.Logger,
	tli *component.TaskLaunchInfo,
) (*TaskInfo, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	// Generate an ID for our task name.
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}

	taskName := fmt.Sprintf("waypoint-odr-task-%s", id.String())

	return nil, nil
}
