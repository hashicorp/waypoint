package singleprocess

import (
	"bufio"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// TODO: test
func (s *service) CreateSnapshot(
	req *empty.Empty,
	srv pb.Waypoint_CreateSnapshotServer,
) error {
	// Always send the open message. In the future we'll send some metadata here.
	if err := srv.Send(&pb.CreateSnapshotResponse{
		Event: &pb.CreateSnapshotResponse_Open_{
			Open: &pb.CreateSnapshotResponse_Open{},
		},
	}); err != nil {
		return err
	}

	// Create the snapshot and write the data
	if err := s.state.CreateSnapshot(bufio.NewWriter(&snapshotWriter{
		srv: srv,
	})); err != nil {
		return err
	}

	return nil
}

type snapshotWriter struct {
	srv pb.Waypoint_CreateSnapshotServer
}

func (w *snapshotWriter) Write(p []byte) (nn int, err error) {
	// Ignore empty data.
	if len(p) == 0 {
		return 0, nil
	}

	// Write our data
	return len(p), w.srv.Send(&pb.CreateSnapshotResponse{
		Event: &pb.CreateSnapshotResponse_Chunk{
			Chunk: p,
		},
	})
}
