package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestServiceConfig(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(testDB(t))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	// Simplify writing tests
	type (
		SReq = pb.ConfigSetRequest
		GReq = pb.ConfigGetRequest
	)

	Var := &pb.ConfigVar{Name: "DATABASE_URL", Value: "postgresql:///"}

	t.Run("set and get", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.SetConfig(ctx, &SReq{Var: Var})
		require.NoError(err)
		require.NotNil(resp)

		// Let's write some data

		grep, err := client.GetConfig(ctx, &GReq{})
		require.NoError(err)
		require.NotNil(grep)

		require.Equal(1, len(grep.Variables))

		require.Equal(Var.Name, grep.Variables[0].Name)
		require.Equal(Var.Value, grep.Variables[0].Value)
	})
}
