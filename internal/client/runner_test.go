// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"
	"testing"

	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/pkg/server"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func Test_remoteOpPreferred(t *testing.T) {
	log := hclog.Default()
	require := require.New(t)

	ctx := context.Background()

	client := singleprocess.TestServer(t)

	project := &pb.Project{
		Name: "test",
	}

	_, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{Project: project})
	require.NoError(err)

	t.Run("Choose local if remote enabled is false for the project.", func(t *testing.T) {
		project = &pb.Project{
			Name:          "test",
			RemoteEnabled: false,
		}
		_, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{Project: project})
		require.NoError(err)

		remote, err := remoteOpPreferred(ctx, client, project, nil, log)
		require.NoError(err)
		require.False(remote)
	})

	t.Run("Choose local if the datasource is not remote-capable.", func(t *testing.T) {
		project = &pb.Project{
			Name:          "test",
			RemoteEnabled: true,
			DataSource: &pb.Job_DataSource{
				Source: &pb.Job_DataSource_Local{},
			},
		}
		_, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{Project: project})
		require.NoError(err)

		remote, err := remoteOpPreferred(ctx, client, project, nil, log)
		require.NoError(err)
		require.False(remote)
	})

	remoteCapableDataSource := &pb.Job_DataSource{
		Source: &pb.Job_DataSource_Git{
			Git: &pb.Job_Git{
				Ref: "main",
				Url: "git.test",
			},
		},
	}

	// Register a remote runner
	_, remoteRunnerClose := server.TestRunner(t, client, &pb.Runner{
		Kind: &pb.Runner_Remote_{Remote: &pb.Runner_Remote{}},
	})
	defer remoteRunnerClose()

	// Register a non-default runner profile
	odrProfileName := "project-specific-ODR-profile"
	_, err = client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: &pb.OnDemandRunnerConfig{
			Name:       odrProfileName,
			PluginType: "docker",
			Default:    false,
		},
	})
	require.NoError(err)

	t.Run("Choose remote if the datasource is good, a remote runner exists, and a runner profile is set for the project", func(t *testing.T) {
		project = &pb.Project{
			Name:          "test",
			RemoteEnabled: true,
			DataSource:    remoteCapableDataSource,
		}
		_, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{Project: project})
		require.NoError(err)

		runnerCfgs := []*configpkg.Runner{{Profile: "test"}}

		remote, err := remoteOpPreferred(ctx, client, project, runnerCfgs, log)
		require.NoError(err)
		require.True(remote)
	})

	t.Run("Choose remote if the app on the project has a runner profile set", func(t *testing.T) {
		project = &pb.Project{
			Name:          "test",
			RemoteEnabled: true,
			DataSource:    remoteCapableDataSource,
			Applications: []*pb.Application{{
				Name: "test-app",
			}},
		}
		_, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{Project: project})
		require.NoError(err)

		runnerCfgs := []*configpkg.Runner{{Profile: "test"}}

		remote, err := remoteOpPreferred(ctx, client, project, runnerCfgs, log)
		require.NoError(err)
		require.True(remote)
	})

	t.Run("Choose local if no runner profile is set for the project, and there is no default", func(t *testing.T) {
		project = &pb.Project{
			Name:          "test",
			RemoteEnabled: true,
			DataSource:    remoteCapableDataSource,
		}
		_, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{Project: project})
		require.NoError(err)

		remote, err := remoteOpPreferred(ctx, client, project, nil, log)
		require.NoError(err)
		require.False(remote)
	})

	t.Run("Choose remote if the project is good and the default runner is set", func(t *testing.T) {
		// Register a default runner profile
		_, err = client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
			Config: &pb.OnDemandRunnerConfig{
				Name:       "the-default",
				PluginType: "docker",
				Default:    true,
			},
		})
		require.NoError(err)

		project = &pb.Project{
			Name:          "test",
			RemoteEnabled: true,
			DataSource:    remoteCapableDataSource,
		}
		_, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{Project: project})
		require.NoError(err)

		remote, err := remoteOpPreferred(ctx, client, project, nil, log)
		require.NoError(err)
		require.True(remote)
	})
}
