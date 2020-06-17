package runner

import (
	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// Runners in Waypoint execute operations. These can be local (the CLI)
// or they can be remote (triggered by some webhook). In either case, they
// share this same underlying implementation.
//
// To use a runner:
//
//   1. Initialize it with New. This will setup some initial state but
//      will not register with the server or run jobs.
//
//   2. Start the runner with "Start". This will register the runner and
//      kick off some management goroutines. This will not execute any jobs.
//
//   3. Run a single job with "Accept". This is named to be similar to a
//      network listener "accepting" a connection. This will request a single
//      job from the Waypoint server, block until one is available, and execute
//      it. Repeat this call for however many jobs you want to execute.
//
//   4. Clean up with "Close". This will gracefully exit the runner, waiting
//      for any running jobs to finish.
//
type Runner struct {
	logger hclog.Logger
	client pb.WaypointClient
}

// New initializes a new runner.
//
// You must call Start to start the runner and register with the Waypoint
// server. See the Runner struct docs for more details.
func New(opts ...Option) (*Runner, error) {
	return nil, nil
}

// Start starts the runner by registering the runner with the Waypoint
// server. This will spawn goroutines for management. This will return after
// registration so this should not be executed in a goroutine.
func (r *Runner) Start() error {
	return nil
}

// Accept will accept and execute a single job. This will block until
// a job is available.
//
// This is safe to be called concurrently which can be used to execute
// multiple jobs in parallel as a runner.
func (r *Runner) Accept() error {
	return nil
}

type config struct{}

type Option func(*Runner, *config) error
