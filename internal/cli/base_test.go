package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
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
