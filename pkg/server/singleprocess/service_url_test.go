// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package singleprocess

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/serverconfig"
)

func TestServerURLService(t *testing.T) {
	t.Run("with URL service disabled", func(t *testing.T) {
		// Our default TestImpl doesn't have a URL service.
		impl := testServiceImpl(TestImpl(t))
		require.Nil(t, impl.urlClient())
	})

	t.Run("with URL service enabled", func(t *testing.T) {
		impl := testServiceImpl(TestImpl(t, TestWithURLService(t, nil)))
		require.NotNil(t, impl.urlClient())
	})

	t.Run("with URL service enabled, but down", func(t *testing.T) {
		ln, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		ln.Close()

		impl := testServiceImpl(TestImpl(t, WithConfig(&serverconfig.Config{
			URL: &serverconfig.URL{
				Enabled:              true,
				APIAddress:           ln.Addr().String(),
				APIInsecure:          true,
				APIToken:             "",
				ControlAddress:       fmt.Sprintf("dev://%s", ln.Addr().String()),
				AutomaticAppHostname: true,
			},
		})))

		require.Nil(t, impl.urlClient())
	})
}
