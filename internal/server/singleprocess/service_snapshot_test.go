package singleprocess

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestServiceRestoreSnapshot_badOpen(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Start exec with a bad starting message
	stream, err := client.RestoreSnapshot(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.RestoreSnapshotRequest{
		Event: &pb.RestoreSnapshotRequest_Chunk{
			Chunk: []byte("Hello"),
		},
	}))

	// Wait for data
	resp, err := stream.CloseAndRecv()
	require.Error(err)
	require.Equal(codes.FailedPrecondition, status.Code(err))
	require.Nil(resp)
}

func TestServiceRestoreSnapshot_full(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Take a snapshot and write the contents to a buf
	var snapshotBuf bytes.Buffer
	{
		stream, err := client.CreateSnapshot(ctx, &empty.Empty{})
		require.NoError(err)

		// Should get the open message
		resp, err := stream.Recv()
		require.NoError(err)
		require.IsType((*pb.CreateSnapshotResponse_Open_)(nil), resp.Event)

		// Get all the data
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			require.NoError(err)
			require.IsType((*pb.CreateSnapshotResponse_Chunk)(nil), resp.Event)

			_, err = io.Copy(&snapshotBuf, bytes.NewReader(
				resp.Event.(*pb.CreateSnapshotResponse_Chunk).Chunk))
			require.NoError(err)
		}
	}
	t.Logf("snapshot length: %d", snapshotBuf.Len())

	// Restore
	stream, err := client.RestoreSnapshot(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.RestoreSnapshotRequest{
		Event: &pb.RestoreSnapshotRequest_Open_{
			Open: &pb.RestoreSnapshotRequest_Open{},
		},
	}))

	for {
		var buf [1024]byte
		n, err := snapshotBuf.Read(buf[:])
		if err == io.EOF {
			err = nil
		}
		require.NoError(err)
		if n == 0 {
			t.Log("finished writing restore data")
			break
		}

		t.Logf("writing log data, len: %d", n)
		require.NoError(stream.Send(&pb.RestoreSnapshotRequest{
			Event: &pb.RestoreSnapshotRequest_Chunk{
				Chunk: buf[:n],
			},
		}))
	}

	resp, err := stream.CloseAndRecv()
	require.NoError(err)
	require.NotNil(resp)
}
