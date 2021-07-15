package clicontext

import (
	"testing"

	"github.com/hashicorp/waypoint/internal/serverconfig"
	"github.com/stretchr/testify/require"
)

func TestConfigFromURL(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Expected Config
	}{
		{
			"host only",
			"foo.com",
			Config{
				Server: serverconfig.Client{
					Address:       "foo.com",
					Tls:           true,
					TlsSkipVerify: true,
				},
			},
		},

		{
			"host with port",
			"foo.com:1234",
			Config{
				Server: serverconfig.Client{
					Address:       "foo.com:1234",
					Tls:           true,
					TlsSkipVerify: true,
				},
			},
		},

		{
			"IP only",
			"127.1.2.3",
			Config{
				Server: serverconfig.Client{
					Address:       "127.1.2.3",
					Tls:           true,
					TlsSkipVerify: true,
				},
			},
		},

		{
			"IP with port",
			"127.0.0.1:9701",
			Config{
				Server: serverconfig.Client{
					Address:       "127.0.0.1:9701",
					Tls:           true,
					TlsSkipVerify: true,
				},
			},
		},

		{
			"http",
			"http://foo.com:1234",
			Config{
				Server: serverconfig.Client{
					Address:       "foo.com:1234",
					Tls:           false,
					TlsSkipVerify: true,
				},
			},
		},

		{
			"https with path",
			"https://foo.com:1234/whocares/about/this/at/all",
			Config{
				Server: serverconfig.Client{
					Address:       "foo.com:1234",
					Tls:           true,
					TlsSkipVerify: true,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			var actual Config
			err := actual.FromURL(tt.Input)
			require.NoError(err)
			require.Equal(actual, tt.Expected)
		})
	}
}
