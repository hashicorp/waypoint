package singleprocess

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server/boltdbstate"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/grpcmetadata"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	"github.com/hashicorp/waypoint/pkg/server/logbuffer"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

const (
	// defaultLogLimitBacklog is the default backlog amount to send down.
	defaultLogLimitBacklog = 100

	// maxEntriesPerRead is how many log entries we request at a time.
	maxEntriesPerRead = 60
)

func (s *Service) spawnLogPlugin(
	ctx context.Context,
	log hclog.Logger,
	deployment *pb.Deployment,
) (*boltdbstate.InstanceLogs, string, error) {
	// TODO(mitchellh): We only support logs if we're using the in-memory
	// state store. We will add support for our other stores later.
	inmemstate, ok := s.state(ctx).(*boltdbstate.State)
	if !ok {
		return nil, "", status.Errorf(codes.Unimplemented,
			"state storage doesn't support log streaming")
	}

	instId, err := server.Id()
	if err != nil {
		return nil, "", err
	}

	log.Info("spawning logs plugin via job",
		"instance-id", instId, "deployment", deployment.Id)

	// Create an InstanceLogs entry for EntrypointLogStream to detect and
	// write logs to. In this way, we can easily coordinate the log entries
	// from the logs plugin to here, the reading half.
	var lo boltdbstate.InstanceLogs
	lo.LogBuffer = logbuffer.New()

	err = inmemstate.InstanceLogsCreate(instId, &lo)
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

// streamRec is a single "stream record" that represents a single instance
// with logs. This is used by the sendInstanceLogs function to get all the
// required information for streaming logs.
type streamRec struct {
	InstanceId   string
	DeploymentId string
	LogBuffer    *logbuffer.Buffer

	InstanceLogsId int64
	JobId          string
}

func (s *Service) GetLogStream(
	req *pb.GetLogStreamRequest,
	srv pb.Waypoint_GetLogStreamServer,
) error {
	ctx := srv.Context()
	log := hclog.FromContext(srv.Context())

	// TODO(mitchellh): We only support logs if we're using the in-memory
	// state store. We will add support for our other stores later.
	inmemstate, ok := s.state(ctx).(*boltdbstate.State)
	if !ok {
		return status.Errorf(codes.Unimplemented,
			"state storage doesn't support log streaming")
	}

	// Default the limit
	if req.LimitBacklog == 0 {
		req.LimitBacklog = defaultLogLimitBacklog
	}

	// instanceFunc will be the function that sendInstanceLogs calls in order
	// to grab the list of instances. This is expected to setup the WatchSet
	// to notify the caller when the set of instances changes.
	var instanceFunc func(ws memdb.WatchSet) ([]*streamRec, error)

	switch scope := req.Scope.(type) {
	case *pb.GetLogStreamRequest_DeploymentId:
		deployment, err := s.state(ctx).DeploymentGet(ctx, &pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{
				Id: scope.DeploymentId,
			},
		})
		if err != nil {
			return hcerr.Externalize(
				log,
				err,
				"failed to get deployment in get log stream",
				"deployment_id",
				scope.DeploymentId,
			)
		}

		log = log.With("deployment_id", scope.DeploymentId)

		// This flag is set when we create the Deployment value by detecting if the plugin
		// had a LogsFunc defined.
		if deployment.HasLogsPlugin {
			log.Debug("deployment supports log plugin. spawning log plugin")
			inst, jobId, err := s.spawnLogPlugin(srv.Context(), log, deployment)
			if err != nil {
				return hcerr.Externalize(
					log,
					fmt.Errorf("error spawning log plugin: %w", err),
					"error spawning log plugin",
				)
			}

			// Be sure to cleanup when we, the job creator, are finished.
			defer s.state(ctx).JobCancel(ctx, jobId, false)

			// Because we spawned the writer, we can safely delete the whole thing
			// when the reader is done.
			defer inmemstate.InstanceLogsDelete(inst.Id)

			log.Debug("log plugin spawned", "job_id", jobId)
			instanceFunc = func(ws memdb.WatchSet) ([]*streamRec, error) {
				return []*streamRec{{
					InstanceId:   inst.InstanceId,
					DeploymentId: scope.DeploymentId,
					LogBuffer:    inst.LogBuffer,
				}}, nil
			}
		} else {
			log.Debug("deployment log will watch connected instances")
			instanceLog := log.Named("instancefunc")
			instanceFunc = func(ws memdb.WatchSet) ([]*streamRec, error) {
				instances, err := s.state(ctx).InstancesByDeployment(ctx, scope.DeploymentId, ws)
				if err != nil {
					return nil, hcerr.Externalize(
						log,
						err,
						"failed to get instance by deployment",
						"deployment_id",
						scope.DeploymentId,
					)
				}
				instanceLog.Trace("instances loaded for deployment",
					"instances_len", len(instances))

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
		log.Debug("application-scope log stream requested")

		// We don't want to respawn plugins if we don't need to, so this
		// memoizes the instance created by launching a plugin for a deployment
		deploymentToInstance := map[string]*streamRec{}

		// Be sure to cleanup all our detritus when done!
		defer func() {
			for _, il := range deploymentToInstance {
				inmemstate.InstanceLogsDelete(il.InstanceLogsId)
				s.state(ctx).JobCancel(ctx, il.JobId, false)
			}
		}()

		instanceLog := log.Named("instancefunc")
		instanceFunc = func(ws memdb.WatchSet) ([]*streamRec, error) {
			// The old version of this code just asked for all the instances (meaning
			// live connected entrypoint processes) for the application. Because we need
			// to change behavior based on Deployment, it had to be changed.
			// We filter and only consider successful and created deployments because
			// that most accurate matches the only scope, so we shouldn't miss anything.

			deployments, err := s.state(ctx).DeploymentList(ctx, scope.Application.Application,
				serverstate.ListWithPhysicalState(pb.Operation_CREATED),
				serverstate.ListWithWorkspace(scope.Application.Workspace),
				serverstate.ListWithStatusFilter(&pb.StatusFilter{
					Filters: []*pb.StatusFilter_Filter{
						{
							Filter: &pb.StatusFilter_Filter_State{
								State: pb.Status_SUCCESS,
							},
						},
					},
				}),
				serverstate.ListWithWatchSet(ws),
			)
			if err != nil {
				return nil, hcerr.Externalize(
					log,
					err,
					"failed to list successful deployments",
				)
			}
			instanceLog.Trace(
				"deployments refreshed for application",
				"deployments_len", len(deployments),
			)

			var streams []*streamRec
			for _, dep := range deployments {
				// If this deployment uses a logs plugin, either used the previously spawn
				// instance or spawn a new instance by invoking the logs pluign.
				if dep.HasLogsPlugin {
					if inst, ok := deploymentToInstance[dep.Id]; ok {
						streams = append(streams, inst)
					} else {
						instanceLog.Trace("deployment supports log plugin, spawning log plugin instance")
						inst, jobId, err := s.spawnLogPlugin(srv.Context(), log, dep)
						if err != nil {
							return nil, hcerr.Externalize(
								log,
								err,
								"failed to launch log plugin for deployment",
							)
						}
						instanceLog.Trace("log plugin spawned", "job_id", jobId)

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
					instanceLog.Trace("tracking instances for deployment", "deployment_id", dep.Id)
					depInstances, err := s.state(ctx).InstancesByDeployment(ctx, dep.Id, ws)
					if err != nil {
						return nil, hcerr.Externalize(
							log,
							err,
							"failed to get instance by deployment",
							"deployment_id",
							dep.Id,
						)
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

	return s.sendInstanceLogs(srv.Context(), log, srv, instanceFunc, req.LimitBacklog)
}

// Used to reduce the functional surface area of sendInstanceLogs. This is
// implemented by the GetLogStream service stream.
type batchSender interface {
	Send(batch *pb.LogBatch) error
}

// sendInstanceLogs calls instanceFunc and coordinates sending both the known
// log entries (ie ones that are already stored on the server) to the waiting
// client, as well as blocking for all new generate log entries that are created
// currently known as well as future instances that are returned by instanceFunc.
func (s *Service) sendInstanceLogs(
	ctx context.Context,
	log hclog.Logger,
	sender batchSender,
	instanceFunc func(ws memdb.WatchSet) ([]*streamRec, error),
	backlog int32,
) error {
	// We keep track of what instances we already have readers for here.
	var instanceSetLock sync.Mutex
	instanceSet := make(map[string]struct{})

	// For values returned by LogMerge, we want to be able to map back to
	// the stream that the reader was for, so we can include the instanceid.
	// This map lets us do that.
	readerToInstance := make(map[*logbuffer.Reader]*streamRec)

	// Step 1: we use log merge to read all the known entries from existing
	// instances. This will never block, it will just weave the log entries
	// that are already known together and then send them back to the client.

	// Get all current records
	ws := memdb.NewWatchSet()
	records, err := instanceFunc(ws)
	if err != nil {
		return err
	}
	log.Trace("initial instances loaded", "len", len(records))

	var readers []logbuffer.MergeReader
	for _, record := range records {
		r := record.LogBuffer.Reader(backlog)
		readerToInstance[r] = record

		readers = append(readers, r)
	}

	// Read out all the log entries from LogMerge. This never blocks waiting
	// for new log entries, it will simply let each reader output all known
	// entries and then loop exits.
	lm := logbuffer.NewMerger(readers...)
	lines := make([]*pb.LogBatch_Entry, maxEntriesPerRead)
	for {
		entries, err := lm.Read(len(lines))
		if err != nil {
			return err
		}

		// When there are no more buffered entries to read, means
		// we'll switch to the on-demand reading logic below.
		if len(entries) == 0 {
			break
		}

		log.Trace("sending known log data", "entries", len(entries))

		var (
			prev *streamRec
			idx  int
		)

		// We batch up the lines so long as the previous line is from the
		// same instance as the current one. When we detect a change, we flush
		// lines and begin buffering again.
		for _, v := range entries {
			rec := readerToInstance[v.Reader.(*logbuffer.Reader)]

			if prev != nil && prev != rec {
				// Flush current lines that were all the same instance
				sender.Send(&pb.LogBatch{
					DeploymentId: prev.DeploymentId,
					InstanceId:   prev.InstanceId,
					Lines:        lines[:idx],
				})

				idx = 0
			}

			val := v.Value().(*pb.LogBatch_Entry)
			lines[idx] = val

			prev = rec
			idx++
		}

		// Flush any unsent lines
		if idx > 0 {
			// Flush current lines that were all the same instance
			sender.Send(&pb.LogBatch{
				DeploymentId: prev.DeploymentId,
				InstanceId:   prev.InstanceId,
				Lines:        lines[:idx],
			})
		}
	}

	// Step 2: startup background forwarders for all the readers we spawned above.

	// We lock around this section because if one of the launched goroutines exits
	// very quickly, we might still be adding the others when it does exit.
	{
		instanceSetLock.Lock()

		for r, rec := range readerToInstance {
			instanceSet[rec.InstanceId] = struct{}{}

			instanceLog := log.With("instance_id", rec.InstanceId)
			instanceLog.Debug("instance log stream starting", "instance-id", rec.InstanceId)

			go func(r *logbuffer.Reader, rec *streamRec) {
				defer func() {
					instanceSetLock.Lock()
					defer instanceSetLock.Unlock()
					delete(instanceSet, rec.InstanceId)
				}()

				s.forwardLogBatches(ctx, instanceLog, sender, r, rec)
			}(r, rec)
		}

		instanceSetLock.Unlock()
	}

	// Step 3: Now we setup goroutines to read on-demand log entries. Additionally
	// this will pickup any new instances and begin streaming their on-demand log entries.

	// We loop forever so that we can automatically get any new instances that
	// join as we have an open log stream.
	for {
		// Wait for changes or to be done
		if err := ws.WatchCtx(ctx); err != nil {
			// If our context ended, exit with that
			if err := ctx.Err(); err != nil {
				return err
			}

			return err
		}

		// Get all current records
		ws = memdb.NewWatchSet()
		records, err := instanceFunc(ws)
		if err != nil {
			return err
		}
		log.Trace("instances reloaded", "len", len(records))

		// For each record, start a goroutine that reads the log entries and sends them.
		// This will skip any instance we already know about.
		for _, record := range records {
			instanceId := record.InstanceId

			// If we already have a reader for this, then do nothing.
			instanceSetLock.Lock()
			_, exit := instanceSet[instanceId]
			instanceSet[instanceId] = struct{}{}
			instanceSetLock.Unlock()
			if exit {
				continue
			}

			// Start our reader up
			r := record.LogBuffer.Reader(backlog)
			instanceLog := log.With("instance_id", instanceId)
			instanceLog.Info("instance log stream starting", "instance-id", instanceId)

			go func(record *streamRec) {
				defer func() {
					instanceSetLock.Lock()
					defer instanceSetLock.Unlock()
					delete(instanceSet, record.InstanceId)
				}()

				s.forwardLogBatches(ctx, instanceLog, sender, r, record)
			}(record)
		}
	}
}

// forwardLogBatches reads entries from the reader and spews them at the sender,
// nothing more. This function wires the reader up to the given context, such that
// when the context is finished, the reader is forced closed so this function can
// return.
func (s *Service) forwardLogBatches(
	ctx context.Context,
	log hclog.Logger,
	sender batchSender,
	r *logbuffer.Reader,
	record *streamRec,
) {
	go r.CloseContext(ctx)
	defer log.Debug("instance log stream ending")

	for {
		entries := r.Read(64, true)
		if entries == nil {
			log.Debug("exitting logs loop, Read returned nil")
			return
		}

		lines := make([]*pb.LogBatch_Entry, len(entries))
		for i, v := range entries {
			lines[i] = v.(*pb.LogBatch_Entry)
		}

		log.Trace("sending instance log data", "entries", len(entries))
		sender.Send(&pb.LogBatch{
			DeploymentId: record.DeploymentId,
			InstanceId:   record.InstanceId,
			Lines:        lines,
		})
	}
}
