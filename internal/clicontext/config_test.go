// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clicontext

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/pkg/serverconfig"
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
					Address:       "foo.com:" + serverconfig.DefaultGRPCPort,
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
					Address:       "127.1.2.3:" + serverconfig.DefaultGRPCPort,
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
