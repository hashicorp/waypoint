package client

import (
	"github.com/hashicorp/waypoint/internal/runner"
)

// startRunner initializes and starts a local runner. If the returned
// runner is non-nil, you must call Close on it to clean up resources properly.
func (c *Client) startRunner() (*runner.Runner, error) {
	// Initialize our runner
	r, err := runner.New(
		runner.WithClient(c.client),
		runner.WithLogger(c.logger),
	)
	if err != nil {
		return nil, err
	}

	// Start the runner
	if err := r.Start(); err != nil {
		return nil, err
	}

	return r, nil
}
