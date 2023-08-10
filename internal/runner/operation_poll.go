// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package runner

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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
	if err != nil {
		return nil, err
	}

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
	newRef, ignore, err := sourcer.Changes(ctx, log, ui, job.DataSource, currentRef, r.tempDir)
	if err != nil {
		return nil, err
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

	// NOTE(briancain): We set a singleton ID for a poll project operation to ensure that the
	// poll handler does not fire off many operations of the same kind more than once,
	// clogging up the job system. By setting a singleton ID that is unique to this
	// application or project, we can ensure only 1 operation will be active at once rather than
	// many operations (such as in the case where a poll interval is shorter than it
	// takes to run the operation)

	// We assume a project and workspace is set given this is Project polling
	singletonId := strings.ToLower(fmt.Sprintf(
		"poll-trigger/%s/%s",
		job.Workspace.Workspace,
		job.Application.Project,
	))
	// Not all jobs set an application
	if job.Application.Application != "" {
		singletonId += "/" + strings.ToLower(job.Application.Application)
	}

	log.Debug("queueing job")
	queueResp, err := r.client.QueueJob(ctx, &pb.QueueJobRequest{
		Job: &pb.Job{
			// We set a singleton ID to verify that we only setup
			// a queue operation once (in case it is taking longer to
			// process than the poll interval).
			SingletonId: singletonId,

			// Target only our project, we don't need an app for this.
			Application: &pb.Ref_Application{
				Project: resp.Project.Name,
			},

			// Copy all of these fields from the job template since we
			// want to execute the same way.
			TargetRunner:        jobTemplate.TargetRunner,
			Workspace:           jobTemplate.Workspace,
			DataSource:          jobTemplate.DataSource,
			DataSourceOverrides: jobTemplate.DataSourceOverrides,

			// Doing a plain old "up"
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
			OldRef: currentRef,
			NewRef: newRef,
		},
	}, nil
}
