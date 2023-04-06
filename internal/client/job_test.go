// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/grpcmetadata"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"

	"github.com/stretchr/testify/require"
)

func Test_setupLocalJobSystem(t *testing.T) {
	hclog.Default().SetLevel(hclog.Debug)
	require := require.New(t)

	ctx := context.Background()
	var err error

	// Validates the side effects of running setupLocalJobSystem
	validateLocalSetupSideEffects := func(ctx context.Context, c *Project, expectLocal bool) {
		// Validate saved locality setting for future operations
		require.NotNil(c.useLocalRunner)
		require.Equal(*c.useLocalRunner, expectLocal)

		// Validate started local runner (or not)
		hasActiveLocalRunner := c.activeRunner != nil
		require.Equal(hasActiveLocalRunner, expectLocal)

		// Validate set GRPC metadata
		_, hasRunnerId := grpcmetadata.OutgoingRunnerId(ctx)
		require.Equal(hasRunnerId, expectLocal)

		// TODO(izaak): check the context
	}

	t.Run("Respects project settings", func(t *testing.T) {
		for _, requestLocal := range []bool{true, false} {
			c := TestProject(t,
				WithClient(singleprocess.TestServer(t)),
				WithUseLocalRunner(requestLocal),
			)
			defer c.Close()

			// We don't need to upsert the project - it shouldn't need to make any API calls to
			// choose a locality if we set it explicitly on the client.

			isLocal, newCtx, err := c.setupLocalJobSystem(ctx)
			require.Nil(err)
			require.Equal(requestLocal, isLocal)

			// Check that running setupLocalJobSystem had the right side effects
			validateLocalSetupSideEffects(newCtx, c, requestLocal)
		}
	})

	t.Run("Automatically determines locality if unset", func(t *testing.T) {

		// Simple setup, uses a
		c := TestProject(t,
			WithClient(singleprocess.TestServer(t)),
		)
		defer c.Close()

		project := &pb.Project{
			Name:          c.project.Project,
			RemoteEnabled: false,
		}

		_, err = c.Client().UpsertProject(ctx, &pb.UpsertProjectRequest{Project: project})
		require.Nil(err)

		isLocal, newCtx, err := c.setupLocalJobSystem(ctx)
		require.Nil(err)

		// we don't care what the value of isLocal is for this test - just that it picked _something_,
		// saved it, and performed its side-effects.
		validateLocalSetupSideEffects(newCtx, c, isLocal)
	})
}
