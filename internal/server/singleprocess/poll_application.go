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
	project, pollTime, err := a.state.ApplicationPollPeek(ws)
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

// GeneratePollJobs will generate a QueuedJobRequest to generate a status report
// for each application defined in the given project.
func (a *applicationPoll) GeneratePollJobs(
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
		job, err := a.PollJob(log, app)
		if err != nil {
			return nil, err
		}

		jobList = append(jobList, job)
	}

	return jobList, nil
}

// PollJob will generate a job to queue a project on.
// When determining what to generate a report on, either a Deployment or Release,
// this job assumes that the _Release_ was the last operation taken on the application.
// If there's a release, this job will queue a status report genreation on that.
// Otherwise if there's just a deployment, it returns a job to generate a
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
	if latestRelease != nil && !latestRelease.Unimplemented {
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

		// TODO define timeout interval based on projects app poll interval
	}

	return jobRequest, nil
}

// Complete will mark the job that was queued as complete, if it
// fails to do so, it will return false with the err to continue the loop
func (a *applicationPoll) Complete(
	log hclog.Logger,
	p interface{},
) error {
	project, ok := p.(*pb.Project)
	if !ok || project == nil {
		log.Error("could not mark application poll as complete, incorrect type passed in")
		return status.Error(codes.FailedPrecondition, "incorrect type passed into Application Complete")
	}
	log = log.Named(project.Name)

	// Mark this as complete so the next poll gets rescheduled.
	log.Trace("marking app poll as complete")
	if err := a.state.ApplicationPollComplete(project, time.Now()); err != nil {
		return err
	}
	return nil
}
