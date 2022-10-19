package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) UpsertTrigger(
	ctx context.Context,
	req *pb.UpsertTriggerRequest,
) (*pb.UpsertTriggerResponse, error) {
	if err := serverptypes.ValidateUpsertTriggerRequest(req); err != nil {
		return nil, err
	}

	result := req.Trigger
	if err := s.state(ctx).TriggerPut(result); err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to upsert trigger",
			"trigger_id",
			result.GetId(),
		)
	}

	return &pb.UpsertTriggerResponse{Trigger: result}, nil
}

// GetTrigger returns a Trigger based on ID
func (s *Service) GetTrigger(
	ctx context.Context,
	req *pb.GetTriggerRequest,
) (*pb.GetTriggerResponse, error) {
	if err := serverptypes.ValidateGetTriggerRequest(req); err != nil {
		return nil, err
	}

	t, err := s.state(ctx).TriggerGet(req.Ref)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to get trigger",
			"trigger_id",
			req.Ref.GetId(),
		)
	}

	return &pb.GetTriggerResponse{Trigger: t}, nil
}

// DeleteTrigger deletes a Trigger based on ID
func (s *Service) DeleteTrigger(
	ctx context.Context,
	req *pb.DeleteTriggerRequest,
) (*empty.Empty, error) {
	if err := serverptypes.ValidateDeleteTriggerRequest(req); err != nil {
		return nil, err
	}

	err := s.state(ctx).TriggerDelete(req.Ref)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to delete trigger",
			"trigger_id",
			req.Ref.GetId(),
		)
	}

	return &empty.Empty{}, nil
}

func (s *Service) ListTriggers(
	ctx context.Context,
	req *pb.ListTriggerRequest,
) (*pb.ListTriggerResponse, error) {
	// NOTE: no ptype validation at the moment, as all Ref fields are optional

	result, err := s.state(ctx).TriggerList(req.Workspace, req.Project, req.Application, req.Tags)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to list triggers",
		)
	}

	return &pb.ListTriggerResponse{Triggers: result}, nil
}

// NoAuthRunTrigger is defined in our auth.go package to not need a token for
// incoming requests. Eventually, Waypoint should have a middleware func
// for incoming requests that can handle API endpoints that can support
// authenticated and non-authenticated requests like this one.
// https://github.com/hashicorp/waypoint/issues/2990
func (s *Service) NoAuthRunTrigger(
	ctx context.Context,
	req *pb.RunTriggerRequest,
) (*pb.RunTriggerResponse, error) {
	if err := serverptypes.ValidateRunTriggerRequest(req); err != nil {
		return nil, err
	}

	log := hclog.FromContext(ctx)
	log.Trace("attempting to find and run trigger from authless func", "trigger_id", req.Ref.Id)

	trigger, err := s.state(ctx).TriggerGet(req.Ref)
	if err != nil {
		log.Error("failed to get requested trigger", "trigger_id", req.Ref.Id, "error", err)
		return nil, status.Errorf(codes.NotFound,
			"trigger id %q not found. check the waypoint server logs for more information", req.Ref.Id)
	}

	if trigger.Authenticated {
		log.Error("requested trigger id requires authentication to run", "trigger_id", trigger.Id)
		return nil, status.Error(codes.PermissionDenied, "trigger requires authentication")
	}

	resp, err := s.RunTrigger(ctx, req)
	return resp, err
}

