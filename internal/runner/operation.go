package runner

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/datadir"
)

// executeJob executes an assigned job. This will source the data (if necessary),
// setup the project, execute the job, and return the outcome.
func (r *Runner) executeJob(ctx context.Context, log hclog.Logger, job *pb.Job) error {
	// Eventually we'll need to extract the data source. For now we're
	// just building for local exec so it is the working directory.
	path := configpkg.Filename

	// Decode the configuration
	var cfg configpkg.Config
	log.Trace("reading configuration", "path", path)
	if err := hclsimple.DecodeFile(path, nil, &cfg); err != nil {
		return err
	}

	// Setup our project data directory.
	projDir, err := datadir.NewProject(".waypoint")
	if err != nil {
		return err
	}

	// Create our project
	log.Trace("initializing project", "project", cfg.Project)
	project, err := core.NewProject(ctx,
		core.WithLogger(log),
		core.WithConfig(&cfg),
		core.WithDataDir(projDir),
		core.WithWorkspace("default"), // TODO(mitchellh): configurable
	)
	if err != nil {
		return err
	}
	defer project.Close()

	// Execute the operation
	log.Info("executing operation", "type", fmt.Sprintf("%T", job.Operation))
	switch job.Operation.(type) {
	case *pb.Job_Noop_:
		log.Debug("noop job success")
		return nil

	default:
		return status.Errorf(codes.Aborted, "unknown operation %T", job.Operation)
	}
}
