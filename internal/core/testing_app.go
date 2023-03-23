// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
)

// TestApp returns the app named n in the project.
func TestApp(t testing.T, p *Project, n string) *App {
	app, err := p.App(n)
	require.NoError(t, err)
	return app
}
