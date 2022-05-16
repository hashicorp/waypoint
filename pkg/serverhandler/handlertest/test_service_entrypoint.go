package handlertest

//
//func initz() {
//	tests["entrypoint"] = []testFunc{
//		TestServiceEntrypointLogStream,
//	}
//}
//
//func TestServiceEntrypointLogStream(t *testing.T, factory Factory) {
//	ctx := context.Background()
//	require := require.New(t)
//
//	_, client := factory(t)
//
//	// Create a deployment
//	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
//		Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
//			Component: &pb.Component{
//				Name: "testapp",
//			},
//		}),
//	})
//	require.NoError(err)
//	dep := resp.Deployment
//
//	// Connect a CEB and create an instance
//
//	// Create the config
//	instanceId, err := server.Id()
//	require.NoError(err)
//	_, err = client.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
//		InstanceId:   instanceId,
//		DeploymentId: dep.Id,
//	})
//	require.NoError(err)
//
//	// Send a log
//	stream, err := client.EntrypointLogStream(ctx)
//	require.NoError(err)
//
//	err = stream.Send(&pb.EntrypointLogBatch{
//		InstanceId: instanceId,
//		Lines: []*pb.LogBatch_Entry{{
//			Source:    pb.LogBatch_Entry_APP,
//			Timestamp: timestamppb.Now(),
//			Line:      "i'm a log message!",
//		}},
//	})
//
//	require.NoError(err)
//
//	// Receive the log
//
//}
