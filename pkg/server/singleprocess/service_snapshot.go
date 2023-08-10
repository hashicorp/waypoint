// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package singleprocess

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
)

func (s *Service) CreateSnapshot(
	req *empty.Empty,
	srv pb.Waypoint_CreateSnapshotServer,
) error {
	ctx := srv.Context()
	// Always send the open message. In the future we'll send some metadata here.
	if err := srv.Send(&pb.CreateSnapshotResponse{
		Event: &pb.CreateSnapshotResponse_Open_{
			Open: &pb.CreateSnapshotResponse_Open{},
		},
	}); err != nil {
		return hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to send create snapshot message",
		)
	}

	// Create the snapshot and write the data
	bw := bufio.NewWriter(&snapshotWriter{srv: srv})
	if err := s.state(ctx).CreateSnapshot(ctx, bw); err != nil {
		return hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"error creating snapshot",
		)
	}

	return bw.Flush()
}

func (s *Service) RestoreSnapshot(
	srv pb.Waypoint_RestoreSnapshotServer,
) error {
	ctx := srv.Context()
	log := hclog.FromContext(srv.Context())

	// Read our first event which must be an Open event.
	log.Trace("waiting for Open message")
	req, err := srv.Recv()
	if err != nil {
		return hcerr.Externalize(
			log,
			err,
			"error reading snapshot events",
		)
	}
	open, ok := req.Event.(*pb.RestoreSnapshotRequest_Open_)
	if !ok {
		return status.Errorf(codes.FailedPrecondition,
			"first message must be Open type")
	}

	// Start our receive loop to read data from the client
	clientEventCh := make(chan *pb.RestoreSnapshotRequest)
	clientCloseCh := make(chan error, 1)
	go func() {
		defer close(clientEventCh)
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

	// Create a pipe for our data. We defer the writer close so we ensure
	// that the reader gets an EOF and ends always. Calling close on a pipe
	// multiple times is safe.
	pr, pw := io.Pipe()
	defer pw.Close()

	// Start our restore goroutine that is actively attempting the restore
	// process. This will ensure our reader end of our pipe is closed.
	restoreCloseCh := make(chan error, 1)
	go func() {
		defer close(restoreCloseCh)
		defer pr.Close()
		restoreCloseCh <- s.state(ctx).StageRestoreSnapshot(ctx, pr)
	}()

	// Buffer our writes so that we store some window of restore data in memory
	bw := bufio.NewWriterSize(pw, 1024*1024) // buffer 1 MB of write data

	// Loop through and read events
	for {
		select {
		case <-srv.Context().Done():
			// The context was closed so we just exit.
			return srv.Context().Err()

		case err := <-clientCloseCh:
			// The client closed the connection so we want to exit the stream.
			if err != nil {
				return hcerr.Externalize(
					log,
					err,
					"client closed connection",
				)
			}

			return srv.SendAndClose(&empty.Empty{})

		case err := <-restoreCloseCh:
			// The restore ended
			if err != nil {
				return hcerr.Externalize(
					log,
					err,
					"error during restore",
				)
			}

			// Restore was successful.
			if open.Open.Exit {
				log.Warn("restore requested exit, closing database and exiting NOW")
				s.state(ctx).Close()
				os.Exit(2) // kind of a weird exit code to note this was manufactured
			}

			return srv.SendAndClose(&empty.Empty{})

		case req, active := <-clientEventCh:
			// If we aren't active anymore, then the client closed the connection
			// so we close our write side and wait for the stage to complete
			if !active {
				log.Debug("event channel closed, waiting for restore to complete")

				// Close the write end so the reader gets an EOF and knows
				// we can continue with the restore.
				if err := bw.Flush(); err != nil {
					return err
				}
				if err := pw.Close(); err != nil {
					return err
				}

				// Set our event channel to nil so it blocks forever
				clientEventCh = nil
				continue
			}

			chunk, ok := req.Event.(*pb.RestoreSnapshotRequest_Chunk)
			if !ok {
				log.Info("restore received unexpected event",
					"type", fmt.Sprintf("%T", req.Event))
				return status.Errorf(codes.FailedPrecondition,
					"all messages after Open must be data chunks")
			}

			_, err := io.Copy(bw, bytes.NewReader(chunk.Chunk))
			if err != nil {
				return hcerr.Externalize(
					log,
					err,
					"error reading request data during restore",
				)
			}
		}
	}
}

type snapshotWriter struct {
	srv pb.Waypoint_CreateSnapshotServer
}

func (w *snapshotWriter) Write(p []byte) (int, error) {
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
