// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ceb

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint/internal/plugin"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"

	"github.com/stretchr/testify/require"
)

func TestRun_reconnect(t *testing.T) {
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
	ceb := testRun(t, context.Background(), &testRunOpts{Client: client})

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

	// Shut down the server
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

	// We should get re-registered
	require.Eventually(func() bool {
		resp, err := impl.ListInstances(ctx, &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: ceb.DeploymentId(),
			},
		})
		require.NoError(err)
		return len(resp.Instances) == 1
	}, 5*time.Second, 10*time.Millisecond)
}

// Test how the CEB behaves when the server is down on startup.
func TestRun_serverDownBasic(t *testing.T) {
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

	// Create a deployment
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, nil),
	})
	require.NoError(err)
	dep := resp.Deployment

	// But actually shut it down
	cancel()

	// Wait to get an unavailable error so we know the server is down
	require.Eventually(func() bool {
		_, err := client.ListInstances(context.Background(), &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: "doesn't matter",
			},
		})
		return status.Code(err) == codes.Unavailable
	}, 5*time.Second, 10*time.Millisecond)

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(td)
	path := filepath.Join(td, "hello")

	// Start the CEB
	ceb := testRun(t, context.Background(), &testRunOpts{
		Client:       client,
		DeploymentId: dep.Id,
		Helper:       "write-file",
		HelperEnv: map[string]string{
			"HELPER_PATH": path,
		},
	})

	// The child should still start up
	require.Eventually(func() bool {
		_, err := ioutil.ReadFile(path)
		return err == nil
	}, 5*time.Second, 10*time.Millisecond)

	// Restart
	restartCh <- struct{}{}

	// We should get re-registered eventually
	require.Eventually(func() bool {
		resp, err := impl.ListInstances(context.Background(), &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: ceb.DeploymentId(),
			},
		})
		require.NoError(err)
		return len(resp.Instances) == 1
	}, 5*time.Second, 10*time.Millisecond)
}

func TestRun_serverDownRequired(t *testing.T) {
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

	// Create a deployment
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, nil),
	})
	require.NoError(err)
	dep := resp.Deployment

	// But actually shut it down
	cancel()

	// Wait to get an unavailable error so we know the server is down
	require.Eventually(func() bool {
		_, err := client.ListInstances(context.Background(), &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: "doesn't matter",
			},
		})
		return status.Code(err) == codes.Unavailable
	}, 5*time.Second, 10*time.Millisecond)

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(td)
	path := filepath.Join(td, "hello")

	// Start the CEB
	testRun(t, context.Background(), &testRunOpts{
		Client:       client,
		DeploymentId: dep.Id,
		Helper:       "write-file",
		HelperEnv: map[string]string{
			envCEBServerRequired: "1",
			"HELPER_PATH":        path,
		},
	})

	// The child should NOT start up
	time.Sleep(1 * time.Second)
	_, err = ioutil.ReadFile(path)
	require.Error(err)

	// Restart
	restartCh <- struct{}{}

	// The child should start up
	require.Eventually(func() bool {
		_, err := ioutil.ReadFile(path)
		return err == nil
	}, 5*time.Second, 10*time.Millisecond)
}

// Test how the CEB behaves when the server is down on startup.
func TestRun_serverDownNoConnect(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start a listener that will refuse connections
	ln, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(err)
	ln.Close()

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(td)
	path := filepath.Join(td, "hello")

	// Start the CEB
	testRun(t, ctx, &testRunOpts{
		ClientDisable: true,
		DeploymentId:  "ABCD1234",
		Helper:        "write-file",
		HelperEnv: map[string]string{
			envServerAddr: ln.Addr().String(),
			"HELPER_PATH": path,
		},
	})

	// The child should still start up
	require.Eventually(func() bool {
		_, err := ioutil.ReadFile(path)
		return err == nil
	}, 10*time.Second, 10*time.Millisecond)
}

// Test CEB disabled with server up. Shouldn't connect at all.
func TestRun_disabledUp(t *testing.T) {
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

	// Create a deployment
	resp, err := client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
		Deployment: serverptypes.TestValidDeployment(t, nil),
	})
	require.NoError(err)
	dep := resp.Deployment

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(td)
	path := filepath.Join(td, "hello")

	// Start the CEB
	ceb := testRun(t, ctx, &testRunOpts{
		Client:       client,
		DeploymentId: dep.Id,
		Helper:       "write-file",
		HelperEnv: map[string]string{
			envCEBDisable: "1",
			"HELPER_PATH": path,
		},
	})

	// The child should start up
	require.Eventually(func() bool {
		_, err := ioutil.ReadFile(path)
		return err == nil
	}, 5*time.Second, 10*time.Millisecond)

	// We should NOT get registered
	{
		time.Sleep(500 * time.Millisecond)
		resp, err := client.ListInstances(ctx, &pb.ListInstancesRequest{
			Scope: &pb.ListInstancesRequest_DeploymentId{
				DeploymentId: ceb.DeploymentId(),
			},
		})
		require.NoError(err)
		require.Empty(resp.Instances)
	}
}

