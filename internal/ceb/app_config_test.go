// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ceb

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	sdkpb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"
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
				Target: &pb.ConfigVar_Target{
					AppScope: &pb.ConfigVar_Target_Application{
						Application: deployment.Application,
					},
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
				Target: &pb.ConfigVar_Target{
					AppScope: &pb.ConfigVar_Target_Application{
						Application: deployment.Application,
					},
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

// Test that our child process is signaled with a file change
func TestConfig_fileChange(t *testing.T) {
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
	fooPath := filepath.Join(td, "foo.txt")

	// Start the CEB
	ceb := testRun(t, context.Background(), &testRunOpts{
		Client: client,
		Helper: "read-file",
		HelperEnv: map[string]string{
			"HELPER_PATH": path,
			"READ_PATH":   fooPath,
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

	client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
		Project:          &pb.Ref_Project{Project: deployment.Application.Project},
		Name:             deployment.Application.Application,
		FileChangeSignal: "USR2",
	})

	// Change our config
	_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{
		Variables: []*pb.ConfigVar{
			{
				Target: &pb.ConfigVar_Target{
					AppScope: &pb.ConfigVar_Target_Application{
						Application: deployment.Application,
					},
				},
				Name:       fooPath,
				NameIsPath: true,
				Value:      &pb.ConfigVar_Static{Static: "via config"},
			},
		},
	})
	require.NoError(err)

	// The child should still start up
	var data []byte
	require.Eventually(func() bool {
		var err error
		data, err = ioutil.ReadFile(path)
		return err == nil && strings.Contains(string(data), "via config")
	}, 5*time.Second, 10*time.Millisecond)

	// Set our config again but to the same value
	_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{
		Variables: []*pb.ConfigVar{
			{
				Target: &pb.ConfigVar_Target{
					AppScope: &pb.ConfigVar_Target_Application{
						Application: deployment.Application,
					},
				},
				Name:       fooPath,
				NameIsPath: true,
				Value:      &pb.ConfigVar_Static{Static: "via config"},
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

	// Short refresh so we can test changing values
	testChConfigRefresh(t, 100*time.Millisecond)

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
			"TEST_VALUE":  "original",
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
				Target: &pb.ConfigVar_Target{
					AppScope: &pb.ConfigVar_Target_Application{
						Application: deployment.Application,
					},
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

	// Change the value and make sure we get it
	testSource.Lock()
	testSource.readValue["key"] = "goodbye"
	testSource.Unlock()

	// The child should change
	require.Eventually(func() bool {
		var err error
		data, err = ioutil.ReadFile(path)
		return err == nil && strings.Contains(string(data), "goodbye")
	}, 5*time.Second, 10*time.Millisecond)

	// We should've called Stop once: exactly for the first read
	testSource.Lock()
	val := testSource.stopCount
	testSource.Unlock()
	require.Equal(1, val)

	// Unset our dynamic config
	_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{
		Variables: []*pb.ConfigVar{
			{
				Target: &pb.ConfigVar_Target{
					AppScope: &pb.ConfigVar_Target_Application{
						Application: deployment.Application,
					},
				},
				Name: "TEST_VALUE",
				Value: &pb.ConfigVar_Unset{
					Unset: &empty.Empty{},
				},
			},
		},
	})
	require.NoError(err)

	// The child should change
	require.Eventually(func() bool {
		var err error
		data, err = ioutil.ReadFile(path)
		return err == nil && strings.Contains(string(data), "original")
	}, 5*time.Second, 10*time.Millisecond)

	// We should call stop once more to end the previous run
	testSource.Lock()
	val = testSource.stopCount
	testSource.Unlock()
	require.Equal(2, val)
}

// Test that we read dynamic config variables where the source
// takes a configuration.
func TestConfig_dynamicConfigurable(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Short refresh so we can test changing values
	testChConfigRefresh(t, 100*time.Millisecond)

	// Create our test config source
	testSource := &testConfigSourcerWithConfig{}

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
			"TEST_VALUE":  "original",
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

	// Set the config for our source
	_, err = client.SetConfigSource(ctx, &pb.SetConfigSourceRequest{
		ConfigSource: &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type: "cloud",

			Config: map[string]string{
				"value": "flower",
			},
		},
	})
	require.NoError(err)

	// Change our config
	_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{
		Variables: []*pb.ConfigVar{
			{
				Target: &pb.ConfigVar_Target{
					AppScope: &pb.ConfigVar_Target_Application{
						Application: deployment.Application,
					},
				},
				Name: "TEST_VALUE",
				Value: &pb.ConfigVar_Dynamic{
					Dynamic: &pb.ConfigVar_DynamicVal{
						From: "cloud",
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
		return err == nil && strings.Contains(string(data), "flower")
	}, 5*time.Second, 10*time.Millisecond)

	// Change the value and make sure we get it
	_, err = client.SetConfigSource(ctx, &pb.SetConfigSourceRequest{
		ConfigSource: &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type: "cloud",

			Config: map[string]string{
				"value": "leaf",
			},
		},
	})
	require.NoError(err)

	// The child should change
	require.Eventually(func() bool {
		var err error
		data, err = ioutil.ReadFile(path)
		return err == nil && strings.Contains(string(data), "leaf")
	}, 5*time.Second, 10*time.Millisecond)

	// We should've called Stop twice: once for the first read and
	// then again when we changed the configuration.
	require.Equal(uint32(2), testSource.StopCount())

	// Unset our dynamic config
	_, err = client.SetConfig(ctx, &pb.ConfigSetRequest{
		Variables: []*pb.ConfigVar{
			{
				Target: &pb.ConfigVar_Target{
					AppScope: &pb.ConfigVar_Target_Application{
						Application: deployment.Application,
					},
				},
				Name: "TEST_VALUE",
				Value: &pb.ConfigVar_Unset{
					Unset: &empty.Empty{},
				},
			},
		},
	})
	require.NoError(err)

	// The child should change
	require.Eventually(func() bool {
		var err error
		data, err = ioutil.ReadFile(path)
		return err == nil && strings.Contains(string(data), "original")
	}, 5*time.Second, 10*time.Millisecond)

	// We should call stop once more to end the previous run
	require.Equal(uint32(3), testSource.StopCount())
}

// When a dynamic source is unused and we set a configuration,
// it should not impact our process.
func TestConfig_dynamicConfigurableUnused(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Short refresh so we can test changing values
	testChConfigRefresh(t, 100*time.Millisecond)

	// Create our test config source
	testSource := &testConfigSourcerWithConfig{}

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
			"TEST_VALUE":  "original",
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
	_, err = client.GetDeployment(ctx, &pb.GetDeploymentRequest{
		Ref: &pb.Ref_Operation{
			Target: &pb.Ref_Operation_Id{
				Id: ceb.DeploymentId(),
			},
		},
	})
	require.NoError(err)

	// Set the config for our source
	_, err = client.SetConfigSource(ctx, &pb.SetConfigSourceRequest{
		ConfigSource: &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type: "cloud",

			Config: map[string]string{
				"value": "flower",
			},
		},
	})
	require.NoError(err)

	// The child should start up
	var pid string
	require.Eventually(func() bool {
		var err error
		data, err := ioutil.ReadFile(path)
		if err == nil && strings.Contains(string(data), "original") {
			parts := strings.Split(string(data), ",")
			pid = parts[0]
			return true
		}

		return false
	}, 5*time.Second, 10*time.Millisecond)

	// Change the value
	_, err = client.SetConfigSource(ctx, &pb.SetConfigSourceRequest{
		ConfigSource: &pb.ConfigSource{
			Scope: &pb.ConfigSource_Global{
				Global: &pb.Ref_Global{},
			},

			Type: "cloud",

			Config: map[string]string{
				"value": "leaf",
			},
		},
	})
	require.NoError(err)

	// Sleep for a bit and ensure we have the same value
	{
		time.Sleep(1 * time.Second)
		data, err := ioutil.ReadFile(path)
		require.NoError(err)
		require.True(strings.HasPrefix(string(data), pid))
	}
}

// testChConfigRefresh changes the amount of time between config refreshes.
// A test cleanup function is automatically registered to revert.
func testChConfigRefresh(t *testing.T, d time.Duration) {
	old := appConfigRefreshPeriod
	appConfigRefreshPeriod = d
	t.Cleanup(func() { appConfigRefreshPeriod = old })
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

// ConfigSourcer that implements Configurable
type testConfigSourcerWithConfig struct {
	config struct {
		Value string `hcl:"value,attr"`
	}

	stopCount uint32
}

func (s *testConfigSourcerWithConfig) Config() (interface{}, error) {
	return &s.config, nil
}

func (s *testConfigSourcerWithConfig) ReadFunc() interface{} {
	return func(reqs []*component.ConfigRequest) ([]*sdkpb.ConfigSource_Value, error) {
		var result []*sdkpb.ConfigSource_Value
		for _, req := range reqs {
			result = append(result, &sdkpb.ConfigSource_Value{
				Name: req.Name,
				Result: &sdkpb.ConfigSource_Value_Value{
					Value: s.config.Value,
				},
			})
		}

		return result, nil
	}
}

func (s *testConfigSourcerWithConfig) StopFunc() interface{} {
	return func() error {
		atomic.AddUint32(&s.stopCount, 1)
		return nil
	}
}

func (s *testConfigSourcerWithConfig) StopCount() uint32 {
	return atomic.LoadUint32(&s.stopCount)
}
