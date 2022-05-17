package singleprocess

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) UpsertPipeline(
	ctx context.Context,
	req *pb.UpsertPipelineRequest,
) (*pb.UpsertPipelineResponse, error) {
	if err := serverptypes.ValidateUpsertPipelineRequest(req); err != nil {
		return nil, err
	}

	result := req.Pipeline
	if err := s.state(ctx).PipelinePut(result); err != nil {
		return nil, err
	}

	return &pb.UpsertPipelineResponse{Pipeline: result}, nil
}

func (s *Service) ListPipelines(
	ctx context.Context,
	req *pb.ListPipelinesRequest,
) (*pb.ListPipelinesResponse, error) {
	if err := serverptypes.ValidateListPipelinesRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state(ctx).PipelineList(req.Project)
	if err != nil {
		return nil, err
	}

	return &pb.ListPipelinesResponse{
		Pipelines: result,
	}, nil
}

func (s *Service) RunPipeline(
	ctx context.Context,
	req *pb.RunPipelineRequest,
) (*pb.RunPipelineResponse, error) {
	if err := serverptypes.ValidateRunPipelineRequest(req); err != nil {
		return nil, err
	}

	// Get the pipeline we should execute
	pipeline, err := s.state(ctx).PipelineGet(req.Pipeline)
	if err != nil {
		return nil, err
	}

	// Generate job IDs for each of the steps. We need to know the IDs in
	// advance to setup the dependency chain.
	stepIds := map[string]string{}
	for name := range pipeline.Steps {
		stepIds[name], err = server.Id()
		if err != nil {
			return nil, err
		}
	}

	// Generate the jobs for each of the steps
	var stepJobs []*pb.QueueJobRequest
	for name, step := range pipeline.Steps {
		var dependsOn []string
		for _, dep := range step.DependsOn {
			dependsOn = append(dependsOn, stepIds[dep])
		}

		job := proto.Clone(req.JobTemplate).(*pb.Job)
		job.Id = stepIds[name]
		job.DependsOn = append(job.DependsOn, dependsOn...)
		job.Operation = &pb.Job_PipelineStep{
			PipelineStep: &pb.Job_PipelineStepOp{
				Step: step,
			},
		}

		stepJobs = append(stepJobs, &pb.QueueJobRequest{
			Job: job,
		})
	}

	// Get the graph for the steps so we can get the root. We enforce a
	// single root so the root is always the first step.
	stepGraph, err := serverptypes.PipelineGraph(pipeline)
	if err != nil {
		return nil, err
	}

	// Get the ordered jobs.
	var jobIds []string
	jobMap := map[string]*pb.Ref_PipelineStep{}
	for _, v := range stepGraph.KahnSort() {
		jobId := stepIds[v.(string)]
		jobIds = append(jobIds, jobId)
		jobMap[jobId] = &pb.Ref_PipelineStep{
			Pipeline: pipeline.Id,
			Step:     v.(string),
		}
	}

	// Queue all the jobs atomically
	if _, err := s.queueJobMulti(ctx, stepJobs); err != nil {
		return nil, err
	}

	return &pb.RunPipelineResponse{
		JobId:     jobIds[0],
		AllJobIds: jobIds,
		JobMap:    jobMap,
	}, nil
}
