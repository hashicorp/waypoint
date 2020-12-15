package snapshot

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

type Config struct {
	// This Client points to a project's client, as set in the baseCommand
	Client pb.WaypointClient
}

func (c *Config) WriteSnapshot(ctx context.Context, w io.Writer) error {
	stream, err := c.Client.CreateSnapshot(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to generate snapshot: %s", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("failed to receive snapshot start message: %s", err)
	}

	if _, ok := resp.Event.(*pb.CreateSnapshotResponse_Open_); !ok {
		return fmt.Errorf("failed to receive snapshot start message: %s", err)
	}

	for {
		ev, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error receiving snapshot data: %s", err)
		}

		chunk, ok := ev.Event.(*pb.CreateSnapshotResponse_Chunk)
		if ok {
			_, err = w.Write(chunk.Chunk)
			if err != nil {
				return fmt.Errorf("error writing snapshot data: %s", err)
			}
		} else {
			return fmt.Errorf("unexpected protocol value: %T", ev.Event)
		}
	}
	return nil
}

func (c *Config) ReadSnapshot(ctx context.Context, r io.Reader, exit bool) error {
	stream, err := c.Client.RestoreSnapshot(ctx)
	if err != nil {
		return fmt.Errorf("failed to restore snapshot: %s", err)
	}

	err = stream.Send(&pb.RestoreSnapshotRequest{
		Event: &pb.RestoreSnapshotRequest_Open_{
			Open: &pb.RestoreSnapshotRequest_Open{
				Exit: exit,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %s", err)
	}

	// Write the data in smaller chunks so we don't overwhelm the grpc stream
	// processing machinary.
	var buf [1024]byte

	for {
		// use ReadFull here because if r is an OS pipe, each bare call to Read()
		// can result in just one or two bytes per call, so we want to batch those
		// up before sending them off for better performance.
		n, err := io.ReadFull(r, buf[:])
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = nil
		}

		if n == 0 {
			break
		}

		err = stream.Send(&pb.RestoreSnapshotRequest{
			Event: &pb.RestoreSnapshotRequest_Chunk{
				Chunk: buf[:n],
			},
		})
		if err != nil {
			return fmt.Errorf("failed to write snapshot data: %s", err)
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil && !exit {
		return fmt.Errorf("failed to receive snapshot close message: %s", err)
	}
	return nil
}
