package singleprocess

import (
	"context"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/grpcmetadata"
	"github.com/hashicorp/waypoint/internal/server/logbuffer"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

// defaultLogLimitBacklog is the default backlog amount to send down.
const defaultLogLimitBacklog = 100

func (s *service) spawnLogPlugin(
	ctx context.Context,
	log hclog.Logger,
	deployment *pb.Deployment,
) (*state.InstanceLogs, string, error) {
	instId, err := server.Id()
	if err != nil {
		return nil, "", err
	}

	log.Info("spawning logs plugin via job",
		"instance-id", instId, "deployment", deployment.Id)

	// Create an InstanceLogs entry for EntrypointLogStream to detect and
	// write logs to. In this way, we can easily coordinate the log entries
	// from the logs plugin to here, the reading half.
	var lo state.InstanceLogs
	lo.LogBuffer = logbuffer.New()

	err = s.state.InstanceLogsCreate(instId, &lo)
	if err != nil {
		return nil, "", err
	}

	job := &pb.Job{
		Workspace:   deployment.Workspace,
		Application: deployment.Application,
		Operation: &pb.Job_Logs{
			Logs: &pb.Job_LogsOp{
				InstanceId: instId,
				Deployment: deployment,
			},
		},
	}

	// Means the client WANTS the job run on itself, so let's target the
	// job back to it.
	if runnerId, ok := grpcmetadata.RunnerId(ctx); ok {
		job.DataSource = &pb.Job_DataSource{
			Source: &pb.Job_DataSource_Local{
				Local: &pb.Job_Local{},
			},
		}

		job.TargetRunner = &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Id{
				Id: &pb.Ref_RunnerId{
					Id: runnerId,
				},
			},
		}

		// Otherwise, the client wants a logs session but doesn't have a runner
		// to use, so we just target any runner.
	} else {
		job.TargetRunner = &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Any{
				Any: &pb.Ref_RunnerAny{},
			},
		}

		// We leave DataSource unset here so that QueueJob will port over the data
		// source from the project.
	}

	qresp, err := s.QueueJob(ctx, &pb.QueueJobRequest{
		Job: job,

		// TODO unknown if this is enough time for when the request is queued
		// by a runner-less client but a user waiting 60 seconds will get impatient
		// regardless.
		ExpiresIn: "60s",
	})
	if err != nil {
		return nil, "", err
	}

	jobId := qresp.JobId

	log.Debug("waiting on job state", "job-id", jobId)

	state, err := s.waitOnJobStarted(ctx, jobId)
	if err != nil {
		return nil, "", err
	}

	switch state {
	case pb.Job_ERROR:
		return nil, "", status.Errorf(codes.FailedPrecondition, "job errored out before starting")
	case pb.Job_SUCCESS:
		return nil, "", status.Errorf(codes.Internal, "job succeeded before running")
	case pb.Job_RUNNING:
		// ok
	default:
		return nil, "", status.Errorf(codes.Internal, "unexpected job status: %s", state.String())
	}

	return &lo, jobId, nil
}

