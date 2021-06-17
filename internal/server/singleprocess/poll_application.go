package singleprocess

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

// applicationPoll accepts a state management interface which provides access
// to a projects current state implementation. Functions like Peek and Complete
// need access to this state interface for peeking at the next available project
// as well as marking a projects poll as complete.
type applicationPoll struct {
	// state is the state management interface that provides functions for
	// safely mutating server state.
	state *state.State

	// the workspace to check for polling applications and running their status on
	workspace string
}

// Peek returns the latest project to poll on
// If there is an error in the ProjectPollPeek, it will return nil
// to allow the outer caller loop to continue and try again
func (a *applicationPoll) Peek(
	log hclog.Logger,
	ws memdb.WatchSet,
) (interface{}, time.Time, error) {
	app, pollTime, err := a.state.ApplicationPollPeek(ws)
	if err != nil {
		return nil, time.Time{}, err // continue loop
	}

	if app != nil {
		log = log.With("application", app.Name)
	}

	return app, pollTime, nil
}

// PollJob will generate a job to queue a project on
func (a *applicationPoll) PollJob(
	log hclog.Logger,
	appl interface{},
) (*pb.QueueJobRequest, error) {
	app, ok := appl.(*pb.Application)
	if !ok || app == nil {
		log.Error("could not generate poll job for application, incorrect type passed in")
		return nil, status.Error(codes.FailedPrecondition, "incorrect type passed into Application PollJob")
	}

	// Determine the latest deployment or release to poll for a status report
	appRef := &pb.Ref_Application{
		Application: app.Name,
		Project:     app.Project.Project,
	}

	latestDeployment, err := a.state.DeploymentLatest(appRef, &pb.Ref_Workspace{Workspace: a.workspace})
	if err != nil {
		return nil, err
	}
	latestRelease, err := a.state.ReleaseLatest(appRef, &pb.Ref_Workspace{Workspace: a.workspace})
	if err != nil {
		return nil, err
	}

	pollOperation := &pb.Job_StatusReport{
		StatusReport: &pb.Job_StatusReportOp{},
	}

	// Default to poll on the "latest" lifecycle operation, so if there's a
	// release, queue up a status on release. If the latest is deploy, then queue that.
	if latestRelease != nil {
		pollOperation.StatusReport.Target = &pb.Job_StatusReportOp_Release{
			Release: latestRelease,
		}
	} else if latestDeployment != nil {
		pollOperation.StatusReport.Target = &pb.Job_StatusReportOp_Deployment{
			Deployment: latestDeployment,
		}
	} else {
		log.Debug("no release or deploy target to run a status report poll against.")
		return nil, nil
	}

	jobRequest := &pb.QueueJobRequest{
		Job: &pb.Job{
			// SingletonId so that we only have one poll operation at
			// any time queued per project.
			SingletonId: fmt.Sprintf("poll/%s", app.Name),

			Application: &pb.Ref_Application{
				Application: app.Name,
				Project:     app.Project.Project,
			},

			Workspace: &pb.Ref_Workspace{Workspace: a.workspace},

			// Poll!
			Operation: pollOperation,

			// Any runner is fine for polling.
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Any{
					Any: &pb.Ref_RunnerAny{},
				},
			},
		},
	}

	return jobRequest, nil
}

// Complete will mark the job that was queued as complete, if it
// fails to do so, it will return false with the err to continue the loop
func (a *applicationPoll) Complete(
	log hclog.Logger,
	appl interface{},
) error {
	app, ok := appl.(*pb.Application)
	if !ok || app == nil {
		log.Error("could not mark application poll as complete, incorrect type passed in")
		return status.Error(codes.FailedPrecondition, "incorrect type passed into Application Complete")
	}

	// Mark this as complete so the next poll gets rescheduled.
	appRef := &pb.Ref_Application{
		Application: app.Name,
		Project:     app.Project.Project,
	}
	if err := a.state.ApplicationPollComplete(appRef, time.Now()); err != nil {
		return err
	}
	return nil
}
