package singleprocess

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

// projectPoll accepts a state management interface which provides access
// to a projects current state implementation. Functions like Peek and Complete
// need access to this state interface for peeking at the next available project
// as well as marking a projects poll as complete.
type projectPoll struct {
	// state is the state management interface that provides functions for
	// safely mutating server state.
	state serverstate.Interface
}

// Peek returns the latest project to poll on
// If there is an error in the ProjectPollPeek, it will return nil
// to allow the outer caller loop to continue and try again
func (pp *projectPoll) Peek(
	log hclog.Logger,
	ws memdb.WatchSet,
) (interface{}, time.Time, error) {
	p, pollTime, err := pp.state.ProjectPollPeek(context.Background(), ws)
	if err != nil {
		return nil, time.Time{}, err // continue loop
	}

	if p != nil {
		log = log.With("project", p.Name)
	}

	return p, pollTime, nil
}

// PollJob will generate a job to queue a project on
func (pp *projectPoll) PollJob(
	log hclog.Logger,
	project interface{},
) ([]*pb.QueueJobRequest, error) {
	p, ok := project.(*pb.Project)
	if !ok || p == nil {
		log.Error("could not generate poll job for project, incorrect type passed in")
		return nil, status.Error(codes.FailedPrecondition, "incorrect type passed into Project PollJob")
	}

	var jobList []*pb.QueueJobRequest

	jobRequest := &pb.QueueJobRequest{
		Job: &pb.Job{
			// SingletonId so that we only have one poll operation at
			// any time queued per project.
			SingletonId: fmt.Sprintf("poll/%s", p.Name),

			// Project polling requires a data source to be configured for the project
			DataSource: p.DataSource,

			Application: &pb.Ref_Application{
				Project: p.Name,
				// No Application set since PollOp is project-oriented
			},

			// Polling always happens on the default workspace even
			// though the PollOp is across every workspace.
			Workspace: &pb.Ref_Workspace{Workspace: "default"},

			// Poll!
			Operation: &pb.Job_Poll{
				Poll: &pb.Job_PollOp{},
			},

			// Any runner is fine for polling.
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Any{
					Any: &pb.Ref_RunnerAny{},
				},
			},
		},
	}

	jobList = append(jobList, jobRequest)

	return jobList, nil
}

// Complete will mark the job that was queued as complete, if it
// fails to do so, it will return false with the err to continue the loop
func (pp *projectPoll) Complete(
	log hclog.Logger,
	project interface{},
) error {
	p, ok := project.(*pb.Project)
	if !ok || p == nil {
		log.Error("could not mark project poll as complete, incorrect type passed in")
		return status.Error(codes.FailedPrecondition, "incorrect type passed into Project Complete")
	}

	// Mark this as complete so the next poll gets rescheduled.
	if err := pp.state.ProjectPollComplete(context.Background(), p, time.Now()); err != nil {
		return err
	}
	return nil
}
