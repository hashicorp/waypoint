package client

import (
	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverclient"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

// Client is the primary structure for interacting with a Waypoint
// server as a client. The client exposes a slightly higher level of
// abstraction over the server API for performing operations locally and
// remotely.
type Client struct {
	client      pb.WaypointClient
	logger      hclog.Logger
	application *pb.Ref_Application
	runner      *pb.Ref_Runner
	ui          terminal.UI
	cleanupFunc func()

	local bool
}

// New initializes a new client.
func New(opts ...Option) (*Client, error) {
	// Our default client
	client := &Client{
		logger: hclog.L(),
		ui:     &terminal.BasicUI{},
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

	return client, nil
}

// APIClient returns the raw Waypoint server API client.
func (c *Client) APIClient() pb.WaypointClient {
	return c.client
}

// AppRef returns the application reference that this client is using.
func (c *Client) AppRef() *pb.Ref_Application {
	return c.application
}

// Close should be called to clean up any resources that the client created.
func (c *Client) Close() error {
	// Run any cleanup necessary
	if f := c.cleanupFunc; f != nil {
		f()
	}

	return nil
}

// cleanup stacks cleanup functions to call when Close is called.
func (c *Client) cleanup(f func()) {
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

type Option func(*Client, *config) error

// WithAppRef sets the application reference for all operations performed.
// This determines what project/application operations such as "Build" will
// target.
func WithAppRef(ref *pb.Ref_Application) Option {
	return func(c *Client, cfg *config) error {
		c.application = ref
		return nil
	}
}

// WithClient sets the client directly. In this case, the runner won't
// attempt any connection at all regardless of other configuration (env
// vars or waypoint config file). This will be used.
func WithClient(client pb.WaypointClient) Option {
	return func(c *Client, cfg *config) error {
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
	return func(c *Client, cfg *config) error {
		cfg.connectOpts = opts
		return nil
	}
}

// WithLocal puts the client in local exec mode. In this mode, the client
// will spin up a per-operation runner locally and reference the local on-disk
// data for all operations.
func WithLocal() Option {
	return func(c *Client, cfg *config) error {
		c.local = true
		return nil
	}
}

// WithLogger sets the logger for the client.
func WithLogger(log hclog.Logger) Option {
	return func(c *Client, cfg *config) error {
		c.logger = log
		return nil
	}
}
