package runner

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/sdk/component"
	"github.com/hashicorp/waypoint/sdk/datadir"
	"github.com/hashicorp/waypoint/sdk/terminal"
)

// executeJob executes an assigned job. This will source the data (if necessary),
// setup the project, execute the job, and return the outcome.
func (r *Runner) executeJob(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	job *pb.Job,
) (*pb.Job_Result, error) {
	// Eventually we'll need to extract the data source. For now we're
	// just building for local exec so it is the working directory.
	path := configpkg.Filename

	// Decode the configuration
	var cfg configpkg.Config
	log.Trace("reading configuration", "path", path)
	// TODO: This also consults env for things like the waypoint url token.
	// This should instead consult the auth config storage instead when that arrives.
	if err := cfg.LoadPath(path); err != nil {
		return nil, err
	}

	// Setup our project data directory.
	projDir, err := datadir.NewProject(".waypoint")
	if err != nil {
		return nil, err
	}

	// Build our job info
	jobInfo := &component.JobInfo{
		Local: r.local,
	}

	// Create our project
	log.Trace("initializing project", "project", cfg.Project)
	project, err := core.NewProject(ctx,
		core.WithLogger(log),
		core.WithUI(ui),
		core.WithComponents(r.factories),
		core.WithClient(r.client),
		core.WithConfig(&cfg),
		core.WithDataDir(projDir),
		core.WithLabels(job.Labels),
		core.WithWorkspace(job.Workspace.Workspace),
		core.WithJobInfo(jobInfo),
	)
	if err != nil {
		return nil, err
	}
	defer project.Close()

	// Execute the operation
	log.Info("executing operation", "type", fmt.Sprintf("%T", job.Operation))
	switch job.Operation.(type) {
	case *pb.Job_Noop_:
		if r.noopCh != nil {
			select {
			case <-r.noopCh:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		log.Debug("noop job success")
		return nil, nil

	case *pb.Job_Build:
		return r.executeBuildOp(ctx, job, project)

	case *pb.Job_Push:
		return r.executePushOp(ctx, job, project)

	case *pb.Job_Deploy:
		return r.executeDeployOp(ctx, job, project)

	case *pb.Job_DestroyDeploy:
		return r.executeDestroyDeployOp(ctx, job, project)

	case *pb.Job_Release:
		return r.executeReleaseOp(ctx, log, job, project)

	case *pb.Job_Validate:
		return r.executeValidateOp(ctx, job, project)

	case *pb.Job_Auth:
		return r.executeAuthOp(ctx, log, job, project)

	default:
		return nil, status.Errorf(codes.Aborted, "unknown operation %T", job.Operation)
	}
}
