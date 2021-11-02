package nomad

import (
	"context"
	"crypto/rand"
	"fmt"
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

const (
	// Build plugins like pack require a decemt amount of memory to build
	// an artifact. This default may seem large, but if we used the default
	// static runner default of 600 MB, it would OOM on a small Go app when
	// buildpack attempts to finish up its build. 2GB was choosen to be a little
	// more than what it might need so that Nomad doesn't OOM the task
	defaultODRMemory = 2000 // in mb
	defaultODRCPU    = 200  // in mhz

	defaultODRRegion     = "global"
	defaultODRDatacenter = "dc1"
	defaultODRNamespace  = "default"

	defaultNomadHost = "http://localhost:4646"
)

// TaskLauncherConfig is the configuration structure for the task plugin.
type TaskLauncherConfig struct {
	// The Datacenter the runner should be created and run in
	Datacenter string `hcl:"datacenter,optional"`

	// The namespace the runner should be created and run in
	Namespace string `hcl:"namespace,optional"`

	// The Nomad region to deploy the task to, defaults to "global"
	Region string `hcl:"region,optional"`

	// Resource request limits for an on-demand runner
	Memory int `hcl:"resources_memory,optional"`
	CPU    int `hcl:"resources_cpu,optional"`

	// The host to connect to for making Nomad API requests
	NomadHost string `hcl:"nomad_host,optional"`
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
		"The Nomad region to deploy the on-demand runner task to.",
		docs.Default(defaultODRRegion),
	)

	doc.SetField(
		"datacenter",
		"The Nomad datacenter to deploy the on-demand runner task to.",
		docs.Default(defaultODRDatacenter),
	)

	doc.SetField(
		"namespace",
		"The Nomad namespace to deploy the on-demand runner task to.",
		docs.Default(defaultODRNamespace),
	)

	doc.SetField(
		"resources_cpu",
		"Amount of CPU in MHz to allocate to this task. This can be overriden with "+
			"the '-nomad-runner-cpu' flag on server install.",
		docs.Default(fmt.Sprint(defaultODRCPU)),
	)

	doc.SetField(
		"resources_memory",
		"Amount of memory in MB to allocate to this task. This can be overriden with "+
			"the '-nomad-runner-memory' flag on server install.",
		docs.Default(fmt.Sprint(defaultODRMemory)),
	)

	doc.SetField(
		"nomad_host",
		"Hostname of the Nomad server to use for launching on-demand tasks.",
		docs.Default(defaultNomadHost),
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
		p.config.Region = defaultODRRegion
	}
	if p.config.Datacenter == "" {
		p.config.Datacenter = defaultODRDatacenter
	}
	if p.config.Namespace == "" {
		p.config.Namespace = defaultODRNamespace
	}
	if p.config.Memory == 0 {
		p.config.Memory = defaultODRMemory
	}
	if p.config.CPU == 0 {
		p.config.CPU = defaultODRCPU
	}
	if p.config.NomadHost == "" {
		p.config.NomadHost = defaultNomadHost
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

	task.Resources = &api.Resources{
		CPU:      &p.config.CPU,
		MemoryMB: &p.config.Memory,
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
	task.Env = env

	// Let the on-demand runner know about the Nomad IP
	task.Env["NOMAD_ADDR"] = p.config.NomadHost

	job.TaskGroups[0].Tasks[0].Env = env

	// On-Demand runner specific configuration to start the task with
	config := map[string]interface{}{
		"image":   tli.OciUrl,
		"args":    tli.Arguments,
		"command": tli.Entrypoint,
	}

	job.TaskGroups[0].Tasks[0].Config = config

	log.Debug("registering on-demand task job", "task-name", taskName)
	_, _, err = jobclient.Register(job, nil)
	if err != nil {
		log.Debug("failed to register job to nomad")
		return nil, err
	}

	log.Debug("finished launching on-demand task for build", "task-name", taskName)
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
