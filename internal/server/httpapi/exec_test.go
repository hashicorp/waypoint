// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package httpapi

import (
	"context"
	"io"
	"net"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wspb"

	"github.com/hashicorp/waypoint/pkg/protocolversion"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/gen/mocks"
)

// This code uses the magic token value 445DHu. This value is a base58 encoded
// empty token. An empty token has the magic 'wp24' at the beginning, so this
// empty token is just base58.Encode([]byte("wp24"))

func TestHandleExec(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Get our gRPC server
	impl := &execImpl{}
	addr := testServer(t, impl)

	// Start up our test HTTP server
	httpServer := httptest.NewServer(HandleExec(addr, false))
	defer httpServer.Close()

	// Dial it up
	conn, _, err := websocket.Dial(ctx, httpServer.URL+"?token=445DHu", nil)
	require.NoError(err)
	defer conn.Close(websocket.StatusInternalError, "early exit")

	// Send our start request
	require.NoError(wspb.Write(ctx, conn, &pb.ExecStreamRequest{
		Event: &pb.ExecStreamRequest_Start_{
			Start: &pb.ExecStreamRequest_Start{
				Args: []string{"foo", "bar"},
			},
		},
	}))

	// We should receive them, eventually.
	var value *pb.ExecStreamRequest
	require.Eventually(func() bool {
		impl.Lock()
		defer impl.Unlock()
		if len(impl.Recv) == 1 {
			value = impl.Recv[0]
			return true
		}

		return false
	}, 5*time.Second, 10*time.Millisecond)

	// It should be our start request
	startReq := value.Event.(*pb.ExecStreamRequest_Start_).Start
	require.Equal([]string{"foo", "bar"}, startReq.Args)
}

func testServer(t *testing.T, impl pb.WaypointServer) string {
	// Create a listener
	ln, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	t.Cleanup(func() { ln.Close() })

	// Register our gRPC service
	s := grpc.NewServer()
	pb.RegisterWaypointServer(s, impl)
	t.Cleanup(s.Stop)
	go s.Serve(ln)

	return ln.Addr().String()
}

type execImpl struct {
	sync.Mutex
	mocks.WaypointServer
	pb.UnsafeWaypointServer

	// Send is the list of responses to send
	Send []*pb.ExecStreamResponse

	// Recv is the list of requests received
	Recv []*pb.ExecStreamRequest
}

// InlineKeepaliveInterceptor may call this
func (v *execImpl) GetVersionInfo(_ context.Context, _ *empty.Empty) (*pb.GetVersionInfoResponse, error) {
	return &pb.GetVersionInfoResponse{
		Info: protocolversion.Current(),
	}, nil
}

func (v *execImpl) StartExecStream(srv pb.Waypoint_StartExecStreamServer) error {
	// Send down all our responses
	for _, resp := range v.Send {
		if err := srv.Send(resp); err != nil {
			return err
		}
	}

	// Receive forever
	for {
		req, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		v.Lock()
		v.Recv = append(v.Recv, req)
		v.Unlock()
	}
}
