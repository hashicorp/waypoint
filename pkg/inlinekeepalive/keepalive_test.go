// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package inlinekeepalive

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint/pkg/protocolversion"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/gen/mocks"
)

// Happy path - server and client both support inline keepalives
func TestCompatibility_NewServerNewClient(t *testing.T) {
	impl := &serverImpl{
		features: []pb.ServerFeaturesFeature{pb.ServerFeatures_FEATURE_INLINE_KEEPALIVES},
	}
	sendInterval := time.Duration(50) * time.Millisecond

	addr := testServer(t, impl, grpc.ChainStreamInterceptor(KeepaliveServerStreamInterceptor(sendInterval)))
	client := testClient(t, addr, grpc.WithChainStreamInterceptor(KeepaliveClientStreamInterceptor(sendInterval)))

	runTest(t, impl, client, sendInterval)
}

// Server should not send inline keepalives to a client that does not have the interceptor configured
func TestCompatibility_NewServerOldClient(t *testing.T) {
	impl := &serverImpl{
		features: []pb.ServerFeaturesFeature{pb.ServerFeatures_FEATURE_INLINE_KEEPALIVES},
	}
	sendInterval := time.Duration(50) * time.Millisecond

	addr := testServer(t, impl, grpc.ChainStreamInterceptor(KeepaliveServerStreamInterceptor(sendInterval)))
	client := testClient(t, addr)

	runTest(t, impl, client, sendInterval)
}

// Client should not send keepalives to server that does not have the interceptor configured
func TestCompatibility_OldServerNewClient(t *testing.T) {
	impl := &serverImpl{}
	sendInterval := time.Duration(50) * time.Millisecond

	addr := testServer(t, impl)
	client := testClient(t, addr, grpc.WithChainStreamInterceptor(KeepaliveClientStreamInterceptor(sendInterval)))

	runTest(t, impl, client, sendInterval)
}

// Doesn't actually exercise any of the keepalive code, but ensures that the test fixture works.
func TestCompatibility_OldServerOldClient(t *testing.T) {
	impl := &serverImpl{}
	sendInterval := time.Duration(50) * time.Millisecond

	addr := testServer(t, impl)
	client := testClient(t, addr)

	runTest(t, impl, client, sendInterval)
}

// Simulates bidirectional communication on a grpc stream.
func runTest(t *testing.T, impl *serverImpl, client pb.WaypointClient, sendInterval time.Duration) {
	require := require.New(t)
	ctx := context.Background()

	resp1 := &pb.RunnerConfigResponse{
		Config: &pb.RunnerConfig{
			ConfigVars: []*pb.ConfigVar{{
				Name: "resp1",
			}},
		},
	}

	resp2 := &pb.RunnerConfigResponse{
		Config: &pb.RunnerConfig{
			ConfigVars: []*pb.ConfigVar{{
				Name: "resp2",
			}},
		},
	}

	req1 := &pb.RunnerConfigRequest{
		Event: &pb.RunnerConfigRequest_Open_{
			Open: &pb.RunnerConfigRequest_Open{
				Runner: &pb.Runner{Id: "req1"},
			},
		},
	}

	req2 := &pb.RunnerConfigRequest{
		Event: &pb.RunnerConfigRequest_Open_{
			Open: &pb.RunnerConfigRequest_Open{
				Runner: &pb.Runner{Id: "req2"},
			},
		},
	}

	rcfgClient, err := client.RunnerConfig(ctx)
	require.NoError(err)

	require.NoError(rcfgClient.Send(req1))

	// Ensure we've sent a keepalive
	time.Sleep(sendInterval * 2)

	impl.Send <- resp1

	// Ensure we've sent a keepalive
	time.Sleep(sendInterval * 2)

	require.NoError(rcfgClient.Send(req2))

	// Ensure we've sent a keepalive
	time.Sleep(sendInterval * 2)

	impl.Send <- resp2

	// Check the responses

	var actualRecv1 *pb.RunnerConfigRequest
	var actualRecv2 *pb.RunnerConfigRequest
	require.Eventually(func() bool {
		impl.Lock()
		defer impl.Unlock()
		if len(impl.Recv) == 2 {
			actualRecv1 = impl.Recv[0].(*pb.RunnerConfigRequest)
			actualRecv2 = impl.Recv[1].(*pb.RunnerConfigRequest)
			return true
		}

		return false
	}, 5*time.Second, 10*time.Millisecond)

	require.Equal(actualRecv1.Event.(*pb.RunnerConfigRequest_Open_).Open.Runner.Id, "req1")
	require.Equal(actualRecv2.Event.(*pb.RunnerConfigRequest_Open_).Open.Runner.Id, "req2")

	actualResp1, err := rcfgClient.Recv()
	require.NoError(err)
	require.Equal(actualResp1.Config.ConfigVars[0].Name, "resp1")

	// Again, ensure we've sent a keepalive
	time.Sleep(sendInterval * 2)

	actualResp2, err := rcfgClient.Recv()
	require.NoError(err)
	require.Equal(actualResp2.Config.ConfigVars[0].Name, "resp2")
}

func testClient(t *testing.T, addr string, opts ...grpc.DialOption) pb.WaypointClient {
	require := require.New(t)
	conn, err := grpc.Dial(addr, append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))...)
	require.NoError(err)

	return pb.NewWaypointClient(conn)
}

func testServer(t *testing.T, impl *serverImpl, opts ...grpc.ServerOption) (addr string) {
	require := require.New(t)

	impl.t = t
	impl.Send = make(chan proto.Message)

	ln, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(err)

	s := grpc.NewServer(opts...)
	pb.RegisterWaypointServer(s, impl)
	t.Cleanup(s.Stop)
	go s.Serve(ln)
	return ln.Addr().String()
}

type serverImpl struct {
	sync.Mutex
	mocks.WaypointServer
	pb.UnsafeWaypointServer
	//pb.UnimplementedWaypointServer

	t *testing.T

	// Features to serve on GetVersionInfo
	features []pb.ServerFeaturesFeature

	// Send is the list of responses to send
	Send chan proto.Message

	// Recv is the list of requests received
	Recv []proto.Message
}

// Using RunnerConfig as an arbitrary bi-directional stream.
func (s *serverImpl) RunnerConfig(
	srv pb.Waypoint_RunnerConfigServer,
) error {
	// Send down responses as we receive them
	go func() {
		for {
			msg := <-s.Send
			require.NoError(s.t, srv.Send(msg.(*pb.RunnerConfigResponse)))
		}
	}()

	for {
		req, err := srv.Recv()
		if err != nil {
			return err
		}

		s.Lock()
		s.Recv = append(s.Recv, req)
		s.Unlock()
	}

}

func (s *serverImpl) GetVersionInfo(_ context.Context, _ *empty.Empty) (*pb.GetVersionInfoResponse, error) {
	return &pb.GetVersionInfoResponse{
		Info:           protocolversion.Current(),
		ServerFeatures: &pb.ServerFeatures{Features: s.features},
	}, nil
}
