// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package statetest

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func init() {
	tests["runner_ondemand"] = []testFunc{
		TestOnDemandRunnerConfig,
		TestOnDemandRunnerConfig_LabelTargeting,
	}
}

func TestOnDemandRunnerConfig(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("Get returns not found error if not exist", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		_, err := s.OnDemandRunnerConfigGet(ctx, &pb.Ref_OnDemandRunnerConfig{
			Id: "foo",
		})
		require.Error(err)
		require.Equal(codes.NotFound, status.Code(err))
	})

	t.Run("Set with no name", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		rec := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
			OciUrl:               "h/w:s",
			EnvironmentVariables: map[string]string{"CONTAINER": "DOCKER", "FOO": "BAR"},
			TargetRunner:         &pb.Ref_Runner{Target: &pb.Ref_Runner_Any{Any: &pb.Ref_RunnerAny{}}},
			PluginConfig:         []byte(`{"foo":"bar"}`),
			ConfigFormat:         pb.Hcl_JSON,
			Default:              true,
		})

		result, err := s.OnDemandRunnerConfigPut(ctx, rec)
		require.NoError(err)
		require.NotEmpty(result.Id)
		require.NotEmpty(result.Name)
	})

	t.Run("Client cannot control the ID", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		{
			// Setting the ID for a new profile
			rec := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
				Id: "client-chooses",
			})

			_, err := s.OnDemandRunnerConfigPut(ctx, rec)
			require.Error(err)
		}

		{
			// Setting the ID for an existing profile
			rec := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
				Name: "client-chooses-name",
			})

			resp, err := s.OnDemandRunnerConfigPut(ctx, rec)
			require.NoError(err)
			require.NotEmpty(resp.Id) // server chose the ID

			updateRec := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
				Name: rec.Name,
				Id:   "trying to overwrite the ID",
			})

			_, err = s.OnDemandRunnerConfigPut(ctx, updateRec)
			require.Error(err)
		}
	})

	t.Run("Put and Get", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		rec := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
			Name:                 "test",
			OciUrl:               "h/w:s",
			EnvironmentVariables: map[string]string{"CONTAINER": "DOCKER", "FOO": "BAR"},
			TargetRunner:         &pb.Ref_Runner{Target: &pb.Ref_Runner_Any{Any: &pb.Ref_RunnerAny{}}},
			PluginConfig:         []byte(`{"foo":"bar"}`),
			ConfigFormat:         pb.Hcl_JSON,
			Default:              true,
		})

		putResp, err := s.OnDemandRunnerConfigPut(ctx, rec)
		require.NoError(err)
		require.NotEmpty(putResp.Id) // must have chosen an id
		require.Equal(rec.Name, putResp.Name)
		require.Equal(rec.TargetRunner.Target, putResp.TargetRunner.Target)
		require.Equal(rec.OciUrl, putResp.OciUrl)
		require.Equal(rec.EnvironmentVariables, putResp.EnvironmentVariables)
		require.Equal(rec.PluginConfig, putResp.PluginConfig)
		require.Equal(rec.ConfigFormat, putResp.ConfigFormat)
		require.Equal(rec.Default, putResp.Default)

		// Get exact
		{
			resp, err := s.OnDemandRunnerConfigGet(ctx, &pb.Ref_OnDemandRunnerConfig{
				Id: rec.Id,
			})
			require.NoError(err)
			require.NotNil(resp)

			// Ensure fields were saved correctly
			require.Equal(putResp.Id, resp.Id)
			require.Equal(rec.Name, resp.Name)
			require.Equal(rec.TargetRunner.Target, resp.TargetRunner.Target)
			require.Equal(rec.OciUrl, resp.OciUrl)
			require.Equal(rec.EnvironmentVariables, resp.EnvironmentVariables)
			require.Equal(rec.PluginConfig, resp.PluginConfig)
			require.Equal(rec.ConfigFormat, resp.ConfigFormat)
			require.Equal(rec.Default, resp.Default)
		}

		// Get case insensitive
		{
			resp, err := s.OnDemandRunnerConfigGet(ctx, &pb.Ref_OnDemandRunnerConfig{
				Id: strings.ToUpper(rec.Id),
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Get by name
		{
			resp, err := s.OnDemandRunnerConfigGet(ctx, &pb.Ref_OnDemandRunnerConfig{
				Name: rec.Name,
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(rec.Name, resp.Name)
			require.Equal(putResp.Id, resp.Id)
		}

		// Get missing (returns not found error)
		{
			_, err := s.OnDemandRunnerConfigGet(ctx, &pb.Ref_OnDemandRunnerConfig{
				Id: strings.ToUpper("unknown"),
			})
			require.Error(err)
			require.Equal(status.Code(err), codes.NotFound)
		}

		// List
		{
			resp, err := s.OnDemandRunnerConfigList(ctx)
			require.NoError(err)
			require.Len(resp, 1)
		}
	})

	t.Run("cannot create a new profile with existing name", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		rec := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
			Name:                 "test",
			OciUrl:               "h/w:s",
			EnvironmentVariables: map[string]string{"CONTAINER": "DOCKER", "FOO": "BAR"},
			TargetRunner:         &pb.Ref_Runner{Target: &pb.Ref_Runner_Any{Any: &pb.Ref_RunnerAny{}}},
			PluginConfig:         []byte(`{"foo":"bar"}`),
			ConfigFormat:         pb.Hcl_JSON,
			Default:              true,
		})

		putResp, err := s.OnDemandRunnerConfigPut(ctx, rec)
		require.NoError(err)
		require.NotEmpty(putResp.Id) // must have chosen an id
		require.Equal(rec.Name, putResp.Name)
		require.Equal(rec.TargetRunner.Target, putResp.TargetRunner.Target)
		require.Equal(rec.OciUrl, putResp.OciUrl)
		require.Equal(rec.EnvironmentVariables, putResp.EnvironmentVariables)
		require.Equal(rec.PluginConfig, putResp.PluginConfig)
		require.Equal(rec.ConfigFormat, putResp.ConfigFormat)
		require.Equal(rec.Default, putResp.Default)

		// Get exact
		{
			resp, err := s.OnDemandRunnerConfigGet(ctx, &pb.Ref_OnDemandRunnerConfig{
				Id: rec.Id,
			})
			require.NoError(err)
			require.NotNil(resp)

			// Ensure fields were saved correctly
			require.Equal(putResp.Id, resp.Id)
			require.Equal(rec.Name, resp.Name)
			require.Equal(rec.TargetRunner.Target, resp.TargetRunner.Target)
			require.Equal(rec.OciUrl, resp.OciUrl)
			require.Equal(rec.EnvironmentVariables, resp.EnvironmentVariables)
			require.Equal(rec.PluginConfig, resp.PluginConfig)
			require.Equal(rec.ConfigFormat, resp.ConfigFormat)
			require.Equal(rec.Default, resp.Default)
		}

		// Try to create a new profile with existing name!
		rec2 := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
			Name:                 "test",
			OciUrl:               "hc/wp:s",
			EnvironmentVariables: map[string]string{"CONTAINER": "DOCKER", "BAZ": "BAR"},
			TargetRunner:         &pb.Ref_Runner{Target: &pb.Ref_Runner_Any{Any: &pb.Ref_RunnerAny{}}},
			PluginConfig:         []byte(`{"baz":"bar"}`),
			ConfigFormat:         pb.Hcl_JSON,
			Default:              false,
		})

		// No error, it just updates the existing one
		resp, err := s.OnDemandRunnerConfigPut(ctx, rec2)
		require.NoError(err)
		// We updated a field, see if it stuck
		require.Equal(rec2.OciUrl, resp.OciUrl)

		// List should only return 1
		{
			resp, err := s.OnDemandRunnerConfigList(ctx)
			require.NoError(err)
			require.Len(resp, 1)
		}

	})

	t.Run("Delete", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		rec := serverptypes.TestOnDemandRunnerConfig(t, serverptypes.TestOnDemandRunnerConfig(t, nil))

		_, err := s.OnDemandRunnerConfigPut(ctx, rec)
		require.NoError(err)

		// Read
		resp, err := s.OnDemandRunnerConfigGet(ctx, &pb.Ref_OnDemandRunnerConfig{
			Id: rec.Id,
		})
		require.NoError(err)
		require.NotNil(resp)

		// Delete
		{
			err := s.OnDemandRunnerConfigDelete(ctx, &pb.Ref_OnDemandRunnerConfig{
				Id: rec.Id,
			})
			require.NoError(err)
		}

		// Read
		{
			_, err := s.OnDemandRunnerConfigGet(ctx, &pb.Ref_OnDemandRunnerConfig{
				Id: rec.Id,
			})
			require.Error(err)
			require.Equal(codes.NotFound, status.Code(err))
		}

		// List
		{
			resp, err := s.OnDemandRunnerConfigList(ctx)
			require.NoError(err)
			require.Len(resp, 0)
		}
	})
}

func TestOnDemandRunnerConfig_LabelTargeting(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	require := require.New(t)

	s := factory(t)
	defer s.Close()

	labels := map[string]string{"name": "testrunner", "env": "test"}

	// Set
	rec := serverptypes.TestOnDemandRunnerConfig(t, &pb.OnDemandRunnerConfig{
		Name:   "test",
		OciUrl: "h/w:s",
		TargetRunner: &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Labels{
				Labels: &pb.Ref_RunnerLabels{
					Labels: labels,
				},
			},
		},
	})

	_, err := s.OnDemandRunnerConfigPut(ctx, rec)
	require.NoError(err)

	resp, err := s.OnDemandRunnerConfigGet(ctx, &pb.Ref_OnDemandRunnerConfig{
		Id: rec.Id,
	})
	require.NoError(err)
	require.NotNil(resp)

	// Ensure the target saved properly
	switch target := rec.TargetRunner.Target.(type) {
	case *pb.Ref_Runner_Labels:
		require.Equal(target.Labels.Labels, labels)
	default:
		t.Fatalf("runner target type %t is not by label", target)
	}
}
