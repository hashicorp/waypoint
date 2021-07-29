package singleprocess

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/logbuffer"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

// TODO: test
func (s *service) GetJob(
	ctx context.Context,
	req *pb.GetJobRequest,
) (*pb.Job, error) {
	job, err := s.state.JobById(req.JobId, nil)
	if err != nil {
		return nil, err
	}
	if job == nil || job.Job == nil {
		return nil, status.Errorf(codes.NotFound, "job not found")
	}

	return job.Job, nil
}

// TODO: test
func (s *service) XListJobs(
	ctx context.Context,
	req *pb.ListJobsRequest,
) (*pb.ListJobsResponse, error) {
	jobs, err := s.state.JobList()
	if err != nil {
		return nil, err
	}

	return &pb.ListJobsResponse{
		Jobs: jobs,
	}, nil
}

func (s *service) CancelJob(
	ctx context.Context,
	req *pb.CancelJobRequest,
) (*empty.Empty, error) {
	if err := s.state.JobCancel(req.JobId, false); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *service) QueueJob(
	ctx context.Context,
	req *pb.QueueJobRequest,
) (*pb.QueueJobResponse, error) {
	job := req.Job

	// Validation
	if job == nil {
		return nil, status.Errorf(codes.FailedPrecondition, "job must be set")
	}
	if job.Operation == nil {
		// We special case this check and return "Unimplemented" because
		// the primary case where operation is nil is if a client is sending
		// us an unsupported operation.
		return nil, status.Errorf(codes.Unimplemented, "operation is nil or unknown")
	}
	if err := serverptypes.ValidateJob(job); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}

	// Verify the project exists and use that to set the default data source
	project, err := s.state.ProjectGet(&pb.Ref_Project{Project: job.Application.Project})
	if status.Code(err) == codes.NotFound {
		return nil, status.Errorf(codes.NotFound,
			"Project %s was not found! Please ensure that 'waypoint init' was run with this project.",
			job.Application.Project,
		)
	}

	if job.DataSource == nil {
		if project.DataSource == nil {
			return nil, status.Errorf(codes.FailedPrecondition,
				"Project %s does not have a data source configured. Remote jobs "+
					"require a data source such as Git to be configured with the project. "+
					"Data sources can be configured via the CLI or UI. For help, see : "+
					"https://www.waypointproject.io/docs/projects/git#configuring-the-project",
				job.Application.Project,
			)
		}

		job.DataSource = project.DataSource
	}

	// Get the next id
	id, err := server.Id()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "uuid generation failed: %s", err)
	}
	job.Id = id

	// Validate expiry if we have one
	job.ExpireTime = nil
	if req.ExpiresIn != "" {
		dur, err := time.ParseDuration(req.ExpiresIn)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition,
				"Invalid expiry duration: %s", err.Error())
		}

		job.ExpireTime, err = ptypes.TimestampProto(time.Now().Add(dur))
		if err != nil {
			return nil, status.Errorf(codes.Aborted, "error configuring expiration: %s", err)
		}
	}

	// Queue the job
	if err := s.state.JobCreate(job); err != nil {
		return nil, err
	}

	return &pb.QueueJobResponse{JobId: job.Id}, nil
}

