// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestServiceHostname(t *testing.T) {
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

		// Should have no hostnames
		{
			resp, err := client.ListHostnames(ctx, &pb.ListHostnamesRequest{})
			require.NoError(err)
			require.Empty(resp.Hostnames)
		}

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
		hostname := resp.Hostname

		// Should have the hostname
		{
			resp, err := client.ListHostnames(ctx, &pb.ListHostnamesRequest{})
			require.NoError(err)
			require.Len(resp.Hostnames, 1)
		}

		// Should be able to filter and have it
		{
			resp, err := client.ListHostnames(ctx, &pb.ListHostnamesRequest{
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
			require.Len(resp.Hostnames, 1)

			resp, err = client.ListHostnames(ctx, &pb.ListHostnamesRequest{
				Target: &pb.Hostname_Target{
					Target: &pb.Hostname_Target_Application{
						Application: &pb.Hostname_TargetApp{
							Application: &pb.Ref_Application{
								Application: "web2",
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
			require.Len(resp.Hostnames, 0)
		}

		// Can delete
		{
			_, err := client.DeleteHostname(ctx, &pb.DeleteHostnameRequest{
				Hostname: hostname.Hostname,
			})
			require.NoError(err)

			resp, err := client.ListHostnames(ctx, &pb.ListHostnamesRequest{})
			require.NoError(err)
			require.Empty(resp.Hostnames)
		}
	})
}
