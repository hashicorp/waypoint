package ecs

import (
	"context"

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
	// ServiceAccount is the name of the AWS service account to apply to the
	// application deployment. This is useful to apply AWS RBAC to the pod.
	ServiceAccount string `hcl:"service_account,optional"`

	// Set an explicit pull policy for this task launching. By default
	// we use "PullIfNotPresent" unless the image tag is "latest" when we
	// use "Always".
	PullPolicy string `hcl:"image_pull_policy,optional"`
}

func (p *TaskLauncher) Documentation() (*docs.Documentation, error) {
	doc, err := docs.New(
		docs.FromConfig(&TaskLauncherConfig{}),
		docs.FromFunc(p.StartTaskFunc()),
	)
	if err != nil {
		return nil, err
	}

	doc.Description(`
Launch an ECS task for on-demand tasks from the Waypoint server.

This will use the standard AWS environment variables and service account to
source authentication information for AWS. If this is running within ECS
itself (typical for a ECS-based installation), it will use the task's
service account unless other auth is explicitly given. This allows the task
launcher to work by default.
`)

	doc.Example(`
task {
	use "ecs" {}
}
`)

	doc.SetField(
		"service_account",
		"service account name to be used for the execution role in the task",
		docs.Summary(
			"service account is the name of the AWS service account to use "+
				"as the execution role for the task.",
		),
	)

	doc.SetField(
		"image_pull_policy",
		"pull policy to use for the task container image",
	)

	return doc, nil
}

// TaskLauncher implements Configurable
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
	return nil, nil
}

var _ component.TaskLauncher = (*TaskLauncher)(nil)
