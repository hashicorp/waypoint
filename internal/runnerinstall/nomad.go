package runnerinstall

import (
	"context"
	"fmt"
	"github.com/hashicorp/nomad/api"
	nomad "github.com/hashicorp/waypoint/internal/installutil/nomad"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"strconv"
	"strings"
)

type NomadRunnerInstaller struct {
	config nomadConfig
}

const (
	runnerName = "waypoint-runner"
)

type nomadConfig struct {
	authSoftFail       bool              `hcl:"auth_soft_fail,optional"`
	image              string            `hcl:"server_image,optional"`
	namespace          string            `hcl:"namespace,optional"`
	serviceAnnotations map[string]string `hcl:"service_annotations,optional"`

	consulService            bool     `hcl:"consul_service,optional"`
	consulServiceUITags      []string `hcl:"consul_service_ui_tags:optional"`
	consulServiceBackendTags []string `hcl:"consul_service_backend_tags:optional"`
	consulDatacenter         string   `hcl:"consul_datacenter,optional"`
	consulDomain             string   `hcl:"consul_datacenter,optional"`

	// If set along with consul, will use this hostname instead of
	// making a consul DNS hostname for the server address in its context
	consulServiceHostname string `hcl:"consul_service_hostname,optional"`

	odrImage string `hcl:"odr_image,optional"`

	region         string   `hcl:"namespace,optional"`
	datacenters    []string `hcl:"datacenters,optional"`
	policyOverride bool     `hcl:"policy_override,optional"`

	serverResourcesCPU    string `hcl:"server_resources_cpu,optional"`
	serverResourcesMemory string `hcl:"server_resources_memory,optional"`
	runnerResourcesCPU    string `hcl:"runner_resources_cpu,optional"`
	runnerResourcesMemory string `hcl:"runner_resources_memory,optional"`

	hostVolume           string `hcl:"host_volume,optional"`
	csiVolumeProvider    string `hcl:"csi_volume_provider,optional"`
	csiVolumeCapacityMin int64  `hcl:"csi_volume_capacity_min,optional"`
	csiVolumeCapacityMax int64  `hcl:"csi_volume_capacity_max,optional"`
	csiFS                string `hcl:"csi_fs,optional"`

	nomadHost string `hcl:"nomad_host,optional"`
}

var (
	// default resources used for both the Server and its runners. Can be overridden
	// through config flags at install
	defaultResourcesCPU    = 200
	defaultResourcesMemory = 600

	// bytes
	defaultCSIVolumeCapacityMin = int64(1073741824)
	defaultCSIVolumeCapacityMax = int64(2147483648)

	defaultCSIVolumeMountFS = "xfs"

	// Defaults to use for setting up Consul
	defaultConsulServiceTag       = "waypoint"
	defaultConsulDatacenter       = "dc1"
	defaultConsulDomain           = "consul"
	waypointConsulBackendName     = "waypoint-server"
	waypointConsulUIName          = "waypoint-ui"
	defaultWaypointConsulHostname = fmt.Sprintf("%s.%s.service.%s.%s",
		defaultConsulServiceTag, waypointConsulBackendName, defaultConsulDatacenter, defaultConsulDomain)

	defaultNomadHost = "http://localhost:4646"
)

func (i *NomadRunnerInstaller) Install(ctx context.Context, opts *InstallOpts) error {
	//TODO implement me
	ui := opts.UI

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing Nomad client...")
	defer func() { s.Abort() }()

	//Build api client from environment
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	// Install the runner
	s.Update("Installing the Waypoint runner")
	_, err = nomad.RunJob(ctx, s, client, waypointRunnerNomadJob(i.config, opts), false)
	if err != nil {
		return err
	}
	s.Update("Waypoint runner installed")
	s.Done()

	return nil
}

// waypointRunnerNomadJob takes in a nomadConfig and returns a Nomad Job
// for the Nomad API to run a Waypoint runner.
func waypointRunnerNomadJob(c nomadConfig, opts *InstallOpts) *api.Job {
	job := api.NewServiceJob(runnerName+opts.Id, runnerName, c.region, 50)
	job.Namespace = &c.namespace
	job.Datacenters = c.datacenters
	job.Meta = c.serviceAnnotations
	tg := api.NewTaskGroup(runnerName, 1)
	tg.Networks = []*api.NetworkResource{
		{
			// Host mode so we can communicate to our server.
			Mode: "host",
		},
	}
	job.AddTaskGroup(tg)

	var image string
	if c.image == "" {
		image = defaultRunnerImage
	} else {
		image = c.image
	}

	task := api.NewTask("runner", "docker")
	task.Config = map[string]interface{}{
		"image": image,
		"args": []string{
			"runner",
			"agent",
			"-id=" + opts.Id,
			"-state-dir=/data/runner",
			"-cookie=" + opts.Cookie,
			"-vv",
		},
		"auth_soft_fail": c.authSoftFail,
	}

	cpu := defaultResourcesCPU
	mem := defaultResourcesMemory

	if c.runnerResourcesCPU != "" {
		cpu, _ = strconv.Atoi(c.runnerResourcesCPU)
	}
	if c.runnerResourcesMemory != "" {
		mem, _ = strconv.Atoi(c.runnerResourcesMemory)
	}
	task.Resources = &api.Resources{
		CPU:      &cpu,
		MemoryMB: &mem,
	}

	task.Env = map[string]string{}
	for _, line := range opts.AdvertiseClient.Env() {
		idx := strings.Index(line, "=")
		if idx == -1 {
			// Should never happen but let's not crash.
			continue
		}

		key := line[:idx]
		value := line[idx+1:]
		task.Env[key] = value
	}

	// Let the runner know about the Nomad IP
	if c.nomadHost == "" {
		c.nomadHost = defaultNomadHost
	}
	task.Env["NOMAD_ADDR"] = c.nomadHost

	tg.AddTask(task)

	return job
}

func (i *NomadRunnerInstaller) InstallFlags(set *flag.Set) {
	set.StringSliceVar(&flag.StringSliceVar{
		Name:    "nomad-dc",
		Target:  &i.config.datacenters,
		Default: []string{"dc1"},
		Usage:   "Datacenters to install to for Nomad.",
	})
}

func (i *NomadRunnerInstaller) Uninstall(ctx context.Context, opts *InstallOpts) error {
	//TODO implement me
	panic("implement me")
}

func (i *NomadRunnerInstaller) UninstallFlags(set *flag.Set) {
	//TODO implement me
	panic("implement me")
}
