package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestServiceWaypointHclFmt(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	t.Run("basic formatting", func(t *testing.T) {
		require := require.New(t)

		const input = `
project="foo"
`

		const output = `
project = "foo"
`

		// Create, should get an ID back
		resp, err := client.WaypointHclFmt(ctx, &pb.WaypointHclFmtRequest{
			WaypointHcl: []byte(input),
		})
		require.NoError(err)
		require.NotNil(resp)

		// Let's write some data
		require.Equal(string(resp.WaypointHcl), output)
	})
}
