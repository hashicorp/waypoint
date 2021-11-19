package cli

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess"
)

func TestCheckFlagsAfterArgs(t *testing.T) {
	var boolVal bool

	cases := []struct {
		Name string
		Flag func(*flag.Sets)
		Args []string
		Err  bool
	}{
		{
			"empty args",
			func(*flag.Sets) {},
			[]string{},
			false,
		},

		{
			"flag with space",
			func(sets *flag.Sets) {
				s := sets.NewSet("test")
				s.BoolVar(&flag.BoolVar{Name: "foo", Target: &boolVal})
			},
			[]string{"-foo", "bar"},
			true,
		},

		{
			"double hyphen",
			func(sets *flag.Sets) {
				s := sets.NewSet("test")
				s.BoolVar(&flag.BoolVar{Name: "foo", Target: &boolVal})
			},
			[]string{"--foo", "bar"},
			true,
		},

		{
			"equals",
			func(sets *flag.Sets) {
				s := sets.NewSet("test")
				s.BoolVar(&flag.BoolVar{Name: "foo", Target: &boolVal})
			},
			[]string{"--foo=bar"},
			true,
		},

		{
			"ignores after double hyphen",
			func(sets *flag.Sets) {
				s := sets.NewSet("test")
				s.BoolVar(&flag.BoolVar{Name: "foo", Target: &boolVal})
			},
			[]string{"hello", "--", "--foo=bar"},
			false,
		},

		{
			"other flag",
			func(sets *flag.Sets) {
				s := sets.NewSet("test")
				s.BoolVar(&flag.BoolVar{Name: "foo", Target: &boolVal})
			},
			[]string{"--bar=bar"},
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			s := flag.NewSets()
			tt.Flag(s)

			err := checkFlagsAfterArgs(tt.Args, s)
			if !tt.Err {
				require.NoError(err)
				return
			}
			require.Error(err)
		})
	}
}

func TestWorkspacePrecedence(t *testing.T) {
	cases := []struct {
		Name             string
		Args             []string
		Env              string
		ContextWorkspace string
		Expected         string
	}{
		{
			"default no inputs",
			[]string{},
			"",
			"",
			defaultWorkspace,
		},
		{
			"workspace flag",
			[]string{"-workspace", "dev"},
			"",
			"",
			"dev",
		},
		{
			"workspace flag and context",
			[]string{"-workspace", "dev"},
			"",
			"lab",
			"dev",
		},
		{
			"context only",
			[]string{},
			"",
			"lab",
			"lab",
		},
		{
			"env only",
			[]string{},
			"test",
			"",
			"test",
		},
		{
			"env and storage",
			[]string{},
			"test",
			"dev",
			"test",
		},
		{
			"everything",
			[]string{"-workspace", "other"},
			"test",
			"dev",
			"other",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)
			c := baseCommand{}

			// Add a context, set the given context workspace value if any
			cfg := &clicontext.Config{}
			if tt.ContextWorkspace != "" {
				cfg.Workspace = tt.ContextWorkspace
			}

			st := clicontext.TestStorage(t)
			require.NoError(st.Set("default", cfg))

			c.contextStorage = st

			// setup flags and arguments of the base command, and set in the
			// base config. This is work typically done in base.Init()
			sets := flag.NewSets()
			set := sets.NewSet("test")
			set.StringVar(&flag.StringVar{Name: "workspace", Target: &c.flagWorkspace})

			baseCfg := baseConfig{
				Flags: sets,
				Args:  tt.Args,
			}

			err := baseCfg.Flags.Parse(baseCfg.Args)
			require.NoError(err)

			c.args = baseCfg.Flags.Args()

			if tt.Env != "" {
				// setup env with the test value. This is unset at the end of
				// each test regardless
				os.Setenv(defaultWorkspaceEnvName, tt.Env)
			}

			// execute the method and test the value
			workspace, err := c.workspace()
			require.NoError(err)
			require.Equal(tt.Expected, workspace)

			// reset the env after every test
			os.Unsetenv(defaultWorkspaceEnvName)
		})
	}
}

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
