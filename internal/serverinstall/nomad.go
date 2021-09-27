package serverinstall

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverconfig"
)

type NomadInstaller struct {
	config nomadConfig
}

type nomadConfig struct {
	authSoftFail             bool              `hcl:"auth_soft_fail,optional"`
	serverImage              string            `hcl:"server_image,optional"`
	namespace                string            `hcl:"namespace,optional"`
	serviceAnnotations       map[string]string `hcl:"service_annotations,optional"`
	consulService            bool              `hcl:"consul_service,optional"`
	consulServiceUITags      []string          `hcl:"consul_service_ui_tags:optional"`
	consulServiceBackendTags []string          `hcl:"consul_service_backend_tags:optional"`

	region         string   `hcl:"namespace,optional"`
	datacenters    []string `hcl:"datacenters,optional"`
	policyOverride bool     `hcl:"policy_override,optional"`

	serverResourcesCPU    string `hcl:"server_resources_cpu,optional"`
	serverResourcesMemory string `hcl:"server_resources_memory,optional"`
	runnerResourcesCPU    string `hcl:"runner_resources_cpu,optional"`
	runnerResourcesMemory string `hcl:"runner_resources_memory,optional"`

	volumeType        string `hcl:"volume_type,optional"`
	hostVolume        string `hcl:"host_volume,optional"`
	csiVolumeProvider string `hcl:"csi_volume_provider"`
}

var (
	// default resources used for both the Server and its runners. Can be overridden
	// through config flags at install
	defaultResourcesCPU    = 200
	defaultResourcesMemory = 600
)

// Install is a method of NomadInstaller and implements the Installer interface to
// register a waypoint-server job with a Nomad cluster
func (i *NomadInstaller) Install(
	ctx context.Context,
	opts *InstallOpts,
) (*InstallResults, error) {
	ui := opts.UI

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing Nomad client...")
	defer func() { s.Abort() }()

	// Build api client from environment
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	s.Update("Checking for existing Waypoint server...")

	// Check if waypoint-server has already been deployed
	jobs, _, err := client.Jobs().PrefixList(serverName)
	if err != nil {
		return nil, err
	}
	var serverDetected bool
	for _, j := range jobs {
		if j.Name == serverName {
			if j.Status != "running" {
				return nil, fmt.Errorf("waypoint-server job found but not running")
			}
			serverDetected = true
			break
		}
	}

	var (
		clicfg   clicontext.Config
		addr     pb.ServerConfig_AdvertiseAddr
		httpAddr string
	)

	clicfg.Server = serverconfig.Client{
		Tls:           true,
		TlsSkipVerify: true,
		Platform:      "nomad",
	}

	addr.Tls = true
	addr.TlsSkipVerify = true

	if serverDetected {
		allocs, _, err := client.Jobs().Allocations(serverName, false, nil)
		if err != nil {
			return nil, err
		}

		var activeAllocs []*api.AllocationListStub
		for _, alloc := range allocs {
			if alloc.DesiredStatus == "run" {
				activeAllocs = append(activeAllocs, alloc)
			}
		}
		if len(allocs) == 0 || len(activeAllocs) == 0 {
			return nil, fmt.Errorf("waypoint-server job found but no running allocations available")
		}
		serverAddr, err := getAddrFromAllocID(allocs[0].ID, client)
		if err != nil {
			return nil, err
		}

		s.Update("Detected existing Waypoint server")
		s.Status(terminal.StatusWarn)
		s.Done()

		clicfg.Server.Address = serverAddr
		addr.Addr = serverAddr
		httpAddr = serverAddr
		return &InstallResults{
			Context:       &clicfg,
			AdvertiseAddr: &addr,
			HTTPAddr:      httpAddr,
		}, nil
	}

	if strings.ToLower(i.config.volumeType) == "csi" {
		if !c.config.csiVolumeProvider {
			return nil, fmt.Errorf("please include '-nomad-csi-volume-provider' flag")
		}

		s.Update("Creating persistent volume")

		vol := api.CSIVolume{
			ID:   "waypoint",
			Name: "waypoint",
			RequestedCapabilities: []*api.CSIVolumeCapability{
				{
					AccessMode:     "single-node-writer",
					AttachmentMode: "file-system",
				},
			},
			MountOptions: &api.CSIMountOptions{
				FSType:     "xfs",
				MountFlags: []string{"noatime"},
			},
			RequestedCapacityMin: 1073741824,
			RequestedCapacityMax: 2147483648,
			PluginID:             i.config.csiVolumeProvider,
		}

		_, _, err = client.CSIVolumes().Create(&vol, &api.WriteOptions{})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed creating Nomad persistent volume ID %s: %s", vol.ID, err)
		}
	} else if strings.ToLower(i.config.volumeType) == "host" && i.config.hostVolume == "" {
		return nil, fmt.Errorf("please include '-nomad-host-volume' flag")
	}

	s.Update("Installing Waypoint server to Nomad")
	allocID, err := i.runJob(ctx, s, client, waypointNomadJob(i.config, opts.ServerRunFlags))
	if err != nil {
		return nil, err
	}

	serverAddr, err := getAddrFromAllocID(allocID, client)
	if err != nil {
		return nil, err
	}
	hAddr, err := getHTTPFromAllocID(allocID, client)
	if err != nil {
		return nil, err
	}
	httpAddr = hAddr
	addr.Addr = serverAddr
	clicfg = clicontext.Config{
		Server: serverconfig.Client{
			Address:       addr.Addr,
			Tls:           true,
			TlsSkipVerify: true, // always for now
			Platform:      "nomad",
		},
	}

	s.Update("Waypoint server ready")
	s.Done()

	return &InstallResults{
		Context:       &clicfg,
		AdvertiseAddr: &addr,
		HTTPAddr:      httpAddr,
	}, nil
}

