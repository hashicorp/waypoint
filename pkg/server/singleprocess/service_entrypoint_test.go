// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package singleprocess

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

func TestServiceEntrypointConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("deployment info", func(t *testing.T) {
		require := require.New(t)

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)
		client := server.TestServer(t, impl)

		// Create a deployment
		resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
				Component: &pb.Component{
					Name: "testapp",
				},
			}),
		})
		require.NoError(err)
		dep := resp.Deployment

		// Create the config
		instanceId, err := server.Id()
		require.NoError(err)
		stream, err := client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
			InstanceId:   instanceId,
			DeploymentId: dep.Id,
		})
		require.NoError(err)

		// Wait for the first config so that we know we're registered
		cfgResp, err := stream.Recv()
		require.NoError(err)

		// Validate config
		require.NotNil(cfgResp.Config.Deployment)
		require.Equal("testapp", cfgResp.Config.Deployment.Component.Name)
	})

	t.Run("no URL service", func(t *testing.T) {
		require := require.New(t)

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)
		client := server.TestServer(t, impl)

		// Create a deployment
		resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
				Component: &pb.Component{
					Name: "testapp",
				},
			}),
		})
		require.NoError(err)
		dep := resp.Deployment

		// Create the config
		instanceId, err := server.Id()
		require.NoError(err)
		stream, err := client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
			InstanceId:   instanceId,
			DeploymentId: dep.Id,
		})
		require.NoError(err)

		// Wait for the first config so that we know we're registered
		cfgResp, err := stream.Recv()
		require.NoError(err)

		// Validate config
		require.Nil(cfgResp.Config.UrlService)
	})

	t.Run("URL service", func(t *testing.T) {
		require := require.New(t)

		// Create our server
		impl, err := New(WithDB(testDB(t)), TestWithURLService(t, nil))
		require.NoError(err)
		client := server.TestServer(t, impl)

		// Create a deployment
		resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
				Component: &pb.Component{
					Name: "testapp",
				},

				Labels: map[string]string{
					"hello": "world",
				},
			}),
		})
		require.NoError(err)
		dep := resp.Deployment

		// Create the config
		instanceId, err := server.Id()
		require.NoError(err)
		stream, err := client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
			InstanceId:   instanceId,
			DeploymentId: dep.Id,
		})
		require.NoError(err)

		// Wait for the first config so that we know we're registered
		cfgResp, err := stream.Recv()
		require.NoError(err)

		// Validate config
		require.NotNil(cfgResp.Config.UrlService)
		require.NotEmpty(cfgResp.Config.UrlService.Token)
		require.NotEmpty(cfgResp.Config.UrlService.Labels)
	})

	t.Run("URL service with guest account", func(t *testing.T) {
		require := require.New(t)

		// Create our server
		impl, err := New(
			WithDB(testDB(t)),
			TestWithURLService(t, nil),
			TestWithURLServiceGuestAccount(t),
		)
		require.NoError(err)
		client := server.TestServer(t, impl)

		// Create a deployment
		resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
				Component: &pb.Component{
					Name: "testapp",
				},

				Labels: map[string]string{
					"hello": "world",
				},
			}),
		})
		require.NoError(err)
		dep := resp.Deployment

		// Create the config
		instanceId, err := server.Id()
		require.NoError(err)
		stream, err := client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
			InstanceId:   instanceId,
			DeploymentId: dep.Id,
		})
		require.NoError(err)

		// Wait for the first config so that we know we're registered
		cfgResp, err := stream.Recv()
		require.NoError(err)

		// Validate config
		require.NotNil(cfgResp.Config.UrlService)
		require.NotEmpty(cfgResp.Config.UrlService.Token)
		require.NotEmpty(cfgResp.Config.UrlService.Labels)
	})

	t.Run("config sources", func(t *testing.T) {
		require := require.New(t)

		// Create our server
		impl, err := New(WithDB(testDB(t)))
		require.NoError(err)
		client := server.TestServer(t, impl)

		// Create a deployment
		resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
				Component: &pb.Component{
					Name: "testapp",
				},
			}),
		})
		require.NoError(err)
		dep := resp.Deployment

		// Create a config source
		{
			_, err := client.SetConfigSource(ctx, &pb.SetConfigSourceRequest{
				ConfigSource: &pb.ConfigSource{
					Scope: &pb.ConfigSource_Global{
						Global: &pb.Ref_Global{},
					},

					Type: "foo",
					Config: map[string]string{
						"value": "42",
					},
				},
			})
			require.NoError(err)
		}

		// Create a static config
		{
			_, err := client.SetConfig(ctx, &pb.ConfigSetRequest{
				Variables: []*pb.ConfigVar{
					{
						Target: &pb.ConfigVar_Target{
							AppScope: &pb.ConfigVar_Target_Application{
								Application: dep.Application,
							},
						},

						Name:  "DATABASE_URL",
						Value: &pb.ConfigVar_Static{Static: "postgresql:///"},
					},
				},
			})
			require.NoError(err)
		}

		// Create the config
		instanceId, err := server.Id()
		require.NoError(err)
		stream, err := client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
			InstanceId:   instanceId,
			DeploymentId: dep.Id,
		})
		require.NoError(err)

		{
			// Wait for the first config so that we know we're registered
			cfgResp, err := stream.Recv()
			require.NoError(err)

			// Validate config
			require.NotEmpty(cfgResp.Config.EnvVars)
			require.Empty(cfgResp.Config.ConfigSources)
		}

		// Create a dynamic config
		{
			_, err := client.SetConfig(ctx, &pb.ConfigSetRequest{
				Variables: []*pb.ConfigVar{
					{
						Target: &pb.ConfigVar_Target{
							AppScope: &pb.ConfigVar_Target_Application{
								Application: dep.Application,
							},
						},

						Name: "DATABASE_URL",
						Value: &pb.ConfigVar_Dynamic{
							Dynamic: &pb.ConfigVar_DynamicVal{
								From: "foo",
							},
						},
					},
				},
			})
			require.NoError(err)
		}

		{
			// Next config
			cfgResp, err := stream.Recv()
			require.NoError(err)

			// Validate config
			require.NotEmpty(cfgResp.Config.EnvVars)
			require.NotEmpty(cfgResp.Config.ConfigSources)
		}
	})
}

