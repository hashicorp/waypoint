package client

import (
	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
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

	return client, nil
}

type config struct{}

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

// WithLocal puts the client in local exec mode. In this mode, the client
// will spin up a per-operation runner locally and reference the local on-disk
// data for all operations.
func WithLocal() Option {
	return func(c *Client, cfg *config) error {
		c.local = true
		return nil
	}
}
