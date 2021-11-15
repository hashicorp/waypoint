package runner

import (
	"context"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/pkg/errors"
)

// LocalRunner is a thin wrapper around Runner that helps manage
// starting and stopping the runner safely
type LocalRunner struct {
	runner   *Runner
	wg       sync.WaitGroup
	bg       context.Context
	bgCancel func()
}

// NewLocalRunner creates a new local runner
func NewLocalRunner(
	client pb.WaypointClient,
	log hclog.Logger,
	ui terminal.UI,
) (*LocalRunner, error) {
	runner, err := New(
		WithClient(client),
		WithLogger(log.Named("runner")),
		WithLocal(ui),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create local runner")
	}

	if err := runner.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start local runner")
	}

	lr := LocalRunner{
		runner: runner,
	}

	lr.bg, lr.bgCancel = context.WithCancel(context.Background())

	return &lr, nil
}

func (r *LocalRunner) RunnerId() string {
	return r.runner.Id()
}

// Start starts the runner. It blocks until the runner's context is cancelled or an error occurs.
func (r *LocalRunner) Start() error {
	r.wg.Add(1)
	r.runner.AcceptMany(r.bg) // This blocks
	r.wg.Done()
	return nil
}

// Close cancels the local runner and waits for it to exit.
func (r *LocalRunner) Close() {
	r.bgCancel()
	r.wg.Wait()
	return
}
