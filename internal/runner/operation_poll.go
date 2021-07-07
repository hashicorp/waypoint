package runner

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/waypoint/internal/datasource"

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
	sourcer, err := r.dataSourcer(ctx, log, job.DataSource, job.DataSourceOverrides)
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
	if err != nil || resp == nil {
		return nil, err
	}

	// For each workspace this project has been deployed to, check if there is a ScopedSetting that
	// overrides the original DataSource

	log.Trace("finding latest ref")
	var queueResults []*pb.Job_PollResult
	for _, wsProject := range resp.Workspaces {
		var dataSource *pb.Job_DataSource
		dataSourceRef := wsProject.DataSourceRef
		wsName := wsProject.Workspace.Workspace

		// Check if we have a custom data source
		v, ok := resp.Project.ScopedSettings[wsName]
		if wsName == "default" || !ok {
			// Workspace default or custom data source not found -> default Data Source
			dataSource = resp.Project.DataSource
		} else if v != nil {
			dataSource = v.DataSource
		}

		if dataSource == nil {
			log.Warn("data source is nil", "workspace", wsName, "project", resp.Project)
			continue
		}

		jobResult, err := r.detectChanges(log, wsName, resp.Project, dataSource, dataSourceRef, sourcer, ui, ctx, job)
		if err != nil {
			log.Error("unable to detect changes",
				"err", err,
				"project", resp.Project.Name,
				"workspace", wsName,
			)
		} else {
			poll := jobResult.Poll
			if poll == nil {
				log.Error("returned result is not a poll result",
					"jobResult", poll,
					"project", resp.Project.Name,
					"workspace", wsName,
				)
				continue
			}
			queueResults = append(queueResults, poll)
		}
	}

	return &pb.Job_Result{
		MultiPoll: &pb.Job_MultiPollResult{
			Results: queueResults,
		},
	}, nil
}

func (r *Runner) detectChanges(log hclog.Logger,
	workspace string,
	project *pb.Project,
	dataSource *pb.Job_DataSource,
	ref *pb.Job_DataSource_Ref,
	sourcer datasource.Sourcer,
	ui terminal.UI,
	ctx context.Context,
	job *pb.Job,
) (*pb.Job_Result, error) {
	log.Debug("current ref for poll operation", "ref", ref)

	// Get any change
	newRef, ignore, err := sourcer.Changes(ctx, log, ui, dataSource, ref, r.tempDir)
	if err != nil {
		return nil, fmt.Errorf("unable to get changes")
	}
	log.Debug("result of Changes, nil means no changes", "result", newRef, "ignore", ignore)

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

	// Setup our job template. This will be used with the QueueProject operation
	// to queue an "up" for each app within the project.
	jobTemplate := &pb.Job{
		TargetRunner: &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Any{
				Any: &pb.Ref_RunnerAny{},
			},
		},

		Workspace: &pb.Ref_Workspace{Workspace: workspace},

		// Reuse the same data source and bring in our overrides to set the ref
		DataSource:          job.DataSource,
		DataSourceOverrides: overrides,

		// Doing a plain old "up"
		Operation: &pb.Job_Up{
			Up: &pb.Job_UpOp{},
		},
	}

	// If we're ignoring, we change the job to a noop job. This will
	// still trigger the machinery to update the ref associated with
	// the project/app and avoids the poll job from having to have too
	// much access or require new APIs to do this.
	if ignore {
		log.Debug("changes marked as ignorable, scheduling a noop job to update our data ref")
		jobTemplate.Operation = &pb.Job_Noop_{
			Noop: &pb.Job_Noop{},
		}
	}

	log.Debug("queueing job")
	queueResp, err := r.client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: &pb.Job{
			// We set a singleton ID to verify that we only setup
			// a queue operation once (in case it is taking longer to
			// process than the poll interval).
			SingletonId: strings.ToLower(fmt.Sprintf(
				"poll-trigger/%s/%s/%s",
				job.Application.Project,
				job.Application.Application,
				job.Workspace.Workspace,
			)),

			// Target only our project, we don't need an app for this.
			Application: &pb.Ref_Application{
				Project: project.Name,
			},

			// Copy all of these fields from the job template since we
			// want to execute the same way.
			TargetRunner:        jobTemplate.TargetRunner,
			Workspace:           jobTemplate.Workspace,
			DataSource:          jobTemplate.DataSource,
			DataSourceOverrides: jobTemplate.DataSourceOverrides,

			// Queue the job
			Operation: &pb.Job_QueueProject{
				QueueProject: &pb.Job_QueueProjectOp{
					JobTemplate: jobTemplate,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	log.Debug("job queued", "job_id", queueResp.JobId)
	return &pb.Job_Result{
		Poll: &pb.Job_PollResult{
			JobId:  queueResp.JobId,
			OldRef: ref,
			NewRef: newRef,
		},
	}, nil
}
