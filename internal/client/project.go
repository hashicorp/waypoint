package client

import (
	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

// Project is the primary structure for interacting with a Waypoint
// server as a client. The client exposes a slightly higher level of
// abstraction over the server API for performing operations locally and
// remotely.
type Project struct {
	UI terminal.UI

	client      pb.WaypointClient
	logger      hclog.Logger
	project     *pb.Ref_Project
	workspace   *pb.Ref_Workspace
	runner      *pb.Ref_Runner
	labels      map[string]string
	cleanupFunc func()

	local bool
}

// New initializes a new client.
func New(opts ...Option) (*Project, error) {
	// Our default client
	client := &Project{
		UI:     &terminal.BasicUI{},
		logger: hclog.L(),
	}

	// Build our config
	var cfg config
	for _, o := range opts {
		err := o(client, &cfg)
		if err != nil {
			return nil, err
		}
	}

	// If a client was explicitly provided, we use that. Otherwise, we
	// have to establish a connection either through the serverclient
	// package or spinning up an in-process server.
	if client.client == nil {
		client.logger.Trace("no API client provided, initializing connection if possible")
		conn, err := client.initServerClient(&cfg)
		if err != nil {
			return nil, err
		}
		client.client = pb.NewWaypointClient(conn)
	}

	// Default workspace if not specified
	if client.workspace == nil {
		client.workspace = &pb.Ref_Workspace{Workspace: "default"}
	}

	return client, nil
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

// Close should be called to clean up any resources that the client created.
func (c *Project) Close() error {
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
