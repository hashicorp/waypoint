// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package singleprocess

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/internal/pkg/graph"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
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
	if err := s.state(ctx).PipelinePut(ctx, result); err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error upserting pipeline",
		)
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

	p, err := s.state(ctx).PipelineGet(ctx, req.Pipeline)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error getting pipeline",
		)
	}

	// Get the graph for the steps so we can get the root. We enforce a
	// single root so the root is always the first step.
	stepGraph, err := serverptypes.PipelineGraph(p)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error generating pipline graph",
			"pipeline_id",
			p.Id,
		)
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

	result, err := s.state(ctx).PipelineList(ctx, req.Project)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error listing piplines",
		)
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
	pipeline, err := s.state(ctx).PipelineGet(ctx, req.Pipeline)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error getting pipeline",
		)
	}

	// Get the graph for the steps so we can get the root. We enforce a
	// single root so the root is always the first step.
	stepGraph, nodeToStepRef, err := s.pipelineGraphFull(ctx, log, nil, "", nil,
		make(map[string]string), nil, pipeline)
	if err != nil {
		log.Error("server failed to build full pipeline graph to determine cycles", "err", err)
		return nil, hcerr.Externalize(
			log,
			fmt.Errorf("server failed to build full pipeline graph to determine cycles: %w", err),
			"server failed to build full pipeline graph to determine cycles",
		)
	}

	// Initialize a pipeline run
	if err = s.state(ctx).PipelineRunPut(ctx, &pb.PipelineRun{
		Pipeline: &pb.Ref_Pipeline{
			Ref: &pb.Ref_Pipeline_Id{
				Id: pipeline.Id,
			},
		},
		State: pb.PipelineRun_PENDING,
	}); err != nil {
		return nil, hcerr.Externalize(
			log,
			err,
			"error initializing pipeline run",
			"pipeline_id",
			pipeline.Id,
		)
	}

	pipelineRun, err := s.state(ctx).PipelineRunGetLatest(ctx, pipeline.Id)
	if err != nil {
		return nil, hcerr.Externalize(
			log,
			err,
			"error getting latest pipeline run",
			"pipeline_id",
			pipeline.Id,
		)
	}

	// Build out all of the queued job requests for running this pipeline's steps
	stepJobs, pipelineRun, stepIds, err := s.buildStepJobs(ctx, log, req,
		make(map[string]interface{}), nodeToStepRef, make(map[string][]string), pipeline, pipelineRun, pipeline)
	if err != nil {
		return nil, hcerr.Externalize(
			log,
			err,
			"error building jobs for pipeline",
			"pipeline_id",
			pipeline.Id,
		)
	}

	// Get the ordered jobs.
	var jobIds []string
	jobMap := map[string]*pb.Ref_PipelineStep{}
	order := stepGraph.KahnSort()

	for _, v := range order {
		// Look up step name and ref by the assigned node ID from graph generation
		nodeId := v.(string)
		stepRef, ok := nodeToStepRef.nodeStepRefs[nodeId]
		if !ok {
			return nil, hcerr.Externalize(
				log,
				fmt.Errorf("could not get pipeline step ref for node id: %q", nodeId),
				"error getting pipeline step",
				"pipeline_id",
				pipeline.Id,
			)
		} else if stepRef == nil {
			return nil, hcerr.Externalize(
				log,
				fmt.Errorf("node id %q returned a nil pipeline step ref", nodeId),
				"error getting pipeline steps reference",
				"pipeline_id",
				pipeline.Id,
			)
		}

		// get the generated queued job request
		jobId, ok := stepIds[nodeId]
		if !ok {
			// NOTE(briancain):
			// This could be better. It's basically here because we want to keep track
			// of an embedded pipeline step ref within the step graph, but it doesn't have
			// a "job id" because the root of the embedded pipeline is the actual ID whereas
			// this is simply a reference.
			// We don't want to add it as a job id because it doesn't actually create a job.
			// jobId = "embedded-pipeline-ref-" + nodeId
			continue
		}

		jobIds = append(jobIds, jobId)
		stepRef.RunSequence = pipelineRun.Sequence
		jobMap[jobId] = stepRef
	}

	pipelineRun.State = pb.PipelineRun_STARTING
	if err = s.state(ctx).PipelineRunPut(ctx, pipelineRun); err != nil {
		return nil, hcerr.Externalize(
			log,
			err,
			"error updating pipeline to starting state",
			"pipeline_id",
			pipeline.Id,
		)
	}

	// Queue all the jobs atomically
	if _, err := s.queueJobMulti(ctx, stepJobs); err != nil {
		return nil, hcerr.Externalize(
			log,
			err,
			"error queueing jobs for pipeline",
			"pipeline_id",
			pipeline.Id,
		)
	}

	return &pb.RunPipelineResponse{
		JobId:     jobIds[0],
		AllJobIds: jobIds,
		JobMap:    jobMap,
		Sequence:  pipelineRun.Sequence,
	}, nil
}