// Upgrade is a method of NomadInstaller and implements the Installer interface to
// upgrade a waypoint-server in a Nomad cluster
func (i *NomadInstaller) Upgrade(
	ctx context.Context, opts *InstallOpts, serverCfg serverconfig.Client) (
	*InstallResults, error,
) {
	ui := opts.UI

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing Nomad client...")
	defer func() { s.Abort() }()

	// Build api client from environment
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	s.Update("Checking for existing Waypoint server...")

	// Check if waypoint-server has already been deployed
	jobs, _, err := client.Jobs().PrefixList(serverName)
	if err != nil {
		return nil, err
	}

	var (
		serverDetected bool

		clicfg   clicontext.Config
		addr     pb.ServerConfig_AdvertiseAddr
		httpAddr string
	)

	for _, j := range jobs {
		if j.Name == serverName {
			serverDetected = true
			break
		}
	}

	clicfg.Server = serverconfig.Client{
		Tls:           true,
		TlsSkipVerify: true,
	}

	addr.Tls = true
	addr.TlsSkipVerify = true

	if !serverDetected {
		s.Update("No existing Waypoint server detected")
		s.Status(terminal.StatusError)
		s.Done()
		return nil, fmt.Errorf("No waypoint server job named %q detected in Nomad", serverName)
	} else {
		allocs, _, err := client.Jobs().Allocations(serverName, false, nil)
		if err != nil {
			return nil, err
		}
		if len(allocs) == 0 {
			return nil, fmt.Errorf("waypoint server job %q found but no running allocations available", serverName)
		}
		serverAddr, err := getAddrFromAllocID(allocs[0].ID, client)
		if err != nil {
			return nil, err
		}

		s.Update("Detected existing Waypoint server")
		s.Done()

		clicfg.Server.Address = serverAddr
		addr.Addr = serverAddr
		httpAddr = serverAddr
	}

	s = sg.Add("Upgrading Waypoint server on Nomad to %q", i.config.serverImage)
	job := waypointNomadJob(i.config, opts.ServerRunFlags)
	jobOpts := &api.RegisterOptions{
		PolicyOverride: i.config.policyOverride,
	}

	resp, _, err := client.Jobs().RegisterOpts(job, jobOpts, nil)
	if err != nil {
		return nil, err
	}

	s.Update("Waiting for allocation to be scheduled")
EVAL:
	qopts := &api.QueryOptions{
		WaitIndex: resp.EvalCreateIndex,
	}

	eval, meta, err := client.Evaluations().Info(resp.EvalID, qopts)
	if err != nil {
		return nil, err
	}
	qopts.WaitIndex = meta.LastIndex
	switch eval.Status {
	case "pending":
		goto EVAL
	case "complete":
		s.Update("Nomad allocation created")
	case "failed", "canceled", "blocked":
		s.Update("Nomad failed to schedule the waypoint server ", serverName)
		s.Status(terminal.StatusError)
		return nil, fmt.Errorf("nomad evaluation did not transition to 'complete'")
	default:
		return nil, fmt.Errorf("unknown eval status: %q", eval.Status)
	}

	var allocID string

	for {
		// We look for allocations by serverName here instead of the recent
		// evaluations ID from eval, because if the upgrade job is identical to what is
		// currently running, we won't get back a list of allocations, which will
		// fail the upgrade with no allocations running
		s.Update("Getting allocations for nomad server job")
		allocs, qmeta, err := client.Jobs().Allocations(serverName, false, qopts)
		s.Update("Got allocations for server install job")

		if err != nil {
			return nil, err
		}

		qopts.WaitIndex = qmeta.LastIndex
		if len(allocs) == 0 {
			return nil, fmt.Errorf("no allocations found after evaluation completed")
		}

		switch allocs[0].ClientStatus {
		case "running":
			allocID = allocs[0].ID
			s.Update("Nomad allocation running")
		case "pending":
			s.Update(fmt.Sprintf("Waiting for allocation %q to start", allocs[0].ID))
			// retry
		default:
			return nil, fmt.Errorf("allocation failed")

		}

		if allocID != "" {
			break
		}

		select {
		case <-time.After(500 * time.Millisecond):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	serverAddr, err := getAddrFromAllocID(allocID, client)
	if err != nil {
		return nil, err
	}
	hAddr, err := getHTTPFromAllocID(allocID, client)
	if err != nil {
		return nil, err
	}
	httpAddr = hAddr
	addr.Addr = serverAddr
	clicfg = clicontext.Config{
		Server: serverconfig.Client{
			Address:       addr.Addr,
			Tls:           true,
			TlsSkipVerify: true, // always for now
		},
	}

	s.Update("Upgrade of Waypoint server on Nomad complete!")
	s.Done()

	return &InstallResults{
		Context:       &clicfg,
		AdvertiseAddr: &addr,
		HTTPAddr:      httpAddr,
	}, nil
}

// Unnstall is a method of NomadInstaller and implements the Installer interface to
// stop and purge the waypoint-server job on a Nomad cluster
func (i *NomadInstaller) Uninstall(ctx context.Context, opts *InstallOpts) error {
	ui := opts.UI

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing Nomad client...")
	defer func() { s.Abort() }()

	// Build api client from environment
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	s.Update("Checking for existing Waypoint server...")

	// Find waypoint-server job
	jobs, _, err := client.Jobs().PrefixList(serverName)
	if err != nil {
		return err
	}
	var serverDetected bool
	for _, j := range jobs {
		if j.Name == serverName {
			serverDetected = true
			break
		}
	}
	if !serverDetected {
		return fmt.Errorf("No job with server name %q found; cannot uninstall", serverName)
	}

	s.Update("Removing Waypoint server from Nomad...")

	_, _, err = client.Jobs().Deregister(serverName, true, &api.WriteOptions{})
	if err != nil {
		ui.Output(
			"Error deregistering waypoint server job: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return err
	}
	allocs, _, err := client.Jobs().Allocations(serverName, true, nil)
	if err != nil {
		return err
	}
	for _, alloc := range allocs {
		if alloc.DesiredStatus != "stop" {
			a, _, err := client.Allocations().Info(alloc.ID, &api.QueryOptions{})
			if err != nil {
				return err
			}
			_, err = client.Allocations().Stop(a, &api.QueryOptions{})
			if err != nil {
				return err
			}
		}
	}

	s.Update("Waypoint job and allocations purged")
	s.Done()

	return nil
}

// InstallRunner implements Installer.
func (i *NomadInstaller) InstallRunner(
	ctx context.Context,
	opts *InstallRunnerOpts,
) error {
	ui := opts.UI

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing Nomad client...")
	defer func() { s.Abort() }()

	// Build api client from environment
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	// Install the runner
	s.Update("Installing the Waypoint runner")
	_, err = i.runJob(ctx, s, client, waypointRunnerNomadJob(i.config, opts))
	if err != nil {
		return err
	}
	s.Update("Waypoint runner installed")
	s.Done()

	return nil
}

// UninstallRunner implements Installer.
func (i *NomadInstaller) UninstallRunner(
	ctx context.Context,
	opts *InstallOpts,
) error {
	ui := opts.UI

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing Nomad client...")
	defer func() { s.Abort() }()

	// Build api client from environment
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	s.Update("Checking for existing Waypoint runner...")
	jobs, _, err := client.Jobs().PrefixList(runnerName)
	if err != nil {
		return err
	}
	var detected bool
	for _, j := range jobs {
		if j.Name == runnerName {
			detected = true
			break
		}
	}
	if !detected {
		s.Update("No Waypoint runner detected.")
		s.Done()
		return nil
	}

	s.Update("Removing Waypoint runner...")
	_, _, err = client.Jobs().Deregister(runnerName, true, &api.WriteOptions{})
	if err != nil {
		ui.Output(
			"Error deregistering Waypoint runner job: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return err
	}

	allocs, _, err := client.Jobs().Allocations(serverName, true, nil)
	if err != nil {
		return err
	}
	for _, alloc := range allocs {
		if alloc.DesiredStatus != "stop" {
			a, _, err := client.Allocations().Info(alloc.ID, &api.QueryOptions{})
			if err != nil {
				return err
			}
			_, err = client.Allocations().Stop(a, &api.QueryOptions{})
			if err != nil {
				return err
			}
		}
	}

	s.Update("Waypoint runner job and allocations purged")
	s.Done()

	return nil
}

// HasRunner implements Installer.
func (i *NomadInstaller) HasRunner(
	ctx context.Context,
	opts *InstallOpts,
) (bool, error) {
	// Build api client from environment
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return false, err
	}

	jobs, _, err := client.Jobs().PrefixList(runnerName)
	if err != nil {
		return false, err
	}
	for _, j := range jobs {
		if j.Name == runnerName {
			return true, nil
		}
	}

	return false, nil
}

func (i *NomadInstaller) runJob(
	ctx context.Context,
	s terminal.Step,
	client *api.Client,
	job *api.Job,
) (string, error) {
	jobOpts := &api.RegisterOptions{
		PolicyOverride: i.config.policyOverride,
	}

	resp, _, err := client.Jobs().RegisterOpts(job, jobOpts, nil)
	if err != nil {
		return "", err
	}

	s.Update("Waiting for allocation to be scheduled")
EVAL:
	qopts := &api.QueryOptions{
		WaitIndex: resp.EvalCreateIndex,
	}

	eval, meta, err := client.Evaluations().Info(resp.EvalID, qopts)
	if err != nil {
		return "", err
	}
	qopts.WaitIndex = meta.LastIndex
	switch eval.Status {
	case "pending":
		goto EVAL
	case "complete":
		s.Update("Nomad allocation created")
	case "failed", "canceled", "blocked":
		s.Update("Nomad failed to schedule the job")
		s.Status(terminal.StatusError)
		return "", fmt.Errorf("nomad evaluation did not transition to 'complete'")
	default:
		return "", fmt.Errorf("unknown eval status: %q", eval.Status)
	}

	var allocID string
	for {
		allocs, qmeta, err := client.Evaluations().Allocations(eval.ID, qopts)
		if err != nil {
			return "", err
		}
		qopts.WaitIndex = qmeta.LastIndex
		if len(allocs) == 0 {
			return "", fmt.Errorf("no allocations found after evaluation completed")
		}

		switch allocs[0].ClientStatus {
		case "running":
			allocID = allocs[0].ID
			s.Update("Nomad allocation running")
		case "pending":
			s.Update(fmt.Sprintf("Waiting for allocation %q to start", allocs[0].ID))
			// retry
		default:
			return "", fmt.Errorf("allocation failed")
		}

		if allocID != "" {
			return allocID, nil
		}

		select {
		case <-time.After(500 * time.Millisecond):
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

// waypointNomadJob takes in a nomadConfig and returns a Nomad Job per the
// Nomad API
func waypointNomadJob(c nomadConfig, rawRunFlags []string) *api.Job {
	job := api.NewServiceJob(serverName, serverName, c.region, 50)
	job.Namespace = &c.namespace
	job.Datacenters = c.datacenters
	job.Meta = c.serviceAnnotations
	tg := api.NewTaskGroup(serverName, 1)

	grpcPort, _ := strconv.Atoi(defaultGrpcPort)
	httpPort, _ := strconv.Atoi(defaultHttpPort)

	// Include services to be registered in Consul, if specified in the server install
	// One service added for Waypoint UI, and one for Waypoint backend port
	if c.consulService {
		tg.Services = []*api.Service{
			{
				Name:      "waypoint-ui",
				PortLabel: "ui",
				Tags:      c.consulServiceUITags,
			},
			{
				Name:      "waypoint-server",
				PortLabel: "server",
				Tags:      c.consulServiceBackendTags,
			},
		}
	}

	tg.Networks = []*api.NetworkResource{
		{
			Mode: "host",
			DynamicPorts: []api.Port{
				{
					Label: "server",
					To:    grpcPort,
				},
			},
			// currently set to static; when ui command can be dynamic - update this
			ReservedPorts: []api.Port{
				{
					Label: "ui",
					Value: httpPort,
					To:    httpPort,
				},
			},
		},
	}

	// Preserve disk, otherwise upgrades will destroy previous allocation and the disk along with it
	volumeRequest := api.VolumeRequest{
		Type:     c.volumeType,
		ReadOnly: false,
	}

	if strings.ToLower(c.volumeType) == "csi" {
		volumeRequest.Source = "waypoint"
		volumeRequest.AccessMode = "single-node-writer"
		volumeRequest.AttachmentMode = "file-system"
	} else {
		volumeRequest.Source = c.hostVolume
	}

	tg.Volumes = map[string]*api.VolumeRequest{
		"waypoint-server": &volumeRequest,
	}

	job.AddTaskGroup(tg)

	ras := []string{"server", "run", "-accept-tos", "-vv", "-db=/alloc/data/data.db", fmt.Sprintf("-listen-grpc=0.0.0.0:%s", defaultGrpcPort), fmt.Sprintf("-listen-http=0.0.0.0:%s", defaultHttpPort)}
	ras = append(ras, rawRunFlags...)
	task := api.NewTask("server", "docker")
	task.Config = map[string]interface{}{
		"image":          c.serverImage,
		"ports":          []string{"server", "ui"},
		"args":           ras
		"auth_soft_fail": c.authSoftFail,
	}
	task.Env = map[string]string{
		"PORT": defaultGrpcPort,
	}
	
	readOnly := false
	volume := "waypoint-server"
	destination := "/data"
	volumeMounts := []*api.VolumeMount{
		{
			Volume:      &volume,
			Destination: &destination,
			ReadOnly:    &readOnly,
		},
	}

	cpu := defaultResourcesCPU
	mem := defaultResourcesMemory

	preTask := api.NewTask("pre_task", "docker")
	// Observed WP user and group IDs in the published container, update if those ever change
	waypointUserID := 100
	waypointGroupID := 1000
	preTask.Config = map[string]interface{}{
		// TODO(xx): pin busybox image
		"image":   "busybox:latest",
		"command": "sh",
		"args":    []string{"-c", fmt.Sprintf("chown -R %d:%d /data/", waypointUserID, waypointGroupID)},
	}
	preTask.VolumeMounts = volumeMounts
	preTask.Resources = &api.Resources{
		// TODO(xx): make pointer vars for smaller cpu and mem
		CPU:      &cpu,
		MemoryMB: &mem,
	}
	preTask.Lifecycle = &api.TaskLifecycle{
		Hook:    "prestart",
		Sidecar: false,
	}
	tg.AddTask(preTask)

	task := api.NewTask("server", "docker")
	task.Config = map[string]interface{}{
		"image":          c.serverImage,
		"ports":          []string{"server", "ui"},
		"args":           []string{"server", "run", "-accept-tos", "-vvv", "-db=/data/data.db", fmt.Sprintf("-listen-grpc=0.0.0.0:%s", defaultGrpcPort), fmt.Sprintf("-listen-http=0.0.0.0:%s", defaultHttpPort)},
		"auth_soft_fail": c.authSoftFail,
	}
	task.Env = map[string]string{
		"PORT": defaultGrpcPort,
	}
	task.VolumeMounts = volumeMounts

	if c.serverResourcesCPU != "" {
		cpu, _ = strconv.Atoi(c.serverResourcesCPU)
	}

	if c.serverResourcesMemory != "" {
		mem, _ = strconv.Atoi(c.serverResourcesMemory)
	}
	task.Resources = &api.Resources{
		CPU:      &cpu,
		MemoryMB: &mem,
	}
	tg.AddTask(task)

	return job
}

// waypointRunnerNomadJob takes in a nomadConfig and returns a Nomad Job
// for the Nomad API to run a Waypoint runner.
func waypointRunnerNomadJob(c nomadConfig, opts *InstallRunnerOpts) *api.Job {
	job := api.NewServiceJob(runnerName, runnerName, c.region, 50)
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

	task := api.NewTask("runner", "docker")
	task.Config = map[string]interface{}{
		"image": c.serverImage,
		"args": []string{
			"runner",
			"agent",
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
	tg.AddTask(task)

	return job
}

// getAddrFromAllocID takes in an allocID and a Nomad Client and returns
// the address for the server
func getAddrFromAllocID(allocID string, client *api.Client) (string, error) {
	alloc, _, err := client.Allocations().Info(allocID, nil)
	if err != nil {
		return "", err
	}

	for _, port := range alloc.AllocatedResources.Shared.Ports {
		if port.Label == "server" {
			return fmt.Sprintf("%s:%d", port.HostIP, port.Value), nil
		}
	}

	return "", nil
}

// getHTTPFromAllocID takes in an allocID and a Nomad Client and returns
// the http address
func getHTTPFromAllocID(allocID string, client *api.Client) (string, error) {
	alloc, _, err := client.Allocations().Info(allocID, nil)
	if err != nil {
		return "", err
	}

	for _, port := range alloc.AllocatedResources.Shared.Ports {
		if port.Label == "ui" {
			return fmt.Sprintf(port.HostIP + ":9702"), nil
		}
	}

	return "", nil
}

func (i *NomadInstaller) InstallFlags(set *flag.Set) {
	set.StringMapVar(&flag.StringMapVar{
		Name:   "nomad-annotate-service",
		Target: &i.config.serviceAnnotations,
		Usage:  "Annotations for the Service generated.",
	})

	set.BoolVar(&flag.BoolVar{
		Name:    "nomad-auth-soft-fail",
		Target:  &i.config.authSoftFail,
		Default: false,
		Usage: "Don't fail the Nomad task on an auth failure obtaining server " +
			"image container. Attempt to continue without auth.",
	})

	set.StringSliceVar(&flag.StringSliceVar{
		Name:    "nomad-dc",
		Target:  &i.config.datacenters,
		Default: []string{"dc1"},
		Usage:   "Datacenters to install to for Nomad.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-namespace",
		Target:  &i.config.namespace,
		Default: "default",
		Usage:   "Namespace to install the Waypoint server into for Nomad.",
	})

	set.BoolVar(&flag.BoolVar{
		Name:    "nomad-policy-override",
		Target:  &i.config.policyOverride,
		Default: false,
		Usage:   "Override the Nomad sentinel policy for enterprise Nomad.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-region",
		Target:  &i.config.region,
		Default: "global",
		Usage:   "Region to install to for Nomad.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-server-cpu",
		Target:  &i.config.serverResourcesCPU,
		Usage:   "CPU required to run this task in MHz.",
		Default: strconv.Itoa(defaultResourcesCPU),
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-server-memory",
		Target:  &i.config.serverResourcesMemory,
		Usage:   "MB of Memory to allocate to the Server job task.",
		Default: strconv.Itoa(defaultResourcesMemory),
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-runner-cpu",
		Target:  &i.config.runnerResourcesCPU,
		Usage:   "CPU required to run this task in MHz.",
		Default: strconv.Itoa(defaultResourcesCPU),
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-runner-memory",
		Target:  &i.config.runnerResourcesMemory,
		Usage:   "MB of Memory to allocate to the runner job task.",
		Default: strconv.Itoa(defaultResourcesMemory),
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: defaultServerImage,
	})

	set.BoolVar(&flag.BoolVar{
		Name:    "nomad-consul-service",
		Target:  &i.config.consulService,
		Usage:   "Create service for Waypoint UI in Consul.",
		Default: false,
	})

	set.StringSliceVar(&flag.StringSliceVar{
		Name:   "nomad-consul-service-ui-tags",
		Target: &i.config.consulServiceUITags,
		Usage:  "Tags for the Waypoint UI service generated in Consul.",
	})

	set.StringSliceVar(&flag.StringSliceVar{
		Name:   "nomad-consul-service-backend-tags",
		Target: &i.config.consulServiceBackendTags,
		Usage:  "Tags for the Waypoint backend service generated in Consul.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-volume-type",
		Target: &i.config.volumeType,
		Usage:  "Nomad volume type for the Waypoint server.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-host-volume",
		Target: &i.config.hostVolume,
		Usage:  "Nomad host volume name.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-csi-volume-provider",
		Target: &i.config.csiVolumeProvider,
		Usage:  "Nomad CSI volume provider.",
	})
}

func (i *NomadInstaller) UpgradeFlags(set *flag.Set) {
	set.StringMapVar(&flag.StringMapVar{
		Name:   "nomad-annotate-service",
		Target: &i.config.serviceAnnotations,
		Usage:  "Annotations for the Service generated.",
	})

	set.BoolVar(&flag.BoolVar{
		Name:    "nomad-auth-soft-fail",
		Target:  &i.config.authSoftFail,
		Default: false,
		Usage: "Don't fail the Nomad task on an auth failure obtaining server " +
			"image container. Attempt to continue without auth.",
	})

	set.StringSliceVar(&flag.StringSliceVar{
		Name:    "nomad-dc",
		Target:  &i.config.datacenters,
		Default: []string{"dc1"},
		Usage:   "Datacenters to install to for Nomad.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-namespace",
		Target:  &i.config.namespace,
		Default: "default",
		Usage:   "Namespace to install the Waypoint server into for Nomad.",
	})

	set.BoolVar(&flag.BoolVar{
		Name:    "nomad-policy-override",
		Target:  &i.config.policyOverride,
		Default: false,
		Usage:   "Override the Nomad sentinel policy for enterprise Nomad.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-region",
		Target:  &i.config.region,
		Default: "global",
		Usage:   "Region to install to for Nomad.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-server-cpu",
		Target:  &i.config.serverResourcesCPU,
		Usage:   "CPU required to run this task in MHz.",
		Default: strconv.Itoa(defaultResourcesCPU),
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-server-memory",
		Target:  &i.config.serverResourcesMemory,
		Usage:   "MB of Memory to allocate to the server job task.",
		Default: strconv.Itoa(defaultResourcesMemory),
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-runner-cpu",
		Target:  &i.config.runnerResourcesCPU,
		Usage:   "CPU required to run this task in MHz.",
		Default: strconv.Itoa(defaultResourcesCPU),
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-runner-memory",
		Target:  &i.config.runnerResourcesMemory,
		Usage:   "MB of Memory to allocate to the runner job task.",
		Default: strconv.Itoa(defaultResourcesMemory),
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: defaultServerImage,
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-volume-type",
		Target: &i.config.volumeType,
		Usage:  "Nomad volume type for the Waypoint server.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-host-volume",
		Target: &i.config.hostVolume,
		Usage:  "Nomad host volume name.",
	})
}

func (i *NomadInstaller) UninstallFlags(set *flag.Set) {
	// Purposely empty, no flags
}
