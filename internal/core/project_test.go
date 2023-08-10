// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/config"
)

func TestNewProject(t *testing.T) {
	require := require.New(t)

	p := TestProject(t,
		WithConfig(config.TestConfig(t, testNewProjectConfig)),
	)

	// App that exists
	app, err := p.App("test")
	require.NoError(err)
	require.NotNil(app)

	// App that doesn't exist
	app, err = p.App("NO")
	require.Error(err)
	require.Nil(app)
	require.Equal(status.Code(err), codes.NotFound)
}

const testNewProjectConfig = `
project = "test"

app "test" {
	build {
		use "test" {}
	}

	deploy {
		use "test" {}
	}
}
`