// TODO: test
func (s *service) GetLogStream(
	req *pb.GetLogStreamRequest,
	srv pb.Waypoint_GetLogStreamServer,
) error {

	// Used to coordinate the data from either Instance or InstanceLogs entries
	// and then funnel them all to the waiting client.
	type streamRec struct {
		InstanceId   string
		DeploymentId string
		LogBuffer    *logbuffer.Buffer

		InstanceLogsId int64
		JobId          string
	}

	log := hclog.FromContext(srv.Context())

	// Default the limit
	if req.LimitBacklog == 0 {
		req.LimitBacklog = defaultLogLimitBacklog
	}

	var instanceFunc func(ws memdb.WatchSet) ([]*streamRec, error)
	switch scope := req.Scope.(type) {
	case *pb.GetLogStreamRequest_DeploymentId:
		deployment, err := s.state.DeploymentGet(&pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{
				Id: scope.DeploymentId,
			},
		})
		if err != nil {
			return err
		}

		log = log.With("deployment_id", scope.DeploymentId)

		// This flag is set when we create the Deployment value by detecting if the plugin
		// had a LogsFunc defined.
		if deployment.HasLogsPlugin {
			inst, jobId, err := s.spawnLogPlugin(srv.Context(), log, deployment)
			if err != nil {
				return err
			}

			// Be sure to cleanup when we, the job creator, are finished.
			defer s.state.JobCancel(jobId, false)

			// Because we spawned the writer, we can safely delete the whole thing
			// when the reader is done.
			go s.state.InstanceLogsDelete(inst.Id)

			instanceFunc = func(ws memdb.WatchSet) ([]*streamRec, error) {
				return []*streamRec{{
					InstanceId:   inst.InstanceId,
					DeploymentId: scope.DeploymentId,
					LogBuffer:    inst.LogBuffer,
				}}, nil
			}
		} else {
			instanceFunc = func(ws memdb.WatchSet) ([]*streamRec, error) {
				instances, err := s.state.InstancesByDeployment(scope.DeploymentId, ws)
				if err != nil {
					return nil, err
				}

				var bufs []*streamRec

				for _, i := range instances {
					bufs = append(bufs, &streamRec{
						InstanceId:   i.Id,
						DeploymentId: scope.DeploymentId,
						LogBuffer:    i.LogBuffer,
					})
				}

				return bufs, nil
			}
		}

	case *pb.GetLogStreamRequest_Application_:
		if scope.Application == nil ||
			scope.Application.Application == nil ||
			scope.Application.Workspace == nil {
			return status.Errorf(
				codes.FailedPrecondition,
				"application scope requires the application and workspace fields to be set",
			)
		}

		log = log.With(
			"project", scope.Application.Application.Project,
			"application", scope.Application.Application.Application,
			"workspace", scope.Application.Workspace.Workspace,
		)

		// We don't want to respawn plugins if we don't need to, so this
		// memoizes the instance created by launching a plugin for a deployment
		deploymentToInstance := map[string]*streamRec{}

		// Be sure to cleanup all our detritus when done!
		defer func() {
			for _, il := range deploymentToInstance {
				s.state.InstanceLogsDelete(il.InstanceLogsId)
				s.state.JobCancel(il.JobId, false)
			}
		}()

		instanceFunc = func(ws memdb.WatchSet) ([]*streamRec, error) {
			// The old version of this code just asked for all the instances (meaning
			// live connected entrypoint processes) for the application. Because we need
			// to change behavior based on Deployment, it had to be changed.
			// We filter and only consider successful and created deployments because
			// that most accurate matches the only scope, so we shouldn't miss anything.

			deployments, err := s.state.DeploymentList(scope.Application.Application,
				state.ListWithPhysicalState(pb.Operation_CREATED),
				state.ListWithStatusFilter(&pb.StatusFilter{
					Filters: []*pb.StatusFilter_Filter{
						{
							Filter: &pb.StatusFilter_Filter_State{
								State: pb.Status_SUCCESS,
							},
						},
					},
				}),
				state.ListWithWatchSet(ws),
			)
			if err != nil {
				return nil, err
			}

			var streams []*streamRec

			for _, dep := range deployments {
				// If this deployment uses a logs plugin, either used the previously spawn
				// instance or spawn a new instance by invoking the logs pluign.
				if dep.HasLogsPlugin {
					if inst, ok := deploymentToInstance[dep.Id]; ok {
						streams = append(streams, inst)
					} else {
						inst, jobId, err := s.spawnLogPlugin(srv.Context(), log, dep)
						if err != nil {
							return nil, err
						}

						rec := &streamRec{
							InstanceId:     inst.InstanceId,
							DeploymentId:   dep.Id,
							LogBuffer:      inst.LogBuffer,
							InstanceLogsId: inst.Id,
							JobId:          jobId,
						}

						deploymentToInstance[dep.Id] = rec

						streams = append(streams, rec)
					}
				} else {
					depInstances, err := s.state.InstancesByDeployment(dep.Id, ws)
					if err != nil {
						return nil, err
					}

					for _, i := range depInstances {
						streams = append(streams, &streamRec{
							InstanceId:   i.Id,
							DeploymentId: i.DeploymentId,
							LogBuffer:    i.LogBuffer,
						})
					}
				}
			}

			return streams, nil
		}

	default:
		return status.Errorf(
			codes.FailedPrecondition,
			"invalid scope supplied: %T",
			req.Scope,
		)
	}

	// We keep track of what instances we already have readers for here.
	var instanceSetLock sync.Mutex
	instanceSet := make(map[string]struct{})

	// We loop forever so that we can automatically get any new instances that
	// join as we have an open log stream.
	for {
		// Get all our records
		ws := memdb.NewWatchSet()
		records, err := instanceFunc(ws)
		if err != nil {
			return err
		}
		log.Trace("instances loaded", "len", len(records))

		// For each record, start a goroutine that reads the log entries and sends them.
		for _, record := range records {
			instanceId := record.InstanceId
			deploymentId := record.DeploymentId

			// If we already have a reader for this, then do nothing.
			instanceSetLock.Lock()
			_, exit := instanceSet[instanceId]
			instanceSet[instanceId] = struct{}{}
			instanceSetLock.Unlock()
			if exit {
				continue
			}

			// Start our reader up
			r := record.LogBuffer.Reader(req.LimitBacklog)
			instanceLog := log.With("instance_id", instanceId)
			instanceLog.Debug("instance log stream starting", "instance-id", instanceId)

			go r.CloseContext(srv.Context())
			go func() {
				defer instanceLog.Debug("instance log stream ending")
				defer func() {
					instanceSetLock.Lock()
					defer instanceSetLock.Unlock()
					delete(instanceSet, instanceId)
				}()

				for {
					entries := r.Read(64, true)
					if entries == nil {
						instanceLog.Debug("exitting logs loop, Read returned nil")
						return
					}

					lines := make([]*pb.LogBatch_Entry, len(entries))
					for i, v := range entries {
						lines[i] = v.(*pb.LogBatch_Entry)
					}

					instanceLog.Trace("sending instance log data", "entries", len(entries))
					srv.Send(&pb.LogBatch{
						DeploymentId: deploymentId,
						InstanceId:   instanceId,
						Lines:        lines,
					})
				}
			}()
		}

		// Wait for changes or to be done
		if err := ws.WatchCtx(srv.Context()); err != nil {
			// If our context ended, exit with that
			if err := srv.Context().Err(); err != nil {
				return err
			}

			return err
		}
	}
}
