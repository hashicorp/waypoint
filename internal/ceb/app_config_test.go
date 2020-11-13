package ceb

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	sdkpb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"
)

// Test that our child process is restarted with an env var change.
func TestConfig_envVarChange(t *testing.T) {
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

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(td)
	path := filepath.Join(td, "hello")

	// Start the CEB
	ceb := testRun(t, context.Background(), &testRunOpts{
		Client: client,
		Helper: "write-env",
		HelperEnv: map[string]string{
			"HELPER_PATH": path,
			"TEST_VALUE":  "",
		},
	})

	// The child should still start up
	require.Eventually(func() bool {
		_, err := ioutil.ReadFile(path)
		return err == nil
	}, 5*time.Second, 10*time.Millisecond)

	// Get our deployment
	deployment, err := client.GetDeployment(ctx, &pb.GetDeploymentRequest{
		Ref: &pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{
				Id: ceb.DeploymentId(),
			},
		},
	})
	require.NoError(err)

	// Change our config
	_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{
		Variables: []*pb.ConfigVar{
			{
				Scope: &pb.ConfigVar_Application{
					Application: deployment.Application,
				},
				Name:  "TEST_VALUE",
				Value: &pb.ConfigVar_Static{Static: "hello"},
			},
		},
	})
	require.NoError(err)

	// The child should still start up
	var data []byte
	require.Eventually(func() bool {
		var err error
		data, err = ioutil.ReadFile(path)
		return err == nil && strings.Contains(string(data), "hello")
	}, 5*time.Second, 10*time.Millisecond)

	// Set our config again but to the same value
	_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{
		Variables: []*pb.ConfigVar{
			{
				Scope: &pb.ConfigVar_Application{
					Application: deployment.Application,
				},
				Name:  "TEST_VALUE",
				Value: &pb.ConfigVar_Static{Static: "hello"},
			},
		},
	})
	require.NoError(err)

	// The child should still start up
	time.Sleep(1 * time.Second)
	data2, err := ioutil.ReadFile(path)
	require.NoError(err)
	require.Equal(data, data2)
}

// Test that we read dynamic config variables.
func TestConfig_dynamicSuccess(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create our test config source
	testSource := &testConfigSourcer{
		readValue: map[string]string{"key": "hello"},
	}

	// Start up the server
	impl := singleprocess.TestImpl(t)
	client := server.TestServer(t, impl,
		server.TestWithContext(ctx),
	)

	// Create a temporary directory for our test
	td, err := ioutil.TempDir("", "test")
	require.NoError(err)
	defer os.RemoveAll(td)
	path := filepath.Join(td, "hello")

	// Start the CEB
	ceb := testRun(t, context.Background(), &testRunOpts{
		Client: client,
		Helper: "write-env",
		HelperEnv: map[string]string{
			"HELPER_PATH": path,
			"TEST_VALUE":  "",
		},
		ConfigPlugins: map[string]component.ConfigSourcer{
			"cloud": testSource,
		},
	})

	// The child should still start up
	require.Eventually(func() bool {
		_, err := ioutil.ReadFile(path)
		return err == nil
	}, 5*time.Second, 10*time.Millisecond)

	// Get our deployment
	deployment, err := client.GetDeployment(ctx, &pb.GetDeploymentRequest{
		Ref: &pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{
				Id: ceb.DeploymentId(),
			},
		},
	})
	require.NoError(err)

	// Change our config
	_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{
		Variables: []*pb.ConfigVar{
			{
				Scope: &pb.ConfigVar_Application{
					Application: deployment.Application,
				},
				Name: "TEST_VALUE",
				Value: &pb.ConfigVar_Dynamic{
					Dynamic: &pb.ConfigVar_DynamicVal{
						From: "cloud",
						Config: map[string]string{
							"key": "key",
						},
					},
				},
			},
		},
	})
	require.NoError(err)

	// The child should still start up
	var data []byte
	require.Eventually(func() bool {
		var err error
		data, err = ioutil.ReadFile(path)
		return err == nil && strings.Contains(string(data), "hello")
	}, 5*time.Second, 10*time.Millisecond)

	// We should've called Stop once: exactly for the first read
	testSource.Lock()
	val := testSource.stopCount
	testSource.Unlock()
	require.Equal(1, val)
}

type testConfigSourcer struct {
	sync.Mutex

	stopCount int
	readValue map[string]string
}

func (s *testConfigSourcer) ReadFunc() interface{} {
	return func(reqs []*component.ConfigRequest) ([]*sdkpb.ConfigSource_Value, error) {
		s.Lock()
		defer s.Unlock()

		var result []*sdkpb.ConfigSource_Value
		for _, req := range reqs {
			result = append(result, &sdkpb.ConfigSource_Value{
				Name: req.Name,
				Result: &sdkpb.ConfigSource_Value_Value{
					Value: s.readValue[req.Config["key"]],
				},
			})
		}

		return result, nil
	}
}

func (s *testConfigSourcer) StopFunc() interface{} {
	return func() error {
		s.Lock()
		defer s.Unlock()
		s.stopCount++
		return nil
	}
}
