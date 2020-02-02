package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mitchellh/devflow/internal/config"
)

func TestNewProject(t *testing.T) {
	require := require.New(t)

	// TODO(mitchellh): need something more robust
	cfg := &config.Config{
		Apps: []*config.App{
			&config.App{
				Name: "test",
				Build: &config.Component{
					Type: "pack",
				},
			},
		},
	}

	p, err := NewProject(context.Background(), WithConfig(cfg))
	require.NoError(err)

	app, err := p.App("test")
	require.NoError(err)
	require.NotNil(app.Builder)
}