func TestServiceEntrypointExecStream_badOpen(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Start exec with a bad starting message
	stream, err := client.EntrypointExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Output_{
			Output: &pb.EntrypointExecRequest_Output{},
		},
	}))

	// Wait for data
	resp, err := stream.Recv()
	require.Error(err)
	require.Equal(codes.FailedPrecondition, status.Code(err))
	require.Nil(resp)
}

func TestServiceEntrypointExecStream_invalidInstanceId(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	// Start exec with a bad starting message
	stream, err := client.EntrypointExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Open_{
			Open: &pb.EntrypointExecRequest_Open{
				InstanceId: "nope",
				Index:      0,
			},
		},
	}))

	// Wait for data
	resp, err := stream.Recv()
	require.Error(err)
	require.Equal(codes.NotFound, status.Code(err))
	require.Nil(resp)
}

func TestServiceEntrypointExecStream_invalidSessionId(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	exec, closer := testRegisterExec(ctx, t, client, impl)
	defer closer()

	// Start exec with a bad starting message
	stream, err := client.EntrypointExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Open_{
			Open: &pb.EntrypointExecRequest_Open{
				InstanceId: exec.InstanceId,
				Index:      exec.Id + 4,
			},
		},
	}))

	// Wait for data
	resp, err := stream.Recv()
	require.Error(err)
	require.Equal(codes.NotFound, status.Code(err))
	require.Nil(resp)
}

func TestServiceEntrypointExecStream_closeSend(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	exec, closer := testRegisterExec(ctx, t, client, impl)
	defer closer()

	// Start exec with a bad starting message
	stream, err := client.EntrypointExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Open_{
			Open: &pb.EntrypointExecRequest_Open{
				InstanceId: exec.InstanceId,
				Index:      exec.Id,
			},
		},
	}))

	// Wait to hear we opened
	testEntrypointExecOpened(t, stream)

	// Close our sending side
	require.NoError(stream.CloseSend())

	// Wait for data
	resp, err := stream.Recv()
	require.Error(err)
	require.Equal(io.EOF, err)
	require.Nil(resp)
}

func TestServiceEntrypointExecStream_doubleStart(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(err)
	client := server.TestServer(t, impl)

	exec, closer := testRegisterExec(ctx, t, client, impl)
	defer closer()

	// Start exec
	stream, err := client.EntrypointExecStream(ctx)
	require.NoError(err)
	require.NoError(stream.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Open_{
			Open: &pb.EntrypointExecRequest_Open{
				InstanceId: exec.InstanceId,
				Index:      exec.Id,
			},
		},
	}))
	defer stream.CloseSend()

	// Wait to hear we opened
	testEntrypointExecOpened(t, stream)

	// Start a second exec
	stream2, err := client.EntrypointExecStream(ctx)
	require.NoError(err)
	require.NoError(stream2.Send(&pb.EntrypointExecRequest{
		Event: &pb.EntrypointExecRequest_Open_{
			Open: &pb.EntrypointExecRequest_Open{
				InstanceId: exec.InstanceId,
				Index:      exec.Id,
			},
		},
	}))
	defer stream2.CloseSend()

	// Wait for data
	resp, err := stream2.Recv()
	require.Error(err)
	require.Equal(codes.FailedPrecondition, status.Code(err))
	require.Nil(resp)
}

func testRegisterExec(ctx context.Context, t *testing.T, client pb.WaypointClient, impl pb.WaypointServer) (*serverstate.InstanceExec, func()) {
	// Create an instance
	instanceId, deploymentId, closer := TestEntrypoint(t, client)
	defer closer()

	// Start exec
	stream, err := client.StartExecStream(context.Background())
	require.NoError(t, err)
	require.NoError(t, stream.Send(&pb.ExecStreamRequest{
		Event: &pb.ExecStreamRequest_Start_{
			Start: &pb.ExecStreamRequest_Start{
				Target: &pb.ExecStreamRequest_Start_DeploymentId{
					DeploymentId: deploymentId,
				},
				Args: []string{"foo", "bar"},
			},
		},
	}))

	// Wait for the registered exec
	ws := memdb.NewWatchSet()
	list, err := testStateInmem(impl).InstanceExecListByInstanceId(ctx, instanceId, ws)
	require.NoError(t, err)
	if len(list) == 0 {
		ws.Watch(time.After(1 * time.Second))
		list, err = testStateInmem(impl).InstanceExecListByInstanceId(ctx, instanceId, ws)
		require.NoError(t, err)
	}
	require.Len(t, list, 1)

	return list[0], func() {
		stream.CloseSend()
	}
}

func testEntrypointExecOpened(t *testing.T, stream pb.Waypoint_EntrypointExecStreamClient) {
	resp, err := stream.Recv()
	require.NoError(t, err)
	require.IsType(t, resp.Event, (*pb.EntrypointExecResponse_Opened)(nil))
}
