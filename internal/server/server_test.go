// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package server

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"

	serverpkg "github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	pbmocks "github.com/hashicorp/waypoint/pkg/server/gen/mocks"
)

func TestComponentEnum(t *testing.T) {
	for idx, name := range pb.Component_Type_name {
		// skip the invalid value
		if idx == 0 {
			continue
		}

		typ := component.Type(idx)
		require.Equal(t, strings.ToUpper(typ.String()), strings.ToUpper(name))
	}
}

type mockServer struct {
	pbmocks.WaypointServer

	pb.UnsafeWaypointServer
}

func TestRun_reconnect(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := &mockServer{}
	m.On("BootstrapToken", mock.Anything, mock.Anything).Return(&pb.NewTokenResponse{Token: "hello"}, nil)
	m.On("GetVersionInfo", mock.Anything, mock.Anything).Return(serverpkg.TestVersionInfoResponse(), nil)
	m.On("GetWorkspace", mock.Anything, mock.Anything).Return(&pb.GetWorkspaceResponse{}, nil)

	// Create the server
	restartCh := make(chan struct{})
	client := serverpkg.TestServer(t, m,
		serverpkg.TestWithContext(ctx),
		serverpkg.TestWithRestart(restartCh),
	)

	// Request should work
	_, err := client.GetWorkspace(ctx, &pb.GetWorkspaceRequest{
		Workspace: &pb.Ref_Workspace{
			Workspace: "test",
		},
	})
	require.NoError(err)

	// Stop it
	cancel()

	// Should not work
	require.Eventually(func() bool {
		_, err := client.GetWorkspace(context.Background(), &pb.GetWorkspaceRequest{
			Workspace: &pb.Ref_Workspace{
				Workspace: "test",
			},
		})
		t.Logf("error: %s", err)
		return status.Code(err) == codes.Unavailable
	}, 2*time.Second, 10*time.Millisecond)

	// Restart
	restartCh <- struct{}{}

	// Should work
	require.Eventually(func() bool {
		_, err := client.GetWorkspace(context.Background(), &pb.GetWorkspaceRequest{
			Workspace: &pb.Ref_Workspace{
				Workspace: "test",
			},
		})
		return err == nil
	}, 5*time.Second, 10*time.Millisecond)
}
