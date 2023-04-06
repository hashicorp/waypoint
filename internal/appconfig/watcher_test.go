// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appconfig

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	sdkpb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func init() {
	hclog.L().SetLevel(hclog.Debug)
}

func TestWatcher_initialBlock(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	w, err := NewWatcher(WithRefreshInterval(10 * time.Millisecond))
	require.NoError(err)
	defer w.Close()

	// We should timeout because we haven't configured any variables.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, _, err = w.Next(ctx, 0)
	require.Error(err)
	require.Equal(err, ctx.Err())

	// Update with empty vars
	w.UpdateVars(context.Background(), nil)

	// We should get an empty result
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	env, _, err := w.Next(ctx, 0)
	require.NoError(err)
	require.Empty(env.EnvVars)
}

func TestWatcher_static(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	w, err := NewWatcher(WithRefreshInterval(10 * time.Millisecond))
	require.NoError(err)
	defer w.Close()

	// Update with some static vars
	w.UpdateVars(ctx, []*pb.ConfigVar{
		{
			Name:  "TEST_VALUE",
			Value: &pb.ConfigVar_Static{Static: "hello"},
		},
	})

	// We should get the static vars back
	env, iter, err := w.Next(ctx, 0)
	require.NoError(err)
	require.Equal(env.EnvVars, []string{"TEST_VALUE=hello"})

	// Update with some other static vars
	w.UpdateVars(ctx, []*pb.ConfigVar{
		{
			Name:  "VALUE2",
			Value: &pb.ConfigVar_Static{Static: "goodbye"},
		},
		{
			Name:  "VALUE3",
			Value: &pb.ConfigVar_Static{Static: "hello"},
		},
	})

	// We should get the static vars back
	env, iter, err = w.Next(ctx, iter)
	require.NoError(err)
	require.Equal(env.EnvVars, []string{"VALUE2=goodbye", "VALUE3=hello"})
	require.Equal(env.DeletedEnvVars, []string{"TEST_VALUE"})
}

func TestWatcher_staticChange(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	w, err := NewWatcher(WithRefreshInterval(10 * time.Millisecond))
	require.NoError(err)
	defer w.Close()

	// Update with some static vars
	w.UpdateVars(ctx, []*pb.ConfigVar{
		{
			Name:  "TEST_VALUE",
			Value: &pb.ConfigVar_Static{Static: "hello"},
		},
	})

	// We should get the static vars back
	env, iter, err := w.Next(ctx, 0)
	require.NoError(err)
	require.Equal(env.EnvVars, []string{"TEST_VALUE=hello"})

	// Update with some other static vars
	w.UpdateVars(ctx, []*pb.ConfigVar{
		{
			Name:  "TEST_VALUE",
			Value: &pb.ConfigVar_Static{Static: "goodbye"},
		},
	})

	// We should get the static vars back
	env, iter, err = w.Next(ctx, iter)
	require.NoError(err)
	require.Equal(env.EnvVars, []string{"TEST_VALUE=goodbye"})
	require.Empty(env.DeletedEnvVars)
}

func TestWatcher_staticOriginalEnv(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	w, err := NewWatcher(
		WithRefreshInterval(10*time.Millisecond),
		WithOriginalEnv([]string{"TEST_VALUE=original"}),
	)
	require.NoError(err)
	defer w.Close()

	// Update with some static vars
	w.UpdateVars(ctx, []*pb.ConfigVar{
		{
			Name:  "TEST_VALUE",
			Value: &pb.ConfigVar_Static{Static: "hello"},
		},
	})

	// We should get the static vars back
	env, iter, err := w.Next(ctx, 0)
	require.NoError(err)
	require.Equal(env.EnvVars, []string{"TEST_VALUE=hello"})

	// Unset our vars
	w.UpdateVars(ctx, []*pb.ConfigVar{})

	// We should get the static vars back
	env, iter, err = w.Next(ctx, iter)
	require.NoError(err)
	require.Equal(env.EnvVars, []string{"TEST_VALUE=original"})
	require.Empty(env.DeletedEnvVars)
}

