// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ceb

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"
)

func TestLogs_basic(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start up the server
	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := server.TestServer(t, impl,
		server.TestWithContext(ctx),
		server.TestWithRestart(restartCh),
	)

	// Start the CEB
	ceb := testRun(t, context.Background(), &testRunOpts{
		Client: client,
		Helper: "logs-stdout",
	})

	// We should get registered
	require.Eventually(func() bool {
		resp, err := client.ListInstances(ctx, &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: ceb.DeploymentId(),
			},
		})
		require.NoError(err)
		return len(resp.Instances) == 1
	}, 2*time.Second, 10*time.Millisecond)

	// Get the log stream
	stream, err := client.GetLogStream(ctx, &pb.GetLogStreamRequest{
		Scope: &pb.GetLogStreamRequest_DeploymentId{
			DeploymentId: ceb.DeploymentId(),
		},
	})
	require.NoError(err)

	// Get a batch
	require.Eventually(func() bool {
		batch, err := stream.Recv()
		require.NoError(err)
		return len(batch.Lines) > 0
	}, 1*time.Second, 50*time.Millisecond)
}

func TestLogs_reconnect(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start up the server
	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := server.TestServer(t, impl,
		server.TestWithContext(ctx),
		server.TestWithRestart(restartCh),
	)

	// Start the CEB
	runCtx, runCancel := context.WithCancel(context.Background())
	defer runCancel()
	ceb := testRun(t, runCtx, &testRunOpts{
		Client: client,
		Helper: "logs-stdout",
	})

	// We should get registered
	require.Eventually(func() bool {
		resp, err := client.ListInstances(ctx, &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: ceb.DeploymentId(),
			},
		})
		require.NoError(err)
		return len(resp.Instances) == 1
	}, 2*time.Second, 10*time.Millisecond)

	// Get the log stream
	stream, err := client.GetLogStream(ctx, &pb.GetLogStreamRequest{
		Scope: &pb.GetLogStreamRequest_DeploymentId{
			DeploymentId: ceb.DeploymentId(),
		},
	})
	require.NoError(err)

	// Get a batch
	require.Eventually(func() bool {
		batch, err := stream.Recv()
		require.NoError(err)
		return len(batch.Lines) > 0
	}, 1*time.Second, 50*time.Millisecond)

	// Shutdown the server
	cancel()
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// We should get deregistered
	require.Eventually(func() bool {
		resp, err := impl.ListInstances(ctx, &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: ceb.DeploymentId(),
			},
		})
		require.NoError(err)
		return len(resp.Instances) == 0
	}, 5*time.Second, 10*time.Millisecond)

	// Restart
	restartCh <- struct{}{}

	// Get the log stream
	stream, err = client.GetLogStream(ctx, &pb.GetLogStreamRequest{
		Scope: &pb.GetLogStreamRequest_DeploymentId{
			DeploymentId: ceb.DeploymentId(),
		},
	}, grpc.WaitForReady(true))
	require.NoError(err)

	// Get a batch
	require.Eventually(func() bool {
		batch, err := stream.Recv()
		require.NoError(err)
		return len(batch.Lines) > 0
	}, 1*time.Second, 50*time.Millisecond)
}
