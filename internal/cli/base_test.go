package cli

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	pbmocks "github.com/hashicorp/waypoint/internal/server/gen/mocks"
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

func Test_remoteIsPossible(t *testing.T) {
	log := hclog.Default()

	type args struct {
		project           *pb.Project
		runnersResp       *pb.ListRunnersResponse
		runnerConfigsResp *pb.ListOnDemandRunnerConfigsResponse
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Choose local if remote enabled is false for the project.",
			args: args{
				project: &pb.Project{
					RemoteEnabled: false,
				},
			},
			want:    false,
			wantErr: false,
		},

		{
			name: "Choose local if the datasource is not remote-capable",
			args: args{
				project: &pb.Project{
					RemoteEnabled: true,
					DataSource: &pb.Job_DataSource{
						Source: &pb.Job_DataSource_Local{},
					},
				},
			},
			want:    false,
			wantErr: false,
		},

		{
			name: "Choose local if there are no remote runners",
			args: args{
				project: &pb.Project{
					RemoteEnabled: true,
					DataSource: &pb.Job_DataSource{
						Source: &pb.Job_DataSource_Git{},
					},
				},
				runnersResp: &pb.ListRunnersResponse{Runners: []*pb.Runner{{Odr: true}}},
			},
			want:    false,
			wantErr: false,
		},

		{
			name: "Choose remote if the datasource is good, a remote runner exists, and a runner profile is set",
			args: args{
				project: &pb.Project{
					RemoteEnabled: true,
					DataSource: &pb.Job_DataSource{
						Source: &pb.Job_DataSource_Git{},
					},
					OndemandRunner: &pb.Ref_OnDemandRunnerConfig{},
				},
				runnersResp: &pb.ListRunnersResponse{Runners: []*pb.Runner{{Odr: false}}},
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "Choose local if no runner profile is set for the project, and there is no default",
			args: args{
				project: &pb.Project{
					RemoteEnabled: true,
					DataSource: &pb.Job_DataSource{
						Source: &pb.Job_DataSource_Git{},
					},
					OndemandRunner: nil,
				},
				runnerConfigsResp: &pb.ListOnDemandRunnerConfigsResponse{
					Configs: []*pb.OnDemandRunnerConfig{{
						Default: false,
					}},
				},
				runnersResp: &pb.ListRunnersResponse{Runners: []*pb.Runner{{Odr: false}}},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Choose remote if the project is good and the default runner is set",
			args: args{
				project: &pb.Project{
					RemoteEnabled: true,
					DataSource: &pb.Job_DataSource{
						Source: &pb.Job_DataSource_Git{},
					},
					OndemandRunner: nil,
				},
				runnerConfigsResp: &pb.ListOnDemandRunnerConfigsResponse{
					Configs: []*pb.OnDemandRunnerConfig{{
						Default: true,
					}},
				},
				runnersResp: &pb.ListRunnersResponse{Runners: []*pb.Runner{{Odr: false}}},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			var client pb.WaypointClient
			if tt.args.runnerConfigsResp != nil || tt.args.runnersResp != nil {
				m := &pbmocks.WaypointServer{}
				// Called when initializing the client
				m.On("BootstrapToken", mock.Anything, mock.Anything).Return(&pb.NewTokenResponse{Token: "hello"}, nil)
				m.On("GetVersionInfo", mock.Anything, mock.Anything).Return(server.TestVersionInfoResponse(), nil)

				m.On("ListOnDemandRunnerConfigs", mock.Anything, mock.Anything).Return(tt.args.runnerConfigsResp, nil)
				m.On("ListRunners", mock.Anything, mock.Anything).Return(tt.args.runnersResp, nil)
				client = server.TestServer(t, m, server.TestWithContext(ctx))
			}

			got, err := remoteIsPossible(ctx, client, tt.args.project, log)
			if (err != nil) != tt.wantErr {
				t.Errorf("remoteIsPossible() error = %v, wantErr %v", err, tt.wantErr)
				cancel()
				return
			}
			if got != tt.want {
				t.Errorf("remoteIsPossible() got = %v, want %v", got, tt.want)
			}
			cancel()
		})
	}
}
