package singleprocess

import (
	"context"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/grpcmetadata"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	"github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func (s *Service) StartExecStream(
	srv pb.Waypoint_StartExecStreamServer,
) error {
	ctx := srv.Context()
	log := hclog.FromContext(srv.Context())

	// Instance exec support is optional to state, so we need to check before we offer
	// to start the exec stream.
	iexec, ok := s.state(ctx).(serverstate.InstanceExecHandler)
	if !ok {
		return hcerr.Externalize(
			log,
			status.Errorf(codes.Unimplemented,
				"state storage doesn't support exec streaming"),
			"state storage doesn't support exec streaming",
		)
	}

	// Read our first event which must be a Start event.
	log.Trace("waiting for Start message")
	req, err := srv.Recv()
	if err != nil {
		return hcerr.Externalize(
			log,
			err,
			"failed to receive entrypoint exex stream",
		)
	}
	start, ok := req.Event.(*pb.ExecStreamRequest_Start_)
	if !ok {
		return hcerr.Externalize(
			log,
			status.Errorf(codes.FailedPrecondition,
				"error reading entrypoint exec stream, first message must be start type"),
			"error reading entrypoint exec stream, first message must be start type",
		)
	}
	if err := ptypes.ValidateExecStreamRequestStart(start.Start); err != nil {
		return err
	}

	// Create our exec. We have to populate everything here first because
	// once we register, this will trigger any watchers to be notified of
	// a change and the instance should try to connect to us.
	clientEventCh := make(chan *pb.ExecStreamRequest)
	eventCh := make(chan *pb.EntrypointExecRequest)
	execRec := &serverstate.InstanceExec{
		Args:              start.Start.Args,
		Pty:               start.Start.Pty,
		ClientEventCh:     clientEventCh,
		EntrypointEventCh: eventCh,
		Context:           srv.Context(),
	}

	// Register the exec session
	switch t := start.Start.Target.(type) {
	case *pb.ExecStreamRequest_Start_InstanceId:
		log = log.With("instance_id", t.InstanceId)
		err = iexec.InstanceExecCreateByTargetedInstance(t.InstanceId, execRec)
		if err != nil {
			return hcerr.Externalize(
				log,
				err,
				"failed to create exec session with target instance",
			)
		}
	case *pb.ExecStreamRequest_Start_DeploymentId:
		log = log.With("deployment_id", t.DeploymentId)

		deployment, err := s.state(ctx).DeploymentGet(&pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{
				Id: t.DeploymentId,
			},
		})
		if err != nil {
			return hcerr.Externalize(
				log,
				err,
				"failed to get deployment for exec",
				"deployment_id",
				t.DeploymentId,
			)
		}

		// We need to spawn a job that will in turn spawn a virtual CEB
		// that will connect back and create an instance exec record for us
		// to use.
		if deployment.HasExecPlugin {
			instId, err := server.Id()
			if err != nil {
				return hcerr.Externalize(
					log,
					err,
					"failed to connect with server for exec",
				)
			}

			log.Info("spawning exec plugin via job", "instance-id", instId)

			job := &pb.Job{
				Workspace:   deployment.Workspace,
				Application: deployment.Application,
				Operation: &pb.Job_Exec{
					Exec: &pb.Job_ExecOp{
						InstanceId: instId,
						Deployment: deployment,
					},
				},
			}

			// Means the client WANTS the job run on itself, so let's target the
			// job back to it.
			if runnerId, ok := grpcmetadata.RunnerId(srv.Context()); ok {
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

				// Otherwise, the client wants an exec session but doesn't have a runner
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

			qresp, err := s.QueueJob(srv.Context(), &pb.QueueJobRequest{
				Job: job,

				// TODO unknown if this is enough time for when the request is queued
				// by a runner-less client but a user waiting 60 seconds will get impatient
				// regardless.
				ExpiresIn: "60s",
			})
			if err != nil {
				return hcerr.Externalize(
					log,
					err,
					"failed to connect with server for exec",
				)
			}

			jobId := qresp.JobId

			// Be sure that if we decide things aren't going well, the job doesn't outlive
			// its usefulness.
			defer s.state(ctx).JobCancel(jobId, false)

			log.Debug("waiting on job state", "job-id", jobId)

			state, err := s.waitOnJobStarted(srv.Context(), jobId)
			if err != nil {
				return hcerr.Externalize(
					log,
					err,
					"failed waiting for exec job",
				)
			}

			switch state {
			case pb.Job_ERROR:
				return status.Errorf(codes.FailedPrecondition, "job errored out before starting")
			case pb.Job_SUCCESS:
				return status.Errorf(codes.Internal, "job succeeded before running")
			case pb.Job_RUNNING:
				// ok
			default:
				return status.Errorf(codes.Internal, "unexpected job status: %s", state.String())
			}

			// If the virtual instance doesn't show up in 60 seconds, just time out and return
			// an error.
			ctx, cancel := context.WithTimeout(srv.Context(), 60*time.Second)
			defer cancel()

			err = iexec.InstanceExecCreateForVirtualInstance(ctx, instId, execRec)
			if err != nil {
				return hcerr.Externalize(
					log,
					err,
					"failed to create exec session",
				)
			}
		} else {
			err = iexec.InstanceExecCreateByDeployment(t.DeploymentId, execRec)
			if err != nil {
				return hcerr.Externalize(
					log,
					err,
					"failed to create exec session with deployment",
					"deployment_id",
					t.DeploymentId,
				)
			}
		}
	default:
		log.Error("exec request sent neither instance id nor deployment id")

		return status.Errorf(codes.FailedPrecondition,
			"request sent neither instance id nor deployment id")
	}

	log.Debug("exec requested", "args", start.Start.Args)

	// Make sure we always deregister it
	defer iexec.InstanceExecDelete(execRec.Id)

	// Always send the open message. In the future we'll send some metadata here.
	if err := srv.Send(&pb.ExecStreamResponse{
		Event: &pb.ExecStreamResponse_Open_{
			Open: &pb.ExecStreamResponse_Open{},
		},
	}); err != nil {
		return hcerr.Externalize(
			log,
			err,
			"failed to send open exec message",
		)
	}

	err = iexec.InstanceExecWaitConnected(ctx, execRec)
	if err != nil {
		return hcerr.Externalize(
			log,
			err,
			"exec session failed while waiting for connection",
		)
	}

	// Start our receive loop to read data from the client
	clientCloseCh := make(chan error, 1)
	go func() {
		defer close(clientEventCh)
		defer close(clientCloseCh)
		for {
			resp, err := srv.Recv()
			if err == io.EOF {
				// This means our client closed the stream. if the client
				// closed the stream, we want to end the exec stream completely.
				return
			}

			if err != nil {
				// Non EOF errors we will just send the error down and exit.
				clientCloseCh <- err
				return
			}

			clientEventCh <- resp
		}
	}()

	// Loop through and read events
	for {
		select {
		case <-srv.Context().Done():
			// The context was closed so we just exit. This will trigger
			// the EOF in the recv goroutine which will end the entrypoint
			// side as well.
			return nil

		case err := <-clientCloseCh:
			// The client closed the connection so we want to exit the stream.
			return hcerr.Externalize(
				log,
				err,
				"client closed exec session",
			)

		case entryReq, active := <-eventCh:
			// We got an event, exit out of the select and determine our action
			if !active {
				log.Debug("event channel closed, exiting")
				return nil
			}

			exit, err := s.handleEntrypointExecRequest(log, srv, entryReq)
			if exit || err != nil {
				return hcerr.Externalize(
					log,
					err,
					"error handling exec session",
				)
			}
		}
	}
}

func (s *Service) handleEntrypointExecRequest(
	log hclog.Logger,
	srv pb.Waypoint_StartExecStreamServer,
	entryReq *pb.EntrypointExecRequest,
) (bool, error) {
	log.Trace("event received from entrypoint", "event", entryReq.Event)
	var send *pb.ExecStreamResponse
	exit := false
	switch event := entryReq.Event.(type) {
	case *pb.EntrypointExecRequest_Output_:
		send = &pb.ExecStreamResponse{
			Event: &pb.ExecStreamResponse_Output_{
				Output: &pb.ExecStreamResponse_Output{
					Channel: pb.ExecStreamResponse_Output_Channel(event.Output.Channel),
					Data:    event.Output.Data,
				},
			},
		}

	case *pb.EntrypointExecRequest_Exit_:
		exit = true
		send = &pb.ExecStreamResponse{
			Event: &pb.ExecStreamResponse_Exit_{
				Exit: &pb.ExecStreamResponse_Exit{
					Code: event.Exit.Code,
				},
			},
		}
	case *pb.EntrypointExecRequest_Error_:
		log.Warn("error observed processing entrypoint exec stream", "error", event.Error.Error)
		exit = true
		send = &pb.ExecStreamResponse{
			Event: &pb.ExecStreamResponse_Exit_{
				Exit: &pb.ExecStreamResponse_Exit{
					Code: 1,
				},
			},
		}
	default:
		log.Warn("unimplemented exec entrypoint message seen", "event", hclog.Fmt("%T", event))
	}

	// Send our response
	if send != nil {
		if err := srv.Send(send); err != nil {
			log.Warn("stream error", "err", err)
			return false, err
		}
	}

	return exit, nil
}

// Wait for the given job to reach a state where it has been been acted upon in some manner.
func (s *Service) waitOnJobStarted(ctx context.Context, jobId string) (pb.Job_State, error) {
	log := hclog.FromContext(ctx)

	// Get the job
	ws := memdb.NewWatchSet()
	job, err := s.state(ctx).JobById(jobId, ws)
	if err != nil {
		return 0, err
	}
	if job == nil {
		return 0, status.Errorf(codes.NotFound, "job not found for ID: %s", jobId)
	}

	log = log.With("job_id", job.Id)

	for {
		switch job.State {
		case pb.Job_ERROR, pb.Job_RUNNING, pb.Job_SUCCESS:
			return job.State, nil
		}

		// Wait for the job to update
		if err := ws.WatchCtx(ctx); err != nil {
			if ctx.Err() != nil {
				return 0, ctx.Err()
			}

			return 0, err

		}

		// Updated job, requery it
		ws = memdb.NewWatchSet()
		job, err = s.state(ctx).JobById(job.Id, ws)
		if err != nil {
			return 0, err
		}
		if job == nil {
			return 0, status.Errorf(codes.Internal, "job disappeared for ID: %s", jobId)
		}
	}
}
