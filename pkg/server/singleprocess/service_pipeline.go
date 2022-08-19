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
	stepGraph, nodeIdMap, err := s.pipelineGraphFull(ctx, log, nil, "",
		make(map[string]string), make(map[string]*pb.Ref_PipelineStep), pipeline)
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

	// Build out all of the queued job requests for running this pipelines steps
	stepJobs, pipelineRun, stepIds, err := s.buildStepJobs(ctx, log, req,
		make(map[string]interface{}), nodeIdMap, pipeline, pipelineRun)
	if err != nil {
		return nil, err
	}

	// Get the ordered jobs.
	var jobIds []string
	jobMap := map[string]*pb.Ref_PipelineStep{}
	for _, v := range stepGraph.KahnSort() {
		// Look up step name and ref by the assigned node ID from graph generation
		nodeId := v.(string)
		stepRef, ok := nodeIdMap[nodeId]
		if !ok {
			return nil, status.Errorf(codes.Internal,
				"could not get pipeline step ref for node id %q", nodeId)
		} else if stepRef == nil {
			return nil, status.Errorf(codes.Internal,
				"node id %q returned a nil pipeline step ref", nodeId)
		}

		// get the generated queued job request
		jobId, ok := stepIds[nodeId]
		if !ok {
			// NOTE(briancain):
			// This could be better. It's basically here because we want to keep track
			// of an embedded pipeline step ref within the step graph, but it doesn't have
			// a "job id" because the root of the embedded pipeline is the actual ID where
			// this is simply a reference.
			// We don't want to add it as a job id because it doesn't actually create a job.
			// jobId = "embedded-pipeline-ref-" + nodeId
			continue
		}

		jobIds = append(jobIds, jobId)
		jobMap[jobId] = stepRef
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
	nodeIdMap map[string]*pb.Ref_PipelineStep,
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
			nodeId, ok := s.stepToNodeId(ctx, log, pipeline.Id, name, nodeIdMap)
			if !ok {
				return nil, nil, nil, status.Errorf(codes.Internal,
					"failed to get node ID from pipeline %q and step name %q",
					pipeline.Id, name)
			}

			var err error
			stepIds[nodeId], err = server.Id()
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
			nodeDepId, ok := s.stepToNodeId(ctx, log, pipeline.Id, dep, nodeIdMap)
			if !ok {
				return nil, nil, nil, status.Errorf(codes.Internal,
					"failed to get node ID from pipeline %q and step name %q",
					pipeline.Id, dep)
			}
			if nodeDepId == "" {
				return nil, nil, status.Errorf(codes.Internal,
					"node ID was blank from pipeline %q and step name %q!!",
					pipeline.Id, dep)
			}

			// TODO(briancain): not sure if this is the right solution here....
			// maybe we gotta look to see if the step is a pipeline ref and don't
			// add it as a DependsOn because it won't have a job id
			// Committing this for now, will be fixed next week.
			d, ok := stepIds[nodeDepId]
			if !ok {
				log.Info("No step id found for nodeDepId", "pipeline", pipeline.Name,
					"step", step.Name, "dep_id", nodeDepId)
				continue
			}
			dependsOn = append(dependsOn, d)
		}

		nodeId, ok := s.stepToNodeId(ctx, log, pipeline.Id, name, nodeIdMap)
		if !ok {
			return nil, nil, nil, status.Errorf(codes.Internal,
				"failed to get node ID from pipeline %q and step name %q",
				pipeline.Id, name)
		}

		job := proto.Clone(req.JobTemplate).(*pb.Job)
		job.Id = stepIds[nodeId]
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
				visitedPipelines, nodeIdMap, embeddedPipeline, pipelineRun)
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
					return nil, nil, nil, status.Errorf(codes.Internal,
						"an embedded pipeline step matches a parent step name: %s", k)
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
// any nested pipeline steps. It keeps track of visisted steps so that we don't
// get stuck in a loop building the graph. This means step names must be unique
// across pipelines.
func (s *Service) pipelineGraphFull(
	ctx context.Context,
	log hclog.Logger,
	g *graph.Graph,
	parentStep string,
	visistedNodes map[string]string,
	nodeIdMap map[string]*pb.Ref_PipelineStep,
	pipeline *pb.Pipeline,
) (*graph.Graph, map[string]*pb.Ref_PipelineStep, error) {
	var stepGraph *graph.Graph
	if g != nil {
		// We're handling an embedded pipeline graph
		stepGraph = g
	} else if stepGraph == nil {
		stepGraph = &graph.Graph{}
	}

	// Note that pipeline.Steps is not an ordered list of steps. It's a map of key val
	// steps so the order will not match the order steps are defined in a waypoint.hcl.
	for _, step := range pipeline.Steps {
		if len(visistedNodes) != 0 {
			if pipeName, ok := visistedNodes[step.Name]; ok && pipeName == pipeline.Name {
				log.Trace("we've cycled to a node we've already visisted!", "pipeline", pipeName, "step", step.Name)
				return nil, nil, status.Error(codes.FailedPrecondition,
					"we've already visisted this node, that means we've got a cycle")
			}
		}

		// Look up step and pipeline id in case we generated the id when we added
		// dependencies from a different step
		nodeId, ok := s.stepToNodeId(ctx, log, pipeline.Id, step.Name, nodeIdMap)
		if !ok {
			var err error
			// unique node graph id for full graph
			nodeId, err = server.Id()
			if err != nil {
				return nil, nil, err
			}
		}

		nodeIdMap[nodeId] = &pb.Ref_PipelineStep{
			Pipeline: pipeline.Id,
			Step:     step.Name,
		}

		// Add our step to the graph as a vertex
		stepGraph.Add(nodeId)

		// Keep track of the fact that we've visited this step node in this pipeline
		// to prevent infinite cycles as we build embedded pipeline graphs
		visistedNodes[step.Name] = pipeline.Name

		if g != nil {
			if parentStep == "" {
				// if we're building an embedded graph and g is not nil but parentStep
				// is then that's an internal error on the caller
				return nil, nil, status.Errorf(codes.Internal,
					"parentStep cannot be empty string if building an embedded graph for pipeline %q!",
					pipeline.Name)
			}

			// Add an edge to the parent step as an implicit dependency
			// Embedded pipeline steps have an implicit dependency on the parent step
			// from the parent pipeline.
			stepGraph.AddEdge(parentStep, nodeId)
		}

		// Add any dependencies as defined by the current Step
		for _, dep := range step.DependsOn {
			// Look up the dependency step by name in case we've already generated
			// a node ID for it for our stepGraph
			depId, ok := s.stepToNodeId(ctx, log, pipeline.Id, dep, nodeIdMap)
			if !ok {
				var err error
				// We haven't reached this node yet, but we're adding it to the graph so
				// we generate an id here too so we can add it to the graph as a vertex
				// and create an edge to the given step with the depencny
				depId, err = server.Id() // unique node graph id for full graph
				if err != nil {
					return nil, nil, err
				}

				// add node id to map
				nodeIdMap[depId] = &pb.Ref_PipelineStep{
					Pipeline: pipeline.Id,
					Step:     dep,
				}
			}

			// Add the dependency as a vertex and draw an edge to *this* steps node ID
			stepGraph.Add(depId)
			stepGraph.AddEdge(depId, nodeId)

			// This only checks for steps inside *this* pipeline. Plain steps cannot
			// reference other steps in other pipelines as a depenency without being
			// an embedded pipeline reference.
			if _, ok := pipeline.Steps[dep]; !ok {
				return nil, nil, fmt.Errorf(
					"step %q depends on non-existent step %q", step, dep)
			}
		}

		// Build a graph for any nested pipelines
		if embedRef, ok := step.Kind.(*pb.Pipeline_Step_Pipeline_); ok {
			// This is only a "ref" to the pipeline, we have to
			// look it up here to get the actual steps.
			embeddedPipeline, err := s.state(ctx).PipelineGet(embedRef.Pipeline.Ref)
			if err != nil {
				return nil, nil, err
			}

			// Build the nested pipelines graph
			parentStep := nodeId
			embeddedGraph, n, err := s.pipelineGraphFull(ctx, log, stepGraph,
				parentStep, visistedNodes, nodeIdMap, embeddedPipeline)
			if err != nil {
				return nil, nil, err
			}

			nodeIdMap = n
			stepGraph = embeddedGraph
		}
	}

	if cycles := stepGraph.Cycles(); len(cycles) > 0 {
		return nil, nil, status.Errorf(codes.FailedPrecondition,
			"step dependencies in pipeline %q contain one or more cycles: %s",
			pipeline.Name, cycles)
	}

	return stepGraph, nodeIdMap, nil
}

func (s *Service) stepToNodeId(
	ctx context.Context,
	log hclog.Logger,
	pipelineName string,
	stepName string,
	nodeIdMap map[string]*pb.Ref_PipelineStep,
) (string, bool) {
	for nodeId, stepRef := range nodeIdMap {
		if stepRef != nil &&
			stepRef.Pipeline == pipelineName && stepRef.Step == stepName {
			return nodeId, true
		}
	}

	return "", false
}
