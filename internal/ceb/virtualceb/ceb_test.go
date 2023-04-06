// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package virtualceb

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/internal/server/execclient"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testVirtualHandler struct {
	info   *ExecInfo
	closed int

	stdin bytes.Buffer
}

func (t *testVirtualHandler) Run(ctx context.Context) error {
	io.Copy(&t.stdin, t.info.Input)
	fmt.Fprintf(t.info.Output, "test output\n")
	fmt.Fprintf(t.info.Error, "test error\n")
	time.Sleep(time.Second)

	return nil
}

func (t *testVirtualHandler) Close() error {
	t.closed++
	return nil
}

func (t *testVirtualHandler) PTYResize(_ *pb.ExecStreamRequest_WindowSize) error {
	return nil
}

func (t *testVirtualHandler) CreateSession(ctx context.Context, sess *ExecInfo) (ExecSession, error) {
	t.info = sess
	return t, nil
}

func TestVirtual_exec(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	L := hclog.New(&hclog.LoggerOptions{
		Level:           hclog.Trace,
		IncludeLocation: true,
	})

	// Start up the server
	restartCh := make(chan struct{})
	impl := singleprocess.TestImpl(t)
	client := server.TestServer(t, impl,
		server.TestWithContext(ctx),
		server.TestWithRestart(restartCh),
	)

	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, nil),
	})
	require.NoError(err)

	virt, err := New(L, Config{
		DeploymentId: resp.Deployment.Id,
		InstanceId:   "A",
		Client:       client,
	})

	require.NoError(err)

	var th testVirtualHandler

	go virt.RunExec(ctx, &th, 1)

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

	var stderr, stdout bytes.Buffer

	ec := &execclient.Client{
		Logger:  L,
		Context: ctx,
		Client:  client,
		Args:    []string{"date"},
		Stdin:   ioutil.NopCloser(strings.NewReader("input data\n")),
		Stdout:  &stdout,
		Stderr:  &stderr,

		InstanceId: "A",
	}

	code, err := ec.Run()
	require.NoError(err)

	assert.Equal(0, code)

	assert.Equal("input data\n", th.stdin.String())
	assert.Equal("test output\n", stdout.String())
	assert.Equal("test error\n", stderr.String())

	assert.Equal([]string{"date"}, th.info.Arguments)

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