// Test CEB disabled via no server addr being set.
func TestRun_disabledNoServerAddr(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(td)
	path := filepath.Join(td, "hello")

	// Start the CEB
	testRun(t, ctx, &testRunOpts{
		ClientDisable: true,
		DeploymentId:  "who cares",
		Helper:        "write-file",
		HelperEnv: map[string]string{
			envServerAddr: "",
			"HELPER_PATH": path,
		},
	})

	// The child should start up
	require.Eventually(func() bool {
		_, err := ioutil.ReadFile(path)
		return err == nil
	}, 5*time.Second, 10*time.Millisecond)
}

var (
	testExec      = os.Args[0]
	envHelperMode = "TEST_HELPER_MODE"
)

func TestMain(m *testing.M) {
	switch os.Getenv(envHelperMode) {
	case "":
		// Log
		hclog.L().SetLevel(hclog.Debug)

		// Normal test mode
		os.Exit(m.Run())

	case "sleep":
		time.Sleep(10 * time.Minute)

	case "logs-stdout":
		for {
			time.Sleep(250 * time.Millisecond)
			fmt.Println(time.Now().String())
		}

	case "write-file":
		path := os.Getenv("HELPER_PATH")
		if path == "" {
			panic("bad")
		}

		ioutil.WriteFile(path, []byte("hello"), 0600)
		time.Sleep(10 * time.Minute)

	case "write-env":
		path := os.Getenv("HELPER_PATH")
		if path == "" {
			panic("bad")
		}

		ioutil.WriteFile(path, []byte(fmt.Sprintf("%d,%s", os.Getpid(), os.Getenv("TEST_VALUE"))), 0600)
		time.Sleep(10 * time.Minute)

	case "read-file":
		path := os.Getenv("HELPER_PATH")
		if path == "" {
			panic("bad")
		}

		rp := os.Getenv("READ_PATH")
		if rp == "" {
			panic("bad")
		}

		sig := make(chan os.Signal, 1)

		signal.Notify(sig, syscall.SIGUSR2)
		go func() {
			<-sig
			data, _ := ioutil.ReadFile(rp)

			ioutil.WriteFile(path, []byte(fmt.Sprintf("%d,%s", os.Getpid(), string(data))), 0600)
		}()

		data, _ := ioutil.ReadFile(rp)

		ioutil.WriteFile(path, []byte(fmt.Sprintf("%d,%s", os.Getpid(), string(data))), 0600)
		time.Sleep(10 * time.Minute)

	default:
		panic("invalid helperfunc")
	}
}

func testRun(t *testing.T, ctx context.Context, opts *testRunOpts) *CEB {
	if opts == nil {
		opts = &testRunOpts{}
	}

	if opts.Client == nil && !opts.ClientDisable {
		opts.Client = singleprocess.TestServer(t)
	}

	if opts.Helper == "" {
		opts.Helper = "sleep"
	}

	if opts.DeploymentId == "" {
		resp, err := opts.Client.UpsertDeployment(ctx, &pb.UpsertDeploymentRequest{
			Deployment: serverptypes.TestValidDeployment(t, nil),
		})
		require.NoError(t, err)
		opts.DeploymentId = resp.Deployment.Id
	}

	// Setup our deployment
	require.NoError(t, os.Setenv(envDeploymentId, opts.DeploymentId))
	t.Cleanup(func() { os.Setenv(envDeploymentId, "") })

	// Setup our helper
	require.NoError(t, os.Setenv(envHelperMode, opts.Helper))
	t.Cleanup(func() { os.Setenv(envHelperMode, "") })

	// Setup helper env vars
	for k, v := range opts.HelperEnv {
		kcopy := k
		require.NoError(t, os.Setenv(k, v))
		t.Cleanup(func() { os.Setenv(kcopy, "") })
	}

	// This is so we can wait to get it set
	cebCh := make(chan *CEB, 1)

	// Run it
	go Run(ctx,
		WithExec([]string{testExec}),
		WithClient(opts.Client),
		WithEnvDefaults(),
		withCEBValue(cebCh),
	)
	ceb := <-cebCh

	// Register our config plugins. NOTE(mitchellh): This is nasty cause we're
	// just poking at internal state, so we should clean this up one day.
	for k, v := range opts.ConfigPlugins {
		ceb.configPlugins[k] = &plugin.Instance{
			Component: v,
		}
	}

	return ceb
}

type testRunOpts struct {
	Client        pb.WaypointClient
	ClientDisable bool
	Helper        string
	HelperEnv     map[string]string
	DeploymentId  string
	ConfigPlugins map[string]component.ConfigSourcer
}
