package singleprocess

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/waypoint/internal/pkg/graph"
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

// GetPipeline returns a Pipeline based on a pipeline ref id
func (s *Service) GetPipeline(
	ctx context.Context,
	req *pb.GetPipelineRequest,
) (*pb.GetPipelineResponse, error) {
	if err := serverptypes.ValidateGetPipelineRequest(req); err != nil {
		return nil, err
	}

	p, err := s.state(ctx).PipelineGet(req.Pipeline)
	if err != nil {
		return nil, err
	}

	// Get the graph for the steps so we can get the root. We enforce a
	// single root so the root is always the first step.
	stepGraph, err := serverptypes.PipelineGraph(p)
	if err != nil {
		return nil, err
	}

	orderedStep := stepGraph.KahnSort()
	rootStepName := orderedStep[0].(string)

	return &pb.GetPipelineResponse{
		Pipeline: p,
		RootStep: rootStepName,
		// TODO: Leaving this out intentionally for now, need to convert stepGraph into mermaid
		//Graph:    graph,
	}, nil
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
	log := hclog.FromContext(ctx)
	if err := serverptypes.ValidateRunPipelineRequest(req); err != nil {
		return nil, err
	}

	// Get the pipeline we should execute
	pipeline, err := s.state(ctx).PipelineGet(req.Pipeline)
	if err != nil {
		return nil, err
	}

	// Get the graph for the steps so we can get the root. We enforce a
	// single root so the root is always the first step.
	//stepGraph, err := serverptypes.PipelineGraph(pipeline) // the old pipeline graph func. Does not account for nested pipes
	stepGraph, err := s.pipelineGraphFull(ctx, log, nil, "", make(map[string]interface{}), pipeline)
	if err != nil {
		return nil, err
	}

	// Initialize a pipeline run
	if err = s.state(ctx).PipelineRunPut(&pb.PipelineRun{
		Pipeline: &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: &pb.Ref_PipelineId{
					Id: pipeline.Id,
				},
			},
		},
		State: pb.PipelineRun_PENDING,
	}); err != nil {
		return nil, err
	}

	pipelineRun, err := s.state(ctx).PipelineRunGetLatest(pipeline.Id)
	if err != nil {
		return nil, err
	}

	stepJobs, pipelineRun, stepIds, err := s.buildStepJobs(ctx, log, req, make(map[string]interface{}), pipeline, pipelineRun)
	if err != nil {
		return nil, err
	}

	// Get the ordered jobs.
	// TODO(briancain): probably move this to buildStepJobs to get the proper Pipeline Ids
	var jobIds []string
	jobMap := map[string]*pb.Ref_PipelineStep{}
	for _, v := range stepGraph.KahnSort() {
		// NOTE(briancain):
		// This could be better. It's basically here because we want to keep track
		// of an embedded pipeline step ref within the step graph, but it doesn't have
		// a "job id" because the root of the embedded pipeline is the actual ID where
		// this is simply a reference
		jobId, ok := stepIds[v.(string)]
		if !ok {
			jobId = "embedded-pipeline-ref-" + v.(string)
		}

		jobIds = append(jobIds, jobId)
		jobMap[jobId] = &pb.Ref_PipelineStep{
			Pipeline:    pipeline.Id, // TODO(briancain): Fix me for embedded pipelines
			Step:        v.(string),
			RunSequence: pipelineRun.Sequence,
		}
	}

	pipelineRun.State = pb.PipelineRun_STARTING
	if err = s.state(ctx).PipelineRunPut(pipelineRun); err != nil {
		return nil, err
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

func (s *Service) buildStepJobs(
	ctx context.Context,
	log hclog.Logger,
	req *pb.RunPipelineRequest,
	visitedPipelines map[string]interface{},
	pipeline *pb.Pipeline,
	pipelineRun *pb.PipelineRun,
) ([]*pb.QueueJobRequest, *pb.PipelineRun, map[string]string, error) {
	if len(visitedPipelines) != 0 {
		// Determine if we've already visisted this pipeline and included its jobs.
		// Otherwise we'll get stuck in a cycle. This only really works because
		// pipeline names are unique for a project. If we ever start allowing for
		// pipelines across projects we'll need to namespace this value.
		if _, ok := visitedPipelines[pipeline.Name]; ok {
			return nil, nil, nil, nil
		}
	}

	// Mark that we've visited this pipeline already
	visitedPipelines[pipeline.Name] = struct{}{}

	// Generate job IDs for each of the steps. We need to know the IDs in
	// advance to setup the dependency chain.
	stepIds := map[string]string{}
	for name, step := range pipeline.Steps {
		if _, ok := step.Kind.(*pb.Pipeline_Step_Pipeline_); !ok {
			// TODO(briancain) should step id keys be project/pipeline/name to avoid embedded collisions?
			var err error
			stepIds[name], err = server.Id()
			if err != nil {
				return nil, nil, nil, err
			}
		}
	}

	// For every step in a pipeline, generate the job and the kind of operation
	// the job should run based on the pipeline proto config.
	var stepJobs []*pb.QueueJobRequest
	for name, step := range pipeline.Steps {
		var dependsOn []string
		for _, dep := range step.DependsOn {
			dependsOn = append(dependsOn, stepIds[dep])
		}

		job := proto.Clone(req.JobTemplate).(*pb.Job)
		job.Id = stepIds[name]
		job.DependsOn = append(job.DependsOn, dependsOn...)
		job.Pipeline = &pb.Ref_PipelineStep{
			Pipeline:    pipeline.Id,
			Step:        step.Name,
			RunSequence: pipelineRun.Sequence,
		}

		// Queue the right job depending on the Step type. We will queue a Waypoint
		// operation if the type is a reserved built in step.
		switch o := step.Kind.(type) {
		case *pb.Pipeline_Step_Build_:
			job.Operation = &pb.Job_Build{
				Build: &pb.Job_BuildOp{
					DisablePush: o.Build.DisablePush,
				},
			}
		case *pb.Pipeline_Step_Deploy_:
			job.Operation = &pb.Job_Deploy{
				Deploy: &pb.Job_DeployOp{},
			}

			if o.Deploy.Release {
				// NOTE(briancain): Unclear if this is really a behavior we want to
				// encourage, unlike the CLI. If users want to release after a deploy
				// they can just add a Release step to their pipeline

				// Queue a release job too
				log.Warn("Waypoint server current does not support queueing an automatic " +
					"release job via a deploy....Sorry!!!")
			}
		case *pb.Pipeline_Step_Release_:
			var deployment *pb.Deployment
			if o.Release.Deployment != nil {
				switch d := o.Release.Deployment.Ref.(type) {
				case *pb.Ref_Deployment_Latest:
					// Nothing, keep the Deployment proto nil
					log.Trace("using nil deployment to queue job, which is latest deployment")
				case *pb.Ref_Deployment_Sequence:
					// Look up deployment sequence here and set proto?
					deployment, err := s.GetDeployment(ctx, &pb.GetDeploymentRequest{
						Ref: &pb.Ref_Operation{
							Target: &pb.Ref_Operation_Sequence{
								Sequence: &pb.Ref_OperationSeq{
									Application: job.Application,
									Number:      d.Sequence,
								},
							},
						},
						LoadDetails: pb.Deployment_ARTIFACT,
					})
					if err != nil {
						return nil, nil, nil, err
					}
					if deployment == nil {
						log.Debug("could not find deploy sequence, using latest instead",
							"seq", d.Sequence)
					}
				default:
					// return an error
					log.Error("invalid deployment ref received", "ref", d)
					return nil, nil, nil, status.Errorf(codes.Internal,
						"invalid deployment ref received: %T", d)
				}
			}

			job.Operation = &pb.Job_Release{
				Release: &pb.Job_ReleaseOp{
					Deployment:          deployment,
					Prune:               o.Release.Prune,
					PruneRetain:         o.Release.PruneRetain,
					PruneRetainOverride: o.Release.PruneRetainOverride,
				},
			}
		case *pb.Pipeline_Step_Up_:
			job.Operation = &pb.Job_Up{
				Up: &pb.Job_UpOp{
					Release: &pb.Job_ReleaseOp{
						Prune:               o.Up.Prune,
						PruneRetain:         o.Up.PruneRetain,
						PruneRetainOverride: o.Up.PruneRetainOverride,
					},
				},
			}
		case *pb.Pipeline_Step_Pipeline_:
			embeddedPipeline, err := s.state(ctx).PipelineGet(o.Pipeline.Ref)
			if err != nil {
				return nil, nil, nil, err
			}

			embedJobs, embedRun, embedStepIds, err := s.buildStepJobs(ctx, log, req,
				visitedPipelines, embeddedPipeline, pipelineRun)
			if err != nil {
				return nil, nil, nil, err
			}

			// add the nested jobs
			stepJobs = append(stepJobs, embedJobs...)
			pipelineRun = embedRun

			// Include nested pipeline steps in stepId map
			for k, v := range embedStepIds {
				if _, ok := stepIds[k]; !ok {
					stepIds[k] = v
				} else {
					// Embedded pipeline steps match an existing step id
					return nil, nil, nil, status.Errorf(codes.Internal, "an embedded pipeline step matches a parent step name: %s", k)
				}
			}
		default:
			job.Operation = &pb.Job_PipelineStep{
				PipelineStep: &pb.Job_PipelineStepOp{
					Step: step,
				},
			}
		}

		pipelineRun.Jobs = append(pipelineRun.Jobs, &pb.Ref_Job{Id: job.Id})
		stepJobs = append(stepJobs, &pb.QueueJobRequest{
			Job: job,
		})
	}

	return stepJobs, pipelineRun, stepIds, nil
}

// pipelineGraphFull takes a pipeline, and optionally accepts an existing Graph
// and parent step, and attempts to build a full graph for a given Pipeline including
// any nested pipeline steps.
func (s *Service) pipelineGraphFull(
	ctx context.Context,
	log hclog.Logger,
	g *graph.Graph,
	parentStep string,
	visistedNodes map[string]interface{},
	v *pb.Pipeline,
) (*graph.Graph, error) {
	var stepGraph *graph.Graph
	if g != nil {
		stepGraph = g
	} else if stepGraph == nil {
		stepGraph = &graph.Graph{}
	}

	// Note that v.Steps is not an ordered list of steps. It's a map of key val
	// steps so the order will not match the order steps are defined.
	for _, step := range v.Steps {
		if len(visistedNodes) != 0 {
			if _, ok := visistedNodes[step.Name]; ok {
				// we've been here
				continue
			}
		}
		// Add our job
		stepGraph.Add(step.Name)

		// Keep track of the fact that we've visited this step node to prevent
		// infinite cycles as we build embedded pipeline graphs
		visistedNodes[step.Name] = struct{}{}

		if g != nil {
			if parentStep == "" {
				return nil, status.Error(codes.FailedPrecondition,
					"parentStep cannot be empty string")
			}

			// Add an edge to the parent step as a dependency
			stepGraph.AddEdge(parentStep, step.Name)
		}

		// Add any dependencies
		for _, dep := range step.DependsOn {
			stepGraph.Add(dep)
			stepGraph.AddEdge(dep, step.Name)

			if _, ok := v.Steps[dep]; !ok {
				return nil, fmt.Errorf(
					"step %q depends on non-existent step %q", step, dep)
			}
		}

		// Build a graph for any nested pipelines
		if embedRef, ok := step.Kind.(*pb.Pipeline_Step_Pipeline_); ok {
			// This is only a "ref" to the pipeline, we have to
			// look it up here to get the actual steps.
			embeddedPipeline, err := s.state(ctx).PipelineGet(embedRef.Pipeline.Ref)
			if err != nil {
				return nil, err
			}

			// NOTE(briancain):
			// One issue with this is any nested pipeline Step names can't be the same
			// as parent step names. We should namespace step names by the pipeline
			// they are in. Pipeline mames are unique within a project. If we ever start
			// embedding pipelines across projects we'll have to figure out something else

			// build the nested pipelines graph
			parentStep := step.Name
			embeddedGraph, err := s.pipelineGraphFull(ctx, log, stepGraph,
				parentStep, visistedNodes, embeddedPipeline)
			if err != nil {
				return nil, err
			}

			stepGraph = embeddedGraph
		}
	}

	if cycles := stepGraph.Cycles(); len(cycles) > 0 {
		return nil, fmt.Errorf(
			"step dependencies contain one or more cycles: %s", cycles)
	}

	return stepGraph, nil
}
