package runner

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (r *Runner) executePollOp(
	ctx context.Context,
	log hclog.Logger,
	ui terminal.UI,
	job *pb.Job,
) (*pb.Job_Result, error) {
	sourcer, err := r.dataSourcer(ctx, log, ui, job.DataSource, job.DataSourceOverrides)
	if err != nil {
		return nil, err
	}

	// Query this project. We're mainly trying to get all the pb.Workspace_Project
	// values for a project so that we can get the data ref that we last polled
	// for each project.
	log.Trace("calling GetProject to get list of workspaces for project")
	resp, err := r.client.GetProject(ctx, &pb.GetProjectRequest{
		Project: &pb.Ref_Project{
			Project: job.Application.Project,
		},
	})
	if err != nil {
		return nil, err
	}
	project := resp.Project

	// Get the current ref for the default workspace.
	//
	// NOTE(mitchellh): for now, we only support the default workspace. We
	// will expand support to polling non-default workspaces in the future.
	log.Trace("finding latest ref")
	var currentRef *pb.Job_DataSource_Ref
	if resp != nil {
		for _, p := range resp.Workspaces {
			if p.Workspace.Workspace == "default" {
				currentRef = p.DataSourceRef
				break
			}
		}
	}
	log.Debug("current ref for poll operation", "ref", currentRef)

	// Get any changes
	newRef, err := sourcer.Changes(ctx, log, ui, job.DataSource, currentRef)
	if err != nil {
		return nil, err
	}
	log.Debug("result of Changes, nil means no changes", "result", newRef)

	// If we have no changes, then we're done.
	if newRef == nil {
		return &pb.Job_Result{}, nil
	}

	// Setup our overrides. Overrides are used to set the exact ref that
	// the job will use.
	overrides, err := sourcer.RefToOverride(newRef)
	if err != nil {
		return nil, err
	}

	// Setup our base job that we'll queue an "up" for. We'll setup a new
	// job for each application within the project but this will be the
	// common fields.
	baseJob := func() *pb.Job {
		return &pb.Job{
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Any{
					Any: &pb.Ref_RunnerAny{},
				},
			},

			// NOTE(mitchellh): default workspace only for now
			Workspace: &pb.Ref_Workspace{Workspace: "default"},

			// Reuse the same data source and bring in our overrides to set the ref
			DataSource:          job.DataSource,
			DataSourceOverrides: overrides,

			// Doing a plain old "up"
			Operation: &pb.Job_Up{
				Up: &pb.Job_UpOp{},
			},
		}
	}

	// Go through each application registered with the project and run up.
	for _, app := range project.Applications {
		job := baseJob()
		job.Application = &pb.Ref_Application{
			Project:     project.Name,
			Application: app.Name,
		}

		// We set a singleton ID so that future polling operations will
		// override this if it already exists.
		job.SingletonId = strings.ToLower(fmt.Sprintf(
			"poll-trigger/%s/%s/%s",
			job.Application.Project,
			job.Application.Application,
			job.Workspace.Workspace,
		))

		log := log.With("project", project.Name, "app", app.Name)

		// Queue it
		log.Debug("queueing job")
		resp, err := r.client.QueueJob(ctx, &pb.QueueJobRequest{
			Job: job,
		})
		if err != nil {
			return nil, err
		}
		log.Debug("job queued", "job_id", resp.JobId)
	}

	return &pb.Job_Result{}, nil
}
