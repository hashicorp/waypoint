// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package clisnapshot provides access for our CLI commands to create and
// restore snapshots
package clisnapshot

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// WriteSnapshot requests a snapshot from the client and writes it to the
// provided writer. Cancelling the context will prematurely cancel the snapshot.
// This may result in partial writes to the writer.
func WriteSnapshot(ctx context.Context, client pb.WaypointClient, w io.Writer) error {
	stream, err := client.CreateSnapshot(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to generate snapshot: %s", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		return status.Error(status.Code(err), fmt.Sprintf("failed to receive snapshot start message: %s", err))
	}

	if _, ok := resp.Event.(*pb.CreateSnapshotResponse_Open_); !ok {
		return status.Error(status.Code(err), fmt.Sprintf("failed to receive snapshot start message: %s", err))
	}

	for {
		ev, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return status.Error(status.Code(err), fmt.Sprintf("error receiving snapshot data: %s", err))
		}

		chunk, ok := ev.Event.(*pb.CreateSnapshotResponse_Chunk)
		if ok {
			_, err = w.Write(chunk.Chunk)
			if err != nil {
				return status.Error(status.Code(err), fmt.Sprintf("error writing snapshot data: %s", err))
			}
		} else {
			return status.Error(status.Code(err), fmt.Sprintf("unexpected protocol value: %T", ev.Event))
		}
	}
	return nil
}

// ReadSnapshot stages a snapshot for restore from the provided reader, and
// sends an exit signal to the server if 'exit' is true. Cancelling the context
// will prematurely cancel the snapshot restore. This may result in a partial
// restore from the reader being staged.
func ReadSnapshot(ctx context.Context, client pb.WaypointClient, r io.Reader, exit bool) error {
	stream, err := client.RestoreSnapshot(ctx)
	if err != nil {
		return status.Error(status.Code(err), fmt.Sprintf("failed to restore snapshot: %s", err))
	}

	err = stream.Send(&pb.RestoreSnapshotRequest{
		Event: &pb.RestoreSnapshotRequest_Open_{
			Open: &pb.RestoreSnapshotRequest_Open{
				Exit: exit,
			},
		},
	})
	if err != nil {
		return status.Error(status.Code(err), fmt.Sprintf("failed to send start message: %s", err))
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
			return status.Error(status.Code(err), fmt.Sprintf("failed to write snapshot data: %s", err))
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil && !exit {
		return status.Error(status.Code(err), fmt.Sprintf("failed to receive snapshot close message: %s", err))
	}
	return nil
}
