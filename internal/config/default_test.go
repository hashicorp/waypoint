package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// This test isn't meant to be exhaustive but to test certain corner cases
// or bugs as they emerge. We aren't exhaustive since that ends up mostly just
// testing mergo.
func TestConfigDefault(t *testing.T) {
	cases := []struct {
		Name     string
		Config   *Config
		Expected *Config
	}{
		{
			"app: no URL block set",
			&Config{
				Apps: []*App{
					&App{},
				},
			},
			&Config{
				Apps: []*App{
					&App{
						URL: &AppURL{
							AutoHostname: nil,
						},
					},
				},
			},
		},

		{
			"app: URL block with no hostname",
			&Config{
				Apps: []*App{
					&App{
						URL: &AppURL{},
					},
				},
			},
			&Config{
				Apps: []*App{
					&App{
						URL: &AppURL{
							AutoHostname: nil,
						},
					},
				},
			},
		},

		{
			"app: URL block with hostname set",
			&Config{
				Apps: []*App{
					&App{
						URL: &AppURL{
							AutoHostname: boolPtr(false),
						},
					},
				},
			},
			&Config{
				Apps: []*App{
					&App{
						URL: &AppURL{
							AutoHostname: boolPtr(false),
						},
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)
			require.NoError(tt.Config.Default())
			require.Equal(tt.Expected, tt.Config)
		})
	}
}
