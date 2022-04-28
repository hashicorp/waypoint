package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func init() {
	tests["waypoint_hcl"] = []testFunc{
		TestServiceWaypointHclFmt,
	}
}

func TestServiceWaypointHclFmt(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()

	// Create our server
	impl := factory(t)
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
