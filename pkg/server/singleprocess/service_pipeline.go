package singleprocess

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/go-hclog"

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

	// Generate job IDs for each of the steps. We need to know the IDs in
	// advance to set up the dependency chain.
	stepIds := map[string]string{}
	for name := range pipeline.Steps {
		stepIds[name], err = server.Id()
		if err != nil {
			return nil, err
		}
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
	run, err := s.state(ctx).PipelineRunGetLatest(pipeline.Id)
	if err != nil {
		return nil, err
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
		job.Pipeline = &pb.Ref_PipelineStep{
			Pipeline:    pipeline.Id,
			Step:        step.Name,
			RunSequence: run.Sequence,
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
				// TODO(briancain): do it
				// copy `job` and update Operation to be release I think. then append
				// job to stepJobs

				// Queue a release job too
				log.Warn("Currently not queueing a release job yet....sry!!!")
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
					deployment, err = s.GetDeployment(ctx, &pb.GetDeploymentRequest{
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
						return nil, err
					}
					if deployment == nil {
						log.Debug("could not find deploy sequence, using latest instead", "seq", d.Sequence)
					}
				default:
					// return an error
					log.Error("invalid deployment ref received", "ref", d)
					return nil, status.Errorf(codes.Internal, "invalid deployment ref received: %T", d)
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
			// TODO(briancain): Look up pipeline by Owner and re-call this func
			// to receieve its step jobs. Might need to validate that there are no
			// cycles here too.
			return nil, status.Error(codes.Unimplemented, "running an embedded pipeline is not yet supported")
		default:
			job.Operation = &pb.Job_PipelineStep{
				PipelineStep: &pb.Job_PipelineStepOp{
					Step: step,
				},
			}
		}

		run.Jobs = append(run.Jobs, &pb.Ref_Job{Id: job.Id})
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
			Pipeline:    pipeline.Id,
			Step:        v.(string),
			RunSequence: run.Sequence,
		}
	}

	run.State = pb.PipelineRun_STARTING
	if err = s.state(ctx).PipelineRunPut(run); err != nil {
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
