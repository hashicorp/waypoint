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
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func init() {
	hclog.L().SetLevel(hclog.Trace)
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
	require.Empty(env)
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
	require.Equal(env, []string{"TEST_VALUE=hello"})

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
	require.Equal(env, []string{"VALUE2=goodbye", "VALUE3=hello"})
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
	require.Equal(env, []string{"TEST_VALUE=hello"})

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
	require.Equal(env, []string{"TEST_VALUE=goodbye"})
}

// Test that we read dynamic config variables.
func TestWatcher_dynamicSuccess(t *testing.T) {
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
	env, iter, err := w.Next(ctx, 0)
	require.NoError(err)
	require.Equal(env, []string{"TEST_VALUE=hello"})

	// Change the value and make sure we get it
	testSource.Lock()
	testSource.readValue["key"] = "goodbye"
	testSource.Unlock()

	// We should get the static vars back
	env, iter, err = w.Next(ctx, iter)
	require.NoError(err)
	require.Equal(env, []string{"TEST_VALUE=goodbye"})

	// We should've called Stop once: exactly for the first read
	testSource.Lock()
	val := testSource.stopCount
	testSource.Unlock()
	require.Equal(1, val)

	// Unset our dynamic config
	w.UpdateVars(ctx, []*pb.ConfigVar{})

	// We should get the static vars back
	env, iter, err = w.Next(ctx, iter)
	require.NoError(err)
	require.Empty(env)

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
	require.Equal(env, []string{"TEST_VALUE=flower"})

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
	require.Equal(env, []string{"TEST_VALUE=leaf"})

	// We should've called Stop twice: once for the first read and
	// then again when we changed the configuration.
	require.Equal(uint32(2), testSource.StopCount())

	// Unset our dynamic config
	w.UpdateVars(ctx, []*pb.ConfigVar{})

	// We should get the static vars back
	env, iter, err = w.Next(ctx, iter)
	require.NoError(err)
	require.Empty(env)

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
	require.Empty(env)

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