// Test that we read dynamic config variables.
func TestWatcher_dynamicSuccess(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create our test config source
	testSource := &testConfigSourcer{
		readValue: map[string]string{
			"envVarKey": "helloEnv",
			"fileKey":   "helloFile",
		},
	}

	w, err := NewWatcher(
		WithRefreshInterval(10*time.Millisecond),
		testWithConfigSourcer("cloud", testSource),
	)
	require.NoError(err)
	defer w.Close()

	// Change our config
	w.UpdateVars(ctx, []*pb.ConfigVar{
		{
			Name: "TEST_VALUE",
			Value: &pb.ConfigVar_Dynamic{
				Dynamic: &pb.ConfigVar_DynamicVal{
					From: "cloud",
					Config: map[string]string{
						"key": "envVarKey",
					},
				},
			},
		},
		{
			Name:       "/tmp/test_file.txt",
			NameIsPath: true,
			Value: &pb.ConfigVar_Dynamic{
				Dynamic: &pb.ConfigVar_DynamicVal{
					From: "cloud",
					Config: map[string]string{
						"key": "fileKey",
					},
				},
			},
		},
	})

	// We should get the static vars and files back
	cfg, iter, err := w.Next(ctx, 0)
	require.NoError(err)
	require.Equal(cfg.EnvVars, []string{"TEST_VALUE=helloEnv"})
	require.Equal(cfg.Files, []*FileContent{{
		Path: "/tmp/test_file.txt", Data: []byte("helloFile"),
	}})

	// Change the values and make sure we get them
	testSource.Lock()
	testSource.readValue["envVarKey"] = "goodbyeEnv"
	testSource.readValue["fileKey"] = "goodbyeFile"
	testSource.Unlock()

	// We should get the static vars and files back
	cfg, iter, err = w.Next(ctx, iter)
	require.NoError(err)
	require.Equal(cfg.EnvVars, []string{"TEST_VALUE=goodbyeEnv"})
	require.Equal(cfg.Files, []*FileContent{{
		Path: "/tmp/test_file.txt", Data: []byte("goodbyeFile"),
	}})

	// We should've called Stop once: exactly for the first read
	testSource.Lock()
	val := testSource.stopCount
	testSource.Unlock()
	require.Equal(1, val)

	// Unset our dynamic config
	w.UpdateVars(ctx, []*pb.ConfigVar{})

	// We should get the static vars and files back
	cfg, iter, err = w.Next(ctx, iter)
	require.NoError(err)
	require.Empty(cfg.EnvVars)
	require.Empty(cfg.Files)

	// We should call stop once more to end the previous run
	testSource.Lock()
	val = testSource.stopCount
	testSource.Unlock()
	require.Equal(2, val)
}

// Test that we read dynamic config variables where the source
// takes a configuration.
func TestWatcher_dynamicConfigurable(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create our test config source
	testSource := &testConfigSourcerWithConfig{}

	w, err := NewWatcher(
		WithRefreshInterval(10*time.Millisecond),
		testWithConfigSourcer("cloud", testSource),
	)
	require.NoError(err)
	defer w.Close()

	// Change our config source
	w.UpdateSources(ctx, []*pb.ConfigSource{
		{
			Type: "cloud",
			Config: map[string]string{
				"value": "flower",
			},
		},
	})

	// Change our config
	w.UpdateVars(ctx, []*pb.ConfigVar{
		{
			Name: "TEST_VALUE",
			Value: &pb.ConfigVar_Dynamic{
				Dynamic: &pb.ConfigVar_DynamicVal{
					From: "cloud",
				},
			},
		},
	})

	// We should get the static vars back
	env, iter, err := w.Next(ctx, 0)
	require.NoError(err)
	require.Equal(env.EnvVars, []string{"TEST_VALUE=flower"})

	// Change our config source
	w.UpdateSources(ctx, []*pb.ConfigSource{
		{
			Type: "cloud",
			Config: map[string]string{
				"value": "leaf",
			},

			// Gotta update the hash manually. Usually the Waypoint
			// server handles this but we're not using the server for
			// these tests.
			Hash: 1,
		},
	})

	// We should get the static vars back
	env, iter, err = w.Next(ctx, iter)
	require.NoError(err)
	require.Equal(env.EnvVars, []string{"TEST_VALUE=leaf"})

	// We should've called Stop twice: once for the first read and
	// then again when we changed the configuration.
	require.Equal(uint32(2), testSource.StopCount())

	// Unset our dynamic config
	w.UpdateVars(ctx, []*pb.ConfigVar{})

	// We should get the static vars back
	env, iter, err = w.Next(ctx, iter)
	require.NoError(err)
	require.Empty(env.EnvVars)

	// We should call stop once more to end the previous run
	require.Equal(uint32(3), testSource.StopCount())
}

// When a dynamic source is unused and we set a configuration,
// it should not impact our process.
func TestConfig_dynamicConfigurableUnused(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create our test config source
	testSource := &testConfigSourcerWithConfig{}

	w, err := NewWatcher(
		WithRefreshInterval(15*time.Millisecond),
		testWithConfigSourcer("cloud", testSource),
	)
	require.NoError(err)
	defer w.Close()

	// Change our config source
	w.UpdateSources(ctx, []*pb.ConfigSource{
		{
			Type: "cloud",
			Config: map[string]string{
				"value": "flower",
			},
		},
	})

	// Update with empty vars
	w.UpdateVars(ctx, nil)

	// We should get an empty result
	env, iter, err := w.Next(ctx, 0)
	require.NoError(err)
	require.Empty(env.EnvVars)

	// Change our config source
	w.UpdateSources(ctx, []*pb.ConfigSource{
		{
			Type: "cloud",
			Config: map[string]string{
				"value": "leaf",
			},

			// Gotta update the hash manually. Usually the Waypoint
			// server handles this but we're not using the server for
			// these tests.
			Hash: 1,
		},
	})

	{
		// We should timeout because nothing changed.
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		_, _, err = w.Next(ctx, iter)
		require.Error(err)
		require.Equal(err, ctx.Err())
	}
}

