package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/waypoint/internal/server/execclient"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	componentmocks "github.com/hashicorp/waypoint-plugin-sdk/component/mocks"
	"github.com/hashicorp/waypoint/internal/config"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestAppExec_happy(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	// Our mock platform, which must also implement Destroyer
	mock := struct {
		*componentmocks.Platform
		*componentmocks.Execer
	}{
		&componentmocks.Platform{},
		&componentmocks.Execer{},
	}

	// Make our factory for platforms
	factory := TestFactory(t, component.PlatformType)
	TestFactoryRegister(t, factory, "test", mock)

	// Make our app
	app := TestApp(t, TestProject(t,
		WithConfig(config.TestConfig(t, testPlatformConfig)),
		WithFactory(component.PlatformType, factory),
	), "test")

	client := app.client

	ctx := context.Background()

	// We're using GetVersionInfoResponse here just because it is a proto message
	// that can be converted to an any.Any easily. We never use it, it's just to keep
	// the tests from blowing up with a nil reference.
	mockPluginArtifact := &pb.GetVersionInfoResponse{}

	anyval, err := ptypes.MarshalAny(mockPluginArtifact)

	aresp, err := client.UpsertPushedArtifact(ctx, &pb.UpsertPushedArtifactRequest{
		Artifact: serverptypes.TestValidArtifact(t, &pb.PushedArtifact{
			Artifact: &pb.Artifact{
				Artifact: anyval,
			},
		}),
	})
	require.NoError(err)

	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, &pb.Deployment{
			ArtifactId: aresp.Artifact.Id,
		}),
	})
	require.NoError(err)

	var stdin bytes.Buffer

	anyd, err := ptypes.MarshalAny(&empty.Empty{})
	require.NoError(err)

	// Expect to have the destroy function called
	require.NoError(err)
	mock.Execer.On("ExecFunc").Return(func(d *any.Any, esi *component.ExecSessionInfo) error {
		app.logger.Info("called mock ExecFunc")

		if d == nil || d != anyd {
			return fmt.Errorf("value didn't match")
		}

		io.Copy(&stdin, esi.Input)

		fmt.Fprintf(esi.Output, "from the mock\n")

		return nil
	})

	instanceId := "A"

	resp.Deployment.Deployment = anyd

	// Exec
	go func() {
		app.Exec(context.Background(), instanceId, resp.Deployment, true)
	}()

	// We should get registered
	require.Eventually(func() bool {
		resp, err := client.ListInstances(ctx, &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: resp.Deployment.Id,
			},
		})
		require.NoError(err)
		return len(resp.Instances) == 1
	}, 2*time.Second, 25*time.Millisecond)

	// Make sure that with all the exec stream tracking we don't leak
	// goroutines
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	var stderr, stdout bytes.Buffer

	ec := &execclient.Client{
		Logger:  app.logger,
		Context: ctx,
		Client:  client,
		Args:    []string{"date"},
		Stdin:   ioutil.NopCloser(strings.NewReader("input data\n")),
		Stdout:  &stdout,
		Stderr:  &stderr,

		InstanceId: "A",
	}

	app.logger.Info("connecting execclient")

	code, err := ec.Run()
	require.NoError(err)

	assert.Equal(0, code)

	assert.Equal("from the mock\n", stdout.String())
	assert.Equal("input data\n", stdin.String())

	// We should get deregistered
	require.Eventually(func() bool {
		resp, err := client.ListInstances(ctx, &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: resp.Deployment.Id,
			},
		})
		require.NoError(err)
		return len(resp.Instances) == 0
	}, time.Second, 100*time.Millisecond)

}
