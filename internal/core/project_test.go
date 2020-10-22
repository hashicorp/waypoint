package core

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/config2"
)

func TestNewProject(t *testing.T) {
	require := require.New(t)

	p := TestProject(t,
		WithConfig(config.TestConfig(t, testNewProjectConfig)),
	)

	app, err := p.App("test")
	require.NoError(err)
	require.NotNil(app)
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
