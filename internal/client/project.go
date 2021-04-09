package client

import (
	"context"
	"sync"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/runner"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
)

// Project is the primary structure for interacting with a Waypoint
// server as a client. The client exposes a slightly higher level of
// abstraction over the server API for performing operations locally and
// remotely.
type Project struct {
	UI terminal.UI

	client              pb.WaypointClient
	logger              hclog.Logger
	project             *pb.Ref_Project
	workspace           *pb.Ref_Workspace
	runner              *pb.Ref_Runner
	labels              map[string]string
	dataSourceOverrides map[string]string
	cleanupFunc         func()

	local bool

	localServer bool // True when a local server is created

	// These are used to manage a local runner and it's job processing
	// in a goroutine.
	wg           sync.WaitGroup
	bg           context.Context
	bgCancel     func()
	activeRunner *runner.Runner
}

// New initializes a new client.
func New(ctx context.Context, opts ...Option) (*Project, error) {
	// Our default client
	client := &Project{
		UI:     terminal.ConsoleUI(ctx),
		logger: hclog.L(),
		runner: &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Any{
				Any: &pb.Ref_RunnerAny{},
			},
		},
	}

	// Build our config
	var cfg config
	for _, o := range opts {
		err := o(client, &cfg)
		if err != nil {
			return nil, err
		}
	}

	// Used by any background goroutines that we'll spawn (like runner job processing)
	client.bg, client.bgCancel = context.WithCancel(context.Background())

	// If a client was explicitly provided, we use that. Otherwise, we
	// have to establish a connection either through the serverclient
	// package or spinning up an in-process server.
	if client.client == nil {
		client.logger.Trace("no API client provided, initializing connection if possible")
		conn, err := client.initServerClient(ctx, &cfg)
		if err != nil {
			return nil, err
		}
		client.client = pb.NewWaypointClient(conn)
	}

	// Negotiate the version
	if err := client.negotiateApiVersion(ctx); err != nil {
		return nil, err
	}

	// Default workspace if not specified
	if client.workspace == nil {
		client.workspace = &pb.Ref_Workspace{Workspace: "default"}
	}

	if client.local {
		client.logger.Debug("starting runner to process local jobs")
		r, err := client.startRunner()
		if err != nil {
			return nil, err
		}

		client.activeRunner = r

		// We spin up the job processing here. Anything that spawns jobs (either locally spawned
		// or server spawned) will be processed by this runner ONLY if the runner is directly targeted.
		// Because this runner's lifetime is bound to a CLI context and therefore transient, we don't
		// want to accept jobs that aren't related to local activities (job's queued or RPCs made)
		// because they'll hang the CLI randomly as those jobs run (it's also a security issue).
		client.wg.Add(1)
		go func() {
			defer client.wg.Done()
			r.AcceptMany(client.bg)
		}()
	}

	return client, nil
}

// LocalRunnerId returns the id of the runner that this project started
// This is used to target jobs specifically at this runner.
func (c *Project) LocalRunnerId() (string, bool) {
	if c.activeRunner == nil {
		return "", false
	}

	return c.activeRunner.Id(), true
}

// Ref returns the raw Waypoint server API client.
func (c *Project) Ref() *pb.Ref_Project {
	return c.project
}

// Client returns the raw Waypoint server API client.
func (c *Project) Client() pb.WaypointClient {
	return c.client
}

// WorkspaceRef returns the application reference that this client is using.
func (c *Project) WorkspaceRef() *pb.Ref_Workspace {
	return c.workspace
}

// Local is true if the server is an in-process just-in-time server.
func (c *Project) Local() bool {
	return c.localServer
}

// Close should be called to clean up any resources that the client created.
func (c *Project) Close() error {
	// Stop the runner early so that it we block here waiting for any outstanding jobs to finish
	// before closing down the rest of the resources.
	if c.activeRunner != nil {
		if err := c.activeRunner.Close(); err != nil {
			c.logger.Error("error stopping runner", "error", err)
		}
	}

	// Forces any background goroutines to stop
	c.bgCancel()

	// Now wait on those goroutines to finish up.
	c.wg.Wait()

	// Run any cleanup necessary
	if f := c.cleanupFunc; f != nil {
		f()
	}

	return nil
}

// cleanup stacks cleanup functions to call when Close is called.
func (c *Project) cleanup(f func()) {
	oldF := c.cleanupFunc
	c.cleanupFunc = func() {
		defer f()
		if oldF != nil {
			oldF()
		}
	}
}

type config struct {
	connectOpts []serverclient.ConnectOption
}

type Option func(*Project, *config) error

// WithProjectRef sets the project reference for all operations performed.
func WithProjectRef(ref *pb.Ref_Project) Option {
	return func(c *Project, cfg *config) error {
		c.project = ref
		return nil
	}
}

// WithWorkspaceRef sets the workspace reference for all operations performed.
// If this isn't set, the default workspace will be used.
func WithWorkspaceRef(ref *pb.Ref_Workspace) Option {
	return func(c *Project, cfg *config) error {
		c.workspace = ref
		return nil
	}
}

// WithClient sets the client directly. In this case, the runner won't
// attempt any connection at all regardless of other configuration (env
// vars or waypoint config file). This will be used.
func WithClient(client pb.WaypointClient) Option {
	return func(c *Project, cfg *config) error {
		c.client = client
		return nil
	}
}

// WithClientConnect specifies the options for connecting to a client.
// If WithClient is specified, that client is always used.
//
// If WithLocal is set and no client is specified and no server creds
// can be found, then an in-process server will be created.
func WithClientConnect(opts ...serverclient.ConnectOption) Option {
	return func(c *Project, cfg *config) error {
		cfg.connectOpts = opts
		return nil
	}
}

// WithLabels sets the labels or any operations.
func WithLabels(m map[string]string) Option {
	return func(c *Project, cfg *config) error {
		c.labels = m
		return nil
	}
}

// WithSourceOverrides sets the data source overrides for queued jobs.
func WithSourceOverrides(m map[string]string) Option {
	return func(c *Project, cfg *config) error {
		c.dataSourceOverrides = m
		return nil
	}
}

// WithLocal puts the client in local exec mode. In this mode, the client
// will spin up a per-operation runner locally and reference the local on-disk
// data for all operations.
func WithLocal() Option {
	return func(c *Project, cfg *config) error {
		c.local = true
		return nil
	}
}

// WithLogger sets the logger for the client.
func WithLogger(log hclog.Logger) Option {
	return func(c *Project, cfg *config) error {
		c.logger = log
		return nil
	}
}

// WithUI sets the UI to use for the client.
func WithUI(ui terminal.UI) Option {
	return func(c *Project, cfg *config) error {
		c.UI = ui
		return nil
	}
}
