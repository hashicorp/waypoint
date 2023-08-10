// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package singleprocess

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"github.com/mitchellh/copystructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

// applicationPoll accepts a state management interface which provides access
// to a projects current state implementation. Functions like Peek and Complete
// need access to this state interface for peeking at the next available project
// as well as marking a projects poll as complete.
type applicationPoll struct {
	// state is the state management interface that provides functions for
	// safely mutating server state.
	state serverstate.Interface

	// the workspace to check for polling applications and running their status on
	// Currently, the application poll handler is only expected to work for default
	// workspaces due to the project poller only working on default workspaces.
	workspace string
}

// Peek returns the latest project to poll on
// If there is an error in the ApplicationPollPeek, it will return nil
// to allow the outer caller loop to continue and try again
func (a *applicationPoll) Peek(
	log hclog.Logger,
	ws memdb.WatchSet,
) (interface{}, time.Time, error) {
	ctx := context.Background()
	project, pollTime, err := a.state.ApplicationPollPeek(ctx, ws)
	if err != nil {
		log.Warn("error peeking for next application to poll", "err", err)
		return nil, time.Time{}, err // continue loop
	}

	if project != nil {
		log = log.With("project", project.Name)
		log.Trace("returning peek for apps")
	} else {
		log.Trace("no application returned from peek")
	}

	return project, pollTime, nil
}

// PollJob will generate a slice of QueuedJobRequests to generate a status report
// for each application defined in the given project.
func (a *applicationPoll) PollJob(
	log hclog.Logger,
	p interface{},
) ([]*pb.QueueJobRequest, error) {
	project, ok := p.(*pb.Project)
	if !ok || project == nil {
		log.Error("could not generate poll jobs for projects applications, incorrect type passed in")
		return nil, status.Error(codes.FailedPrecondition,
			"incorrect type passed into Application GeneratePollJobs")
	}
	log = log.Named(project.Name)
	var jobList []*pb.QueueJobRequest

	for _, app := range project.Applications {
		jobs, err := a.buildPollJobs(log, app)
		if err != nil {
			return nil, err
		}

		jobList = append(jobList, jobs...)
	}

	return jobList, nil
}

