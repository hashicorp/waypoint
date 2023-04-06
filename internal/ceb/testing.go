// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ceb

import (
	"context"
	"os"

	hzntest "github.com/hashicorp/horizon/pkg/testutils/central"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/pkg/server/singleprocess"
)

// TestCEB runs a CEB for test purposes. This will start an in-memory
// server automatically if connection information is not configured.
func TestCEB(t testing.T, opts ...Option) *TestCEBData {
	// Create a context that will stop our CEB
	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(func() { cancelFunc() })

	// Change the options to have some internal ones that spin up
	// a server automatically and set some exec args if necessary.
	// We also perform a copy here so if options are reused we don't
	// modify the caller.
	var data TestCEBData
	doneCh := make(chan struct{})
	opts = append(
		append([]Option{
			WithEnvDefaults(),
		}, opts...),
		withTestDefaults(t, &data, doneCh),
	)

	// Start the CEB
	go func() {
		err := Run(ctx, opts...)
		if err != nil {
			t.Logf("error running CEB: %s", err)
		}
	}()

	<-doneCh
	return &data
}

type TestCEBData struct {
	CEB     *CEB
	Horizon hzntest.DevSetup
}

func withTestDefaults(t testing.T, data *TestCEBData, doneCh chan<- struct{}) Option {
	return func(ceb *CEB, cfg *config) error {
		defer close(doneCh)

		// Store a reference to the CEB
		data.CEB = ceb

		// If we have no server configured, create that.
		if ceb.client == nil && cfg.ServerAddr == "" {
			if err := testCEBServer(t, ceb, cfg, data); err != nil {
				return err
			}
		}

		// If we have no exec specified, set that up to block.
		if len(cfg.ExecArgs) == 0 {
			cfg.ExecArgs = []string{"sleep", "10000"}
		}

		return nil
	}
}

func testCEBServer(t testing.T, ceb *CEB, cfg *config, data *TestCEBData) error {
	ceb.client = singleprocess.TestServer(t,
		singleprocess.TestWithURLService(t, &data.Horizon),
	)
	return nil
}

func testChenv(t testing.T, k, v string) {
	old := os.Getenv(k)
	require.NoError(t, os.Setenv(k, v))
	t.Cleanup(func() { os.Setenv(k, old) })
}
