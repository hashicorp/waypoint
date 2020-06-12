package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidation(t *testing.T) {
	cases := []struct {
		Name   string
		Modify func(*Config)
		Err    string
	}{
		{
			"label map with a waypoint/ key",
			func(c *Config) {
				c.Labels = map[string]string{}
				c.Labels["waypoint/foo"] = "bar"
			},
			"reserved for system",
		},

		{
			"label map with hostname key",
			func(c *Config) {
				c.Labels = map[string]string{}
				c.Labels["foo"] = "bar"
			},
			"",
		},

		{
			"label map with hostname key before '/'",
			func(c *Config) {
				c.Labels = map[string]string{
					"foo/baz":   "bar",
					"foo.com/a": "bar",
					"foo.com/b": "bar",
				}
			},
			"",
		},

		{
			"invalid labels on app",
			func(c *Config) {
				c.Apps = append(c.Apps, &App{
					Labels: map[string]string{
						"waypoint/foo": "bar",
					},
				})
			},
			"reserved for system",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			c := TestConfig(t, testConfigValidSrc)
			tt.Modify(c)
			err := c.Validate()
			if tt.Err == "" {
				require.NoError(err)
				return
			}

			t.Logf(err.Error())
			require.Error(err)
			require.Contains(err.Error(), tt.Err)
		})
	}
}

var testConfigValidSrc = `
project = "test"
`
