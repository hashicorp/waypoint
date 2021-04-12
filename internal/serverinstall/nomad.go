package serverinstall

import (
	"context"
	"fmt"
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
	authSoftFail       bool              `hcl:"auth_soft_fail,optional"`
	serverImage        string            `hcl:"server_image,optional"`
	namespace          string            `hcl:"namespace,optional"`
	serviceAnnotations map[string]string `hcl:"service_annotations,optional"`

	region         string   `hcl:"namespace,optional"`
	datacenters    []string `hcl:"datacenters,optional"`
	policyOverride bool     `hcl:"policy_override,optional"`

	serverPurge bool `hcl:"serverPurge,optional"`
}

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

	s.Update("Installing Waypoint server to Nomad")
	allocID, err := i.runJob(ctx, s, client, waypointNomadJob(i.config))
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
	job := waypointNomadJob(i.config)
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
		allocs, qmeta, err := client.Jobs().Allocations(serverName, false, qopts)
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
// stop, and optionally purge, the waypoint-server job on a Nomad cluster
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

	_, _, err = client.Jobs().Deregister(serverName, i.config.serverPurge, &api.WriteOptions{})
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

	if i.config.serverPurge {
		s.Update("Waypoint job and allocations purged")
	} else {
		s.Update("Waypoint job and allocations stopped")
	}
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
	_, _, err = client.Jobs().Deregister(runnerName, i.config.serverPurge, &api.WriteOptions{})
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

	if i.config.serverPurge {
		s.Update("Waypoint runner job and allocations purged")
	} else {
		s.Update("Waypoint runner job and allocations stopped")
	}
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
func waypointNomadJob(c nomadConfig) *api.Job {
	job := api.NewServiceJob(serverName, serverName, c.region, 50)
	job.Namespace = &c.namespace
	job.Datacenters = c.datacenters
	job.Meta = c.serviceAnnotations
	tg := api.NewTaskGroup(serverName, 1)

	grpcPort, _ := strconv.Atoi(defaultGrpcPort)
	httpPort, _ := strconv.Atoi(defaultHttpPort)

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
	// Preserve disk, otherwise upgrades will destroy previous allocation and the
	// disk along with it
	tg.EphemeralDisk = &api.EphemeralDisk{
		Sticky:  &[]bool{true}[0],
		Migrate: &[]bool{true}[0],
	}
	job.AddTaskGroup(tg)

	task := api.NewTask("server", "docker")
	task.Config = map[string]interface{}{
		"image":          c.serverImage,
		"ports":          []string{"server", "ui"},
		"args":           []string{"server", "run", "-accept-tos", "-vvv", "-db=/alloc/data/data.db", fmt.Sprintf("-listen-grpc=0.0.0.0:%s", defaultGrpcPort), fmt.Sprintf("-listen-http=0.0.0.0:%s", defaultHttpPort)},
		"auth_soft_fail": c.authSoftFail,
	}
	task.Env = map[string]string{
		"PORT": defaultGrpcPort,
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
			"-vvv",
		},
		"auth_soft_fail": c.authSoftFail,
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
		Name:    "nomad-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: defaultServerImage,
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
		Name:    "nomad-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: defaultServerImage,
	})
}

func (i *NomadInstaller) UninstallFlags(set *flag.Set) {
	set.BoolVar(&flag.BoolVar{
		Name:   "nomad-purge",
		Target: &i.config.serverPurge,
		Usage:  "Purge the Waypoint server job after deregistering as part of uninstall.",
	})
}