// RunTrigger takes a configured trigger URL, and crafts all of the appropriate
// data into a QueuedJobRequest for Waypoint to execute. It returns a list of
// job ids that have been queued from the result of the request.
func (s *Service) RunTrigger(
	ctx context.Context,
	req *pb.RunTriggerRequest,
) (*pb.RunTriggerResponse, error) {
	if err := serverptypes.ValidateRunTriggerRequest(req); err != nil {
		return nil, err
	}
	log := hclog.FromContext(ctx)

	runTrigger, err := s.state(ctx).TriggerGet(req.Ref)
	if err != nil {
		return nil, hcerr.Externalize(
			log,
			err,
			"failed to run trigger",
			"trigger_id",
			req.Ref.GetId(),
		)
	}

	log = log.With("run_trigger", runTrigger.Id)

	log.Debug("building run trigger job")

	// Build the job(s)
	job := &pb.Job{
		Workspace: runTrigger.Workspace,
		Labels:    map[string]string{"trigger/id": runTrigger.Id},
	}

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
	case *pb.Trigger_StatusReport:
		job.Operation = &pb.Job_StatusReport{StatusReport: op.StatusReport}
	default:
		return nil, status.Errorf(codes.Internal,
			"trigger %q is configured with an unsupported operation %T", runTrigger.Id, op)
	}

	if len(req.VariableOverrides) > 0 {
		log.Debug("variable overrides have been requested for trigger job")
		for i, v := range req.VariableOverrides {
			switch vType := v.Source.(type) {
			case *pb.Variable_Cli:
				continue
			default:
				if vType == nil {
					return nil, status.Errorf(codes.FailedPrecondition,
						"No Variable type for %q given. Expected \"variable_cli\" type for override.", req.VariableOverrides[i].Name)
				} else {
					return nil, status.Errorf(codes.FailedPrecondition,
						"Incorrect Variable type for %q given. Got %T, but expected \"variable_cli\" type.", req.VariableOverrides[i].Name, vType)
				}
			}
		}

		job.Variables = req.VariableOverrides
	}

	// TODO(briancain): look up a target runner config at the project/app level and apply it to job requests
	job.TargetRunner = &pb.Ref_Runner{Target: &pb.Ref_Runner_Any{}}

	// generate job requests
	var jobList []*pb.QueueJobRequest
	var ids []string
	if runTrigger.Application == nil || runTrigger.Application.Application == "" {
		// we're gonna queue multiple jobs for every application in a project
		log.Debug("building multi-jobs for all apps in project", "project", runTrigger.Project.Project)
		jobList, err = s.state(ctx).JobProjectScopedRequest(ctx, runTrigger.Project, job)
		if err != nil {
			return nil, hcerr.Externalize(
				hclog.FromContext(ctx),
				err,
				"error creating job scoped to project when running a trigger",
				"trigger_id",
				req.Ref.GetId(),
				"project",
				runTrigger.Project.GetProject(),
			)
		}
	} else {
		log.Debug("building a single job for target", "project",
			runTrigger.Application.Project, "app", runTrigger.Application.Application)
		// we're only targetting a specific application, so queue 1 job
		job.Application = runTrigger.Application
		j := &pb.QueueJobRequest{Job: job}
		jobList = append(jobList, j)
	}

	// NOTE(briancain): This loops is to set full messages on an operation. Certain
	// operations don't take references and require the entire message from the database
	// to properly perform its operation. We set those here before queueing the job.
	// Ideally, the executeXOperation would do the lookup after receiving a Ref_X rather than
	// here when we are queueing a job.
	// NOTE(briancain): See https://github.com/hashicorp/waypoint/issues/2884
	// for why we must attach the full PushedArtifact message for a deploy operation
	// We have to set the full artifact message on Deployment operations
	// This is true for other operations too like Push and Release.
	for i, qJob := range jobList {
		switch op := qJob.Job.Operation.(type) {
		case *pb.Job_Push:
			if op.Push.Build.Sequence == 0 {
				buildLatest, err := s.state(ctx).BuildLatest(qJob.Job.Application, qJob.Job.Workspace)
				if err != nil {
					return nil, hcerr.Externalize(
						log,
						err,
						"failed to obtain latest build",
					)
				}

				jobList[i].Job.Operation = &pb.Job_Push{
					Push: &pb.Job_PushOp{
						Build: buildLatest,
					},
				}
			} else {
				build, err := s.state(ctx).BuildGet(&pb.Ref_Operation{
					Target: &pb.Ref_Operation_Sequence{
						Sequence: &pb.Ref_OperationSeq{
							Application: qJob.Job.Application,
							Number:      op.Push.Build.Sequence,
						},
					},
				})
				if err != nil {
					return nil, hcerr.Externalize(
						log,
						err,
						"failed to obtain build by id",
						"build_id",
						op.Push.Build.Sequence,
					)
				}

				jobList[i].Job.Operation = &pb.Job_Push{
					Push: &pb.Job_PushOp{
						Build: build,
					},
				}
			}
		case *pb.Job_Destroy:
			switch destroyTarget := op.Destroy.Target.(type) {
			case *pb.Job_DestroyOp_Deployment:
				if destroyTarget.Deployment.Sequence == 0 {
					// get latest deployment
					deployLatest, err := s.state(ctx).DeploymentLatest(destroyTarget.Deployment.Application, destroyTarget.Deployment.Workspace)
					if err != nil {
						return nil, hcerr.Externalize(
							log,
							err,
							"failed to obtain latest deployment for destroying deployment operation trigger",
							"application",
							destroyTarget.Deployment.GetApplication(),
							"workspace",
							destroyTarget.Deployment.GetWorkspace(),
						)
					}

					jobList[i].Job.Operation = &pb.Job_Destroy{
						Destroy: &pb.Job_DestroyOp{
							Target: &pb.Job_DestroyOp_Deployment{
								Deployment: deployLatest,
							},
						},
					}
				} else {
					// get deployment by id seq
					deploy, err := s.state(ctx).DeploymentGet(&pb.Ref_Operation{
						Target: &pb.Ref_Operation_Sequence{
							Sequence: &pb.Ref_OperationSeq{
								Application: qJob.Job.Application,
								Number:      destroyTarget.Deployment.Sequence,
							},
						},
					})
					if err != nil {
						return nil, hcerr.Externalize(
							log,
							err,
							"failed to obtain deployment by id for destroying deployment operation trigger",
							"application",
							qJob.Job.GetApplication(),
							"sequence",
							destroyTarget.Deployment.GetSequence(),
						)
					}

					jobList[i].Job.Operation = &pb.Job_Destroy{
						Destroy: &pb.Job_DestroyOp{
							Target: &pb.Job_DestroyOp_Deployment{
								Deployment: deploy,
							},
						},
					}
				}
			default:
				// We don't need any setup for destroying workspaces at the moment
				break
			}
		case *pb.Job_Deploy:
			if op.Deploy.Artifact == nil {
				// get latest pushed artifact, then set it on the operation
				artifactLatest, err := s.state(ctx).ArtifactLatest(ctx, qJob.Job.Application, qJob.Job.Workspace)
				if err != nil {
					return nil, hcerr.Externalize(
						log,
						err,
						"failed to obtain latest pushed artifact",
						"application",
						qJob.Job.GetApplication(),
						"workspace",
						qJob.Job.GetWorkspace(),
					)
				}

				jobList[i].Job.Operation = &pb.Job_Deploy{
					Deploy: &pb.Job_DeployOp{
						Artifact: artifactLatest,
					},
				}
			} else {
				// Set the actual pushed artifact on the operation
				buildSeq := op.Deploy.Artifact.Sequence
				artifact, err := s.state(ctx).ArtifactGet(ctx, &pb.Ref_Operation{
					Target: &pb.Ref_Operation_Sequence{
						Sequence: &pb.Ref_OperationSeq{
							Application: qJob.Job.Application,
							Number:      buildSeq,
						},
					},
				})
				if err != nil {
					return nil, hcerr.Externalize(
						log,
						err,
						"failed to obtain pushed artifact for id",
						"application",
						qJob.Job.GetApplication(),
						"sequence",
						buildSeq,
					)
				}

				jobList[i].Job.Operation = &pb.Job_Deploy{
					Deploy: &pb.Job_DeployOp{
						Artifact: artifact,
					},
				}
			}
		case *pb.Job_Release:
			// We have to set the full Deployment message on Release, it does not
			// take a ref. This is a similar issue that Deployments have with artifacts
			if op.Release.Deployment.Sequence == 0 {
				// get latest deployment
				deployLatest, err := s.state(ctx).DeploymentLatest(op.Release.Deployment.Application, op.Release.Deployment.Workspace)
				if err != nil {
					return nil, hcerr.Externalize(
						log,
						err,
						"failed to obtain latest deployment for running release operation trigger",
						"application",
						op.Release.Deployment.GetApplication(),
						"workspace",
						op.Release.Deployment.GetWorkspace(),
					)
				}

				jobList[i].Job.Operation = &pb.Job_Release{
					Release: &pb.Job_ReleaseOp{
						Deployment:          deployLatest,
						Prune:               op.Release.Prune,
						PruneRetain:         op.Release.PruneRetain,
						PruneRetainOverride: op.Release.PruneRetainOverride,
					},
				}
			} else {
				// get deployment by id seq
				deploy, err := s.state(ctx).DeploymentGet(&pb.Ref_Operation{
					Target: &pb.Ref_Operation_Sequence{
						Sequence: &pb.Ref_OperationSeq{
							Application: qJob.Job.Application,
							Number:      op.Release.Deployment.Sequence,
						},
					},
				})
				if err != nil {
					return nil, hcerr.Externalize(
						log,
						err,
						"failed to obtain deployment by id for running release operation trigger",
						"application",
						op.Release.Deployment.GetApplication(),
						"sequence",
						op.Release.Deployment.GetSequence(),
					)
				}

				jobList[i].Job.Operation = &pb.Job_Release{
					Release: &pb.Job_ReleaseOp{
						Deployment:          deploy,
						Prune:               op.Release.Prune,
						PruneRetain:         op.Release.PruneRetain,
						PruneRetainOverride: op.Release.PruneRetainOverride,
					},
				}
			}
		case *pb.Job_StatusReport:
			// determine target, then get either deployment/release latest or by seq id
			switch srTarget := op.StatusReport.Target.(type) {
			case *pb.Job_StatusReportOp_Deployment:
				if srTarget.Deployment.Sequence == 0 {
					// get latest deployment
					deployLatest, err := s.state(ctx).DeploymentLatest(srTarget.Deployment.Application, srTarget.Deployment.Workspace)
					if err != nil {
						return nil, hcerr.Externalize(
							log,
							err,
							"failed to obtain latest deployment for running a status report operation trigger",
							"application",
							srTarget.Deployment.GetApplication(),
							"workspace",
							srTarget.Deployment.GetWorkspace(),
						)
					}

					jobList[i].Job.Operation = &pb.Job_StatusReport{
						StatusReport: &pb.Job_StatusReportOp{
							Target: &pb.Job_StatusReportOp_Deployment{
								Deployment: deployLatest,
							},
						},
					}
				} else {
					// get deployment by id seq
					deploy, err := s.state(ctx).DeploymentGet(&pb.Ref_Operation{
						Target: &pb.Ref_Operation_Sequence{
							Sequence: &pb.Ref_OperationSeq{
								Application: qJob.Job.Application,
								Number:      srTarget.Deployment.Sequence,
							},
						},
					})
					if err != nil {
						return nil, hcerr.Externalize(
							log,
							err,
							"failed to obtain deployment by id for running status report operation trigger",
							"application",
							srTarget.Deployment.GetApplication(),
							"sequence",
							srTarget.Deployment.GetSequence(),
						)
					}

					jobList[i].Job.Operation = &pb.Job_StatusReport{
						StatusReport: &pb.Job_StatusReportOp{
							Target: &pb.Job_StatusReportOp_Deployment{
								Deployment: deploy,
							},
						},
					}
				}
			case *pb.Job_StatusReportOp_Release:
				if srTarget.Release.Sequence == 0 {
					releaseLatest, err := s.state(ctx).ReleaseLatest(srTarget.Release.Application, srTarget.Release.Workspace)
					if err != nil {
						return nil, hcerr.Externalize(
							log,
							err,
							"failed to obtain latest release for running a status report operation trigger",
							"application",
							srTarget.Release.GetApplication(),
							"workspace",
							srTarget.Release.GetWorkspace(),
						)
					}

					jobList[i].Job.Operation = &pb.Job_StatusReport{
						StatusReport: &pb.Job_StatusReportOp{
							Target: &pb.Job_StatusReportOp_Release{
								Release: releaseLatest,
							},
						},
					}
				} else {
					// get deployment by id seq
					release, err := s.state(ctx).ReleaseGet(&pb.Ref_Operation{
						Target: &pb.Ref_Operation_Sequence{
							Sequence: &pb.Ref_OperationSeq{
								Application: qJob.Job.Application,
								Number:      srTarget.Release.Sequence,
							},
						},
					})
					if err != nil {
						return nil, hcerr.Externalize(
							log,
							err,
							"failed to obtain release by id for running status report operation trigger",
							"application",
							qJob.Job.GetApplication(),
							"sequence",
							srTarget.Release.GetSequence(),
						)
					}

					jobList[i].Job.Operation = &pb.Job_StatusReport{
						StatusReport: &pb.Job_StatusReportOp{
							Target: &pb.Job_StatusReportOp_Release{
								Release: release,
							},
						},
					}
				}
			default:
				// This shouldn't happen, but let's check anyway and return an error
				return nil, status.Errorf(codes.Internal,
					"incorrect status report target given for running a trigger: %T", srTarget)
			}
		default:
			// We assume all jobs have the same operation, so if none match, don't loop over all jobs
			break
		}
	}

	if len(jobList) > 0 {
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
		log.Warn("the RunTrigger job list was empty, no jobs to queue")
		return nil, nil
	}

	log.Debug("run trigger job(s) have been queued")

	// Trigger has been requested to queue jobs, update active time
	runTrigger.ActiveTime = timestamppb.Now()
	err = s.state(ctx).TriggerPut(runTrigger)
	if err != nil {
		return nil, err
	}

	return &pb.RunTriggerResponse{JobIds: ids}, nil
}
