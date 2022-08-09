package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/clicontext"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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

func TestInitConfigLoad(t *testing.T) {
	cases := []struct {
		Name                    string
		FileName                string
		ExpectedValidateResults bool
		ExpectedErr             bool
	}{
		{
			"valid file",
			"valid.hcl",
			false,
			false,
		},
		{
			"invalid file",
			"invalid.hcl",
			true,
			true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)
			c := baseCommand{}

			// Add a context, set the given context workspace value if any
			cfg := &clicontext.Config{
				Workspace: defaultWorkspace,
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
			}

			err := baseCfg.Flags.Parse(baseCfg.Args)
			require.NoError(err)
			c.refWorkspace = &pb.Ref_Workspace{Workspace: defaultWorkspace}

			c.args = baseCfg.Flags.Args()
			_, vr, err := c.initConfig(filepath.Join("testdata", tt.FileName))
			require.Equal((err != nil), tt.ExpectedErr, "error")
			require.Equal(len(vr) > 0, tt.ExpectedValidateResults, "validation results")
		})
	}
}
