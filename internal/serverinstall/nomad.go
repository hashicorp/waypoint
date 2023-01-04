package serverinstall

import (
	"context"
	json "encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/waypoint/internal/installutil"

	"github.com/hashicorp/waypoint/internal/installutil/nomad"
	"github.com/hashicorp/waypoint/internal/runnerinstall"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/nomad/api"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
)

const runnerJobName string = "waypoint-static-runner"

type NomadInstaller struct {
	config nomadConfig
}

type nomadConfig struct {
	authSoftFail       bool              `hcl:"auth_soft_fail,optional"`
	serverImage        string            `hcl:"server_image,optional"`
	namespace          string            `hcl:"namespace,optional"`
	serviceAnnotations map[string]string `hcl:"service_annotations,optional"`

	serviceProvider          string   `hcl:"service_provider,optional"`
	serviceUITags            []string `hcl:"service_ui_tags:optional"`
	serviceBackendTags       []string `hcl:"service_backend_tags:optional"`
	consulService            bool     `hcl:"consul_service,optional"`
	consulServiceUITags      []string `hcl:"consul_service_ui_tags:optional"`
	consulServiceBackendTags []string `hcl:"consul_service_backend_tags:optional"`
	consulDatacenter         string   `hcl:"consul_datacenter,optional"`
	consulDomain             string   `hcl:"consul_datacenter,optional"`
	consulToken              string   `hcl:"consul_token,optional"`

	// If set along with consul, will use this hostname instead of
	// making a consul DNS hostname for the server address in its context
	consulServiceHostname string `hcl:"consul_service_hostname,optional"`

	odrImage string `hcl:"odr_image,optional"`

	region         string   `hcl:"namespace,optional"`
	datacenters    []string `hcl:"datacenters,optional"`
	policyOverride bool     `hcl:"policy_override,optional"`

	serverResourcesCPU    string `hcl:"server_resources_cpu,optional"`
	serverResourcesMemory string `hcl:"server_resources_memory,optional"`

	hostVolume           string            `hcl:"host_volume,optional"`
	csiVolumeProvider    string            `hcl:"csi_volume_provider,optional"`
	csiVolumeCapacityMin int64             `hcl:"csi_volume_capacity_min,optional"`
	csiVolumeCapacityMax int64             `hcl:"csi_volume_capacity_max,optional"`
	csiFS                string            `hcl:"csi_fs,optional"`
	csiPluginId          string            `hcl:"csi_plugin_id,optional"`
	csiExternalId        string            `hcl:"nomad_csi_external_id,optional"`
	csiTopologies        map[string]string `hcl:"nomad_csi_topologies,optional"`
	csiSecrets           map[string]string `hcl:"nomad_csi_secrets,optional"`
	csiParams            map[string]string `hcl:"csi_parameters,optional"`
	csiVolume            string            `hcl:"nomad_csi_volume,optional"`
	nomadHost            string            `hcl:"nomad_host,optional"`

	runnerResourcesCPU         string `hcl:"runner_resources_cpu,optional"`
	runnerResourcesMemory      string `hcl:"runner_resources_memory,optional"`
	runnerHostVolume           string `hcl:"runner_host_volume,optional"`
	runnerCsiVolumeProvider    string `hcl:"runner_csi_volume_provider,optional"`
	runnerCsiVolumeCapacityMin int64  `hcl:"runner_csi_volume_capacity_min,optional"`
	runnerCsiVolumeCapacityMax int64  `hcl:"runner_csi_volume_capacity_max,optional"`
	runnerCsiVolume            string `hcl:"runner_csi_volume,optional"`
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
	defaultServiceTag             = "waypoint"
	defaultConsulDatacenter       = "dc1"
	defaultConsulDomain           = "consul"
	waypointBackendServiceName    = "waypoint-server"
	waypointUIServiceName         = "waypoint-ui"
	defaultWaypointConsulHostname = fmt.Sprintf("%s.%s.service.%s.%s",
		defaultServiceTag, waypointBackendServiceName, defaultConsulDatacenter, defaultConsulDomain)

	defaultNomadHost = "http://localhost:4646"
)

