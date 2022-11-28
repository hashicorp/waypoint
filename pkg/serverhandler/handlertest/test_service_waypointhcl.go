package handlertest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func init() {
	tests["waypoint_hcl"] = []testFunc{
		TestServiceWaypointHclFmt,
	}
}

func TestServiceWaypointHclFmt(t *testing.T, factory Factory) {
	ctx := context.Background()

	// Create our server
	client, _ := factory(t)

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
