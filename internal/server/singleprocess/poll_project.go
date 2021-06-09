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

type projectPoll struct {
	// state is the state management interface that provides functions for
	// safely mutating server state.
	state *state.State
}

// Peek returns the latest project to poll on
// If there is an error in the ProjectPollPeek, it will return nil
// to allow the outer caller loop to continue and try again
func (pp *projectPoll) Peek(
	ws memdb.WatchSet,
	log hclog.Logger,
) (interface{}, time.Time, error) {
	p, pollTime, err := pp.state.ProjectPollPeek(ws)
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
	project interface{},
	log hclog.Logger,
) (*pb.QueueJobRequest, error) {
	p, ok := project.(*pb.Project)
	if !ok {
		log.Error("could not generate poll job for project, incorrect type passed in")
		return nil, status.Error(codes.FailedPrecondition, "incorrect type passed into Project PollJob")
	}

	jobRequest := &pb.QueueJobRequest{
		Job: &pb.Job{
			// SingletonId so that we only have one poll operation at
			// any time queued per project.
			SingletonId: fmt.Sprintf("poll/%s", p.Name),

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

	return jobRequest, nil
}

// Complete will mark the job that was queued as complete, if it
// fails to do so, it will return false with the err to continue the loop
func (pp *projectPoll) Complete(
	project interface{},
	log hclog.Logger,
) error {
	p, ok := project.(*pb.Project)
	if !ok {
		log.Error("could not mark project poll as complete, incorrect type passed in")
		return status.Error(codes.FailedPrecondition, "incorrect type passed into Project Complete")
	}

	// Mark this as complete so the next poll gets rescheduled.
	if err := pp.state.ProjectPollComplete(p, time.Now()); err != nil {
		return err
	}
	return nil
}
