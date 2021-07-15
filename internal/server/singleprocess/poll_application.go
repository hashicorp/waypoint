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
		log.Warn("error peeking for next application to poll", "err", err)
		return nil, time.Time{}, err // continue loop
	}

	if app != nil {
		log = log.With("application", app.Name)
		log.Trace("returning peek for app")
	} else {
		log.Trace("no application returned from peek")
	}

	return app, pollTime, nil
}

// PollJob will generate a job to queue a project on.
// When determining what to generate a report on, either a Deployment or Release,
// this job assumes that the _Release_ was the last operation taken on the application.
// If there's a release, this job will queue a status report genreation on that.
// Otherwise if there's just a deployment, it return a job to generate a
// status report on the deployment.
func (a *applicationPoll) PollJob(
	log hclog.Logger,
	appl interface{},
) (*pb.QueueJobRequest, error) {
	app, ok := appl.(*pb.Application)
	if !ok || app == nil {
		log.Error("could not generate poll job for application, incorrect type passed in")
		return nil, status.Error(codes.FailedPrecondition, "incorrect type passed into Application PollJob")
	}
	log = log.Named(app.Name)

	// Determine the latest deployment or release to poll for a status report
	appRef := &pb.Ref_Application{
		Application: app.Name,
		Project:     app.Project.Project,
	}

	log.Trace("looking at latest deployment and release to generate status report on")
	latestDeployment, err := a.state.DeploymentLatest(appRef, &pb.Ref_Workspace{Workspace: a.workspace})
	if err != nil {
		return nil, err
	}
	latestRelease, err := a.state.ReleaseLatest(appRef, &pb.Ref_Workspace{Workspace: a.workspace})
	// Some platforms don't release, so we shouldn't error here if we at least got a deployment
	if err != nil && latestDeployment == nil {
		log.Error("no deployment or release found, cannot generate a poll job")
		return nil, err
	}

	statusReportJob := &pb.Job_StatusReport{
		StatusReport: &pb.Job_StatusReportOp{},
	}

	log.Trace("Determining which target to generate a status report on")

	// Default to poll on the "latest" lifecycle operation, so if there's a
	// release, queue up a status on release. If the latest is deploy, then queue that.
	if latestRelease != nil { // ignore me for now
		//TODO(briancain): the issue right now is that waypoint always stores a release,
		// even if the plugin didn't do one. We don't know if it's a real or fake release
		// and we can't run a status report on a fake release especially if the plugin
		// has not implemented a releaser.

		// Returned latestRelease for a plugin that hasn't implemented a releaser
		//panic: application:{application:"web" project:"nginx-project"} workspace:{workspace:"default"} sequence:1 id:"01FANWNRRVVCX9KKE0AV08X3HR" status:{state:SUCCESS start_time:{seconds:1626379838 nanos
		//:234574394} complete_time:{seconds:1626379838 nanos:237627450}} state:CREATED component:{type:PLATFORM name:"docker"} deployment_id:"01FANWNPJHST0XQ7G8GC06PXN0" labels:{key:"waypoint/workspace" va
		//lue:"default"} job_id:"01FANWNJY5NQWYQA00HQ2817P8" preload:{}

		log.Trace("using latest release as a status report target")
		statusReportJob.StatusReport.Target = &pb.Job_StatusReportOp_Release{
			Release: latestRelease,
		}
	} else if latestDeployment.Deployment != nil {
		log.Trace("using latest deployment as a status report target")
		statusReportJob.StatusReport.Target = &pb.Job_StatusReportOp_Deployment{
			Deployment: latestDeployment,
		}
	} else {
		// Unclear if we'll even reach this. DeploymentLatest and ReleaseLatest will
		// return an error if there's no deployment or release given an app name.
		log.Debug("no release or deploy target to run a status report poll against.")
		return nil, nil
	}

	// App polling needs the parent project to obtain its datasource
	project, err := a.state.ProjectGet(&pb.Ref_Project{Project: app.Project.Project})
	if err != nil {
		return nil, err
	}

	// Application polling requires a remote data source, otherwise a status report
	// cannot be generated without a project and its hcl context. This returns
	// an error so we fail early instead of queueing an already broken job
	if project.DataSource == nil {
		log.Debug("cannot poll a job without a remote data source configured.")
		return nil, status.Error(codes.FailedPrecondition, "application polling requires a remote data source")
	}

	log.Trace("building queue job request for generating status report")
	jobRequest := &pb.QueueJobRequest{
		Job: &pb.Job{
			// SingletonId so that we only have one poll operation at
			// any time queued per application.
			SingletonId: fmt.Sprintf("appl-poll/%s", app.Name),

			Application: &pb.Ref_Application{
				Application: app.Name,
				Project:     app.Project.Project,
			},

			// Applicatioon polling requires a data source to be configured for the project
			// Otherwise a status report can't properly eval the projects hcl context
			// needed to query the deploy or release
			DataSource: project.DataSource,

			Workspace: &pb.Ref_Workspace{Workspace: a.workspace},

			// Generate a status report
			Operation: statusReportJob,

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
	log = log.Named(app.Name)

	// Mark this as complete so the next poll gets rescheduled.
	log.Trace("marking app poll as complete")
	if err := a.state.ApplicationPollComplete(app, time.Now()); err != nil {
		return err
	}
	return nil
}