// Test that we read dynamic config variables.
func TestWatcher_variableReferences(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create our test config source
	testSource := &testConfigSourcer{
		readValue: map[string]string{"key": "hello"},
	}

	w, err := NewWatcher(
		WithRefreshInterval(10*time.Millisecond),
		testWithConfigSourcer("cloud", testSource),
	)
	require.NoError(err)
	defer w.Close()

	// Change our config
	w.UpdateVars(ctx, []*pb.ConfigVar{
		{
			Name: "V1",
			Value: &pb.ConfigVar_Static{
				Static: "connect://${config.env.TEST_VALUE}",
			},
		},
		{
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
	})

	// We should get the static vars back
	env, _, err := w.Next(ctx, 0)
	require.NoError(err)
	require.Equal([]string{"TEST_VALUE=hello", "V1=connect://hello"}, env.EnvVars)
}

// Test that we process escaped variables.
func TestWatcher_variableEscape(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create our test config source
	testSource := &testConfigSourcer{
		readValue: map[string]string{"key": "hello"},
	}

	w, err := NewWatcher(
		WithRefreshInterval(10*time.Millisecond),
		testWithConfigSourcer("cloud", testSource),
	)
	require.NoError(err)
	defer w.Close()

	// Change our config
	w.UpdateVars(ctx, []*pb.ConfigVar{
		{
			Name: "V1",
			Value: &pb.ConfigVar_Static{
				Static: "connect://$${get_hostname()}",
			},
		},
	})

	// We should get the static vars back
	env, _, err := w.Next(ctx, 0)
	require.NoError(err)
	require.Equal([]string{"V1=connect://${get_hostname()}"}, env.EnvVars)
}

// TestWatcher_sourcerJson tests that the configsourcer can handle
// plugins returning json values, not just string values.
// We expect configsourcers to return json typed data to feed
// dynamic default variables, and not for app or runner config.
func TestWatcher_sourcerJson(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create our test config source
	testSource := &testConfigSourcer{
		readJson: map[string][]byte{
			"structuredValKey": []byte("{'foo':'bar'}"),
		},
	}

	w, err := NewWatcher(
		WithRefreshInterval(10*time.Millisecond),
		testWithConfigSourcer("cloud", testSource),
	)
	require.NoError(err)
	defer w.Close()

	// Change our config
	w.UpdateVars(ctx, []*pb.ConfigVar{
		{
			Name:       "structuredVal",
			NameIsPath: true,
			Value: &pb.ConfigVar_Dynamic{
				Dynamic: &pb.ConfigVar_DynamicVal{
					From: "cloud",
					Config: map[string]string{
						"key": "structuredValKey",
					},
				},
			},
		},
	})

	cfg, _, err := w.Next(ctx, 0)
	require.NoError(err)

	// It is undeniably weird that we're returning json data as raw file contents.
	// We've very much outgrown the current app config system that returns either
	// files or env vars. IMO, we should refactor this system to return structured data,
	// and the caller can decide to use it for files, env vars, or some other purpose.
	//
	// For now, we're doubling down on the existing hack wherein the variable system
	// requests its dynamic config vars to be treated as "files", and it
	// can figure out on its own whether to treat it as json or a string based
	// on the variable type the user specified in the HCL.
	require.Equal(cfg.Files, []*FileContent{{
		Path: "structuredVal", Data: []byte("{'foo':'bar'}"),
	}})
}

type testConfigSourcer struct {
	sync.Mutex

	stopCount int
	readValue map[string]string

	// json values
	readJson map[string][]byte
}

func (s *testConfigSourcer) ReadFunc() interface{} {
	return func(reqs []*component.ConfigRequest) ([]*sdkpb.ConfigSource_Value, error) {
		s.Lock()
		defer s.Unlock()

		var result []*sdkpb.ConfigSource_Value
		for _, req := range reqs {

			if stringVal, ok := s.readValue[req.Config["key"]]; ok {
				result = append(result, &sdkpb.ConfigSource_Value{
					Name: req.Name,
					Result: &sdkpb.ConfigSource_Value_Value{
						Value: stringVal,
					},
				})
			} else if jsonVal, ok := s.readJson[req.Config["key"]]; ok {
				result = append(result, &sdkpb.ConfigSource_Value{
					Name: req.Name,
					Result: &sdkpb.ConfigSource_Value_Json{
						Json: jsonVal,
					},
				})
			} else {
				panic("invalid test - request for unset key " + req.Config["key"])
			}
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
