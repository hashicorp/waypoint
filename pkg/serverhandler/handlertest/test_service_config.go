// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func init() {
	tests["workspace"] = []testFunc{
		TestServiceConfig,
		TestServiceConfigSource,
	}
}

func TestServiceConfig(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Simplify writing tests
	type (
		SReq = pb.ConfigSetRequest
		GReq = pb.ConfigGetRequest
	)

	Var := &pb.ConfigVar{
		Target: &pb.ConfigVar_Target{
			AppScope: &pb.ConfigVar_Target_Application{
				Application: &pb.Ref_Application{
					Application: "foo",
					Project:     "bar",
				},
			},
		},

		Name:  "DATABASE_URL",
		Value: &pb.ConfigVar_Static{Static: "postgresql:///"},
	}

	Var2 := &pb.ConfigVar{
		Target: &pb.ConfigVar_Target{
			Workspace: &pb.Ref_Workspace{
				Workspace: "prod",
			},
		},

		Name:  "VAULT_URL",
		Value: &pb.ConfigVar_Static{Static: "example.com"},
	}

	t.Run("set and get", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.SetConfig(ctx, &SReq{Variables: []*pb.ConfigVar{Var}})
		require.NoError(err)
		require.NotNil(resp)

		// Let's write some data

		grep, err := client.GetConfig(ctx, &GReq{
			Scope: &pb.ConfigGetRequest_Application{
				Application: &pb.Ref_Application{
					Application: "foo",
					Project:     "bar",
				},
			},
		})
		require.NoError(err)
		require.NotNil(grep)

		require.Equal(1, len(grep.Variables))

		require.Equal(Var.Name, grep.Variables[0].Name)
		require.Equal(Var.Value, grep.Variables[0].Value)
	})

	t.Run("delete", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.SetConfig(ctx, &SReq{Variables: []*pb.ConfigVar{Var, Var2}})
		require.NoError(err)
		require.NotNil(resp)

		// Unset via static type with empty string
		_, err = client.DeleteConfig(ctx, &pb.ConfigDeleteRequest{
			Variables: []*pb.ConfigVar{{
				Name:  "DATABASE_URL",
				Value: &pb.ConfigVar_Static{Static: ""},
			}},
		})

		// It's gone
		grep, err := client.GetConfig(ctx, &GReq{
			Scope: &pb.ConfigGetRequest_Application{
				Application: &pb.Ref_Application{
					Application: "foo",
					Project:     "bar",
				},
			},
		})
		require.Error(err)
		require.Nil(grep)

		// Unset via Unset protobuf type
		_, err = client.DeleteConfig(ctx, &pb.ConfigDeleteRequest{
			Variables: []*pb.ConfigVar{{
				Name:  "VAULT_URL",
				Value: &pb.ConfigVar_Unset{},
			}},
		})

		// It's gone
		grep, err = client.GetConfig(ctx, &GReq{
			Workspace: &pb.Ref_Workspace{Workspace: "prod"},
		})
		require.Error(err)
		require.Nil(grep)
	})
}

func TestServiceConfigSource(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

	// Simplify writing tests
	type (
		SReq = pb.SetConfigSourceRequest
		GReq = pb.GetConfigSourceRequest
	)

	v := &pb.ConfigSource{
		Scope: &pb.ConfigSource_Global{
			Global: &pb.Ref_Global{},
		},

		Type: "foo",

		Config: map[string]string{
			"value": "42",
		},
	}

	t.Run("set and get", func(t *testing.T) {
		require := require.New(t)

		// Create
		resp, err := client.SetConfigSource(ctx, &SReq{ConfigSource: v})
		require.NoError(err)
		require.NotNil(resp)

		// Read
		{
			resp, err := client.GetConfigSource(ctx, &GReq{
				Scope: &pb.GetConfigSourceRequest_Global{
					Global: &pb.Ref_Global{},
				},

				Type: "foo",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(1, len(resp.ConfigSources))
		}
	})
}