func (s *service) ValidateJob(
	ctx context.Context,
	req *pb.ValidateJobRequest,
) (*pb.ValidateJobResponse, error) {
	var err error
	result := &pb.ValidateJobResponse{Valid: true}

	// Struct validation
	if err := serverptypes.ValidateJob(req.Job); err != nil {
		result.Valid = false
		result.ValidationError = status.New(codes.FailedPrecondition, err.Error()).Proto()
		return result, nil
	}

	// Check assignability
	result.Assignable, err = s.state.JobIsAssignable(ctx, req.Job)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetJobStream(
	req *pb.GetJobStreamRequest,
	server pb.Waypoint_GetJobStreamServer,
) error {
	log := hclog.FromContext(server.Context())
	ctx := server.Context()

	// Get the job
	ws := memdb.NewWatchSet()
	job, err := s.state.JobById(req.JobId, ws)
	if err != nil {
		return err
	}
	if job == nil {
		return status.Errorf(codes.NotFound, "job not found for ID: %s", req.JobId)
	}
	log = log.With("job_id", job.Id)

	// We always send the open message as confirmation the job was found.
	if err := server.Send(&pb.GetJobStreamResponse{
		Event: &pb.GetJobStreamResponse_Open_{
			Open: &pb.GetJobStreamResponse_Open{},
		},
	}); err != nil {
		return err
	}

	// Start a goroutine that watches for job changes
	jobCh := make(chan *state.Job, 1)
	errCh := make(chan error, 1)
	go func() {
		for {
			// Send the job
			select {
			case jobCh <- job:
			case <-ctx.Done():
				return
			}

			// Wait for the job to update
			if err := ws.WatchCtx(ctx); err != nil {
				if ctx.Err() == nil {
					errCh <- err
				}

				return
			}

			// Updated job, requery it
			ws = memdb.NewWatchSet()
			job, err = s.state.JobById(job.Id, ws)
			if err != nil {
				errCh <- err
				return
			}
			if job == nil {
				errCh <- status.Errorf(codes.Internal, "job disappeared for ID: %s", req.JobId)
				return
			}
		}
	}()

	// Track that we only send these events once. We could use a bitmask for
	// this if we cared about that level of optimization but it hurts readability
	// and we don't need the performance yet.
	var cancelSent bool
	var downloadSent bool

	// Enter the event loop
	var lastState pb.Job_State
	var eventsCh <-chan []*pb.GetJobStreamResponse_Terminal_Event
	for {
		select {
		case <-ctx.Done():
			return nil

		case err := <-errCh:
			return err

		case job := <-jobCh:
			log.Debug("job state change", "state", job.State)

			// If we have a state change, send that event down. We also send
			// down a state change if we enter a "cancelled" scenario.
			canceling := job.CancelTime != nil
			if lastState != job.State || cancelSent != canceling {
				if err := server.Send(&pb.GetJobStreamResponse{
					Event: &pb.GetJobStreamResponse_State_{
						State: &pb.GetJobStreamResponse_State{
							Previous:  lastState,
							Current:   job.State,
							Job:       job.Job,
							Canceling: canceling,
						},
					},
				}); err != nil {
					return err
				}

				lastState = job.State
				cancelSent = canceling
			}

			// If we have a data source ref set, then we need to send the download event.
			if !downloadSent && job.DataSourceRef != nil {
				if err := server.Send(&pb.GetJobStreamResponse{
					Event: &pb.GetJobStreamResponse_Download_{
						Download: &pb.GetJobStreamResponse_Download{
							DataSourceRef: job.DataSourceRef,
						},
					},
				}); err != nil {
					return err
				}

				downloadSent = true
			}

			// If we haven't initialized output streaming and the output buffer
			// is now non-nil, initialize that. This will send any buffered
			// data down.
			if eventsCh == nil && job.OutputBuffer != nil {
				eventsCh, err = s.getJobStreamOutputInit(ctx, job, server)
				if err != nil {
					return err
				}
			}

			switch job.State {
			case pb.Job_SUCCESS, pb.Job_ERROR:
				// TODO(mitchellh): we should drain the output buffer

				// Job is done. For success, error will be nil, so this
				// populates the event with the proper values.
				return server.Send(&pb.GetJobStreamResponse{
					Event: &pb.GetJobStreamResponse_Complete_{
						Complete: &pb.GetJobStreamResponse_Complete{
							Error:  job.Error,
							Result: job.Result,
						},
					},
				})
			}

		case events := <-eventsCh:
			if err := server.Send(&pb.GetJobStreamResponse{
				Event: &pb.GetJobStreamResponse_Terminal_{
					Terminal: &pb.GetJobStreamResponse_Terminal{
						Events: events,
					},
				},
			}); err != nil {
				return err
			}
		}
	}
}

func (s *service) readJobLogBatch(r *logbuffer.Reader, block bool) []*pb.GetJobStreamResponse_Terminal_Event {
	entries := r.Read(64, block)
	if entries == nil {
		return nil
	}

	events := make([]*pb.GetJobStreamResponse_Terminal_Event, len(entries))
	for i, entry := range entries {
		events[i] = entry.(*pb.GetJobStreamResponse_Terminal_Event)
	}

	return events
}

func (s *service) getJobStreamOutputInit(
	ctx context.Context,
	job *state.Job,
	server pb.Waypoint_GetJobStreamServer,
) (<-chan []*pb.GetJobStreamResponse_Terminal_Event, error) {
	// Send down all our buffered lines.
	outputR := job.OutputBuffer.Reader(-1)
	go outputR.CloseContext(ctx)
	for {
		events := s.readJobLogBatch(outputR, false)
		if events == nil {
			break
		}

		if err := server.Send(&pb.GetJobStreamResponse{
			Event: &pb.GetJobStreamResponse_Terminal_{
				Terminal: &pb.GetJobStreamResponse_Terminal{
					Events:   events,
					Buffered: true,
				},
			},
		}); err != nil {
			return nil, err
		}
	}

	// Start a goroutine that reads output
	eventsCh := make(chan []*pb.GetJobStreamResponse_Terminal_Event, 1)
	go func() {
		for {
			events := s.readJobLogBatch(outputR, true)
			if events == nil {
				return
			}

			select {
			case eventsCh <- events:
			case <-ctx.Done():
				return
			}
		}
	}()

	return eventsCh, nil
}