// nodeToStepRef is a helper struct used in both buildStepJobs and pipelineGraphFull.
// Its job is to keep track of each vertex in a pipeline graph where its node
// id is unique to the full graph and is backed by a specific Pipeline Step that
// the node Id is referencing.
type nodeToStepRef struct {
	// Map[NodeId] Pipeline Step Ref
	nodeStepRefs map[string]*pb.Ref_PipelineStep

	// Map[Pipeline + Step] Node Id
	// Use the value of this, to look up the Step Ref in nodeStepRefs
	stepRefs map[nodePipelineStepRef]string
}

// nodePipelineStepRef is a simple struct that is *like* a pb.Ref_PipelineStep.
// We include this as its own struct because it's not safe to compare protobuf
// struct pointers within a maps key.
type nodePipelineStepRef struct {
	pipeline, step string
}

func (s *Service) buildStepJobs(
	ctx context.Context,
	log hclog.Logger,
	req *pb.RunPipelineRequest,
	visitedPipelines map[string]interface{},
	nodeStepRef *nodeToStepRef,
	parentDep map[string][]string,
	pipeline *pb.Pipeline,
	pipelineRun *pb.PipelineRun,
	rootPipeline *pb.Pipeline,
) ([]*pb.QueueJobRequest, *pb.PipelineRun, map[string]string, error) {
	if len(visitedPipelines) != 0 {
		// Determine if we've already visited this pipeline and included its jobs.
		// Otherwise, we'll get stuck in a cycle. This only really works because
		// pipeline names are unique for a project. If we ever allow for
		// pipelines across projects, we'll need to namespace this value for find
		// some other way of tracking our visisted pipelines for the job builder.
		if _, ok := visitedPipelines[pipeline.Name]; ok {
			return nil, nil, nil, nil
		}
	}

	// Mark that we've visited this pipeline already
	visitedPipelines[pipeline.Name] = struct{}{}

	// Generate job IDs for each of the steps. We need to know the IDs in
	// advance to set up the dependency chain.
	stepIds := map[string]string{}
	for name := range pipeline.Steps {
		nodeId, ok := nodeStepRef.stepRefs[nodePipelineStepRef{pipeline: pipeline.Id, step: name}]
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

	// For every step in a pipeline, generate the job and the kind of operation
	// the job should run based on the pipeline proto config.
	var stepJobs []*pb.QueueJobRequest
	for name, step := range pipeline.Steps {
		var dependsOn []string
		for _, dep := range step.DependsOn {
			nodeDepId, ok := nodeStepRef.stepRefs[nodePipelineStepRef{pipeline: pipeline.Id, step: dep}]
			if !ok {
				return nil, nil, nil, status.Errorf(codes.Internal,
					"failed to get node ID from pipeline %q and step name %q",
					pipeline.Id, dep)
			}
			if nodeDepId == "" {
				return nil, nil, nil, status.Errorf(codes.Internal,
					"node ID was blank from pipeline %q and step name %q!!",
					pipeline.Id, dep)
			}

			d, ok := stepIds[nodeDepId]
			if !ok {
				log.Error("No step id found for nodeDepId", "pipeline", pipeline.Name,
					"step", step.Name, "dep_id", nodeDepId)
				return nil, nil, nil, status.Errorf(codes.Internal,
					"no step ID was found for nodeDepId from pipeline %q and step name %q!!",
					pipeline.Name, nodeDepId)
			}

			dependsOn = append(dependsOn, d)
		}

		// Depend on all job IDs from parent step. If we were given a parent pipeline
		// with its step dependencies we're in an embedded pipeline and need to ensure
		// the downstream steps have an implicit dependency on the parent embedded
		// pipeline Ref step
		for _, stepDepends := range parentDep {
			dependsOn = append(dependsOn, stepDepends...)
		}

		nodeId, ok := nodeStepRef.stepRefs[nodePipelineStepRef{pipeline: pipeline.Id, step: name}]
		if !ok {
			return nil, nil, nil, status.Errorf(codes.Internal,
				"failed to get node ID from pipeline %q and step name %q",
				pipeline.Id, name)
		}

		job := proto.Clone(req.JobTemplate).(*pb.Job)
		job.Id = stepIds[nodeId]
		job.DependsOn = append(job.DependsOn, dependsOn...)
		job.Pipeline = &pb.Ref_PipelineStep{
			PipelineId:       pipeline.Id,
			PipelineName:     pipeline.Name,
			Step:             step.Name,
			RootPipelineId:   rootPipeline.Id,
			RootPipelineName: rootPipeline.Name,
			RunSequence:      pipelineRun.Sequence,
		}

		// step has a specific workspace set, update the job to use that
		// workspace
		if step.Workspace != nil {
			job.Workspace = &pb.Ref_Workspace{
				Workspace: step.Workspace.Workspace,
			}
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
			embeddedPipeline, err := s.state(ctx).PipelineGet(ctx, o.Pipeline.Ref)
			if err != nil {
				return nil, nil, nil, err
			}

			// Pass through *this* step's job dependency IDs, so that the child
			// step's jobs aren't scheduled prior to any dependencies.
			parentStepDep := map[string][]string{pipeline.Id: job.DependsOn}

			embedJobs, embedRun, embedStepIds, err := s.buildStepJobs(ctx, log, req,
				visitedPipelines, nodeStepRef, parentStepDep, embeddedPipeline, pipelineRun, rootPipeline)
			if err != nil {
				return nil, nil, nil, err
			}

			// Add the parent step workspace ref and apply it to all embedded
			// pipeline job templates if step was not configured with a
			// workspace ref.
			if step.Workspace != nil {
				for _, jobReq := range embedJobs {
					if jobReq.Job.Workspace.Workspace == "default" {
						jobReq.Job.Workspace = &pb.Ref_Workspace{
							Workspace: step.Workspace.Workspace,
						}
					}
				}
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

			// Steps of type Pipeline Refs are now noop jobs. This is currently a work around to ensure that
			// if a step parent is *Also* an embedded pipeline, we should not run until
			// that pipeline is complete. To accomplish this, we make a Noop job depend
			// on all of the embedded pipeline jobs.
			job.Operation = &pb.Job_Noop_{
				Noop: &pb.Job_Noop{},
			}
			for _, stepJobReq := range embedJobs {
				job.DependsOn = append(job.DependsOn, stepJobReq.Job.Id)
			}
		default:
			job.Operation = &pb.Job_PipelineStep{
				PipelineStep: &pb.Job_PipelineStepOp{
					Step: step,
				},
			}
		}

		// Include a list of all associated jobs for this specific run
		pipelineRun.Jobs = append(pipelineRun.Jobs, &pb.Ref_Job{Id: job.Id})
		stepJobs = append(stepJobs, &pb.QueueJobRequest{
			Job: job,
		})
	}

	return stepJobs, pipelineRun, stepIds, nil
}

// pipelineGraphFull takes a pipeline, and optionally accepts an existing Graph
// and parent step, and attempts to build a full graph for a given Pipeline including
// any nested pipeline steps. It keeps track of visited steps so that we don't
// get stuck in a loop building the graph. This means step names must be unique
// across pipelines.
func (s *Service) pipelineGraphFull(
	ctx context.Context,
	log hclog.Logger,
	g *graph.Graph,
	parentStep string,
	parentStepDeps []string,
	visitedNodes map[string]string,
	nodeStepRef *nodeToStepRef,
	pipeline *pb.Pipeline,
) (*graph.Graph, *nodeToStepRef, error) {
	stepGraph := &graph.Graph{}
	if g != nil {
		// We're handling an embedded pipeline graph
		stepGraph = g
	}

	// Beginning of the graph builder
	if nodeStepRef == nil {
		nodeStepRef = &nodeToStepRef{
			nodeStepRefs: make(map[string]*pb.Ref_PipelineStep),
			stepRefs:     make(map[nodePipelineStepRef]string),
		}
	}

	// Note that pipeline.Steps is not an ordered list of steps. It's a map of key val
	// steps so the order will not match the order steps are defined in a waypoint.hcl.
	for _, step := range pipeline.Steps {
		if len(visitedNodes) != 0 {
			if pipeName, ok := visitedNodes[step.Name]; ok && pipeName == pipeline.Name {
				log.Trace("we've cycled to a node we've already visited!", "pipeline", pipeName, "step", step.Name)
				return nil, nil, status.Errorf(codes.FailedPrecondition,
					"cycle has been detected. Node %q in pipeline %q has already been visisted", step.Name, pipeName)
			}
		}

		// Look up step and pipeline id in case we generated the id when we added
		// dependencies from a different step
		nodeId, ok := nodeStepRef.stepRefs[nodePipelineStepRef{pipeline: pipeline.Id, step: step.Name}]
		if !ok {
			var err error
			// unique node graph id for full graph
			uid, err := server.Id()
			if err != nil {
				return nil, nil, err
			}
			nodeId = fmt.Sprintf("%s.%s-%s", pipeline.Id, step.Name, uid)
		}

		nodeStepRef.nodeStepRefs[nodeId] = &pb.Ref_PipelineStep{
			PipelineId: pipeline.Id,
			Step:       step.Name,
		}
		nodeStepRef.stepRefs[nodePipelineStepRef{pipeline: pipeline.Id, step: step.Name}] = nodeId

		// Add our step to the graph as a vertex
		stepGraph.Add(nodeId)

		// Keep track of the fact that we've visited this step node in this pipeline
		// to prevent infinite cycles as we build embedded pipeline graphs
		visitedNodes[step.Name] = pipeline.Name

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
			for _, dep := range parentStepDeps {
				stepGraph.AddEdge(dep, nodeId)
			}

			// The edge here indicates an order. It says that parentStep should run before
			// nodeId. Ie to traverse the graph of work, you need to travel from the parentStep
			// vertex to the nodeId vertex.
			stepGraph.AddEdge(nodeId, parentStep)
		}

		var myDeps []string

		// Add any dependencies as defined by the current Step
		for _, dep := range step.DependsOn {
			// Look up the dependency step by name in case we've already generated
			// a node ID for it for our stepGraph
			depId, ok := nodeStepRef.stepRefs[nodePipelineStepRef{pipeline: pipeline.Id, step: dep}]
			if !ok {
				var err error
				// We haven't reached this node yet, but we're adding it to the graph so
				// we generate an id here too so we can add it to the graph as a vertex
				// and create an edge to the given step with the depencny
				uid, err := server.Id() // unique node graph id for full graph
				if err != nil {
					return nil, nil, err
				}

				depId = fmt.Sprintf("%s.%s-%s", pipeline.Id, dep, uid)

				// add node id to map
				nodeStepRef.nodeStepRefs[depId] = &pb.Ref_PipelineStep{
					PipelineId: pipeline.Id,
					Step:       dep,
				}
				nodeStepRef.stepRefs[nodePipelineStepRef{pipeline: pipeline.Id, step: dep}] = depId
			}

			myDeps = append(myDeps, depId)

			// Add the dependency as a vertex and draw an edge to *this* steps node ID
			stepGraph.Add(depId)

			// The edge here indicates that to travel the graph properly, you have to
			// go from the depId vertex to the nodeId vertex. For example, if
			// B depends on C, then depId == C, nodeId == B.
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
			embeddedPipeline, err := s.state(ctx).PipelineGet(ctx, embedRef.Pipeline.Ref)
			if err != nil {
				return nil, nil, err
			}

			// Build the nested pipelines graph
			parentStep := nodeId
			embeddedGraph, embedNodeToStepRef, err := s.pipelineGraphFull(ctx, log, stepGraph,
				parentStep, myDeps, visitedNodes, nodeStepRef, embeddedPipeline)
			if err != nil {
				return nil, nil, err
			}

			nodeStepRef = embedNodeToStepRef
			stepGraph = embeddedGraph
		}
	}

	if cycles := stepGraph.Cycles(); len(cycles) > 0 {
		return nil, nil, status.Errorf(codes.FailedPrecondition,
			"step dependencies in pipeline %q contain one or more cycles: %s",
			pipeline.Name, cycles)
	}

	return stepGraph, nodeStepRef, nil
}
