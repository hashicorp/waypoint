package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestServiceCreateHostname(t *testing.T) {
	ctx := context.Background()

	t.Run("with no URL service", func(t *testing.T) {
		require := require.New(t)

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)
		client := server.TestServer(t, impl)

		// Create a hostname
		resp, err := client.CreateHostname(ctx, &pb.CreateHostnameRequest{
			Target: &pb.Hostname_Target{
				Target: &pb.Hostname_Target_Application{
					Application: &pb.Hostname_TargetApp{
						Application: &pb.Ref_Application{
							Application: "web",
							Project:     "test",
						},

						Workspace: &pb.Ref_Workspace{
							Workspace: "default",
						},
					},
				},
			},
		})
		require.Error(err)
		require.Nil(resp)
		require.Equal(codes.FailedPrecondition, status.Code(err))
	})

	t.Run("URL service", func(t *testing.T) {
		require := require.New(t)

		// Create our server
		impl, err := New(WithDB(testDB(t)), TestWithURLService(t, nil))
		require.NoError(err)
		client := server.TestServer(t, impl)

		// Create a hostname
		resp, err := client.CreateHostname(ctx, &pb.CreateHostnameRequest{
			Target: &pb.Hostname_Target{
				Target: &pb.Hostname_Target_Application{
					Application: &pb.Hostname_TargetApp{
						Application: &pb.Ref_Application{
							Application: "web",
							Project:     "test",
						},

						Workspace: &pb.Ref_Workspace{
							Workspace: "default",
						},
					},
				},
			},
		})
		require.NoError(err)
		require.NotNil(resp)
		require.NotEmpty(resp.Hostname)
	})
}
