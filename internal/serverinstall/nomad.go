package serverinstall

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverconfig"
)

type NomadInstaller struct {
	config nomadConfig
}

type nomadConfig struct {
	serverImage        string            `hcl:"server_image,optional"`
	namespace          string            `hcl:"namespace,optional"`
	serviceAnnotations map[string]string `hcl:"service_annotations,optional"`

	region         string   `hcl:"namespace,optional"`
	datacenters    []string `hcl:"datacenters,optional"`
	policyOverride bool     `hcl:"policy_override,optional"`
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
	jobs, _, err := client.Jobs().PrefixList("waypoint-server")
	if err != nil {
		return nil, err
	}
	var serverDetected bool
	for _, j := range jobs {
		if j.Name == "waypoint-server" {
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
	}

	addr.Tls = true
	addr.TlsSkipVerify = true

	if serverDetected {
		allocs, _, err := client.Jobs().Allocations("waypoint-server", false, nil)
		if err != nil {
			return nil, err
		}
		if len(allocs) == 0 {
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
		s.Update("Nomad failed to schedule the waypoint-server")
		s.Status(terminal.StatusError)
		return nil, fmt.Errorf("nomad evaluation did not transition to 'complete'")
	default:
		return nil, fmt.Errorf("unknown eval status: %q", eval.Status)
	}

	var allocID string

	for {
		allocs, qmeta, err := client.Evaluations().Allocations(eval.ID, qopts)
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

	s.Update("Nomad allocation ready")
	s.Done()

	return &InstallResults{
		Context:       &clicfg,
		AdvertiseAddr: &addr,
		HTTPAddr:      httpAddr,
	}, nil
}

// waypointNomadJob takes in a nomadConfig and returns a Nomad Job per the
// Nomad API
func waypointNomadJob(c nomadConfig) *api.Job {
	job := api.NewServiceJob("waypoint-server", "waypoint-server", c.region, 50)
	job.Namespace = &c.namespace
	job.Datacenters = c.datacenters
	job.Meta = c.serviceAnnotations
	tg := api.NewTaskGroup("waypoint-server", 1)
	tg.Networks = []*api.NetworkResource{
		{
			Mode: "host",
			DynamicPorts: []api.Port{
				{
					Label: "server",
					To:    9701,
				},
			},
			// currently set to static; when ui command can be dynamic - update this
			ReservedPorts: []api.Port{
				{
					Label: "ui",
					Value: 9702,
					To:    9702,
				},
			},
		},
	}
	job.AddTaskGroup(tg)

	task := api.NewTask("server", "docker")
	task.Config = map[string]interface{}{
		"image": c.serverImage,
		"ports": []string{"server", "ui"},
		"args":  []string{"server", "run", "-accept-tos", "-vvv", "-db=/alloc/data.db", "-listen-grpc=0.0.0.0:9701", "-listen-http=0.0.0.0:9702"},
	}
	task.Env = map[string]string{
		"PORT": "9701",
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

// InstallRunner implements Installer.
func (i *NomadInstaller) InstallRunner(
	ctx context.Context,
	opts *InstallRunnerOpts,
) error {
	// TODO
	return nil
}

func (i *NomadInstaller) InstallFlags(set *flag.Set) {
	set.StringMapVar(&flag.StringMapVar{
		Name:   "nomad-annotate-service",
		Target: &i.config.serviceAnnotations,
		Usage:  "Annotations for the Service generated.",
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
		Default: "hashicorp/waypoint:latest",
	})
}
