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
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	pbmocks "github.com/hashicorp/waypoint/internal/server/gen/mocks"
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

func TestRun_reconnect(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := &pbmocks.WaypointServer{}
	m.On("BootstrapToken", mock.Anything, mock.Anything).Return(&pb.NewTokenResponse{Token: "hello"}, nil)
	m.On("GetVersionInfo", mock.Anything, mock.Anything).Return(testVersionInfoResponse(), nil)
	m.On("GetWorkspace", mock.Anything, mock.Anything).Return(&pb.GetWorkspaceResponse{}, nil)

	// Create the server
	restartCh := make(chan struct{})
	client := TestServer(t, m,
		TestWithContext(ctx),
		TestWithRestart(restartCh),
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
