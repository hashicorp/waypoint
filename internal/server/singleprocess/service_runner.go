package singleprocess

import (
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

// TODO: test
func (s *service) RunnerConfig(
	req *pb.RunnerConfigRequest,
	srv pb.Waypoint_RunnerConfigServer,
) error {
	log := hclog.FromContext(srv.Context())

	// Create our record
	log = log.With("runner_id", req.Id)
	log.Trace("registering runner")
	record := &state.Runner{
		Id: req.Id,
	}
	if err := s.state.RunnerCreate(record); err != nil {
		return err
	}

	// Defer deleting this.
	// TODO(mitchellh): this is too aggressive and we want to have some grace
	// period for reconnecting clients. We should clean this up.
	defer func() {
		log.Trace("deleting runner")
		if err := s.state.RunnerDelete(record.Id); err != nil {
			log.Error("failed to delete runner data. This should not happen.", "err", err)
		}
	}()

	// Build our config in a loop.
	for {
		// Build our config
		config := &pb.RunnerConfig{}

		// Send new config
		if err := srv.Send(&pb.RunnerConfigResponse{
			Config: config,
		}); err != nil {
			return err
		}

		// Nil out the stuff we used so that if we're waiting awhile we can GC
		config = nil

		// We don't ever currently have config changes so we just block
		// until we're done. But soon we'll have config changes.
		<-srv.Context().Done()
	}
}

// TODO: test
func (s *service) RunnerJobStream(
	server pb.Waypoint_RunnerJobStreamServer,
) error {
	log := hclog.FromContext(server.Context())

	// Receive our opening message so we can determine the runner ID.
	req, err := server.Recv()
	if err != nil {
		return err
	}
	reqEvent, ok := req.Event.(*pb.RunnerJobStreamRequest_Request_)
	if !ok {
		return status.Errorf(codes.FailedPrecondition,
			"first message must be a Request event")
	}

	log = log.With("runner_id", reqEvent.Request.RunnerId)

	return nil
}
