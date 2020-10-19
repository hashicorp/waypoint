package config

import (
	"testing"

	"github.com/mitchellh/pointerstructure"
	"github.com/stretchr/testify/require"
)

// This test isn't meant to be exhaustive but to test certain corner cases
// or bugs as they emerge. We aren't exhaustive since that ends up mostly just
// testing mergo.
func TestConfigDefault(t *testing.T) {
	cases := []struct {
		Name     string
		Config   *Config
		Path     string
		Expected interface{}
	}{
		{
			"app: no URL block set",
			&Config{
				Apps: []*App{
					{},
				},
			},
			"/Apps",
			[]*App{
				{
					URL: &AppURL{
						AutoHostname: nil,
					},
				},
			},
		},

		{
			"app: URL block with no hostname",
			&Config{
				Apps: []*App{
					{
						URL: &AppURL{},
					},
				},
			},
			"/Apps",
			[]*App{
				{
					URL: &AppURL{
						AutoHostname: nil,
					},
				},
			},
		},

		{
			"app: URL block with hostname set",
			&Config{
				Apps: []*App{
					{
						URL: &AppURL{
							AutoHostname: boolPtr(false),
						},
					},
				},
			},
			"/Apps",
			[]*App{
				{
					URL: &AppURL{
						AutoHostname: boolPtr(false),
					},
				},
			},
		},

		{
			"runner: disabled by default",
			&Config{},
			"/Runner/Enabled",
			false,
		},

		{
			"runner: disabled by default",
			&Config{},
			"/Runner/Enabled",
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)
			require.NoError(tt.Config.Default())

			var compare interface{} = tt.Config
			if tt.Path != "" {
				var err error
				compare, err = pointerstructure.Get(compare, tt.Path)
				require.NoError(err)
			}

			require.Equal(tt.Expected, compare)
		})
	}
}