// Install is a method of NomadInstaller and implements the Installer interface to
// register a waypoint-server job with a Nomad cluster
func (i *NomadInstaller) Install(
	ctx context.Context,
	opts *InstallOpts,
) (*InstallResults, string, error) {
	ui := opts.UI

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing Nomad client...")
	defer func() { s.Abort() }()

	// Build api client from environment
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, "", err
	}

	s.Update("Checking for existing Waypoint server...")

	// Check if waypoint-server has already been deployed
	jobs, _, err := client.Jobs().PrefixList(serverName)
	if err != nil {
		return nil, "", err
	}

	var serverDetected bool
	for _, j := range jobs {
		if j.Name == serverName {
			if j.Status != "running" {
				return nil, "", fmt.Errorf("waypoint-server job found but not running")
			}
			serverDetected = true
			break
		}
	}

	if i.config.odrImage == "" {
		var err error
		i.config.odrImage, err = installutil.DeriveDefaultODRImage(i.config.serverImage)
		if err != nil {
			return nil, "", err
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
			return nil, "", err
		}

		var activeAllocs []*api.AllocationListStub
		for _, alloc := range allocs {
			if alloc.DesiredStatus == "run" {
				activeAllocs = append(activeAllocs, alloc)
			}
		}
		if len(allocs) == 0 || len(activeAllocs) == 0 {
			return nil, "", fmt.Errorf("waypoint-server job found but no running allocations available")
		}
		serverAddr, err := getAddrFromAllocID(allocs[0].ID, client)
		if err != nil {
			return nil, "", err
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
		}, "", nil
	}

	if i.config.csiVolumeProvider == "" && i.config.hostVolume == "" {
		return nil, "", fmt.Errorf("please include '-nomad-csi-volume-provider' or '-nomad-host-volume'")
	} else if i.config.csiVolumeProvider != "" {
		if i.config.hostVolume != "" {
			return nil, "", fmt.Errorf("choose either CSI or host volume, not both")
		}

		s.Update("Creating persistent volume")
		err = nomad.CreatePersistentVolume(
			ctx,
			client,
			"waypoint-server",
			i.config.csiVolume,
			i.config.csiPluginId,
			i.config.csiVolumeProvider,
			i.config.csiFS,
			i.config.csiExternalId,
			i.config.csiVolumeCapacityMin,
			i.config.csiVolumeCapacityMax,
			i.config.csiTopologies,
			i.config.csiSecrets,
			i.config.csiParams,
		)
		if err != nil {
			return nil, "", status.Errorf(codes.Internal, "Failed creating Nomad persistent volume: %s", err)
		}
		s.Update("Persistent volume created!")
		s.Status(terminal.StatusOK)
		s.Done()
	}

	s.Update("Installing Waypoint server to Nomad")
	allocID, err := nomad.RunJob(ctx, s, client, waypointNomadJob(i.config, opts.ServerRunFlags, false), i.config.policyOverride)
	if err != nil {
		return nil, "", err
	}

	// If a Consul service was requested, set the consul DNS hostname rather
	// than the direct static IP for the CLI context and server config. Otherwise
	// if Nomad restarts the server allocation, a new IP will be assigned and any
	// configured clients will be invalid
	httpAddr, addr.Addr, err = i.getWaypointAddress(client, allocID)
	if err != nil {
		return nil, "", err
	}

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

	if i.config.serviceProvider == "consul" && i.config.consulService {
		opts.UI.Output("The CLI has been configured to automatically install a Consul service for\n" +
			"the Waypoint service backend and ui service in Nomad.")
	} else if i.config.serviceProvider == "nomad" {
		opts.UI.Output("The CLI has been configured to automatically install a Nomad service for\n" +
			"the Waypoint service backend and ui service in Nomad.")
	} else {
		opts.UI.Output("Waypoint server running on Nomad is being accessed via its allocation IP and port.\n" +
			"This could change in the future if Nomad creates a new allocation for the Waypoint server,\n" +
			"which would break all existing Waypoint contexts.\n\n" +
			"It is recommended to use Consul for determining Waypoint servers IP running on Nomad rather than\n" +
			"relying on the static IP that is initially set up for this allocation.")
	}

	return &InstallResults{
		Context:       &clicfg,
		AdvertiseAddr: &addr,
		HTTPAddr:      httpAddr,
	}, "", nil
}