// buildPollJobs will generate jobs to poll the latest deploy and release operations for a project.
func (a *applicationPoll) buildPollJobs(
	log hclog.Logger,
	appl interface{},
) ([]*pb.QueueJobRequest, error) {
	ctx := context.Background()
	app, ok := appl.(*pb.Application)
	if !ok || app == nil {
		log.Error("could not generate poll job for application, incorrect type passed in")
		return nil, status.Error(codes.FailedPrecondition, "incorrect type passed into Application PollJob")
	}
	log = log.Named(app.Name)

	// App polling needs the parent project to obtain its datasource
	project, err := a.state.ProjectGet(ctx, &pb.Ref_Project{Project: app.Project.Project})
	if err != nil {
		return nil, err
	}

	// Application status polling requires a remote data source, otherwise a status report
	// cannot be generated without a project and its hcl context. This returns
	// an error so we fail early instead of queueing an already broken job
	if project.DataSource == nil {
		log.Debug("cannot build an application poll job without a remote data source configured.")
		return nil, status.Error(codes.FailedPrecondition, "application status polling requires a remote data source")
	}

	// Determine the latest deployment or release to poll for a status report
	appRef := &pb.Ref_Application{
		Application: app.Name,
		Project:     app.Project.Project,
	}

	log.Trace("looking at latest deployment and release to generate status report on")
	latestDeployment, err := a.state.DeploymentLatest(ctx, appRef, &pb.Ref_Workspace{Workspace: a.workspace})
	// If the deployment isn't found, it's ok
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}
	latestRelease, err := a.state.ReleaseLatest(ctx, appRef, &pb.Ref_Workspace{Workspace: a.workspace})
	// If the release isn't found, it's ok.
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, err
	}
	// Some platforms don't release, so we shouldn't error here if we at least got a deployment
	if latestRelease == nil && latestDeployment == nil {
		log.Warn("no deployment or release found, cannot generate a poll job")
		return nil, nil
	}

	baseJob := &pb.QueueJobRequest{
		Job: &pb.Job{
			Application: &pb.Ref_Application{
				Application: app.Name,
				Project:     app.Project.Project,
			},

			// Application polling requires a data source to be configured for the project
			// Otherwise a status report can't properly eval the project's hcl context
			// needed to query the deploy or release
			DataSource: project.DataSource,

			Workspace: &pb.Ref_Workspace{Workspace: a.workspace},

			// Generate a status report
			Operation: &pb.Job_StatusReport{
				StatusReport: &pb.Job_StatusReportOp{},
			},

			// Any runner is fine for polling.
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Any{
					Any: &pb.Ref_RunnerAny{},
				},
			},
		},

		// TODO define timeout interval based on projects app poll interval
	}
	var jobs []*pb.QueueJobRequest

	log.Trace("Determining which target to generate a status report on")

	// Default to poll on the "latest" lifecycle operation, so if there's a
	// deploy, queue up a status on deploy. If there is latest is release, then queue that too.
	if latestDeployment.Deployment != nil {
		baseJobCopy, err := copystructure.Copy(baseJob)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to generate a job to poll deployment status: %s", err)
		}
		deploymentJob := baseJobCopy.(*pb.QueueJobRequest)
		deploymentJob.Job.Operation = &pb.Job_StatusReport{
			StatusReport: &pb.Job_StatusReportOp{
				Target: &pb.Job_StatusReportOp_Deployment{
					Deployment: latestDeployment,
				},
			},
		}
		// SingletonId so that we only have one poll operation at
		// any time queued per app/operation.
		deploymentJob.Job.SingletonId = appStatusPollSingletonId(a.workspace, app.Project.Project, app.Name, appStatusPollOperationTypeDeployment)

		jobs = append(jobs, deploymentJob)
	}
	if latestRelease != nil && !latestRelease.Unimplemented {
		baseJobCopy, err := copystructure.Copy(baseJob)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to generate a job to poll release status: %s", err)
		}
		releaseJob := baseJobCopy.(*pb.QueueJobRequest)
		releaseJob.Job.Operation = &pb.Job_StatusReport{
			StatusReport: &pb.Job_StatusReportOp{
				Target: &pb.Job_StatusReportOp_Release{
					Release: latestRelease,
				},
			},
		}
		// SingletonId so that we only have one poll operation at
		// any time queued per app/operation.
		releaseJob.Job.SingletonId = appStatusPollSingletonId(a.workspace, app.Project.Project, app.Name, appStatusPollOperationTypeRelease)

		jobs = append(jobs, releaseJob)
	}
	if len(jobs) == 0 {
		// Unclear if we'll even reach this. DeploymentLatest and ReleaseLatest will
		// return an error if there's no deployment or release given an app name.
		log.Debug("no release or deploy target to run a status report poll against.")
	}

	return jobs, nil
}

// Complete will mark the job that was queued as complete, if it
// fails to do so, it will return false with the err to continue the loop
func (a *applicationPoll) Complete(
	log hclog.Logger,
	p interface{},
) error {
	ctx := context.Background()
	project, ok := p.(*pb.Project)
	if !ok || project == nil {
		log.Error("could not mark application poll as complete, incorrect type passed in")
		return status.Error(codes.FailedPrecondition, "incorrect type passed into Application Complete")
	}
	log = log.Named(project.Name)

	// Mark this as complete so the next poll gets rescheduled.
	log.Trace("marking app poll as complete")
	if err := a.state.ApplicationPollComplete(ctx, project, time.Now()); err != nil {
		return err
	}
	return nil
}

// The name of an operation type that status polling is possible for
type appStatusPollOperationType string

const (
	appStatusPollOperationTypeDeployment appStatusPollOperationType = "deployment"
	appStatusPollOperationTypeRelease    appStatusPollOperationType = "release"
)

// appStatusPollSingletonId generates an application status polling job singleton ID
// for the given workspace, project, app and operation type.
// NOTE(briancain): We set a singleton ID for a poll application operation to ensure that the
// poll handler does not fire off many operations of the same kind more than once,
// clogging up the job system. By setting a singleton ID that is unique to this
// application, we can ensure only 1 operation will be active at once rather than
// many operations (such as in the case where a poll interval is shorter than it
// takes to run the operation)
func appStatusPollSingletonId(
	workspaceName string,
	projectName string,
	appName string,
	operationType appStatusPollOperationType,
) string {
	return fmt.Sprintf("app-status-poll/%s/%s/%s/%s", workspaceName, projectName, appName, operationType)
}
