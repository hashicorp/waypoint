package nomad

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/docs"
)

// TaskLauncher implements the TaskLauncher plugin interface to support
// launching on-demand tasks for the Waypoint server.
type TaskLauncher struct {
	config TaskLauncherConfig
}

// StartTaskFunc implements component.TaskLauncher.
func (p *TaskLauncher) StartTaskFunc() interface{} {
	return p.StartTask
}

// StopTaskFunc implements component.TaskLauncher.
func (p *TaskLauncher) StopTaskFunc() interface{} {
	return p.StopTask
}

// TaskLauncherConfig is the configuration structure for the task plugin.
type TaskLauncherConfig struct {
	// The Datacenter the runner should be created and run in
	Datacenter string `hcl:"datacenter,optional"`

	// The namespace the runner should be created and run in
	Namespace string `hcl:"namespace,optional"`

	// The Nomad region to deploy the task to, defaults to "global"
	Region string `hcl:"region,optional"`

	// Resource request limits for an on-demand runner
	Memory string `hcl:"resources_memory,optional"`
	CPU    string `hcl:"resources_cpu,optional"`
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
Launch a Nomad job for on-demand tasks from the Waypoint server.

This will use the standard Nomad environment used for with the server install
to launch on demand Nomad jobs for Waypoint server tasks.
	`)

	doc.Example(`
task {
	use "nomad" {}
}
`)

	doc.SetField(
		"region",
		"The Nomad region to deploy the job to.",
		docs.Default("global"),
	)

	doc.SetField(
		"datacenter",
		"The Nomad datacenter to deploy the job to.",
		docs.Default("dc1"),
	)

	doc.SetField(
		"namespace",
		"The Nomad namespace to deploy the job to.",
	)

	doc.SetField(
		"resources_cpu",
		"Amount of CPU in MHz to allocate to this task",
		docs.Default("200"),
	)

	doc.SetField(
		"resources_memory",
		"Amount of memory in MB to allocate to this task.",
		docs.Default("600"),
	)

	return doc, nil
}

// Config implements Configurable.
func (p *TaskLauncher) Config() (interface{}, error) {
	return &p.config, nil
}

// StopTask signals to Nomad to stop the nomad job created previously.
func (p *TaskLauncher) StopTask(
	ctx context.Context,
	log hclog.Logger,
	ti *TaskInfo,
) error {
	client, err := p.getNomadClient()
	if err != nil {
		log.Error("failed to create a Nomad API client to stop an ODR task")
		return err
	}

	_, _, err = client.Jobs().Deregister(ti.Id, true, nil)
	return err
}

// StartTask creates a Nomad job for working on the task.
func (p *TaskLauncher) StartTask(
	ctx context.Context,
	log hclog.Logger,
	tli *component.TaskLaunchInfo,
) (*TaskInfo, error) {
	client, err := p.getNomadClient()
	if err != nil {
		log.Error("failed to create a Nomad API client to start an ODR task")
		return nil, err
	}

	// Generate an ID for our pod name.
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// Generate unique task name
	taskName := strings.ToLower(fmt.Sprintf("waypoint-task-%s", id.String()))

	// Set some defaults
	if p.config.Region == "" {
		p.config.Region = "global"
	}
	if p.config.Datacenter == "" {
		p.config.Region = "dc1"
	}
	if p.config.Namespace == "" {
		p.config.Namespace = "default"
	}
	if p.config.Memory != "" {
		p.config.Memory = "600"
	}
	if p.config.CPU != "" {
		p.config.CPU = "200"
	}

	log.Trace("creating Nomad job for task")
	jobclient := client.Jobs()
	job := api.NewServiceJob(taskName, taskName, p.config.Region, 10)
	job.Datacenters = []string{p.config.Datacenter}
	tg := api.NewTaskGroup(taskName, 1)
	tg.Networks = []*api.NetworkResource{
		{
			Mode: "host",
		},
	}

	job.Namespace = &p.config.Namespace
	job.AddTaskGroup(tg)
	task := &api.Task{
		Name:   taskName,
		Driver: "docker",
	}

	cpu, _ := strconv.Atoi(p.config.CPU)
	mem, _ := strconv.Atoi(p.config.Memory)
	task.Resources = &api.Resources{
		CPU:      &cpu,
		MemoryMB: &mem,
	}

	tg.AddTask(task)

	// Set our ID on the meta.
	job.SetMeta(metaId, taskName)
	job.SetMeta(metaNonce, time.Now().UTC().Format(time.RFC3339Nano))

	// Build our env vars
	env := map[string]string{}
	for k, v := range tli.EnvironmentVariables {
		env[k] = v
	}

	job.TaskGroups[0].Tasks[0].Env = env

	config := map[string]interface{}{
		"image": tli.OciUrl,
	}

	// TODO set auth here for pulling ODR image? not needed? we don't do it on install
	//if p.config.Auth != nil {
	//	config["auth"] = map[string]interface{}{
	//		"username": p.config.Auth.Username,
	//		"password": p.config.Auth.Password,
	//	}
	//}

	job.TaskGroups[0].Tasks[0].Config = config

	log.Debug("registering on-demand task job %q...", taskName)
	_, _, err = jobclient.Register(job, nil)
	if err != nil {
		return nil, err
	}

	// TODO: wait for allocation to be scheduled
	//log.Debug("waiting for allocation to be scheduled...")
	// Wait on the allocation
	//evalID := regResult.EvalID

	return &TaskInfo{
		Id: taskName,
	}, nil
}

// getNomadClient provides the client connection used by resources to interact with Nomad.
func (p *TaskLauncher) getNomadClient() (*api.Client, error) {
	// Get our client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return client, nil
}
