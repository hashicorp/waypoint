package client

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
)

// TODO(izaak): refine description after the dust settles
// Project is the primary structure for interacting with a Waypoint
// server as a client. The client exposes a slightly higher level of
// abstraction over the server API for performing operations locally and
// remotely.
type Project struct {
	UI terminal.UI

	client pb.WaypointClient

	logger              hclog.Logger
	project             *pb.Ref_Project
	workspace           *pb.Ref_Workspace
	runner              *pb.Ref_Runner
	labels              map[string]string
	variables           []*pb.Variable
	dataSourceOverrides map[string]string
	cleanupFunc         func()

	// TODO(izaak): delete?
	serverVersion *pb.VersionInfo

	localServer bool // True when a local server is created

	// These tell the project that all of the jobs it creates should be local, not remote.
	// They're used to template the jobs that it creates, and when interacting with
	// the jobs as they're running.
	executeJobsLocally bool
	localRunnerId      string
}

// TODO(izaak): maybe call this NewProject? It's more than just a client.
// NewProjectClient initializes a new client.
func NewProjectClient(ctx context.Context, opts ...Option) (*Project, error) {
	// Our default projectClient
	projectClient := &Project{
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
		err := o(projectClient, &cfg)
		if err != nil {
			return nil, err
		}
	}

	// TODO(izaak): I get the feeling that by the end of this, projectClient will always be set. It should probably
	//// be given as an arg, or maybe always assumed to be here?
	//// If a client was explicitly provided, we use that. Otherwise, we
	//// have to establish a connection either through the serverclient
	//// package or spinning up an in-process server.
	//if projectClient.client == nil {
	//	projectClient.logger.Trace("no API client provided, initializing connection if possible")
	//	conn, err := projectClient.initServerConnection(ctx, &cfg)
	//	if err != nil {
	//		return nil, err
	//	}
	//	projectClient.client = pb.NewWaypointClient(conn)
	//
	//	// Negotiate the version
	//	ver, err := NegotiateApiVersion(ctx, projectClient.client, projectClient.logger)
	//	if err != nil {
	//		return nil, err
	//	}
	//	projectClient.serverVersion = ver
	//}

	// Default workspace if not specified
	if projectClient.workspace == nil {
		projectClient.workspace = &pb.Ref_Workspace{Workspace: "default"}
	}

	return projectClient, nil
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

// ServerVersion returns the server version that this client is connected to.
func (c *Project) ServerVersion() *pb.VersionInfo {
	return c.serverVersion
}

// Close should be called to clean up any resources that the client created.
func (c *Project) Close() error {

	// TODO(izaak): I don't think we need any of this, or it probably shouldn't live here.

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

// TODO(izaak): probably delete?
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

// WithVariables sets variable values from flags and local env on any operations.
func WithVariables(m []*pb.Variable) Option {
	return func(c *Project, cfg *config) error {
		c.variables = m
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

func WithExecuteJobsLocally(localRunnerId string) Option {
	return func(c *Project, cfg *config) error {
		c.executeJobsLocally = true
		c.localRunnerId = localRunnerId
		return nil
	}
}