func (i *NomadInstaller) getWaypointAddress(client *api.Client, allocID string) (string, string, error) {
	if i.config.serviceProvider == "consul" && i.config.consulService {
		if i.config.consulDatacenter == "" {
			i.config.consulDatacenter = defaultConsulDatacenter
		}
		if i.config.consulDomain == "" {
			i.config.consulDomain = defaultConsulDomain
		}

		grpcPort, _ := strconv.Atoi(defaultGrpcPort)
		httpPort, _ := strconv.Atoi(defaultHttpPort)

		if i.config.consulServiceHostname == "" {
			return fmt.Sprintf("%s.service.%s.%s:%d",
					waypointUIServiceName, i.config.consulDatacenter, i.config.consulDomain, httpPort), fmt.Sprintf("%s.service.%s.%s:%d",
					waypointBackendServiceName, i.config.consulDatacenter, i.config.consulDomain, grpcPort), nil
		} else {
			return fmt.Sprintf("%s:%d", i.config.consulServiceHostname, httpPort),
				fmt.Sprintf("%s:%d", i.config.consulServiceHostname, grpcPort), nil
		}
	} else if i.config.serviceProvider == "nomad" {
		backendService, _, err := client.Services().Get(waypointBackendServiceName, nil)
		if err != nil {
			return "", "", err
		}

		uiService, _, err := client.Services().Get(waypointUIServiceName, nil)
		if err != nil {
			return "", "", err
		}

		return fmt.Sprintf("%s:%d", uiService[0].Address, uiService[0].Port),
			fmt.Sprintf("%s:%d", backendService[0].Address, backendService[0].Port), nil
	} else {
		serverAddr, err := getAddrFromAllocID(allocID, client)
		if err != nil {
			return "", "", err
		}
		hAddr, err := getHTTPFromAllocID(allocID, client)
		if err != nil {
			return "", "", err
		}
		return hAddr, serverAddr, nil
	}
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

	if i.config.serviceProvider == "none" && !i.config.consulService {
		// By default, we don't auto-enable the consul service because prior to Waypoint
		// version 0.6.2, we did not enable it by default.
		proceed, err := opts.UI.Input(&terminal.Input{
			Prompt: "Service discovery is disabled for the Waypoint Nomad job. If you had previously enabled " +
				"it in the last installation, please stop this upgrade and re-run with the -nomad-service-provider flag. " +
				"Otherwise, enter 'yes' to continue the upgrade: ",
			Style:  "",
			Secret: false,
		})
		if err != nil {
			opts.UI.Output(
				"Error upgrading server: %s",
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return nil, err
		} else if strings.ToLower(proceed) != "yes" {
			return nil, errors.New("upgrade aborted")
		}
	}

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
		clicfg         clicontext.Config
		addr           pb.ServerConfig_AdvertiseAddr
		httpAddr       string
	)

	for _, j := range jobs {
		if j.Name == serverName {
			serverDetected = true
			break
		}
	}

	if i.config.odrImage == "" {
		var err error
		i.config.odrImage, err = installutil.DeriveDefaultODRImage(i.config.serverImage)
		if err != nil {
			return nil, err
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
	job := waypointNomadJob(i.config, opts.ServerRunFlags, true)
	jobOpts := &api.RegisterOptions{
		PolicyOverride: i.config.policyOverride,
	}

	resp, _, err := client.Jobs().RegisterOpts(job, jobOpts, nil)
	if err != nil {
		return nil, err
	}

	s.Update("Waiting for allocation to be scheduled")
	qopts := &api.QueryOptions{
		WaitIndex: resp.EvalCreateIndex,
		WaitTime:  time.Duration(500 * time.Millisecond),
	}

	eval, meta, err := i.waitForEvaluation(ctx, s, client, resp, qopts)
	if err != nil {
		s.Update("Nomad failed to schedule the waypoint server ", serverName)
		s.Status(terminal.StatusError)
		return nil, err
	}
	if eval == nil {
		return nil, fmt.Errorf("evaluation status could not be determined")
	}
	qopts.WaitIndex = meta.LastIndex

	var allocID string

	for {
		// We look for allocations by serverName here instead of the recent
		// evaluations ID from eval, because if the upgrade job is identical to what is
		// currently running, we won't get back a list of allocations, which will
		// fail the upgrade with no allocations running
		s.Update("Getting allocations for nomad server job, this may take a while...")
		allocs, qmeta, err := client.Jobs().Allocations(serverName, false, qopts)

		if err != nil {
			return nil, err
		}

		qopts.WaitIndex = qmeta.LastIndex
		if len(allocs) == 0 {
			return nil, fmt.Errorf("no allocations found after evaluation completed")
		}

		s.Update("Got allocations for server install job")

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

	// If a Consul service was requested, set the consul DNS hostname rather
	// than the direct static IP for the CLI context and server config. Otherwise
	// if Nomad restarts the server allocation, a new IP will be assigned and any
	// configured clients will be invalid
	httpAddr, addr.Addr, err = i.getWaypointAddress(client, allocID)
	if err != nil {
		return nil, err
	}

	clicfg = clicontext.Config{
		Server: serverconfig.Client{
			Address:       addr.Addr,
			Tls:           true,
			TlsSkipVerify: true, // always for now
			Platform:      "nomad",
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

	vols, _, err := client.CSIVolumes().List(&api.QueryOptions{Prefix: "waypoint"})
	if err != nil {
		return err
	}
	for _, vol := range vols {
		if vol.ID == "waypoint-server" {
			s.Update("Destroying persistent CSI volume")
			err = client.CSIVolumes().Deregister(vol.ID, false, &api.WriteOptions{})
			if err != nil {
				return err
			}
			s.Update("Successfully destroyed persistent volumes")
			break
		}
	}
	s.Done()

	return nil
}

// InstallRunner implements Installer.
func (i *NomadInstaller) InstallRunner(
	ctx context.Context,
	opts *runnerinstall.InstallOpts,
) error {
	runnerInstaller := runnerinstall.NomadRunnerInstaller{
		Config: runnerinstall.NomadConfig{
			AuthSoftFail:          i.config.authSoftFail,
			RunnerImage:           i.config.serverImage,
			Namespace:             i.config.namespace,
			ServiceAnnotations:    i.config.serviceAnnotations,
			Region:                i.config.region,
			Datacenters:           i.config.datacenters,
			PolicyOverride:        i.config.policyOverride,
			RunnerResourcesCPU:    i.config.runnerResourcesCPU,
			RunnerResourcesMemory: i.config.runnerResourcesMemory,
			HostVolume:            i.config.runnerHostVolume,
			CsiVolumeProvider:     i.config.runnerCsiVolumeProvider,
			CsiVolumeCapacityMin:  i.config.runnerCsiVolumeCapacityMin,
			CsiVolumeCapacityMax:  i.config.runnerCsiVolumeCapacityMax,
			CsiFS:                 i.config.csiFS,
			CsiTopologies:         i.config.csiTopologies,
			CsiExternalId:         i.config.csiExternalId,
			CsiParams:             i.config.csiParams,
			CsiPluginId:           i.config.csiPluginId,
			CsiSecrets:            i.config.csiSecrets,
			CsiVolume:             i.config.runnerCsiVolume,
			NomadHost:             i.config.nomadHost,
		},
	}
	err := runnerInstaller.Install(ctx, opts)
	if err != nil {
		return err
	}
	return nil
}

// UninstallRunner implements Installer.
func (i *NomadInstaller) UninstallRunner(
	ctx context.Context,
	opts *runnerinstall.InstallOpts,
) error {

	runnerInstaller := runnerinstall.NomadRunnerInstaller{
		Config: runnerinstall.NomadConfig{
			Namespace: i.config.namespace,
			Region:    i.config.region,
		},
	}

	err := runnerInstaller.Uninstall(ctx, opts)
	if err != nil {
		return err
	}
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
	var jobs []*api.JobListStub

	// Check for runner with job name pre-0.9
	jobsWithOldRunnerName, _, err := client.Jobs().PrefixList(runnerName)
	if err != nil {
		return false, err
	}
	jobs = append(jobs, jobsWithOldRunnerName...)

	// Check for runner with job name post-0.9
	jobWithNewRunnerName, _, err := client.Jobs().PrefixList(runnerJobName)
	if err != nil {
		return false, err
	}
	jobs = append(jobs, jobWithNewRunnerName...)
	for _, j := range jobs {
		if j.Name == runnerJobName || j.Name == runnerName {
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
	qopts := &api.QueryOptions{
		WaitIndex: resp.EvalCreateIndex,
		WaitTime:  time.Duration(500 * time.Millisecond),
	}

	eval, meta, err := i.waitForEvaluation(ctx, s, client, resp, qopts)
	if err != nil {
		return "", err
	}
	if eval == nil {
		return "", fmt.Errorf("evaluation status could not be determined")
	}
	qopts.WaitIndex = meta.LastIndex

	var allocID string
	retries := 0
	maxRetries := 3
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
			retries++
		case "pending":
			s.Update(fmt.Sprintf("Waiting for allocation %q to start", allocs[0].ID))
			// retry
		default:
			return "", fmt.Errorf("allocation failed")
		}

		if allocID != "" {
			if retries == maxRetries {
				return allocID, nil
			} else {
				s.Update("Ensuring allocation %q has properly started up...", allocs[0].ID)
				time.Sleep(1 * time.Second)
			}
		}

		select {
		case <-time.After(500 * time.Millisecond):
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

func (i *NomadInstaller) waitForEvaluation(
	ctx context.Context,
	s terminal.Step,
	client *api.Client,
	resp *api.JobRegisterResponse,
	qopts *api.QueryOptions,
) (*api.Evaluation, *api.QueryMeta, error) {

	for {
		eval, meta, err := client.Evaluations().Info(resp.EvalID, qopts)
		if err != nil {
			return nil, nil, err
		}

		qopts.WaitIndex = meta.LastIndex

		switch eval.Status {
		case "pending":
			s.Update("Nomad allocation pending...")
		case "complete":
			s.Update("Nomad allocation created")

			return eval, meta, nil
		case "failed", "canceled", "blocked":
			s.Update("Nomad failed to schedule the job")
			s.Status(terminal.StatusError)
			return nil, nil, fmt.Errorf("Nomad evaluation did not transition to 'complete'")
		default:
			return nil, nil, fmt.Errorf("receieved unknown eval status from Nomad: %q", eval.Status)
		}
	}
}

// waypointNomadJob takes in a nomadConfig and returns a Nomad Job per the
// Nomad API
func waypointNomadJob(c nomadConfig, rawRunFlags []string, upgrade bool) *api.Job {
	job := api.NewServiceJob(serverName, serverName, c.region, 50)
	job.Namespace = &c.namespace
	job.Datacenters = c.datacenters
	job.Meta = c.serviceAnnotations
	tg := api.NewTaskGroup(serverName, 1)

	grpcPort, _ := strconv.Atoi(defaultGrpcPort)
	httpPort, _ := strconv.Atoi(defaultHttpPort)

	// Include services to be registered. Currently configured to happen with Consul by default
	// One service added for Waypoint UI, and one for Waypoint backend port
	var services []*api.Service
	if (c.serviceProvider == "consul" && (c.consulService || upgrade)) || (c.consulService && (c.serviceProvider == "consul" || upgrade)) {
		token := ""
		if c.consulToken == "" {
			token = os.Getenv("CONSUL_HTTP_TOKEN")
		} else {
			token = c.consulToken
		}
		job.ConsulToken = &token

		// if the user explicitly set tags with the Consul flags, we use those
		// otherwise, use platform-agnostic tag flags (for backwards-compatibility)
		uiTags := c.consulServiceUITags
		if len(uiTags) == 0 {
			uiTags = c.serviceUITags
		}

		backendTags := c.consulServiceBackendTags
		if len(backendTags) == 0 {
			backendTags = c.serviceUITags
		}

		services = []*api.Service{
			{
				Name:      waypointUIServiceName,
				PortLabel: "ui",
				Tags:      uiTags,
				Provider:  "consul",
			},
			{
				Name:      waypointBackendServiceName,
				PortLabel: "server",
				Tags:      backendTags,
				Provider:  "consul",
			},
		}
	} else if c.serviceProvider == "nomad" {
		services = []*api.Service{
			{
				Name:      waypointUIServiceName,
				PortLabel: "ui",
				Tags:      c.serviceUITags,
				Provider:  "nomad",
			},
			{
				Name:      waypointBackendServiceName,
				PortLabel: "server",
				Tags:      c.serviceBackendTags,
				Provider:  "nomad",
			},
		}
	}
	tg.Services = services

	tg.Networks = []*api.NetworkResource{
		{
			Mode: "host",
			// currently set to static; when ui command can be dynamic - update this
			ReservedPorts: []api.Port{
				{
					Label: "ui",
					Value: httpPort,
					To:    httpPort,
				},
				{
					Label: "server",
					To:    grpcPort,
					Value: grpcPort,
				},
			},
		},
	}

	// Preserve disk, otherwise upgrades will destroy previous allocation and the disk along with it
	volumeRequest := api.VolumeRequest{ReadOnly: false}

	if c.csiVolumeProvider != "" {
		volumeRequest.Type = "csi"
		volumeRequest.Source = "waypoint-server"
		volumeRequest.AccessMode = "single-node-writer"
		volumeRequest.AttachmentMode = "file-system"
	} else {
		volumeRequest.Type = "host"
		volumeRequest.Source = c.hostVolume
	}

	tg.Volumes = map[string]*api.VolumeRequest{
		"waypoint-server": &volumeRequest,
	}

	job.AddTaskGroup(tg)

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

	preTask := nomad.SetupPretask(volumeMounts)

	tg.AddTask(preTask)

	ras := []string{"server", "run", "-accept-tos", "-vv", "-db=/data/data.db", fmt.Sprintf("-listen-grpc=0.0.0.0:%s", defaultGrpcPort), fmt.Sprintf("-listen-http=0.0.0.0:%s", defaultHttpPort)}
	ras = append(ras, rawRunFlags...)
	task := api.NewTask("server", "docker")
	task.Config = map[string]interface{}{
		"image":          c.serverImage,
		"ports":          []string{"server", "ui"},
		"args":           ras,
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

	// Let the runner know about the Nomad IP
	if c.nomadHost == "" {
		c.nomadHost = defaultNomadHost
	}
	task.Env["NOMAD_ADDR"] = c.nomadHost

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

func (i *NomadInstaller) OnDemandRunnerConfig() *pb.OnDemandRunnerConfig {
	// Generate some configuration
	cfgMap := map[string]interface{}{}
	if v := i.config.runnerResourcesCPU; v != "" {
		cfgMap["resources_cpu"] = v
	}
	if v := i.config.runnerResourcesMemory; v != "" {
		cfgMap["resources_memory"] = v
	}
	if v := i.config.datacenters[0]; v != "" {
		cfgMap["datacenter"] = v
	}
	if v := i.config.namespace; v != "" {
		cfgMap["namespace"] = v
	}
	if v := i.config.region; v != "" {
		cfgMap["region"] = v
	}
	if v := i.config.nomadHost; v != "" {
		cfgMap["nomad_host"] = v
	}

	// Marshal our config
	cfgJson, err := json.MarshalIndent(cfgMap, "", "\t")
	if err != nil {
		// This shouldn't happen cause we control our input. If it does,
		// just panic cause this will be in a `server install` CLI and
		// we want the user to report a bug.
		panic(err)
	}

	return &pb.OnDemandRunnerConfig{
		Name:         "nomad",
		OciUrl:       i.config.odrImage,
		PluginType:   "nomad",
		Default:      true,
		PluginConfig: cfgJson,
		ConfigFormat: pb.Hcl_JSON,
	}
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
		Name:    "nomad-host",
		Target:  &i.config.nomadHost,
		Default: defaultNomadHost,
		Usage:   "Hostname of the Nomad server to use, like for launching on-demand tasks.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-namespace",
		Target:  &i.config.namespace,
		Default: "default",
		Usage:   "Namespace to install the Waypoint server into for Nomad.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-odr-image",
		Target: &i.config.odrImage,
		Usage: "Docker image for the on-demand runners. If not specified, it " +
			"defaults to the server image name + '-odr' (i.e. 'hashicorp/waypoint-odr:latest')",
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
		Name:   "nomad-runner-host-volume",
		Target: &i.config.runnerHostVolume,
		Usage:  "Name of the host volume to use for the Waypoint runner.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-runner-csi-volume-provider",
		Target: &i.config.runnerCsiVolumeProvider,
		Usage:  "Name of the CSI volume provider to use for the Waypoint runner.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-runner-csi-volume",
		Target: &i.config.runnerCsiVolume,
		Usage:  "The name of the volume to initialize for the Waypoint runner within the CSI provider.",
	})

	// TODO: Update default values for runner - less space is needed for runner compared to server
	set.Int64Var(&flag.Int64Var{
		Name:    "nomad-runner-csi-volume-capacity-min",
		Target:  &i.config.runnerCsiVolumeCapacityMin,
		Usage:   "Waypoint runner Nomad CSI volume capacity minimum, in bytes.",
		Default: defaultCSIVolumeCapacityMin,
	})

	set.Int64Var(&flag.Int64Var{
		Name:    "nomad-runner-csi-volume-capacity-max",
		Target:  &i.config.runnerCsiVolumeCapacityMax,
		Usage:   "Waypoint runner Nomad CSI volume capacity maximum, in bytes.",
		Default: defaultCSIVolumeCapacityMax,
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: installutil.DefaultServerImage,
	})

	set.EnumSingleVar(&flag.EnumSingleVar{
		Name:    "nomad-service-provider",
		Target:  &i.config.serviceProvider,
		Usage:   "Create service for Waypoint UI and Server in Consul.",
		Values:  []string{"consul", "nomad", "none"},
		Default: "consul", //default to consul for fresh installs
	})

	set.StringSliceVar(&flag.StringSliceVar{
		Name:    "nomad-service-ui-tags",
		Target:  &i.config.serviceUITags,
		Usage:   "Tags for the Waypoint UI service.",
		Default: []string{defaultServiceTag},
	})

	set.StringSliceVar(&flag.StringSliceVar{
		Name:    "nomad-service-backend-tags",
		Target:  &i.config.serviceBackendTags,
		Usage:   "Tags for the Waypoint backend service.",
		Default: []string{defaultServiceTag},
	})

	set.BoolVar(&flag.BoolVar{
		Name:    "nomad-consul-service",
		Target:  &i.config.consulService,
		Usage:   "Create service for Waypoint UI and Server in Consul.",
		Default: true, //default to true for fresh installs
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-consul-service-hostname",
		Target: &i.config.consulServiceHostname,
		Usage: "If set, will use this hostname for Consul DNS rather than the default, " +
			"i.e. \"waypoint-server.service.consul\".",
		Default: "",
	})

	set.StringSliceVar(&flag.StringSliceVar{
		Name:   "nomad-consul-service-ui-tags",
		Target: &i.config.consulServiceUITags,
		Usage:  "Tags for the Waypoint UI service generated in Consul.",
	})

	set.StringSliceVar(&flag.StringSliceVar{
		Name:   "nomad-consul-service-backend-tags",
		Target: &i.config.consulServiceBackendTags,
		Usage: "Tags for the Waypoint backend service generated in Consul. The 'first' tag " +
			"will be used when crafting the Consul DNS hostname for accessing Waypoint.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-consul-datacenter",
		Target:  &i.config.consulDatacenter,
		Usage:   "The datacenter where Consul is located.",
		Default: defaultConsulDatacenter,
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-consul-domain",
		Target:  &i.config.consulDomain,
		Usage:   "The domain where Consul is located.",
		Default: defaultConsulDomain,
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-consul-token",
		Target: &i.config.consulToken,
		Usage: "If set, the passed Consul token is stored in the job " +
			"before sending to the Nomad servers. Overrides the CONSUL_HTTP_TOKEN " +
			"environment variable if set.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-host-volume",
		Target: &i.config.hostVolume,
		Usage:  "Nomad host volume name to use for the Waypoint server, required for volume type 'host'.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-csi-volume-provider",
		Target: &i.config.csiVolumeProvider,
		Usage:  "Nomad CSI volume provider, required for volume type 'csi'.",
	})

	set.Int64Var(&flag.Int64Var{
		Name:    "nomad-csi-volume-capacity-min",
		Target:  &i.config.csiVolumeCapacityMin,
		Usage:   "Nomad CSI volume capacity minimum, in bytes.",
		Default: defaultCSIVolumeCapacityMin,
	})

	set.Int64Var(&flag.Int64Var{
		Name:    "nomad-csi-volume-capacity-max",
		Target:  &i.config.csiVolumeCapacityMax,
		Usage:   "Nomad CSI volume capacity maximum, in bytes.",
		Default: defaultCSIVolumeCapacityMax,
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-csi-fs",
		Target:  &i.config.csiFS,
		Usage:   "Nomad CSI volume mount option file system.",
		Default: defaultCSIVolumeMountFS,
	})

	set.StringMapVar(&flag.StringMapVar{
		Name:   "nomad-csi-secrets",
		Target: &i.config.csiSecrets,
		Usage:  "Secrets to provide for the CSI volume.",
	})

	set.StringMapVar(&flag.StringMapVar{
		Name:   "nomad-csi-parameters",
		Target: &i.config.csiParams,
		Usage:  "Parameters passed directly to the CSI plugin to configure the volume.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-csi-plugin-id",
		Target: &i.config.csiPluginId,
		Usage:  "The ID of the CSI plugin that manages the volume, required for volume type 'csi'.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-csi-external-id",
		Target: &i.config.csiExternalId,
		Usage:  "The ID of the physical volume from the Nomad storage provider.",
	})

	set.StringMapVar(&flag.StringMapVar{
		Name:   "nomad-csi-topologies",
		Target: &i.config.csiTopologies,
		Usage:  "Locations from which the Nomad Volume will be accessible.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-csi-volume",
		Target:  &i.config.csiVolume,
		Usage:   "The name of the volume to initialize for Waypoint server within the CSI provider.",
		Default: "waypoint-server",
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
		Name:    "nomad-host",
		Target:  &i.config.nomadHost,
		Default: "http://localhost:4646",
		Usage:   "Hostname of the Nomad server to use, like for launching on-demand tasks.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-namespace",
		Target:  &i.config.namespace,
		Default: "default",
		Usage:   "Namespace to install the Waypoint server into for Nomad.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-odr-image",
		Target: &i.config.odrImage,
		Usage: "Docker image for the on-demand runners. If not specified, it " +
			"defaults to the server image name + '-odr' (i.e. 'hashicorp/waypoint-odr:latest')",
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
		Name:   "nomad-runner-host-volume",
		Target: &i.config.runnerHostVolume,
		Usage:  "Name of the host volume to use for the Waypoint runner.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-runner-csi-volume-provider",
		Target: &i.config.runnerCsiVolumeProvider,
		Usage:  "Name of the CSI volume provider to use for the Waypoint runner.",
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-runner-csi-volume",
		Target: &i.config.runnerCsiVolume,
		Usage:  "The name of the volume to initialize for the Waypoint runner within the CSI provider.",
	})

	// TODO: Update default values for runner - less space is needed for runner compared to server
	set.Int64Var(&flag.Int64Var{
		Name:    "nomad-runner-csi-volume-capacity-min",
		Target:  &i.config.runnerCsiVolumeCapacityMin,
		Usage:   "Waypoint runner Nomad CSI volume capacity minimum, in bytes.",
		Default: defaultCSIVolumeCapacityMin,
	})

	set.Int64Var(&flag.Int64Var{
		Name:    "nomad-runner-csi-volume-capacity-max",
		Target:  &i.config.runnerCsiVolumeCapacityMax,
		Usage:   "Waypoint runner Nomad CSI volume capacity maximum, in bytes.",
		Default: defaultCSIVolumeCapacityMax,
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-server-image",
		Target:  &i.config.serverImage,
		Usage:   "Docker image for the Waypoint server.",
		Default: installutil.DefaultServerImage,
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-host-volume",
		Target: &i.config.hostVolume,
		Usage:  "Nomad host volume name.",
	})

	set.EnumSingleVar(&flag.EnumSingleVar{
		Name:    "nomad-service-provider",
		Target:  &i.config.serviceProvider,
		Usage:   "Create service for Waypoint UI and Server.",
		Values:  []string{"consul", "nomad", "none"},
		Default: "none", // default none on upgrades, make sure people opt-in to service discovery
	})

	set.StringSliceVar(&flag.StringSliceVar{
		Name:    "nomad-service-ui-tags",
		Target:  &i.config.serviceUITags,
		Usage:   "Tags for the Waypoint UI service.",
		Default: []string{defaultServiceTag},
	})

	set.StringSliceVar(&flag.StringSliceVar{
		Name:    "nomad-service-backend-tags",
		Target:  &i.config.serviceBackendTags,
		Usage:   "Tags for the Waypoint backend service.",
		Default: []string{defaultServiceTag},
	})

	set.BoolVar(&flag.BoolVar{
		Name:    "nomad-consul-service",
		Target:  &i.config.consulService,
		Usage:   "Create service for Waypoint UI and Server in Consul.",
		Default: false, // default to false, make sure people opt into this for upgrades
	})

	set.StringVar(&flag.StringVar{
		Name:   "nomad-consul-service-hostname",
		Target: &i.config.consulServiceHostname,
		Usage: "If set, will use this hostname for Consul DNS rather than the default, " +
			"i.e. \"waypoint-server.service.consul\".",
		Default: "",
	})

	set.StringSliceVar(&flag.StringSliceVar{
		Name:   "nomad-consul-service-ui-tags",
		Target: &i.config.consulServiceUITags,
		Usage:  "Tags for the Waypoint UI service generated in Consul.",
	})

	set.StringSliceVar(&flag.StringSliceVar{
		Name:   "nomad-consul-service-backend-tags",
		Target: &i.config.consulServiceBackendTags,
		Usage: "Tags for the Waypoint backend service generated in Consul. The 'first' tag " +
			"will be used when crafting the Consul DNS hostname for accessing Waypoint.",
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-consul-datacenter",
		Target:  &i.config.consulDatacenter,
		Usage:   "The datacenter where Consul is located.",
		Default: defaultConsulDatacenter,
	})

	set.StringVar(&flag.StringVar{
		Name:    "nomad-consul-domain",
		Target:  &i.config.consulDomain,
		Usage:   "The domain where Consul is located.",
		Default: defaultConsulDomain,
	})
}

func (i *NomadInstaller) UninstallFlags(set *flag.Set) {
	// Purposely empty, no flags
}
