package singleprocess

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *service) UpsertTrigger(
	ctx context.Context,
	req *pb.UpsertTriggerRequest,
) (*pb.UpsertTriggerResponse, error) {
	if err := serverptypes.ValidateUpsertTriggerRequest(req); err != nil {
		return nil, err
	}

	result := req.Trigger
	if err := s.state.TriggerPut(result); err != nil {
		return nil, err
	}

	return &pb.UpsertTriggerResponse{Trigger: result}, nil
}

// GetTrigger returns a Trigger based on ID
func (s *service) GetTrigger(
	ctx context.Context,
	req *pb.GetTriggerRequest,
) (*pb.GetTriggerResponse, error) {
	if err := serverptypes.ValidateGetTriggerRequest(req); err != nil {
		return nil, err
	}

	t, err := s.state.TriggerGet(req.Ref)
	if err != nil {
		return nil, err
	}

	return &pb.GetTriggerResponse{Trigger: t}, nil
}

// DeleteTrigger deletes a Trigger based on ID
func (s *service) DeleteTrigger(
	ctx context.Context,
	req *pb.DeleteTriggerRequest,
) (*empty.Empty, error) {
	if err := serverptypes.ValidateDeleteTriggerRequest(req); err != nil {
		return nil, err
	}

	err := s.state.TriggerDelete(req.Ref)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *service) ListTriggers(
	ctx context.Context,
	req *pb.ListTriggerRequest,
) (*pb.ListTriggerResponse, error) {
	// NOTE: no ptype validation at the moment, as all Ref fields are optional

	result, err := s.state.TriggerList(req.Workspace, req.Project, req.Application, req.Tags)
	if err != nil {
		return nil, err
	}

	return &pb.ListTriggerResponse{Triggers: result}, nil
}

func (s *service) RunTrigger(
	ctx context.Context,
	req *pb.RunTriggerRequest,
) (*pb.RunTriggerResponse, error) {
	if err := serverptypes.ValidateRunTriggerRequest(req); err != nil {
		return nil, err
	}
	log := hclog.FromContext(ctx)

	runTrigger, err := s.state.TriggerGet(req.Ref)
	if err != nil {
		return nil, err
	}

	log = log.With("run_trigger", runTrigger.Id)

	log.Debug("building run trigger job")

	// Build the job(s)
	job := &pb.Job{
		Workspace: runTrigger.Workspace,
	}

	// TODO is there an easy way to convert this without the big switch
	switch op := runTrigger.Operation.(type) {
	case *pb.Trigger_Build:
		job.Operation = &pb.Job_Build{Build: op.Build}
	case *pb.Trigger_Push:
		job.Operation = &pb.Job_Push{Push: op.Push}
	case *pb.Trigger_Deploy:
		job.Operation = &pb.Job_Deploy{Deploy: op.Deploy}
	case *pb.Trigger_Destroy:
		job.Operation = &pb.Job_Destroy{Destroy: op.Destroy}
	case *pb.Trigger_Release:
		job.Operation = &pb.Job_Release{Release: op.Release}
	case *pb.Trigger_Up:
		job.Operation = &pb.Job_Up{Up: op.Up}
	case *pb.Trigger_Init:
		job.Operation = &pb.Job_Init{Init: op.Init}
	default:
		return nil, status.Errorf(codes.Internal,
			"trigger %q is configured with an unsupported operation %T", runTrigger.Id, op)
	}

	// TODO: Config Variable overrides?

	// TODO(briancain): look up a target runner config at the project/app level and apply it to job requests
	job.TargetRunner = &pb.Ref_Runner{Target: &pb.Ref_Runner_Any{}}

	// generate job requests
	var jobList []*pb.QueueJobRequest
	var ids []string
	if runTrigger.Application == nil {
		// we're gonna queue multiple jobs for every application in a project
		project, err := s.state.ProjectGet(runTrigger.Project)
		if err != nil {
			return nil, err
		}

		log.Debug("building multi-jobs for all apps in project", "project", project.Name)
		for _, app := range project.Applications {
			tempJob := &pb.Job{
				Workspace:    job.Workspace,
				Operation:    job.Operation,
				TargetRunner: job.TargetRunner,
			}

			tempJob.Application = &pb.Ref_Application{
				Project:     project.Name,
				Application: app.Name,
			}

			jobReq := &pb.QueueJobRequest{Job: tempJob}
			jobList = append(jobList, jobReq)
		}

		// Queue the job(s)
		log.Debug("queueing jobs", "total_jobs", len(jobList))
		respList, err := s.queueJobMulti(ctx, jobList)
		if err != nil {
			return nil, err
		}
		// Gather queue job request ids
		for _, qJr := range respList {
			ids = append(ids, qJr.JobId)
		}
	} else {
		log.Debug("building a single job for target", "project",
			runTrigger.Application.Project, "app", runTrigger.Application.Application)
		// we're only targetting a specific application, so queue 1 job
		job.Application = runTrigger.Application
		j := &pb.QueueJobRequest{Job: job}
		jobList = append(jobList, j)

		resp, err := s.QueueJob(ctx, j)
		if err != nil {
			return nil, err
		}
		ids = append(ids, resp.JobId)
	}

	log.Debug("run trigger job(s) have been queued")

	// Trigger has been requested to queue jobs, update active time
	runTrigger.ActiveTime = timestamppb.New(time.Now())
	err = s.state.TriggerPut(runTrigger)
	if err != nil {
		return nil, err
	}

	// maybe update to return array of RunTriggerResponses instead?
	return &pb.RunTriggerResponse{JobIds: ids}, nil
}
