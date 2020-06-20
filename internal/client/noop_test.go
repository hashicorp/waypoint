package client

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server/singleprocess"
)

func init() {
	hclog.L().SetLevel(hclog.Trace)
}

func TestClientNoop(t *testing.T) {
	ctx := context.Background()
	require := require.New(t)
	client := singleprocess.TestServer(t)

	// Build our client
	c := TestClient(t, WithClient(client), WithLocal())

	// TODO(mitchellh): once we have an API to list jobs, verify we have
	// no jobs, and then verify we execute a job after.

	// Noop
	require.NoError(c.Noop(ctx))
}
