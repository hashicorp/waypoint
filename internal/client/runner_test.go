package client

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"
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
	require.Nil(err)

	t.Run("Choose local if remote enabled is false for the project.", func(t *testing.T) {
		project = &pb.Project{
			Name:          "test",
			RemoteEnabled: false,
		}
		_, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{Project: project})
		require.Nil(err)

		remote, err := remoteOpPreferred(ctx, client, project, log)
		require.Nil(err)
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
		require.Nil(err)

		remote, err := remoteOpPreferred(ctx, client, project, log)
		require.Nil(err)
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
	_, remoteRunnerClose := singleprocess.TestRunner(t, client, &pb.Runner{Odr: false})
	defer remoteRunnerClose()

	// Register a non-default runner profile
	odrProfileName := "project-specific ODR profile"
	_, err = client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: &pb.OnDemandRunnerConfig{
			Name:       odrProfileName,
			PluginType: "docker",
			Default:    false,
		},
	})
	require.Nil(err)

	t.Run("Choose remote if the datasource is good, a remote runner exists, and a runner profile is set for the project", func(t *testing.T) {
		project = &pb.Project{
			Name:           "test",
			RemoteEnabled:  true,
			DataSource:     remoteCapableDataSource,
			OndemandRunner: &pb.Ref_OnDemandRunnerConfig{Name: odrProfileName},
		}
		_, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{Project: project})
		require.Nil(err)

		remote, err := remoteOpPreferred(ctx, client, project, log)
		require.Nil(err)
		require.True(remote)
	})

	t.Run("Choose local if no runner profile is set for the project, and there is no default", func(t *testing.T) {
		project = &pb.Project{
			Name:          "test",
			RemoteEnabled: true,
			DataSource:    remoteCapableDataSource,
		}
		_, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{Project: project})
		require.Nil(err)

		remote, err := remoteOpPreferred(ctx, client, project, log)
		require.Nil(err)
		require.False(remote)
	})

	// Register a default runner profile
	_, err = client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
		Config: &pb.OnDemandRunnerConfig{
			Name:       "the default",
			PluginType: "docker",
			Default:    true,
		},
	})
	require.Nil(err)

	t.Run("Choose remote if the project is good and the default runner is set", func(t *testing.T) {
		project = &pb.Project{
			Name:          "test",
			RemoteEnabled: true,
			DataSource:    remoteCapableDataSource,
		}
		_, err := client.UpsertProject(ctx, &pb.UpsertProjectRequest{Project: project})
		require.Nil(err)

		remote, err := remoteOpPreferred(ctx, client, project, log)
		require.Nil(err)
		require.True(remote)
	})
}
