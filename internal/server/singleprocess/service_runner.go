package singleprocess

import (
	"github.com/hashicorp/go-hclog"

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
