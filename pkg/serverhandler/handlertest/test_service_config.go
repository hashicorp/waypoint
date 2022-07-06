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
